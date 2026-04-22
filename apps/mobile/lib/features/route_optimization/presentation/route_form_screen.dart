import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:geolocator/geolocator.dart';
import 'package:dio/dio.dart';

import '../../../core/l10n/app_localizations.dart';
import '../application/route_optimization_provider.dart';
import '../data/models/location.dart';
import '../data/models/optimize_route_request.dart';
import '../data/models/waypoint.dart';
import 'widgets/waypoint_selector_dialog.dart';
import 'widgets/route_header.dart';

class RouteFormScreen extends ConsumerStatefulWidget {
  const RouteFormScreen({super.key});

  @override
  ConsumerState<RouteFormScreen> createState() => _RouteFormScreenState();
}

class _RouteFormScreenState extends ConsumerState<RouteFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _routeNameController = TextEditingController();
  
  Location? _startLocation;
  final List<Waypoint> _waypoints = [];

  bool _isLoading = false;
  bool _isGettingLocation = false;

  @override
  void dispose() {
    _routeNameController.dispose();
    super.dispose();
  }

  /// Reverse geocode coordinates to get address using Nominatim (OpenStreetMap)
  /// Improved with better parameters and parsing for more accurate results
  Future<String?> _reverseGeocode(double lat, double lng) async {
    // Retry logic for better reliability
    const maxRetries = 2;
    for (int attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        // Use better parameters for more accurate results:
        // - zoom=18: street level accuracy
        // - addressdetails=1: get detailed address components
        // - accept-language: prefer Indonesian, fallback to English
        final uri = Uri.parse(
          'https://nominatim.openstreetmap.org/reverse?format=json&lat=$lat&lon=$lng&zoom=18&addressdetails=1',
        );

        final dio = Dio(BaseOptions(
          headers: {
            'User-Agent': 'CRM-Healthcare-Mobile/1.0',
            'Accept-Language': 'id,en',
          },
          connectTimeout: const Duration(seconds: 8),
          receiveTimeout: const Duration(seconds: 8),
        ));
        final response = await dio.getUri(uri);

        if (response.statusCode == 200) {
          final data = response.data as Map<String, dynamic>;
          
          // Check if there's an error in the response
          if (data.containsKey('error')) {
            debugPrint('Nominatim error: ${data['error']}');
            if (attempt < maxRetries) {
              await Future.delayed(Duration(seconds: attempt + 1));
              continue;
            }
            return null;
          }

          // Try to build a more accurate address from address components
          final address = data['address'] as Map<String, dynamic>?;
          if (address != null) {
            // Build address in order of specificity (most specific first)
            final addressParts = <String>[];
            
            // Road/Street name
            if (address['road'] != null) {
              addressParts.add(address['road'] as String);
            } else if (address['pedestrian'] != null) {
              addressParts.add(address['pedestrian'] as String);
            }
            
            // Suburb/Neighborhood
            if (address['suburb'] != null) {
              addressParts.add(address['suburb'] as String);
            } else if (address['neighbourhood'] != null) {
              addressParts.add(address['neighbourhood'] as String);
            }
            
            // Village/District
            if (address['village'] != null) {
              addressParts.add(address['village'] as String);
            } else if (address['city_district'] != null) {
              addressParts.add(address['city_district'] as String);
            }
            
            // City/Municipality
            if (address['city'] != null) {
              addressParts.add(address['city'] as String);
            } else if (address['municipality'] != null) {
              addressParts.add(address['municipality'] as String);
            } else if (address['town'] != null) {
              addressParts.add(address['town'] as String);
            }
            
            // State/Province
            if (address['state'] != null) {
              addressParts.add(address['state'] as String);
            } else if (address['province'] != null) {
              addressParts.add(address['province'] as String);
            }
            
            // Country
            if (address['country'] != null) {
              addressParts.add(address['country'] as String);
            }
            
            // If we have address parts, join them
            if (addressParts.isNotEmpty) {
              return addressParts.join(', ');
            }
          }
          
          // Fallback to display_name if address components parsing fails
          final displayName = data['display_name'] as String?;
          if (displayName != null && displayName.isNotEmpty) {
            return displayName;
          }
        } else if (response.statusCode == 429) {
          // Rate limited - wait longer before retry
          if (attempt < maxRetries) {
            await Future.delayed(Duration(seconds: (attempt + 1) * 2));
            continue;
          }
        } else if (response.statusCode != null && response.statusCode! >= 500) {
          // Server error - retry
          if (attempt < maxRetries) {
            await Future.delayed(Duration(seconds: attempt + 1));
            continue;
          }
        }
      } catch (e) {
        debugPrint('Reverse geocoding attempt ${attempt + 1} failed: $e');
        if (attempt < maxRetries) {
          await Future.delayed(Duration(seconds: attempt + 1));
          continue;
        }
      }
    }
    
    return null;
  }

  Future<void> _getCurrentLocation() async {
    setState(() => _isGettingLocation = true);

    try {
      bool serviceEnabled = await Geolocator.isLocationServiceEnabled();
      if (!serviceEnabled) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('Location services are disabled. Please enable location services.'),
            ),
          );
        }
        return;
      }

      LocationPermission permission = await Geolocator.checkPermission();
      if (permission == LocationPermission.denied) {
        permission = await Geolocator.requestPermission();
        if (permission == LocationPermission.denied) {
          if (mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('Location permissions are denied.'),
              ),
            );
          }
          return;
        }
      }

      if (permission == LocationPermission.deniedForever) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('Location permissions are permanently denied. Please enable in settings.'),
            ),
          );
        }
        return;
      }

      // BEST PRACTICE: Use bestForNavigation for most accurate GPS location
      // This is recommended for navigation, tracking, and precise location needs
      Position position = await Geolocator.getCurrentPosition(
        locationSettings: const LocationSettings(
          accuracy: LocationAccuracy.bestForNavigation, // BEST PRACTICE: Most accurate for navigation
          timeLimit: Duration(seconds: 20), // Allow time for GPS to get accurate fix
        ),
      );

      // Log GPS information for debugging
      debugPrint('=== GPS Position Obtained ===');
      debugPrint('  Latitude: ${position.latitude}');
      debugPrint('  Longitude: ${position.longitude}');
      debugPrint('  Accuracy: ${position.accuracy} meters');
      debugPrint('  Altitude: ${position.altitude}');
      debugPrint('  Speed: ${position.speed}');
      debugPrint('  Heading: ${position.heading}');
      debugPrint('  Timestamp: ${position.timestamp}');
      debugPrint('  Speed Accuracy: ${position.speedAccuracy}');
      debugPrint('============================');

      // Check if accuracy is acceptable (less than 20 meters is excellent for navigation)
      if (position.accuracy > 20) {
        debugPrint('⚠️ GPS accuracy is ${position.accuracy} meters');
        debugPrint('💡 Note: If using emulator, inject location manually via Extended Controls');
        debugPrint('💡 For real device: Wait a few seconds for GPS satellite fix');
        
        // For real devices, wait a bit more and try to get better accuracy
        // This helps GPS get a better satellite fix
        if (position.accuracy > 50) {
          debugPrint('⏳ Attempting to get better GPS fix...');
          await Future.delayed(const Duration(seconds: 3));
          
          try {
            final betterPosition = await Geolocator.getCurrentPosition(
              locationSettings: const LocationSettings(
                accuracy: LocationAccuracy.bestForNavigation,
                timeLimit: Duration(seconds: 15),
              ),
            );
            
            if (betterPosition.accuracy < position.accuracy) {
              debugPrint('✅ Got better position: ${betterPosition.accuracy} meters (was ${position.accuracy})');
              position = betterPosition;
            } else {
              debugPrint('ℹ️ Position accuracy unchanged: ${betterPosition.accuracy} meters');
            }
          } catch (e) {
            debugPrint('❌ Failed to get better position: $e');
            // Continue with original position
          }
        }
      } else {
        debugPrint('✅ GPS accuracy is excellent: ${position.accuracy} meters');
      }

      // Reverse geocode to get address
      final address = await _reverseGeocode(
        position.latitude,
        position.longitude,
      );

      setState(() {
        _startLocation = Location(
          lat: position.latitude,
          lng: position.longitude,
          address: address ?? 'Current Location',
        );
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to get location: $e'),
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isGettingLocation = false);
      }
    }
  }

  void _addWaypoint() {
    showDialog(
      context: context,
      builder: (context) => WaypointSelectorDialog(
        onSelect: (selectedWaypoints) {
          setState(() {
            _waypoints.addAll(selectedWaypoints);
          });
        },
      ),
    );
  }

  void _removeWaypoint(int index) {
    setState(() {
      _waypoints.removeAt(index);
    });
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;

    if (_startLocation == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please set start location'),
        ),
      );
      return;
    }

    if (_waypoints.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please add at least one waypoint'),
        ),
      );
      return;
    }

    setState(() => _isLoading = true);

    try {
      final request = OptimizeRouteRequest(
        routeName: _routeNameController.text.trim().isEmpty 
            ? null 
            : _routeNameController.text.trim(),
        startLocation: _startLocation!,
        waypoints: _waypoints,
      );

      await ref.read(routeOptimizationProvider.notifier).optimizeRoute(request);
      
      final optimizedRoute = ref.read(routeOptimizationProvider).optimizedRoute;
      
      if (mounted && optimizedRoute != null) {
        // Refresh route list
        ref.read(routeListProvider.notifier).refresh();
        
        // Navigate back
        Navigator.pop(context);
        
        // Show success message
        final l10n = AppLocalizations.of(context)!;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.routeOptimizedSuccessfully),
            backgroundColor: Colors.green,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        final l10n = AppLocalizations.of(context)!;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('${l10n.failedToOptimizeRoute}: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final optimizationState = ref.watch(routeOptimizationProvider);

    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            // Header
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 16, 16, 0),
              child: RouteHeader(
                title: l10n.createRoute,
                actions: [
                  if (_isLoading)
                    const Padding(
                      padding: EdgeInsets.all(8.0),
                      child: SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      ),
                    )
                  else
                    TextButton(
                      onPressed: _handleSubmit,
                      child: Text(l10n.optimize),
                    ),
                ],
              ),
            ),
            // Form Content
            Expanded(
              child: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Route Name
            TextFormField(
              controller: _routeNameController,
              decoration: InputDecoration(
                labelText: l10n.routeNameOptional,
                hintText: l10n.enterRouteName,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                filled: true,
                fillColor: theme.colorScheme.surfaceContainerHighest,
              ),
            ),
            const SizedBox(height: 24),

            // Start Location
            Text(
              l10n.startLocation,
              style: theme.textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Card(
              child: ListTile(
                leading: const Icon(Icons.location_on),
                title: Text(
                  _startLocation?.address ?? l10n.notSet,
                  style: theme.textTheme.bodyMedium?.copyWith(
                    fontWeight: FontWeight.w500,
                    color: _startLocation == null
                        ? theme.colorScheme.onSurface.withValues(alpha: 0.5)
                        : null,
                  ),
                ),
                subtitle: _startLocation != null
                    ? Text(
                        '${_startLocation!.lat.toStringAsFixed(6)}, ${_startLocation!.lng.toStringAsFixed(6)}',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
                        ),
                      )
                    : null,
                trailing: _isGettingLocation
                    ? const Padding(
                        padding: EdgeInsets.all(12.0),
                        child: SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        ),
                      )
                    : IconButton(
                        icon: const Icon(Icons.my_location),
                        onPressed: _getCurrentLocation,
                        tooltip: l10n.useCurrentLocation,
                      ),
              ),
            ),
            const SizedBox(height: 24),

            // Waypoints
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '${l10n.waypoints} (${_waypoints.length})',
                  style: theme.textTheme.titleMedium,
                ),
                ElevatedButton.icon(
                  onPressed: _addWaypoint,
                  icon: const Icon(Icons.add),
                  label: Text(l10n.add),
                ),
              ],
            ),
            const SizedBox(height: 8),
            if (_waypoints.isEmpty)
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(24.0),
                  child: Center(
                    child: Column(
                      children: [
                        Icon(
                          Icons.location_off,
                          size: 48,
                          color: theme.colorScheme.onSurface.withValues(alpha: 0.3),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          l10n.noWaypointsAdded,
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: theme.colorScheme.onSurface.withValues(alpha: 0.5),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              )
            else
              ...List.generate(_waypoints.length, (index) {
                final waypoint = _waypoints[index];
                return Card(
                  margin: const EdgeInsets.only(bottom: 8),
                  child: ListTile(
                    leading: CircleAvatar(
                      child: Text('${index + 1}'),
                    ),
                    title: Text(
                      waypoint.accountName ?? waypoint.address ?? '${l10n.destination} ${index + 1}',
                      style: theme.textTheme.bodyMedium?.copyWith(
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                    subtitle: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        if (waypoint.address != null)
                          Text(
                            waypoint.address!,
                            style: theme.textTheme.bodySmall,
                          ),
                        Text(
                          '${waypoint.lat.toStringAsFixed(6)}, ${waypoint.lng.toStringAsFixed(6)}',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
                          ),
                        ),
                      ],
                    ),
                    trailing: IconButton(
                      icon: const Icon(Icons.delete),
                      onPressed: () => _removeWaypoint(index),
                    ),
                  ),
                );
              }),

            if (optimizationState.errorMessage != null) ...[
              const SizedBox(height: 16),
              Card(
                color: theme.colorScheme.errorContainer,
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Text(
                    optimizationState.errorMessage!,
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.onErrorContainer,
                    ),
                  ),
                ),
              ),
            ],
          ],
        ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

