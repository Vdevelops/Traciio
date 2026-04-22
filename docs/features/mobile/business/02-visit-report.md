# Business - Visit Report Management

## CRM Healthcare Mobile App - Flutter

**Module**: Business Domain  
**Sprint**: Sprint 2  
**Version**: 1.0  
**Status**: ✅ **Completed**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Ringkasan Fitur](#ringkasan-fitur)
2. [Fitur Utama](#fitur-utama)
3. [Business Rules](#business-rules)
4. [Keputusan Teknis & Trade-offs](#keputusan-teknis--trade-offs)
5. [Struktur Folder](#struktur-folder)
6. [API Endpoints](#api-endpoints)
7. [Data Models](#data-models)
8. [Configuration](#configuration)
9. [Usage Examples](#usage-examples)
10. [Cara Test Manual](#cara-test-manual)
11. [Dependencies](#dependencies)
12. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Fitur **Visit Report Management** memungkinkan sales rep untuk mencatat kunjungan ke accounts (hospitals, clinics, pharmacies) dengan GPS tracking, photo documentation, dan status workflow. Fitur ini adalah core functionality untuk tracking sales activities dan memastikan accountability.

### Goals

- **Visit Documentation**: Catat detail kunjungan dengan notes dan photos
- **GPS Tracking**: Check-in dan check-out dengan location validation
- **Photo Upload**: Dokumentasikan kunjungan dengan foto
- **Workflow Management**: Status tracking dari draft sampai approval
- **Offline Support**: Create visit reports offline dengan sync

---

## Fitur Utama

### 1. Visit Report List

**Features**:

- List visit reports dengan pagination
- Filter by status (draft, submitted, approved, rejected)
- Filter by date range
- Search by account name
- Pull-to-refresh
- Offline indicator

**Status Badges**:

- 🟡 Draft: Belum di-submit
- 🔵 Submitted: Menunggu approval
- 🟢 Approved: Disetujui supervisor
- 🔴 Rejected: Ditolak dengan reason

### 2. Visit Report Creation

**Form Fields**:

- Account Selection (dropdown/search)
- Visit Date & Time
- Notes/Description
- Photo Upload (multiple, max 5)
- GPS Check-in (automatic)

**Workflow**:

1. Sales rep create visit report
2. Check-in dengan GPS (wajib)
3. Upload photos dan notes
4. Check-out dengan GPS (wajib)
5. Submit untuk approval

### 3. Visit Report Detail

**Information Displayed**:

- Account info
- Visit timeline (check-in → activities → check-out)
- Photos gallery
- Status dan approval info
- Supervisor comments
- Actions (edit jika draft, view jika submitted)

### 4. GPS Check-in/out

**Features**:

- Get current location dengan geolocator
- Validate location proximity ke account address
- Timestamp capture
- Location accuracy check
- Offline location caching

**Location Validation**:

- Minimum accuracy: 50 meters
- Distance threshold: 500 meters dari account location (optional)
- Retry mechanism untuk poor GPS signal

### 5. Photo Upload

**Features**:

- Camera integration (image_picker)
- Multiple photo selection (max 5)
- Photo preview dengan carousel
- Photo compression sebelum upload
- Offline photo storage

**Photo Requirements**:

- Max size: 5MB per photo
- Format: JPG, PNG
- Min resolution: 640x480
- Compression: 80% quality

### 6. Approval Workflow

**Status Flow**:

```
Draft → Submitted → Approved
            ↓
        Rejected → Draft (dengan edits)
```

**Supervisor Actions**:

- View submitted reports
- Approve dengan comment
- Reject dengan reason
- Request additional info

---

## Business Rules

### 1. Visit Report Creation Rules

**Required Fields**:

- Account ID
- Visit Date
- Check-in Location
- Check-out Location (saat submit)

**Optional Fields**:

- Notes
- Photos (min 1, max 5)

### 2. GPS Rules

**Check-in**:

- Wajib dilakukan saat tiba di location
- GPS accuracy harus < 50 meters
- Tidak bisa check-in jika terlalu jauh dari account (> 1km, warning)

**Check-out**:

- Wajib dilakukan saat meninggalkan location
- Check-out location harus sama atau dekat dengan check-in
- Duration minimum: 5 menit

**Offline**:

- Cache GPS coordinates saat offline
- Validate saat connection restored
- Allow manual coordinate input dengan alasan

### 3. Photo Rules

**Upload Requirements**:

- Min 1 photo, max 5 photos
- Max file size: 5MB per photo
- Compression required sebelum upload
- Photo metadata (timestamp, location) preserved

**Offline**:

- Save photos locally saat offline
- Auto-upload saat connection restored
- Queue untuk multiple photo uploads

### 4. Approval Rules

**Sales Rep**:

- Hanya bisa edit saat status Draft
- Tidak bisa edit setelah Submitted
- Dapat melihat rejection reason

**Supervisor**:

- Dapat approve/reject submitted reports
- Wajib memberikan comment saat reject
- Dapat view semua reports di team-nya

**Status Changes**:

- Draft → Submitted: oleh sales rep
- Submitted → Approved: oleh supervisor
- Submitted → Rejected: oleh supervisor (dengan reason)
- Rejected → Draft: oleh sales rep (setelah edits)

### 5. Offline Rules

**Create Offline**:

- Create visit report sebagai draft
- Save locally dengan pending sync flag
- Auto-sync saat online

**Edit Offline**:

- Edit draft reports offline
- Queue changes untuk sync

**Sync Rules**:

- Sync saat connection restored
- Retry 3x dengan exponential backoff
- Conflict resolution: server wins

---

## Keputusan Teknis & Trade-offs

### Mengapa Geolocator vs Location Package Lain?

**Keputusan**: Menggunakan geolocator package.

**Alasan**:

- **Popularity**: Most popular location package untuk Flutter
- **Feature Rich**: Background location, geocoding, permission handling
- **Active Maintenance**: Regularly updated
- **Platform Support**: Support Android, iOS, Web

**Trade-off**: Large package size. **Mitigasi**: Use specific features only.

### Photo Compression Strategy

**Keputusan**: Compress photos client-side sebelum upload.

**Alasan**:

- **Bandwidth**: Reduce data usage untuk mobile users
- **Speed**: Faster upload times
- **Cost**: Reduce server storage dan bandwidth costs

**Trade-off**: Additional processing time di client. **Mitigasi**: Compress asynchronously dengan progress indicator.

### Offline-first vs Online-first

**Keputusan**: Offline-first approach untuk visit report creation.

**Alasan**:

- **User Experience**: User dapat create reports meskipun offline
- **Reliability**: Tidak kehilangan data saat connection issues
- **Flexibility**: Sales rep dapat bekerja di remote areas

**Trade-off**: Complex sync logic. **Mitigasi**: Implement robust sync service dengan conflict resolution.

### Sequential vs Parallel Photo Upload

**Keputusan**: Sequential upload dengan queue.

**Alasan**:

- **Reliability**: Better error handling per photo
- **Progress**: Clear progress tracking
- **Memory**: Lower memory usage

**Trade-off**: Slower total upload time. **Mitigasi**: Compress photos sebelum upload untuk reduce size.

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── visit_reports/
│       ├── data/
│       │   ├── models/
│       │   │   ├── visit_report_model.dart     # Visit report entity
│       │   │   ├── visit_location_model.dart   # GPS coordinates
│       │   │   └── visit_photo_model.dart      # Photo metadata
│       │   └── visit_report_repository.dart    # API & sync
│       ├── application/
│       │   ├── visit_report_list_provider.dart
│       │   ├── visit_report_detail_provider.dart
│       │   ├── visit_report_form_provider.dart # Form state
│       │   └── visit_report_sync_provider.dart # Sync queue
│       └── presentation/
│           ├── screens/
│           │   ├── visit_report_list_screen.dart
│           │   ├── visit_report_detail_screen.dart
│           │   ├── visit_report_form_screen.dart
│           │   └── visit_report_approval_screen.dart
│           └── widgets/
│               ├── visit_report_card.dart
│               ├── visit_report_timeline.dart
│               ├── photo_gallery.dart
│               ├── photo_picker.dart
│               ├── gps_check_in_button.dart
│               └── status_badge.dart
├── core/
│   ├── services/
│   │   ├── location_service.dart              # GPS wrapper
│   │   └── photo_service.dart                 # Photo handling
│   └── utils/
│       └── image_compressor.dart              # Image compression
```

---

## API Endpoints

### Visit Report CRUD

#### GET /api/v1/visit-reports

List visit reports dengan filter.

**Query Parameters**:

```
?page=1&limit=20&status=draft&account_id=uuid&start_date=2025-01-01&end_date=2025-01-31
```

**Response**:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "vr-uuid",
        "account_id": "account-uuid",
        "account_name": "RS Medika Hospital",
        "sales_rep_id": "user-uuid",
        "sales_rep_name": "John Doe",
        "visit_date": "2025-01-15",
        "check_in_at": "2025-01-15T09:00:00Z",
        "check_in_location": {
          "latitude": -6.2088,
          "longitude": 106.8456,
          "accuracy": 10.5
        },
        "check_out_at": "2025-01-15T10:30:00Z",
        "check_out_location": {
          "latitude": -6.2088,
          "longitude": 106.8456,
          "accuracy": 12.0
        },
        "notes": "Discussed new product line",
        "photos": [
          {
            "id": "photo-uuid",
            "url": "https://storage.example.com/photo1.jpg",
            "thumbnail_url": "https://storage.example.com/photo1_thumb.jpg"
          }
        ],
        "status": "submitted",
        "supervisor_comment": null,
        "created_at": "2025-01-15T09:00:00Z",
        "updated_at": "2025-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 100
    }
  }
}
```

#### GET /api/v1/visit-reports/:id

Get visit report detail.

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "vr-uuid",
    "account": {
      "id": "account-uuid",
      "name": "RS Medika Hospital",
      "address": "Jl. Sudirman No. 123"
    },
    "sales_rep": {
      "id": "user-uuid",
      "name": "John Doe",
      "email": "john@example.com"
    },
    "visit_date": "2025-01-15",
    "check_in_at": "2025-01-15T09:00:00Z",
    "check_in_location": {
      "latitude": -6.2088,
      "longitude": 106.8456,
      "accuracy": 10.5,
      "address": "Jl. Sudirman, Jakarta"
    },
    "check_out_at": "2025-01-15T10:30:00Z",
    "check_out_location": {
      "latitude": -6.2088,
      "longitude": 106.8456,
      "accuracy": 12.0
    },
    "duration_minutes": 90,
    "notes": "Discussed new product line. Dr. Smith interested in sample.",
    "photos": [
      {
        "id": "photo-uuid",
        "url": "https://storage.example.com/photo1.jpg",
        "thumbnail_url": "https://storage.example.com/photo1_thumb.jpg",
        "uploaded_at": "2025-01-15T09:15:00Z"
      }
    ],
    "status": "submitted",
    "supervisor_id": null,
    "supervisor_comment": null,
    "approved_at": null,
    "created_at": "2025-01-15T09:00:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

#### POST /api/v1/visit-reports

Create new visit report.

**Request**:

```json
{
  "account_id": "account-uuid",
  "visit_date": "2025-01-15",
  "check_in_at": "2025-01-15T09:00:00Z",
  "check_in_location": {
    "latitude": -6.2088,
    "longitude": 106.8456,
    "accuracy": 10.5
  },
  "notes": "Discussed new product line"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "vr-uuid",
    "status": "draft",
    "message": "Visit report created successfully"
  }
}
```

#### PUT /api/v1/visit-reports/:id

Update visit report (hanya saat draft).

**Request**:

```json
{
  "notes": "Updated notes",
  "check_out_at": "2025-01-15T10:30:00Z",
  "check_out_location": {
    "latitude": -6.2088,
    "longitude": 106.8456,
    "accuracy": 12.0
  }
}
```

#### POST /api/v1/visit-reports/:id/check-in

Record check-in.

**Request**:

```json
{
  "latitude": -6.2088,
  "longitude": 106.8456,
  "accuracy": 10.5,
  "timestamp": "2025-01-15T09:00:00Z"
}
```

#### POST /api/v1/visit-reports/:id/check-out

Record check-out.

**Request**:

```json
{
  "latitude": -6.2088,
  "longitude": 106.8456,
  "accuracy": 12.0,
  "timestamp": "2025-01-15T10:30:00Z"
}
```

#### POST /api/v1/visit-reports/:id/submit

Submit visit report untuk approval.

**Response**:

```json
{
  "success": true,
  "data": {
    "status": "submitted",
    "message": "Visit report submitted for approval"
  }
}
```

### Photo Upload

#### POST /api/v1/visit-reports/:id/photos

Upload photo untuk visit report.

**Request**: Multipart form-data

```
file: [binary image data]
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "photo-uuid",
    "url": "https://storage.example.com/photo1.jpg",
    "thumbnail_url": "https://storage.example.com/photo1_thumb.jpg"
  }
}
```

---

## Data Models

### Visit Report Model

```dart
@freezed
class VisitReport with _$VisitReport {
  const factory VisitReport({
    required String id,
    required String accountId,
    String? accountName,
    required String salesRepId,
    String? salesRepName,
    required DateTime visitDate,
    DateTime? checkInAt,
    VisitLocation? checkInLocation,
    DateTime? checkOutAt,
    VisitLocation? checkOutLocation,
    String? notes,
    @Default([]) List<VisitPhoto> photos,
    @Default('draft') String status,
    String? supervisorId,
    String? supervisorComment,
    DateTime? approvedAt,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _VisitReport;

  factory VisitReport.fromJson(Map<String, dynamic> json) =
      _$VisitReportFromJson(json);
}

enum VisitReportStatus {
  draft,
  submitted,
  approved,
  rejected;

  String get displayName {
    switch (this) {
      case VisitReportStatus.draft:
        return 'Draft';
      case VisitReportStatus.submitted:
        return 'Submitted';
      case VisitReportStatus.approved:
        return 'Approved';
      case VisitReportStatus.rejected:
        return 'Rejected';
    }
  }

  Color get color {
    switch (this) {
      case VisitReportStatus.draft:
        return Colors.orange;
      case VisitReportStatus.submitted:
        return Colors.blue;
      case VisitReportStatus.approved:
        return Colors.green;
      case VisitReportStatus.rejected:
        return Colors.red;
    }
  }
}
```

### Visit Location Model

```dart
@freezed
class VisitLocation with _$VisitLocation {
  const factory VisitLocation({
    required double latitude,
    required double longitude,
    double? accuracy,
    String? address,
  }) = _VisitLocation;

  factory VisitLocation.fromJson(Map<String, dynamic> json) =
      _$VisitLocationFromJson(json);

  factory VisitLocation.fromPosition(Position position) {
    return VisitLocation(
      latitude: position.latitude,
      longitude: position.longitude,
      accuracy: position.accuracy,
    );
  }
}
```

### Visit Photo Model

```dart
@freezed
class VisitPhoto with _$VisitPhoto {
  const factory VisitPhoto({
    required String id,
    required String url,
    String? thumbnailUrl,
    DateTime? uploadedAt,
    @Default(false) bool isLocal,
    String? localPath,
    @Default('pending') String syncStatus,
  }) = _VisitPhoto;

  factory VisitPhoto.fromJson(Map<String, dynamic> json) =
      _$VisitPhotoFromJson(json);
}
```

### Visit Report Form State

```dart
@freezed
class VisitReportFormState with _$VisitReportFormState {
  const factory VisitReportFormState({
    String? id,
    String? accountId,
    String? accountName,
    DateTime? visitDate,
    DateTime? checkInAt,
    VisitLocation? checkInLocation,
    DateTime? checkOutAt,
    VisitLocation? checkOutLocation,
    @Default('') String notes,
    @Default([]) List<XFile> localPhotos,
    @Default([]) List<VisitPhoto> uploadedPhotos,
    @Default(false) bool isSubmitting,
    @Default(false) bool isCheckingIn,
    @Default(false) bool isCheckingOut,
    String? error,
    @Default(false) bool isSuccess,
  }) = _VisitReportFormState;
}
```

---

## Configuration

### Location Service

**File**: `core/services/location_service.dart`

```dart
class LocationService {
  Future<LocationResult> getCurrentLocation() async {
    // Check permissions
    bool serviceEnabled = await Geolocator.isLocationServiceEnabled();
    if (!serviceEnabled) {
      return LocationResult.error('Location services disabled');
    }

    LocationPermission permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
      if (permission == LocationPermission.denied) {
        return LocationResult.error('Location permission denied');
      }
    }

    if (permission == LocationPermission.deniedForever) {
      return LocationResult.error('Location permission permanently denied');
    }

    // Get position
    try {
      final position = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
        timeLimit: const Duration(seconds: 10),
      );

      // Validate accuracy
      if (position.accuracy > 50) {
        return LocationResult.warning(
          position: position,
          message: 'Low accuracy: ${position.accuracy.toStringAsFixed(1)}m',
        );
      }

      return LocationResult.success(position: position);
    } catch (e) {
      return LocationResult.error('Failed to get location: $e');
    }
  }

  double calculateDistance(
    double startLatitude,
    double startLongitude,
    double endLatitude,
    double endLongitude,
  ) {
    return Geolocator.distanceBetween(
      startLatitude,
      startLongitude,
      endLatitude,
      endLongitude,
    );
  }
}

class LocationResult {
  final Position? position;
  final String? error;
  final String? warning;
  final bool isSuccess;

  LocationResult._({
    this.position,
    this.error,
    this.warning,
    required this.isSuccess,
  });

  factory LocationResult.success({required Position position}) {
    return LocationResult._(position: position, isSuccess: true);
  }

  factory LocationResult.error(String message) {
    return LocationResult._(error: message, isSuccess: false);
  }

  factory LocationResult.warning({
    required Position position,
    required String message,
  }) {
    return LocationResult._(
      position: position,
      warning: message,
      isSuccess: true,
    );
  }
}
```

### Photo Service

**File**: `core/services/photo_service.dart`

```dart
class PhotoService {
  final ImagePicker _picker = ImagePicker();

  Future<XFile?> pickPhoto() async {
    return _picker.pickImage(
      source: ImageSource.camera,
      maxWidth: 1920,
      maxHeight: 1080,
      imageQuality: 80,
    );
  }

  Future<List<XFile>> pickMultiplePhotos() async {
    return _picker.pickMultiImage(
      maxWidth: 1920,
      maxHeight: 1080,
      imageQuality: 80,
    );
  }

  Future<File> compressPhoto(XFile file) async {
    final dir = await getTemporaryDirectory();
    final targetPath = '${dir.path}/${DateTime.now().millisecondsSinceEpoch}.jpg';

    final result = await FlutterImageCompress.compressAndGetFile(
      file.path,
      targetPath,
      quality: 80,
      minWidth: 1280,
      minHeight: 720,
    );

    return File(result!.path);
  }

  Future<String> savePhotoLocally(XFile file) async {
    final dir = await getApplicationDocumentsDirectory();
    final photosDir = Directory('${dir.path}/visit_photos');
    if (!await photosDir.exists()) {
      await photosDir.create(recursive: true);
    }

    final fileName = 'photo_${DateTime.now().millisecondsSinceEpoch}.jpg';
    final localPath = '${photosDir.path}/$fileName';

    await File(file.path).copy(localPath);
    return localPath;
  }
}
```

---

## Usage Examples

### Visit Report Creation Flow

```dart
class VisitReportFormScreen extends ConsumerWidget {
  const VisitReportFormScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(visitReportFormProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Create Visit Report'),
      ),
      body: Form(
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Account Selection
            AccountSelector(
              selectedAccount: state.accountId != null
                  ? Account(id: state.accountId!, name: state.accountName!)
                  : null,
              onSelect: (account) {
                ref.read(visitReportFormProvider.notifier)
                    .setAccount(account);
              },
            ),

            const SizedBox(height: 16),

            // Check-in Button
            GpsCheckInButton(
              isCheckedIn: state.checkInAt != null,
              location: state.checkInLocation,
              onCheckIn: () async {
                final result = await ref.read(locationServiceProvider)
                    .getCurrentLocation();

                result.when(
                  success: (position) {
                    ref.read(visitReportFormProvider.notifier).checkIn(
                      VisitLocation.fromPosition(position),
                    );
                  },
                  error: (message) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text(message)),
                    );
                  },
                );
              },
            ),

            const SizedBox(height: 16),

            // Photos
            PhotoPickerSection(
              photos: [...state.localPhotos, ...state.uploadedPhotos],
              onAddPhoto: () async {
                final photo = await ref.read(photoServiceProvider)
                    .pickPhoto();
                if (photo != null) {
                  ref.read(visitReportFormProvider.notifier)
                      .addPhoto(photo);
                }
              },
              onRemovePhoto: (index) {
                ref.read(visitReportFormProvider.notifier)
                    .removePhoto(index);
              },
            ),

            const SizedBox(height: 16),

            // Notes
            TextFormField(
              maxLines: 5,
              decoration: const InputDecoration(
                labelText: 'Notes',
                hintText: 'Describe your visit...',
                border: OutlineInputBorder(),
              ),
              onChanged: (value) {
                ref.read(visitReportFormProvider.notifier).setNotes(value);
              },
            ),

            const SizedBox(height: 16),

            // Check-out Button
            GpsCheckOutButton(
              isCheckedIn: state.checkInAt != null,
              isCheckedOut: state.checkOutAt != null,
              onCheckOut: () async {
                final result = await ref.read(locationServiceProvider)
                    .getCurrentLocation();

                result.when(
                  success: (position) {
                    ref.read(visitReportFormProvider.notifier).checkOut(
                      VisitLocation.fromPosition(position),
                    );
                  },
                  error: (message) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text(message)),
                    );
                  },
                );
              },
            ),

            const SizedBox(height: 24),

            // Submit Button
            FilledButton.icon(
              onPressed: state.isSubmitting
                  ? null
                  : () => _submit(context, ref),
              icon: state.isSubmitting
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Icon(Icons.send),
              label: Text(state.isSubmitting ? 'Submitting...' : 'Submit'),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _submit(BuildContext context, WidgetRef ref) async {
    final success = await ref.read(visitReportFormProvider.notifier).submit();

    if (success && context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Visit report submitted successfully')),
      );
      context.pop();
    }
  }
}
```

### Sync Queue untuk Offline

```dart
class VisitReportSyncService {
  final Box<VisitReportHive> _pendingBox;
  final ApiClient _apiClient;
  final ConnectivityService _connectivity;

  VisitReportSyncService(
    this._pendingBox,
    this._apiClient,
    this._connectivity,
  );

  Future<void> queueVisitReport(VisitReport report) async {
    final hiveReport = VisitReportHive()
      ..id = report.id
      ..accountId = report.accountId
      ..checkInLocation = report.checkInLocation
      ..checkOutLocation = report.checkOutLocation
      ..notes = report.notes
      ..localPhotoPaths = report.localPhotos.map((p) => p.path).toList()
      ..syncStatus = 'pending'
      ..createdAt = DateTime.now();

    await _pendingBox.put(report.id, hiveReport);

    // Try sync jika online
    if (_connectivity.isOnline) {
      await syncPendingReports();
    }
  }

  Future<void> syncPendingReports() async {
    if (!_connectivity.isOnline) return;

    final pending = _pendingBox.values
        .where((r) => r.syncStatus == 'pending')
        .toList();

    for (final report in pending) {
      try {
        // Update status
        report.syncStatus = 'syncing';
        await _pendingBox.put(report.id, report);

        // Upload photos dulu
        final photoUrls = await _uploadPhotos(report.localPhotoPaths);

        // Create visit report
        await _apiClient.post('/api/v1/visit-reports', data: {
          'account_id': report.accountId,
          'visit_date': report.visitDate?.toIso8601String(),
          'check_in_location': report.checkInLocation?.toJson(),
          'check_out_location': report.checkOutLocation?.toJson(),
          'notes': report.notes,
          'photo_ids': photoUrls.map((p) => p['id']).toList(),
        });

        // Mark as synced
        report.syncStatus = 'synced';
        await _pendingBox.put(report.id, report);
        await _pendingBox.delete(report.id);
      } catch (e) {
        report.syncStatus = 'failed';
        report.errorMessage = e.toString();
        await _pendingBox.put(report.id, report);
      }
    }
  }

  Future<List<Map<String, dynamic>>> _uploadPhotos(
    List<String> localPaths,
  ) async {
    final uploaded = <Map<String, dynamic>>[];

    for (final path in localPaths) {
      final file = File(path);
      if (await file.exists()) {
        final response = await _apiClient.uploadFile(
          '/api/v1/upload',
          file,
        );
        uploaded.add(response.data['data']);
      }
    }

    return uploaded;
  }
}
```

---

## Cara Test Manual

### Test Visit Report Creation

1. **Create Flow**:
   - Tap "Create Visit Report"
   - Select account
   - Tap "Check-in"
   - Verifikasi: GPS coordinates captured
   - Add photos
   - Add notes
   - Tap "Check-out"
   - Submit
   - Verifikasi: Report created dengan status Draft

2. **GPS Validation**:
   - Check-in dengan poor GPS signal
   - Verifikasi: Warning message muncul
   - Verifikasi: Tetap bisa continue dengan acknowledgment

3. **Photo Upload**:
   - Add 5 photos
   - Verifikasi: Semua photos tersimpan
   - Compress dan upload
   - Verifikasi: Progress indicator muncul

4. **Offline Creation**:
   - Matikan internet
   - Create visit report
   - Verifikasi: Report saved locally
   - Turn on internet
   - Verifikasi: Auto-sync terjadi

### Test Approval Workflow

1. **Submit untuk Approval**:
   - Create visit report lengkap
   - Submit
   - Verifikasi: Status berubah ke Submitted
   - Tidak bisa edit lagi

2. **Supervisor Approval**:
   - Login sebagai supervisor
   - View submitted reports
   - Approve dengan comment
   - Verifikasi: Status berubah ke Approved

3. **Rejection**:
   - Reject visit report dengan reason
   - Verifikasi: Status berubah ke Rejected
   - Verifikasi: Sales rep dapat melihat rejection reason

---

## Dependencies

### Internal

- `core/services/location_service.dart` - GPS functionality
- `core/services/photo_service.dart` - Photo handling
- `core/network/api_client.dart` - API calls
- `core/storage/hive_storage.dart` - Offline storage

### External

- `geolocator: ^10.0.0` - GPS location
- `image_picker: ^1.0.0` - Camera integration
- `flutter_image_compress: ^2.0.0` - Photo compression
- `photo_view: ^0.14.0` - Photo gallery view
- `carousel_slider: ^4.2.0` - Photo carousel

---

## Notes & Improvements

### Known Limitations

1. **GPS Accuracy**: Accuracy tergantung device dan environment.

2. **Photo Storage**: Large storage requirement untuk offline photos.

3. **Sync Complexity**: Complex sync logic untuk offline scenarios.

4. **No Live Tracking**: Tidak ada live location tracking selama visit.

### Future Improvements

1. **Live Tracking**: Track location continuously selama visit dengan user consent

2. **Voice Notes**: Add voice recording sebagai alternative ke text notes

3. **Signature Capture**: Capture digital signature dari contact person

4. **Barcode Scanner**: Scan product barcodes untuk visit reports related ke product orders

5. **Route Optimization**: Optimize visit route based on multiple accounts

6. **Visit Templates**: Pre-defined templates untuk different types of visits

7. **Automated Reports**: Auto-generate weekly/monthly visit summaries

---

**Document Status**: Active  
**Last Updated**: January 2025  
**Maintained By**: Dev3 (Mobile Development Team)
