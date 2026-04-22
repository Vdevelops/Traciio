import 'location.dart';

class CalculateDistanceRequest {
  final Location origin;
  final Location destination;

  CalculateDistanceRequest({
    required this.origin,
    required this.destination,
  });

  Map<String, dynamic> toJson() {
    return {
      'origin': origin.toJson(),
      'destination': destination.toJson(),
    };
  }
}

class CalculateDistanceResponse {
  final double distance;
  final String distanceFormatted;
  final int duration;
  final String durationFormatted;

  CalculateDistanceResponse({
    required this.distance,
    required this.distanceFormatted,
    required this.duration,
    required this.durationFormatted,
  });

  factory CalculateDistanceResponse.fromJson(Map<String, dynamic> json) {
    return CalculateDistanceResponse(
      distance: (json['distance'] as num).toDouble(),
      distanceFormatted: json['distance_formatted'] as String,
      duration: json['duration'] as int,
      durationFormatted: json['duration_formatted'] as String,
    );
  }
}





