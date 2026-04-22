import 'location.dart';
import 'waypoint.dart';

class OptimizeRouteRequest {
  final String? routeName;
  final Location startLocation;
  final List<Waypoint> waypoints;
  final String? optimizationType; // "distance" or "duration"
  final DateTime? startTime; // When the route starts (for time window calculations)

  OptimizeRouteRequest({
    this.routeName,
    required this.startLocation,
    required this.waypoints,
    this.optimizationType,
    this.startTime,
  });

  Map<String, dynamic> toJson() {
    return {
      if (routeName != null) 'route_name': routeName,
      'start_location': startLocation.toJson(),
      'waypoints': waypoints.map((e) => e.toJson()).toList(),
      if (optimizationType != null) 'optimization_type': optimizationType,
      if (startTime != null) 'start_time': startTime!.toIso8601String(),
    };
  }
}





