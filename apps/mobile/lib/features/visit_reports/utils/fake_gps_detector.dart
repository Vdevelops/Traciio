import 'dart:math' as math;
import 'package:geolocator/geolocator.dart';

/// Utility untuk mendeteksi Fake GPS pada mobile app
/// 
/// Teknik deteksi:
/// 1. Cek accuracy GPS (Fake GPS biasanya memiliki accuracy yang tidak realistis)
/// 2. Cek timestamp (Fake GPS mungkin memiliki timestamp yang tidak sinkron)
/// 3. Cek perubahan lokasi yang terlalu cepat (tidak mungkin secara fisik)
/// 4. Cek apakah location berubah terlalu drastis dalam waktu singkat
/// 5. Cek apakah isMockLocation (Android) - jika tersedia

/// Threshold untuk deteksi Fake GPS
class _Thresholds {
  // Accuracy terlalu baik (kurang dari 1 meter biasanya tidak realistis untuk GPS biasa)
  static const double minRealisticAccuracy = 1.0;
  // Accuracy terlalu buruk (lebih dari 100 meter biasanya tidak realistis untuk GPS modern)
  static const double maxRealisticAccuracy = 100.0;
  // Perubahan lokasi maksimal dalam 1 detik (dalam meter) - kecepatan manusia normal ~5 m/s
  static const double maxSpeedMps = 50.0; // 50 m/s = 180 km/h (sangat cepat, tapi masih mungkin dengan kendaraan)
  // Perubahan lokasi maksimal dalam 1 detik untuk deteksi Fake GPS (dalam meter)
  static const double fakeGPSSpeedThreshold = 1000.0; // 1000 m/s = 3600 km/h (tidak mungkin secara fisik)
  // Timestamp tidak boleh lebih dari 5 detik dari waktu sekarang
  static const int maxTimestampDiffSeconds = 5;
}

/// Hasil deteksi Fake GPS
class FakeGPSDetectionResult {
  final bool isFakeGPS;
  final String? reason;
  final String confidence; // "low", "medium", "high"

  FakeGPSDetectionResult({
    required this.isFakeGPS,
    this.reason,
    required this.confidence,
  });
}

/// Cache untuk menyimpan posisi sebelumnya
final List<_CachedPosition> _previousPositions = [];

class _CachedPosition {
  final Position position;
  final int timestamp;

  _CachedPosition({
    required this.position,
    required this.timestamp,
  });
}

/// Menghitung jarak antara dua koordinat menggunakan Haversine formula
double _calculateDistance(
  double lat1,
  double lon1,
  double lat2,
  double lon2,
) {
  const double R = 6371000; // Radius bumi dalam meter
  final double dLat = (lat2 - lat1) * math.pi / 180;
  final double dLon = (lon2 - lon1) * math.pi / 180;
  final double a = math.sin(dLat / 2) * math.sin(dLat / 2) +
      math.cos(lat1 * math.pi / 180) *
          math.cos(lat2 * math.pi / 180) *
          math.sin(dLon / 2) *
          math.sin(dLon / 2);
  final double c = 2 * math.atan2(math.sqrt(a), math.sqrt(1 - a));
  return R * c;
}

/// Mendeteksi apakah GPS position adalah Fake GPS
FakeGPSDetectionResult detectFakeGPS(
  Position position, {
  Position? previousPosition,
}) {
  final now = DateTime.now().millisecondsSinceEpoch ~/ 1000;
  final checks = <_CheckResult>[];

  // Check 1: isMockLocation (Android) - jika tersedia
  // Note: isMockLocation mungkin tidak tersedia di semua platform
  try {
    // Cek apakah position memiliki property isMockLocation
    // Karena Position dari geolocator tidak memiliki isMockLocation secara langsung,
    // kita perlu menggunakan teknik lain
    // Untuk Android, kita bisa cek melalui platform channel, tapi untuk sekarang
    // kita fokus pada teknik deteksi lainnya
  } catch (e) {
    // Ignore jika isMockLocation tidak tersedia
  }

  // Check 2: Accuracy terlalu baik (kurang dari 1 meter biasanya tidak realistis)
  if (position.accuracy > 0 && position.accuracy < _Thresholds.minRealisticAccuracy) {
    checks.add(_CheckResult(
      passed: false,
      reason:
          'GPS accuracy terlalu baik (${position.accuracy.toStringAsFixed(2)}m). Accuracy GPS biasa tidak bisa lebih baik dari 1 meter.',
      confidence: 'high',
    ));
  }

  // Check 3: Accuracy terlalu buruk (lebih dari 100 meter biasanya tidak realistis untuk GPS modern)
  if (position.accuracy > _Thresholds.maxRealisticAccuracy) {
    checks.add(_CheckResult(
      passed: false,
      reason:
          'GPS accuracy terlalu buruk (${position.accuracy.toStringAsFixed(2)}m). Accuracy GPS modern biasanya di bawah 100 meter.',
      confidence: 'medium',
    ));
  }

  // Check 4: Timestamp tidak sinkron dengan waktu sekarang
  final positionTimestamp = position.timestamp.millisecondsSinceEpoch ~/ 1000;
  final timestampDiff = (now - positionTimestamp).abs();
  if (timestampDiff > _Thresholds.maxTimestampDiffSeconds) {
    checks.add(_CheckResult(
      passed: false,
      reason: 'GPS timestamp tidak sinkron. Selisih waktu: $timestampDiff detik.',
      confidence: 'high',
    ));
  }

  // Check 5: Perubahan lokasi terlalu cepat (jika ada previous position)
  if (previousPosition != null) {
    final distance = _calculateDistance(
      previousPosition.latitude,
      previousPosition.longitude,
      position.latitude,
      position.longitude,
    );
    final timeDiff = (positionTimestamp -
            (previousPosition.timestamp.millisecondsSinceEpoch ~/ 1000))
        .abs();

    if (timeDiff > 0) {
      final speed = distance / timeDiff; // meter per second

      if (speed > _Thresholds.fakeGPSSpeedThreshold) {
        checks.add(_CheckResult(
          passed: false,
          reason:
              'Perubahan lokasi terlalu cepat (${speed.toStringAsFixed(2)} m/s = ${(speed * 3.6).toStringAsFixed(2)} km/h). Ini tidak mungkin secara fisik.',
          confidence: 'high',
        ));
      } else if (speed > _Thresholds.maxSpeedMps) {
        checks.add(_CheckResult(
          passed: false,
          reason:
              'Perubahan lokasi sangat cepat (${speed.toStringAsFixed(2)} m/s = ${(speed * 3.6).toStringAsFixed(2)} km/h). Mungkin menggunakan Fake GPS.',
          confidence: 'medium',
        ));
      }
    }
  }

  // Check 6: Cek dengan posisi sebelumnya dari cache
  if (_previousPositions.isNotEmpty) {
    final recentPositions = _previousPositions
        .where((p) => (now - p.timestamp).abs() < 10) // Posisi dalam 10 detik terakhir
        .toList();

    for (final prev in recentPositions) {
      final distance = _calculateDistance(
        prev.position.latitude,
        prev.position.longitude,
        position.latitude,
        position.longitude,
      );
      final timeDiff = (positionTimestamp - prev.timestamp).abs();

      if (timeDiff > 0) {
        final speed = distance / timeDiff;

        if (speed > _Thresholds.fakeGPSSpeedThreshold) {
          checks.add(_CheckResult(
            passed: false,
            reason:
                'Perubahan lokasi terlalu cepat dari posisi sebelumnya (${speed.toStringAsFixed(2)} m/s).',
            confidence: 'high',
          ));
        }
      }
    }
  }

  // Update cache (simpan maksimal 5 posisi terakhir)
  _previousPositions.add(_CachedPosition(
    position: position,
    timestamp: positionTimestamp,
  ));
  if (_previousPositions.length > 5) {
    _previousPositions.removeAt(0);
  }

  // Jika ada check yang gagal dengan confidence tinggi, return fake GPS
  final highConfidenceFailures =
      checks.where((c) => !c.passed && c.confidence == 'high').toList();
  if (highConfidenceFailures.isNotEmpty) {
    return FakeGPSDetectionResult(
      isFakeGPS: true,
      reason: highConfidenceFailures.first.reason,
      confidence: 'high',
    );
  }

  // Jika ada 2+ check yang gagal dengan confidence medium, return fake GPS
  final mediumConfidenceFailures =
      checks.where((c) => !c.passed && c.confidence == 'medium').toList();
  if (mediumConfidenceFailures.length >= 2) {
    return FakeGPSDetectionResult(
      isFakeGPS: true,
      reason: mediumConfidenceFailures.first.reason,
      confidence: 'medium',
    );
  }

  // Jika hanya 1 check medium yang gagal, return suspicious tapi tidak pasti
  if (mediumConfidenceFailures.length == 1) {
    return FakeGPSDetectionResult(
      isFakeGPS: false,
      reason: mediumConfidenceFailures.first.reason,
      confidence: 'low',
    );
  }

  // Semua check passed
  return FakeGPSDetectionResult(
    isFakeGPS: false,
    confidence: 'low',
  );
}

/// Clear cache posisi sebelumnya (untuk testing atau reset)
void clearPositionCache() {
  _previousPositions.clear();
}

class _CheckResult {
  final bool passed;
  final String reason;
  final String confidence;

  _CheckResult({
    required this.passed,
    required this.reason,
    required this.confidence,
  });
}

