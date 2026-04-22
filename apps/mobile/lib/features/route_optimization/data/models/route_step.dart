import 'location.dart';

class RouteStep {
  final int step;
  final double distance;
  final String? distanceFormatted;
  final int duration;
  final String? durationFormatted;
  final String? instruction;
  final String? polyline;
  final String? maneuver;
  final Location startLocation;
  final Location endLocation;

  RouteStep({
    required this.step,
    required this.distance,
    this.distanceFormatted,
    required this.duration,
    this.durationFormatted,
    this.instruction,
    this.polyline,
    this.maneuver,
    required this.startLocation,
    required this.endLocation,
  });

  factory RouteStep.fromJson(Map<String, dynamic> json) {
    return RouteStep(
      step: json['step'] as int,
      distance: (json['distance'] as num).toDouble(),
      distanceFormatted: json['distance_formatted'] as String?,
      duration: json['duration'] as int,
      durationFormatted: json['duration_formatted'] as String?,
      instruction: json['instruction'] as String?,
      polyline: json['polyline'] as String?,
      maneuver: json['maneuver'] as String?,
      startLocation: Location.fromJson(
        json['start_location'] as Map<String, dynamic>,
      ),
      endLocation: Location.fromJson(
        json['end_location'] as Map<String, dynamic>,
      ),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'step': step,
      'distance': distance,
      if (distanceFormatted != null) 'distance_formatted': distanceFormatted,
      'duration': duration,
      if (durationFormatted != null) 'duration_formatted': durationFormatted,
      if (instruction != null) 'instruction': instruction,
      if (polyline != null) 'polyline': polyline,
      if (maneuver != null) 'maneuver': maneuver,
      'start_location': startLocation.toJson(),
      'end_location': endLocation.toJson(),
    };
  }
}





