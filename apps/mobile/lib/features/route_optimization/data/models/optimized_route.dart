import 'waypoint.dart';
import 'route_step.dart';

class OptimizedRoute {
  final String id;
  final String? routeName;
  final String userId;
  final List<Waypoint> waypoints;
  final List<int> optimizedOrder;
  final double? totalDistance;
  final String? totalDistanceFormatted;
  final int? totalDuration;
  final String? totalDurationFormatted;
  final String? routePolyline;
  final List<RouteStep>? routeSteps;
  final DateTime createdAt;
  final DateTime updatedAt;
  final Map<String, dynamic>? user;

  OptimizedRoute({
    required this.id,
    this.routeName,
    required this.userId,
    required this.waypoints,
    required this.optimizedOrder,
    this.totalDistance,
    this.totalDistanceFormatted,
    this.totalDuration,
    this.totalDurationFormatted,
    this.routePolyline,
    this.routeSteps,
    required this.createdAt,
    required this.updatedAt,
    this.user,
  });

  factory OptimizedRoute.fromJson(Map<String, dynamic> json) {
    return OptimizedRoute(
      id: json['id'] as String,
      routeName: json['route_name'] as String?,
      userId: json['user_id'] as String,
      waypoints: (json['waypoints'] as List<dynamic>)
          .map((e) => Waypoint.fromJson(e as Map<String, dynamic>))
          .toList(),
      optimizedOrder: (json['optimized_order'] as List<dynamic>)
          .map((e) => e as int)
          .toList(),
      totalDistance: json['total_distance'] != null
          ? (json['total_distance'] as num).toDouble()
          : null,
      totalDistanceFormatted: json['total_distance_formatted'] as String?,
      totalDuration: json['total_duration'] as int?,
      totalDurationFormatted: json['total_duration_formatted'] as String?,
      routePolyline: json['route_polyline'] as String?,
      routeSteps: json['route_steps'] != null
          ? (json['route_steps'] as List<dynamic>)
              .map((e) => RouteStep.fromJson(e as Map<String, dynamic>))
              .toList()
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
      user: json['user'] as Map<String, dynamic>?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      if (routeName != null) 'route_name': routeName,
      'user_id': userId,
      'waypoints': waypoints.map((e) => e.toJson()).toList(),
      'optimized_order': optimizedOrder,
      if (totalDistance != null) 'total_distance': totalDistance,
      if (totalDistanceFormatted != null)
        'total_distance_formatted': totalDistanceFormatted,
      if (totalDuration != null) 'total_duration': totalDuration,
      if (totalDurationFormatted != null)
        'total_duration_formatted': totalDurationFormatted,
      if (routePolyline != null) 'route_polyline': routePolyline,
      if (routeSteps != null)
        'route_steps': routeSteps!.map((e) => e.toJson()).toList(),
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
      if (user != null) 'user': user,
    };
  }
}





