/**
 * Utility untuk mendeteksi Fake GPS pada web browser
 * 
 * Teknik deteksi:
 * 1. Cek accuracy GPS (Fake GPS biasanya memiliki accuracy yang tidak realistis)
 * 2. Cek timestamp (Fake GPS mungkin memiliki timestamp yang tidak sinkron)
 * 3. Cek perubahan lokasi yang terlalu cepat (tidak mungkin secara fisik)
 * 4. Cek apakah location berubah terlalu drastis dalam waktu singkat
 */

export interface GPSPosition {
  latitude: number;
  longitude: number;
  accuracy?: number;
  timestamp: number;
}

export interface FakeGPSDetectionResult {
  isFakeGPS: boolean;
  reason?: string;
  confidence: "low" | "medium" | "high";
}

// Threshold untuk deteksi Fake GPS
const THRESHOLDS = {
  // Accuracy terlalu baik (kurang dari 1 meter biasanya tidak realistis untuk GPS biasa)
  MIN_REALISTIC_ACCURACY: 1,
  // Accuracy terlalu buruk (lebih dari 100 meter biasanya tidak realistis untuk GPS modern)
  MAX_REALISTIC_ACCURACY: 100,
  // Perubahan lokasi maksimal dalam 1 detik (dalam meter) - kecepatan manusia normal ~5 m/s
  MAX_SPEED_MPS: 50, // 50 m/s = 180 km/h (sangat cepat, tapi masih mungkin dengan kendaraan)
  // Perubahan lokasi maksimal dalam 1 detik untuk deteksi Fake GPS (dalam meter)
  FAKE_GPS_SPEED_THRESHOLD: 1000, // 1000 m/s = 3600 km/h (tidak mungkin secara fisik)
  // Timestamp tidak boleh lebih dari 5 detik dari waktu sekarang
  MAX_TIMESTAMP_DIFF_SECONDS: 5,
} as const;

// Cache untuk menyimpan posisi sebelumnya
let previousPositions: Array<{ position: GPSPosition; timestamp: number }> = [];

/**
 * Menghitung jarak antara dua koordinat menggunakan Haversine formula
 */
function calculateDistance(
  lat1: number,
  lon1: number,
  lat2: number,
  lon2: number
): number {
  const R = 6371000; // Radius bumi dalam meter
  const dLat = ((lat2 - lat1) * Math.PI) / 180;
  const dLon = ((lon2 - lon1) * Math.PI) / 180;
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos((lat1 * Math.PI) / 180) *
      Math.cos((lat2 * Math.PI) / 180) *
      Math.sin(dLon / 2) *
      Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
}

/**
 * Mendeteksi apakah GPS position adalah Fake GPS
 */
export function detectFakeGPS(
  position: GPSPosition,
  previousPosition?: GPSPosition
): FakeGPSDetectionResult {
  const now = Math.floor(Date.now() / 1000);
  const checks: Array<{ passed: boolean; reason: string; confidence: "low" | "medium" | "high" }> = [];

  // Check 1: Accuracy terlalu baik (kurang dari 1 meter biasanya tidak realistis)
  if (position.accuracy !== undefined && position.accuracy < THRESHOLDS.MIN_REALISTIC_ACCURACY) {
    checks.push({
      passed: false,
      reason: `GPS accuracy terlalu baik (${position.accuracy.toFixed(2)}m). Accuracy GPS biasa tidak bisa lebih baik dari 1 meter.`,
      confidence: "high",
    });
  }

  // Check 2: Accuracy terlalu buruk (lebih dari 100 meter biasanya tidak realistis untuk GPS modern)
  if (position.accuracy !== undefined && position.accuracy > THRESHOLDS.MAX_REALISTIC_ACCURACY) {
    checks.push({
      passed: false,
      reason: `GPS accuracy terlalu buruk (${position.accuracy.toFixed(2)}m). Accuracy GPS modern biasanya di bawah 100 meter.`,
      confidence: "medium",
    });
  }

  // Check 3: Timestamp tidak sinkron dengan waktu sekarang
  const timestampDiff = Math.abs(now - position.timestamp);
  if (timestampDiff > THRESHOLDS.MAX_TIMESTAMP_DIFF_SECONDS) {
    checks.push({
      passed: false,
      reason: `GPS timestamp tidak sinkron. Selisih waktu: ${timestampDiff} detik.`,
      confidence: "high",
    });
  }

  // Check 4: Perubahan lokasi terlalu cepat (jika ada previous position)
  if (previousPosition) {
    const distance = calculateDistance(
      previousPosition.latitude,
      previousPosition.longitude,
      position.latitude,
      position.longitude
    );
    const timeDiff = Math.abs(position.timestamp - previousPosition.timestamp);
    
    if (timeDiff > 0) {
      const speed = distance / timeDiff; // meter per second
      
      if (speed > THRESHOLDS.FAKE_GPS_SPEED_THRESHOLD) {
        checks.push({
          passed: false,
          reason: `Perubahan lokasi terlalu cepat (${speed.toFixed(2)} m/s = ${(speed * 3.6).toFixed(2)} km/h). Ini tidak mungkin secara fisik.`,
          confidence: "high",
        });
      } else if (speed > THRESHOLDS.MAX_SPEED_MPS) {
        checks.push({
          passed: false,
          reason: `Perubahan lokasi sangat cepat (${speed.toFixed(2)} m/s = ${(speed * 3.6).toFixed(2)} km/h). Mungkin menggunakan Fake GPS.`,
          confidence: "medium",
        });
      }
    }
  }

  // Check 5: Cek dengan posisi sebelumnya dari cache
  if (previousPositions.length > 0) {
    const recentPositions = previousPositions.filter(
      (p) => now - p.timestamp < 10 // Posisi dalam 10 detik terakhir
    );
    
    for (const prev of recentPositions) {
      const distance = calculateDistance(
        prev.position.latitude,
        prev.position.longitude,
        position.latitude,
        position.longitude
      );
      const timeDiff = Math.abs(position.timestamp - prev.timestamp);
      
      if (timeDiff > 0) {
        const speed = distance / timeDiff;
        
        if (speed > THRESHOLDS.FAKE_GPS_SPEED_THRESHOLD) {
          checks.push({
            passed: false,
            reason: `Perubahan lokasi terlalu cepat dari posisi sebelumnya (${speed.toFixed(2)} m/s).`,
            confidence: "high",
          });
        }
      }
    }
  }

  // Update cache (simpan maksimal 5 posisi terakhir)
  previousPositions.push({ position, timestamp: now });
  if (previousPositions.length > 5) {
    previousPositions.shift();
  }

  // Jika ada check yang gagal dengan confidence tinggi, return fake GPS
  const highConfidenceFailures = checks.filter(
    (c) => !c.passed && c.confidence === "high"
  );
  if (highConfidenceFailures.length > 0) {
    return {
      isFakeGPS: true,
      reason: highConfidenceFailures[0].reason,
      confidence: "high",
    };
  }

  // Jika ada 2+ check yang gagal dengan confidence medium, return fake GPS
  const mediumConfidenceFailures = checks.filter(
    (c) => !c.passed && c.confidence === "medium"
  );
  if (mediumConfidenceFailures.length >= 2) {
    return {
      isFakeGPS: true,
      reason: mediumConfidenceFailures[0].reason,
      confidence: "medium",
    };
  }

  // Jika hanya 1 check medium yang gagal, return suspicious tapi tidak pasti
  if (mediumConfidenceFailures.length === 1) {
    return {
      isFakeGPS: false,
      reason: mediumConfidenceFailures[0].reason,
      confidence: "low",
    };
  }

  // Semua check passed
  return {
    isFakeGPS: false,
    confidence: "low",
  };
}

/**
 * Mendeteksi Fake GPS dari GeolocationPosition browser
 */
export function detectFakeGPSFromPosition(
  position: GeolocationPosition,
  previousPosition?: GeolocationPosition
): FakeGPSDetectionResult {
  const gpsPosition: GPSPosition = {
    latitude: position.coords.latitude,
    longitude: position.coords.longitude,
    accuracy: position.coords.accuracy ?? undefined,
    timestamp: Math.floor(position.timestamp / 1000), // Convert to seconds
  };

  const prevGPSPosition: GPSPosition | undefined = previousPosition
    ? {
        latitude: previousPosition.coords.latitude,
        longitude: previousPosition.coords.longitude,
        accuracy: previousPosition.coords.accuracy ?? undefined,
        timestamp: Math.floor(previousPosition.timestamp / 1000),
      }
    : undefined;

  return detectFakeGPS(gpsPosition, prevGPSPosition);
}

/**
 * Clear cache posisi sebelumnya (untuk testing atau reset)
 */
export function clearPositionCache(): void {
  previousPositions = [];
}

