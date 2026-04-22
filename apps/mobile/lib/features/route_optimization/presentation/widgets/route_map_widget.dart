import 'package:flutter/material.dart';
import 'package:flutter_map/flutter_map.dart';
import 'package:latlong2/latlong.dart';

import '../../data/models/optimized_route.dart';
import '../../data/models/waypoint.dart';

// Helper to calculate bounds from waypoints
LatLngBounds _calculateBounds(List<Waypoint> waypoints) {
  if (waypoints.isEmpty) {
    return LatLngBounds(
      const LatLng(-6.2, 106.8), // Default Jakarta bounds
      const LatLng(-6.3, 106.9),
    );
  }

  double minLat = waypoints[0].lat;
  double maxLat = waypoints[0].lat;
  double minLng = waypoints[0].lng;
  double maxLng = waypoints[0].lng;

  for (final waypoint in waypoints) {
    if (waypoint.lat < minLat) minLat = waypoint.lat;
    if (waypoint.lat > maxLat) maxLat = waypoint.lat;
    if (waypoint.lng < minLng) minLng = waypoint.lng;
    if (waypoint.lng > maxLng) maxLng = waypoint.lng;
  }

  // Add padding
  final latPadding = (maxLat - minLat) * 0.1;
  final lngPadding = (maxLng - minLng) * 0.1;

  return LatLngBounds(
    LatLng(minLat - latPadding, minLng - lngPadding),
    LatLng(maxLat + latPadding, maxLng + lngPadding),
  );
}

/// Decode polyline string (Google Polyline Algorithm / OSRM format)
/// This decodes an encoded polyline string into a list of LatLng coordinates
List<LatLng> decodePolyline(String encoded) {
  final List<LatLng> points = [];

  if (encoded.isEmpty) {
    return points;
  }

  int index = 0;
  int lat = 0;
  int lng = 0;

  try {
    while (index < encoded.length) {
      // Decode latitude
      int shift = 0;
      int result = 0;
      int byte;
      do {
        if (index >= encoded.length) break;
        byte = encoded.codeUnitAt(index++) - 63;
        result |= (byte & 0x1F) << shift;
        shift += 5;
      } while (byte >= 0x20 && index < encoded.length);

      final int deltaLat = ((result & 1) != 0) ? ~(result >> 1) : (result >> 1);
      lat += deltaLat;

      // Decode longitude
      shift = 0;
      result = 0;
      do {
        if (index >= encoded.length) break;
        byte = encoded.codeUnitAt(index++) - 63;
        result |= (byte & 0x1F) << shift;
        shift += 5;
      } while (byte >= 0x20 && index < encoded.length);

      final int deltaLng = ((result & 1) != 0) ? ~(result >> 1) : (result >> 1);
      lng += deltaLng;

      // Add point (divide by 1e5 to get actual coordinates)
      points.add(LatLng(lat / 1e5, lng / 1e5));
    }
  } catch (e) {
    // If decoding fails, return empty list (will trigger fallback)
    debugPrint('Error decoding polyline: $e');
    return [];
  }

  return points;
}

class RouteMapWidget extends StatefulWidget {
  const RouteMapWidget({super.key, required this.route, this.height = 300});

  final OptimizedRoute route;
  final double height;

  @override
  State<RouteMapWidget> createState() => _RouteMapWidgetState();
}

class _RouteMapWidgetState extends State<RouteMapWidget> {
  final MapController _mapController = MapController();
  bool _mapReady = false;

  @override
  void initState() {
    super.initState();
    // Fit bounds after a short delay to ensure map is initialized
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted && widget.route.waypoints.isNotEmpty) {
        _fitBounds();
      }
    });
  }

  void _fitBounds() {
    if (!_mapReady || widget.route.waypoints.isEmpty) return;

    final bounds = _calculateBounds(widget.route.waypoints);
    _mapController.fitCamera(
      CameraFit.bounds(bounds: bounds, padding: const EdgeInsets.all(50)),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    if (widget.route.waypoints.isEmpty) {
      return Container(
        height: widget.height,
        decoration: BoxDecoration(
          color: theme.colorScheme.surfaceContainerHighest,
          borderRadius: BorderRadius.circular(12),
        ),
        child: Center(
          child: Text(
            'No waypoints available',
            style: theme.textTheme.bodyMedium?.copyWith(
              color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
            ),
          ),
        ),
      );
    }

    // Calculate center point (fallback if bounds fitting fails)
    final center = LatLng(
      widget.route.waypoints[0].lat,
      widget.route.waypoints[0].lng,
    );

    // Build polyline positions
    List<LatLng> polylinePositions = [];

    // Try to use route_polyline first (main route polyline from OSRM)
    if (widget.route.routePolyline != null &&
        widget.route.routePolyline!.isNotEmpty) {
      // Decode polyline from OSRM
      polylinePositions = decodePolyline(widget.route.routePolyline!);

      // Debug: Log polyline info
      debugPrint(
        'Route polyline length: ${widget.route.routePolyline!.length}',
      );
      debugPrint('Decoded points count: ${polylinePositions.length}');

      // If decoding failed or returned empty, try route steps
      if (polylinePositions.isEmpty) {
        debugPrint('Polyline decoding returned empty, trying route steps');
        polylinePositions = _buildPolylineFromRouteSteps();

        // If route steps also empty, fallback to straight line
        if (polylinePositions.isEmpty) {
          debugPrint('Route steps also empty, using straight line fallback');
          polylinePositions = _buildStraightLinePolyline();
        }
      } else {
        debugPrint(
          'Using OSRM polyline with ${polylinePositions.length} points',
        );
      }
    } else if (widget.route.routeSteps != null &&
        widget.route.routeSteps!.isNotEmpty) {
      // If no route_polyline, try to build from route steps
      debugPrint('No route polyline, building from route steps');
      polylinePositions = _buildPolylineFromRouteSteps();

      if (polylinePositions.isEmpty) {
        debugPrint('Route steps polyline empty, using straight line fallback');
        polylinePositions = _buildStraightLinePolyline();
      }
    } else {
      debugPrint(
        'No route polyline or steps available, using straight line fallback',
      );
      // Fallback to straight line
      polylinePositions = _buildStraightLinePolyline();
    }

    // Calculate bounds for fitting
    final bounds = _calculateBounds(widget.route.waypoints);

    return Container(
      height: widget.height,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 3,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(12),
        child: Stack(
          children: [
            FlutterMap(
              mapController: _mapController,
              options: MapOptions(
                initialCenter: center,
                initialZoom: 13.0,
                minZoom: 5.0,
                maxZoom: 18.0,
                interactionOptions: const InteractionOptions(
                  flags: InteractiveFlag.all,
                ),
                onMapReady: () {
                  setState(() {
                    _mapReady = true;
                  });
                  // Fit bounds after map is ready
                  _fitBounds();
                },
                cameraConstraint: CameraConstraint.contain(bounds: bounds),
              ),
              children: [
                // Tile Layer - Using CartoDB Positron (more reliable for emulator)
                // Fallback to Stamen Toner if CartoDB fails
                TileLayer(
                  urlTemplate:
                      'https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png',
                  subdomains: const ['a', 'b', 'c', 'd'],
                  userAgentPackageName: 'com.crm.healthcare.mobile',
                  maxZoom: 19,
                  errorTileCallback: (tile, error, stackTrace) {
                    debugPrint('Tile error for ${tile.coordinates}: $error');
                  },
                ),

                // Route Polyline
                if (polylinePositions.length > 1)
                  PolylineLayer(
                    polylines: [
                      // Shadow/glow effect
                      Polyline(
                        points: polylinePositions,
                        strokeWidth: 8,
                        color: Colors.blue.withValues(alpha: 0.2),
                        borderStrokeWidth: 0,
                      ),
                      // Main route line
                      Polyline(
                        points: polylinePositions,
                        strokeWidth: 5,
                        color: const Color(0xFF3b82f6),
                        borderStrokeWidth: 0,
                      ),
                    ],
                  ),

                // Markers
                MarkerLayer(markers: _buildMarkers(theme)),
              ],
            ),

            // Route Info Overlay (top-left)
            Positioned(
              top: 12,
              left: 12,
              child: _buildRouteInfoOverlay(context, theme),
            ),

            // Error overlay if map fails to load
            if (!_mapReady)
              Positioned.fill(
                child: Container(
                  color: theme.colorScheme.surfaceContainerHighest.withValues(
                    alpha: 0.9,
                  ),
                  child: Center(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const CircularProgressIndicator(),
                        const SizedBox(height: 16),
                        Text(
                          'Loading map...',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }

  /// Build polyline from route steps (if available)
  List<LatLng> _buildPolylineFromRouteSteps() {
    final List<LatLng> positions = [];

    if (widget.route.routeSteps == null || widget.route.routeSteps!.isEmpty) {
      debugPrint('No route steps available for polyline');
      return positions;
    }

    debugPrint(
      'Building polyline from ${widget.route.routeSteps!.length} route steps',
    );

    // Decode polyline from each route step and combine them
    for (final step in widget.route.routeSteps!) {
      if (step.polyline != null && step.polyline!.isNotEmpty) {
        debugPrint(
          'Decoding polyline from step ${step.step}, length: ${step.polyline!.length}',
        );
        final stepPolyline = decodePolyline(step.polyline!);
        if (stepPolyline.isNotEmpty) {
          debugPrint(
            'Decoded ${stepPolyline.length} points from step ${step.step}',
          );
          // Add points from this step (avoid duplicates)
          for (final point in stepPolyline) {
            if (positions.isEmpty ||
                positions.last.latitude != point.latitude ||
                positions.last.longitude != point.longitude) {
              positions.add(point);
            }
          }
        } else {
          debugPrint('Failed to decode polyline from step ${step.step}');
        }
      } else {
        debugPrint(
          'Step ${step.step} has no polyline, using start/end locations',
        );
        // If step has no polyline, use start and end locations
        positions.add(LatLng(step.startLocation.lat, step.startLocation.lng));
      }
    }

    // Add last end location if we have steps
    if (widget.route.routeSteps!.isNotEmpty && positions.isNotEmpty) {
      final lastStep = widget.route.routeSteps!.last;
      final lastPoint = LatLng(
        lastStep.endLocation.lat,
        lastStep.endLocation.lng,
      );
      // Only add if different from last position
      if (positions.last.latitude != lastPoint.latitude ||
          positions.last.longitude != lastPoint.longitude) {
        positions.add(lastPoint);
      }
    }

    debugPrint(
      'Built polyline with ${positions.length} points from route steps',
    );
    return positions;
  }

  List<LatLng> _buildStraightLinePolyline() {
    final List<LatLng> positions = [];

    // Backend returns waypoints already in optimized order:
    // [start(order=0), optimized_stop_1(order=1), optimized_stop_2(order=2), ...]
    // Simply iterate them in their natural array order.
    for (int i = 0; i < widget.route.waypoints.length; i++) {
      positions.add(
        LatLng(widget.route.waypoints[i].lat, widget.route.waypoints[i].lng),
      );
    }

    return positions;
  }

  List<Marker> _buildMarkers(ThemeData theme) {
    final List<Marker> markers = [];

    for (int i = 0; i < widget.route.waypoints.length; i++) {
      final waypoint = widget.route.waypoints[i];
      final isStart = waypoint.order == 0 || i == 0;

      markers.add(
        Marker(
          point: LatLng(waypoint.lat, waypoint.lng),
          width: isStart ? 50 : 40,
          height: isStart ? 50 : 40,
          child: _buildMarkerWidget(isStart, i, waypoint, theme),
        ),
      );
    }

    return markers;
  }

  Widget _buildMarkerWidget(
    bool isStart,
    int index,
    Waypoint waypoint,
    ThemeData theme,
  ) {
    if (isStart) {
      // Start marker with gradient and rocket icon
      return Container(
        decoration: BoxDecoration(
          gradient: const LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF10b981), Color(0xFF059669)],
          ),
          shape: BoxShape.circle,
          border: Border.all(color: Colors.white, width: 3),
          boxShadow: [
            BoxShadow(
              color: const Color(0xFF10b981).withValues(alpha: 0.2),
              blurRadius: 3,
              offset: const Offset(0, 1),
            ),
          ],
        ),
        child: const Center(child: Text('🚀', style: TextStyle(fontSize: 20))),
      );
    } else {
      // Numbered destination marker with gradient
      return Container(
        decoration: BoxDecoration(
          gradient: const LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF667eea), Color(0xFF764ba2)],
          ),
          shape: BoxShape.circle,
          border: Border.all(color: Colors.white, width: 3),
          boxShadow: [
            BoxShadow(
              color: const Color(0xFF667eea).withValues(alpha: 0.2),
              blurRadius: 3,
              offset: const Offset(0, 1),
            ),
          ],
        ),
        child: Center(
          child: Text(
            '$index',
            style: const TextStyle(
              color: Colors.white,
              fontWeight: FontWeight.bold,
              fontSize: 14,
            ),
          ),
        ),
      );
    }
  }

  Widget _buildRouteInfoOverlay(BuildContext context, ThemeData theme) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: theme.brightness == Brightness.dark
            ? Colors.grey[900]!
            : Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: theme.colorScheme.outline.withValues(alpha: 0.2),
          width: 2,
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 3,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            children: [
              Icon(Icons.route, size: 16, color: theme.colorScheme.primary),
              const SizedBox(width: 4),
              Text(
                'Route Information',
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: theme.colorScheme.onSurface,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          _buildInfoRow(
            context,
            theme,
            Icons.straighten,
            'Distance',
            widget.route.totalDistanceFormatted ?? 'N/A',
            Colors.blue,
          ),
          const SizedBox(height: 4),
          _buildInfoRow(
            context,
            theme,
            Icons.access_time,
            'Duration',
            widget.route.totalDurationFormatted ?? 'N/A',
            Colors.green,
          ),
          const SizedBox(height: 4),
          _buildInfoRow(
            context,
            theme,
            Icons.location_on,
            'Stops',
            '${widget.route.waypoints.length - 1}',
            Colors.purple,
          ),
        ],
      ),
    );
  }

  Widget _buildInfoRow(
    BuildContext context,
    ThemeData theme,
    IconData icon,
    String label,
    String value,
    Color iconColor,
  ) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 14, color: iconColor),
        const SizedBox(width: 6),
        Text(
          '$label: ',
          style: theme.textTheme.bodySmall?.copyWith(
            fontWeight: FontWeight.w600,
            color: theme.colorScheme.onSurface.withValues(alpha: 0.7),
          ),
        ),
        Flexible(
          child: Text(
            value,
            style: theme.textTheme.bodySmall?.copyWith(
              fontWeight: FontWeight.bold,
              color: theme.colorScheme.onSurface,
            ),
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );
  }
}
