import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:geolocator/geolocator.dart';
import 'package:image_picker/image_picker.dart';
import 'package:intl/intl.dart';

import '../application/visit_report_provider.dart';
import '../utils/fake_gps_detector.dart';
import '../widgets/fake_gps_warning_dialog.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/widgets/error_widget.dart';
import '../../../core/widgets/loading_widget.dart';
import 'widgets/selfie_preview_screen.dart';

class VisitReportDetailScreen extends ConsumerStatefulWidget {
  const VisitReportDetailScreen({super.key, required this.visitReportId});

  final String visitReportId;

  @override
  ConsumerState<VisitReportDetailScreen> createState() =>
      _VisitReportDetailScreenState();
}

class _VisitReportDetailScreenState
    extends ConsumerState<VisitReportDetailScreen> {
  final ImagePicker _imagePicker = ImagePicker();
  Position? _previousGPSPosition;

  /// Format DateTime to local timezone with readable format
  /// Example: "24 Dec 2025, 08:25 AM"
  String _formatDateTime(DateTime dateTime) {
    // Convert to local timezone
    final localDateTime = dateTime.toLocal();
    // Format: "24 Dec 2025, 08:25 AM"
    return DateFormat('dd MMM yyyy, hh:mm a', 'en_US').format(localDateTime);
  }

  /// Format visit date string with locale-aware formatting
  String _formatVisitDate(String visitDate, Locale locale) {
    try {
      // Try parsing ISO 8601 format
      final dateTime = DateTime.parse(visitDate).toLocal();
      
      // Use locale-aware date format
      final localeString = locale.languageCode == 'id' ? 'id_ID' : 'en_US';
      final dateFormat = DateFormat('EEEE, d MMMM yyyy', localeString);
      final timeFormat = DateFormat('HH:mm', localeString);
      
      // Check if time is included
      if (visitDate.contains('T') && visitDate.length > 10) {
        final timePart = visitDate.split('T')[1];
        if (timePart.isNotEmpty && !timePart.startsWith('00:00:00')) {
          return '${dateFormat.format(dateTime)} at ${timeFormat.format(dateTime)}';
        }
      }
      
      return dateFormat.format(dateTime);
    } catch (e) {
      // If parsing fails, return as-is
      return visitDate;
    }
  }

  Future<void> _requestLocationPermission() async {
    final permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      await Geolocator.requestPermission();
    }
  }

  Future<Position?> _getCurrentLocation() async {
    try {
      // Check if location services are enabled
      bool serviceEnabled = await Geolocator.isLocationServiceEnabled();
      if (!serviceEnabled) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text(
                'Location services are disabled. Please enable location services.',
              ),
              backgroundColor: Colors.orange,
            ),
          );
        }
        return null;
      }

      await _requestLocationPermission();
      final permission = await Geolocator.checkPermission();

      if (permission == LocationPermission.denied) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text(
                'Location permission is required for check-in/out. Please grant permission in settings.',
              ),
              backgroundColor: Colors.orange,
            ),
          );
        }
        return null;
      }

      if (permission == LocationPermission.deniedForever) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text(
                'Location permission is permanently denied. Please enable it in app settings.',
              ),
              backgroundColor: Colors.red,
            ),
          );
        }
        return null;
      }

      // BEST PRACTICE: Use bestForNavigation for most accurate GPS location
      final position = await Geolocator.getCurrentPosition(
        locationSettings: const LocationSettings(
          accuracy: LocationAccuracy.bestForNavigation, // BEST PRACTICE: Most accurate for navigation
          timeLimit: Duration(seconds: 20), // Allow time for GPS to get accurate fix
        ),
      );
      return position;
    } catch (e) {
      if (mounted) {
        String errorMessage = 'Failed to get location';
        if (e.toString().contains('MissingPluginException')) {
          errorMessage =
              'Location plugin not initialized. Please stop and restart the app.';
        } else {
          errorMessage =
              'Failed to get location: ${e.toString().replaceFirst('Exception: ', '')}';
        }
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(errorMessage),
            backgroundColor: Colors.red,
            duration: const Duration(seconds: 5),
          ),
        );
      }
      return null;
    }
  }

  Future<void> _handleCheckIn() async {
    final l10n = AppLocalizations.of(context)!;

    // First, get location
    final position = await _getCurrentLocation();
    if (position == null) return;

    // Detect Fake GPS
    final fakeGPSDetection = detectFakeGPS(
      position,
      previousPosition: _previousGPSPosition,
    );

    if (fakeGPSDetection.isFakeGPS) {
      // Fake GPS detected - block check-in and show warning
      if (mounted) {
        showDialog(
          context: context,
          builder: (context) =>
              FakeGPSWarningDialog(reason: fakeGPSDetection.reason),
        );
      }
      return;
    }

    // GPS is valid - save for next check
    _previousGPSPosition = position;

    // Then, capture selfie
    File? photoFile;
    bool shouldProceed = false;

    while (!shouldProceed) {
      try {
        final XFile? image = await _imagePicker.pickImage(
          source: ImageSource.camera,
          preferredCameraDevice: CameraDevice.front, // Front camera for selfie
          imageQuality: 85,
        );

        if (image == null) {
          // User cancelled camera - show message that selfie is required
          if (mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(l10n.selfieRequiredForCheckIn),
                backgroundColor: Colors.orange,
                duration: const Duration(seconds: 3),
              ),
            );
          }
          return;
        }

        photoFile = File(image.path);

        // Show preview screen
        if (mounted) {
          final confirmed = await Navigator.of(context).push<bool>(
            MaterialPageRoute(
              builder: (context) => SelfiePreviewScreen(
                photoFile: photoFile!,
                visitReportId: widget.visitReportId,
                position: position,
              ),
            ),
          );

          if (confirmed == true) {
            shouldProceed = true;
          } else {
            // User wants to retake, continue loop
            continue;
          }
        } else {
          return;
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(
                'Failed to capture photo: ${e.toString().replaceFirst('Exception: ', '')}',
              ),
              backgroundColor: Colors.red,
            ),
          );
        }
        return;
      }
    }

    // Submit check-in with confirmed photo
    if (photoFile != null && shouldProceed) {
      final formNotifier = ref.read(visitReportFormProvider.notifier);
      final result = await formNotifier.checkIn(
        visitReportId: widget.visitReportId,
        latitude: position.latitude,
        longitude: position.longitude,
        photoFile: photoFile,
        accuracy: position.accuracy,
        // Note: EXIF GPS extraction would require additional package (exif package)
        // For now, we'll use device GPS only
        // photoLatitude: extractedFromExif,
        // photoLongitude: extractedFromExif,
        // photoTimestamp: extractedFromExif,
      );

      if (mounted) {
        if (result != null) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(l10n.checkInSuccessful),
              backgroundColor: Colors.green,
            ),
          );
          // Refresh detail
          ref.invalidate(visitReportDetailProvider(widget.visitReportId));
        } else {
          final error = ref.read(visitReportFormProvider).errorMessage;
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(error ?? l10n.failedToCheckIn),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    }
  }

  Future<void> _handleCheckOut() async {
    final l10n = AppLocalizations.of(context)!;

    // Show loading dialog while getting location
    if (mounted) {
      showDialog(
        context: context,
        barrierDismissible: false,
        builder: (context) => AlertDialog(
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const CircularProgressIndicator(),
              const SizedBox(height: 16),
              Text(l10n.gettingLocation),
            ],
          ),
        ),
      );
    }

    // Get current location
    final position = await _getCurrentLocation();

    // Close loading dialog
    if (mounted) {
      Navigator.of(context).pop();
    }

    if (position == null) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.checkOutLocationRequired),
            backgroundColor: Colors.orange,
            duration: const Duration(seconds: 3),
          ),
        );
      }
      return;
    }

    // Detect Fake GPS
    final fakeGPSDetection = detectFakeGPS(
      position,
      previousPosition: _previousGPSPosition,
    );

    if (fakeGPSDetection.isFakeGPS) {
      // Fake GPS detected - block check-out and show warning
      if (mounted) {
        showDialog(
          context: context,
          builder: (context) =>
              FakeGPSWarningDialog(reason: fakeGPSDetection.reason),
        );
      }
      return;
    }

    // GPS is valid - save for next check
    _previousGPSPosition = position;

    // Ask if user wants to add photo (optional, like web version)
    if (!mounted) return;
    final addPhoto = await showDialog<bool>(
      context: context,
      builder: (context) {
        final theme = Theme.of(context);
        return AlertDialog(
          title: Text(l10n.confirmCheckOut),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                l10n.currentLocation,
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  Icon(Icons.location_on,
                      color: theme.colorScheme.primary, size: 20),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      '${position.latitude.toStringAsFixed(6)}, ${position.longitude.toStringAsFixed(6)}',
                      style: theme.textTheme.bodyMedium,
                    ),
                  ),
                ],
              ),
              if (position.accuracy > 0) ...[
                const SizedBox(height: 8),
                Row(
                  children: [
                    Icon(
                      Icons.gps_fixed,
                      color: theme.colorScheme.secondary,
                      size: 20,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      '${l10n.locationAccuracy}: ${position.accuracy.toStringAsFixed(1)}m',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurface.withValues(alpha:0.7),
                      ),
                    ),
                  ],
                ),
              ],
              const SizedBox(height: 16),
              Text(
                'Add photo? (Optional)',
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(false),
              child: Text(l10n.cancel),
            ),
            TextButton(
              onPressed: () => Navigator.of(context).pop(false), // No photo
              child: Text('Check-out without photo'),
            ),
            FilledButton(
              onPressed: () => Navigator.of(context).pop(true), // With photo
              style: FilledButton.styleFrom(backgroundColor: Colors.blue),
              child: Text('Check-out with photo'),
            ),
          ],
        );
      },
    );

    if (addPhoto == null) return; // User cancelled

    File? photoFile;
    if (addPhoto == true) {
      // User wants to add photo - capture selfie
      try {
        final XFile? image = await _imagePicker.pickImage(
          source: ImageSource.camera,
          preferredCameraDevice: CameraDevice.front, // Front camera for selfie
          imageQuality: 85,
        );

        if (image != null) {
          photoFile = File(image.path);
          // Show preview and confirm
          if (!mounted) return;
          final confirmed = await Navigator.of(context).push<bool>(
            MaterialPageRoute(
              builder: (context) => SelfiePreviewScreen(
                photoFile: photoFile!,
                visitReportId: widget.visitReportId,
                position: position,
              ),
            ),
          );

          if (confirmed != true) {
            photoFile = null; // User cancelled photo
          }
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(
                'Failed to capture photo: ${e.toString().replaceFirst('Exception: ', '')}',
              ),
              backgroundColor: Colors.red,
            ),
          );
        }
        return;
      }
    }

    // Perform check-out (with or without photo)
    final formNotifier = ref.read(visitReportFormProvider.notifier);
    final result = await formNotifier.checkOut(
      visitReportId: widget.visitReportId,
      latitude: position.latitude,
      longitude: position.longitude,
      photoFile: photoFile, // Optional photo
      accuracy: position.accuracy,
    );

    if (mounted) {
      if (result != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.checkOutSuccessful),
            backgroundColor: Colors.green,
          ),
        );
        // Refresh detail
        ref.invalidate(visitReportDetailProvider(widget.visitReportId));
      } else {
        final error = ref.read(visitReportFormProvider).errorMessage;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(error ?? l10n.failedToCheckOut),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Future<void> _handleUpdate() async {
    final l10n = AppLocalizations.of(context)!;
    
    // Get visit report data from provider
    final visitReportAsync = ref.read(visitReportDetailProvider(widget.visitReportId));
    
    // Get visit report data - wait if loading
    final visitReport = await visitReportAsync.when(
      data: (data) => Future.value(data),
      loading: () async {
        // Show loading dialog
        if (mounted) {
          showDialog(
            context: context,
            barrierDismissible: false,
            builder: (context) => const Center(child: CircularProgressIndicator()),
          );
        }
        // Wait for provider to finish loading
        await Future.delayed(const Duration(milliseconds: 100));
        final newAsync = ref.read(visitReportDetailProvider(widget.visitReportId));
        if (mounted) Navigator.of(context).pop(); // Close loading dialog
        return newAsync.when(
          data: (data) => data,
          loading: () => throw Exception('Still loading'),
          error: (error, stack) => throw error,
        );
      },
      error: (error, stack) => throw error,
    );
    
    if (!mounted) return;
    final theme = Theme.of(context);
    
    final purposeController = TextEditingController(text: visitReport.purpose ?? '');
    final notesController = TextEditingController(text: visitReport.notes ?? '');
    
    // Parse visit date - handle both "YYYY-MM-DD" and "YYYY-MM-DD HH:mm:ss" formats
    DateTime visitDateTime;
    if (visitReport.visitDate.contains(' ')) {
      visitDateTime = DateTime.parse(visitReport.visitDate.split(' ')[0]);
      final timeParts = visitReport.visitDate.split(' ')[1].split(':');
      visitDateTime = visitDateTime.add(Duration(
        hours: int.parse(timeParts[0]),
        minutes: int.parse(timeParts[1]),
      ));
    } else {
      visitDateTime = DateTime.parse(visitReport.visitDate);
    }
    
    DateTime selectedDate = visitDateTime;
    TimeOfDay selectedTime = TimeOfDay.fromDateTime(visitDateTime);

    if (!mounted) return;
    final scaffoldMessenger = ScaffoldMessenger.of(context);
    final updated = await showDialog<bool>(
      context: context,
      builder: (context) {
        return StatefulBuilder(
          builder: (context, setState) {
            return AlertDialog(
              title: Text(l10n.updateVisitReport),
              content: SingleChildScrollView(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Visit Date & Time
                    Text(
                      '${l10n.visitDate} *',
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        Expanded(
                          child: OutlinedButton.icon(
                            onPressed: () async {
                              final picked = await showDatePicker(
                                context: context,
                                initialDate: selectedDate,
                                firstDate: DateTime.now().subtract(const Duration(days: 365)),
                                lastDate: DateTime.now().add(const Duration(days: 365)),
                              );
                              if (picked != null) {
                                setState(() {
                                  selectedDate = picked;
                                });
                              }
                            },
                            icon: const Icon(Icons.calendar_today, size: 18),
                            label: Text(
                              DateFormat('dd/MM/yyyy').format(selectedDate),
                              style: theme.textTheme.bodySmall,
                            ),
                          ),
                        ),
                        const SizedBox(width: 8),
                        Expanded(
                          child: OutlinedButton.icon(
                            onPressed: () async {
                              final picked = await showTimePicker(
                                context: context,
                                initialTime: selectedTime,
                              );
                              if (picked != null) {
                                setState(() {
                                  selectedTime = picked;
                                });
                              }
                            },
                            icon: const Icon(Icons.access_time, size: 18),
                            label: Text(
                              selectedTime.format(context),
                              style: theme.textTheme.bodySmall,
                            ),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                    // Purpose
                    Text(
                      '${l10n.purpose} *',
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    TextField(
                      controller: purposeController,
                      decoration: InputDecoration(
                        hintText: 'e.g., Product demo, Follow-up',
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                        contentPadding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 8,
                        ),
                      ),
                      maxLines: 2,
                    ),
                    const SizedBox(height: 16),
                    // Notes
                    Text(
                      l10n.notes,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    TextField(
                      controller: notesController,
                      decoration: InputDecoration(
                        hintText: 'Enter additional notes...',
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                        contentPadding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 8,
                        ),
                      ),
                      maxLines: 3,
                    ),
                  ],
                ),
              ),
              actions: [
                TextButton(
                  onPressed: () {
                    purposeController.dispose();
                    notesController.dispose();
                    Navigator.of(context).pop(false);
                  },
                  child: Text(l10n.cancel),
                ),
                FilledButton(
                  onPressed: () {
                    if (purposeController.text.trim().isEmpty) {
                      scaffoldMessenger.showSnackBar(
                        SnackBar(
                          content: Text(l10n.purposeRequired),
                          backgroundColor: Colors.red,
                        ),
                      );
                      return;
                    }
                    Navigator.of(context).pop(true);
                  },
                  style: FilledButton.styleFrom(backgroundColor: Colors.blue),
                  child: Text(l10n.save),
                ),
              ],
            );
          },
        );
      },
    );

    if (updated != true) {
      purposeController.dispose();
      notesController.dispose();
      return;
    }

    // Format visit date
    final dateStr = DateFormat('yyyy-MM-dd').format(selectedDate);
    final timeStr = '${selectedTime.hour.toString().padLeft(2, '0')}:${selectedTime.minute.toString().padLeft(2, '0')}';
    final visitDate = '$dateStr $timeStr';

    final formNotifier = ref.read(visitReportFormProvider.notifier);
    final result = await formNotifier.updateVisitReport(
      id: widget.visitReportId,
      visitDate: visitDate,
      purpose: purposeController.text.trim(),
      notes: notesController.text.trim().isNotEmpty
          ? notesController.text.trim()
          : null,
    );

    purposeController.dispose();
    notesController.dispose();

    if (mounted) {
      if (result != null) {
        scaffoldMessenger.showSnackBar(
          SnackBar(
            content: Text(l10n.visitReportUpdatedSuccessfully),
            backgroundColor: Colors.green,
          ),
        );
        // Refresh detail
        ref.invalidate(visitReportDetailProvider(widget.visitReportId));
      } else {
        final error = ref.read(visitReportFormProvider).errorMessage;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(error ?? l10n.failedToUpdateVisitReport),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Future<void> _handleDelete() async {
    final l10n = AppLocalizations.of(context)!;

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text(l10n.deleteVisitReport),
          content: Text(l10n.deleteVisitReportConfirmation),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(false),
              child: Text(l10n.cancel),
            ),
            FilledButton(
              onPressed: () => Navigator.of(context).pop(true),
              style: FilledButton.styleFrom(backgroundColor: Colors.red),
              child: Text(l10n.delete),
            ),
          ],
        );
      },
    );

    if (confirmed != true) return;

    final formNotifier = ref.read(visitReportFormProvider.notifier);
    final success = await formNotifier.deleteVisitReport(widget.visitReportId);

    if (mounted) {
      if (success) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.visitReportDeletedSuccessfully),
            backgroundColor: Colors.green,
          ),
        );
        Navigator.of(context).pop(); // Go back to list
      } else {
        final error = ref.read(visitReportFormProvider).errorMessage;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(error ?? l10n.failedToDeleteVisitReport),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Future<void> _handleSubmit() async {
    final l10n = AppLocalizations.of(context)!;

    String? selectedOutcome;
    String? nextSteps;

    final submitted = await showDialog<bool>(
      context: context,
      builder: (context) {
        final theme = Theme.of(context);
        return StatefulBuilder(
          builder: (context, setState) {
            return AlertDialog(
              title: Text(l10n.submitVisitReport),
              content: SingleChildScrollView(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Outcome selection
                    Text(
                      '${l10n.outcome} *',
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    DropdownButtonFormField<String>(
                      initialValue: selectedOutcome,
                      decoration: InputDecoration(
                        hintText: l10n.selectOutcome,
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                        contentPadding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 8,
                        ),
                      ),
                      items: [
                        DropdownMenuItem(
                          value: 'positive',
                          child: Text('Positive'),
                        ),
                        DropdownMenuItem(
                          value: 'very_positive',
                          child: Text('Very Positive'),
                        ),
                        DropdownMenuItem(
                          value: 'neutral',
                          child: Text('Neutral'),
                        ),
                        DropdownMenuItem(
                          value: 'negative',
                          child: Text('Negative'),
                        ),
                      ],
                      onChanged: (value) {
                        setState(() {
                          selectedOutcome = value;
                        });
                      },
                    ),
                    const SizedBox(height: 16),
                    // Next Steps (optional)
                    Text(
                      l10n.nextSteps,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    TextField(
                      decoration: InputDecoration(
                        hintText: l10n.enterNextSteps,
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                        contentPadding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 8,
                        ),
                      ),
                      maxLines: 3,
                      onChanged: (value) {
                        nextSteps = value.trim().isNotEmpty ? value.trim() : null;
                      },
                    ),
                  ],
                ),
              ),
              actions: [
                TextButton(
                  onPressed: () => Navigator.of(context).pop(false),
                  child: Text(l10n.cancel),
                ),
                FilledButton(
                  onPressed: selectedOutcome == null
                      ? null
                      : () {
                          Navigator.of(context).pop(true);
                        },
                  style: FilledButton.styleFrom(backgroundColor: Colors.green),
                  child: Text(l10n.submitVisitReport),
                ),
              ],
            );
          },
        );
      },
    );

    if (submitted != true || selectedOutcome == null) return;

    final formNotifier = ref.read(visitReportFormProvider.notifier);
    final result = await formNotifier.submitVisitReport(
      id: widget.visitReportId,
      outcome: selectedOutcome,
      nextSteps: nextSteps,
    );

    if (mounted) {
      if (result != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.visitReportSubmittedSuccessfully),
            backgroundColor: Colors.green,
          ),
        );
        // Refresh detail
        ref.invalidate(visitReportDetailProvider(widget.visitReportId));
      } else {
        final error = ref.read(visitReportFormProvider).errorMessage;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(error ?? l10n.failedToSubmitVisitReport),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final visitReportAsync = ref.watch(
      visitReportDetailProvider(widget.visitReportId),
    );
    final formState = ref.watch(visitReportFormProvider);
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(title: Text(l10n.visitReportDetails), elevation: 0),
      body: visitReportAsync.when(
        data: (visitReport) => SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Header: Type Badge & Status Badge
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  _TypeBadge(
                    type: visitReport.type,
                    theme: theme,
                    colorScheme: colorScheme,
                    l10n: l10n,
                  ),
                  _StatusBadge(
                    status: visitReport.status,
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
                ],
              ),
              const SizedBox(height: 16),
              // Purpose Section
              _SectionTitle(
                title: l10n.purpose,
                theme: theme,
                colorScheme: colorScheme,
              ),
              const SizedBox(height: 8),
              _InfoCard(
                theme: theme,
                colorScheme: colorScheme,
                children: [
                  if (visitReport.purpose != null && visitReport.purpose!.isNotEmpty)
                    Text(
                      visitReport.purpose!,
                      style: theme.textTheme.bodyLarge?.copyWith(
                        color: colorScheme.onSurface,
                      ),
                    )
                  else
                    Text(
                      'No purpose specified',
                      style: theme.textTheme.bodyLarge?.copyWith(
                        color: colorScheme.onSurface.withValues(alpha:0.5),
                        fontStyle: FontStyle.italic,
                      ),
                    ),
                ],
              ),
              const SizedBox(height: 16),
              // Notes Section
              if (visitReport.notes != null && visitReport.notes!.isNotEmpty) ...[
                _SectionTitle(
                  title: l10n.notes,
                  theme: theme,
                  colorScheme: colorScheme,
                ),
                const SizedBox(height: 8),
                _InfoCard(
                  theme: theme,
                  colorScheme: colorScheme,
                  children: [
                    Text(
                      visitReport.notes!,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: colorScheme.onSurface,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
              ],
              // Visit Information
              _SectionTitle(
                title: l10n.visitInformation,
                theme: theme,
                colorScheme: colorScheme,
              ),
              const SizedBox(height: 8),
              _InfoCard(
                theme: theme,
                colorScheme: colorScheme,
                children: [
                  // Type
                  _InfoRow(
                    icon: Icons.category_outlined,
                    label: 'Type',
                    value: visitReport.type.toUpperCase(),
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
                  // Account/Deal/Lead
                  if (visitReport.account != null) ...[
                    _InfoRow(
                      icon: Icons.business_outlined,
                      label: l10n.accounts,
                      value: visitReport.account!.name,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                    if (visitReport.contact != null)
                      _InfoRow(
                        icon: Icons.person_outline,
                        label: l10n.contacts,
                        value: visitReport.contact!.name,
                        theme: theme,
                        colorScheme: colorScheme,
                      ),
                  ] else if (visitReport.deal != null) ...[
                    _InfoRow(
                      icon: Icons.handshake_outlined,
                      label: l10n.deal,
                      value: visitReport.deal!.title,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                    if (visitReport.deal!.account != null)
                      _InfoRow(
                        icon: Icons.business_outlined,
                        label: l10n.accounts,
                        value: visitReport.deal!.account!.name,
                        theme: theme,
                        colorScheme: colorScheme,
                      ),
                    if (visitReport.contact != null)
                      _InfoRow(
                        icon: Icons.person_outline,
                        label: l10n.contacts,
                        value: visitReport.contact!.name,
                        theme: theme,
                        colorScheme: colorScheme,
                      ),
                  ] else if (visitReport.lead != null) ...[
                    _InfoRow(
                      icon: Icons.person_outline,
                      label: l10n.lead,
                      value: '${visitReport.lead!.firstName} ${visitReport.lead!.lastName ?? ''}'.trim(),
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                    if (visitReport.lead!.companyName != null)
                      _InfoRow(
                        icon: Icons.business_outlined,
                        label: 'Company',
                        value: visitReport.lead!.companyName!,
                        theme: theme,
                        colorScheme: colorScheme,
                      ),
                  ],
                  // Visit Date
                  _InfoRow(
                    icon: Icons.calendar_today_outlined,
                    label: l10n.visitDate,
                    value: _formatVisitDate(visitReport.visitDate, l10n.locale),
                    theme: theme,
                    colorScheme: colorScheme,
                  ),
                ],
              ),
              const SizedBox(height: 16),
              // Check-in/out Information
              _SectionTitle(
                title: l10n.checkInOutStatus,
                theme: theme,
                colorScheme: colorScheme,
              ),
              const SizedBox(height: 8),
              _InfoCard(
                theme: theme,
                colorScheme: colorScheme,
                children: [
                  if (visitReport.checkInTime != null) ...[
                    _InfoRow(
                      icon: Icons.login,
                      label: l10n.checkInTime,
                      value: _formatDateTime(visitReport.checkInTime!),
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                    if (visitReport.checkInLocation != null)
                      _InfoRow(
                        icon: Icons.location_on_outlined,
                        label: l10n.checkInLocation,
                        value:
                            '${visitReport.checkInLocation!.latitude.toStringAsFixed(6)}, ${visitReport.checkInLocation!.longitude.toStringAsFixed(6)}',
                        theme: theme,
                        colorScheme: colorScheme,
                      ),
                  ] else
                    _InfoRow(
                      icon: Icons.login,
                      label: l10n.checkIn,
                      value: l10n.notCheckedIn,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                  if (visitReport.checkOutTime != null) ...[
                    _InfoRow(
                      icon: Icons.logout,
                      label: l10n.checkOutTime,
                      value: _formatDateTime(visitReport.checkOutTime!),
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                    if (visitReport.checkOutLocation != null)
                      _InfoRow(
                        icon: Icons.location_on_outlined,
                        label: l10n.checkOutLocation,
                        value:
                            '${visitReport.checkOutLocation!.latitude.toStringAsFixed(6)}, ${visitReport.checkOutLocation!.longitude.toStringAsFixed(6)}',
                        theme: theme,
                        colorScheme: colorScheme,
                      ),
                  ] else if (visitReport.checkInTime != null)
                    _InfoRow(
                      icon: Icons.logout,
                      label: l10n.checkOut,
                      value: l10n.notCheckedOut,
                      theme: theme,
                      colorScheme: colorScheme,
                    ),
                ],
              ),
              const SizedBox(height: 16),
              // Photos (including selfie from check-in)
              if (visitReport.photoUrls != null &&
                  visitReport.photoUrls!.isNotEmpty) ...[
                _SectionTitle(
                  title: l10n.photos,
                  theme: theme,
                  colorScheme: colorScheme,
                ),
                const SizedBox(height: 8),
                SizedBox(
                  height: 140,
                  child: ListView.builder(
                    scrollDirection: Axis.horizontal,
                    itemCount: visitReport.photoUrls!.length,
                    itemBuilder: (context, index) {
                      // First photo is typically the selfie from check-in
                      final isCheckInSelfie =
                          visitReport.checkInTime != null && index == 0;
                      return GestureDetector(
                        onTap: () {
                          // Show fullscreen image
                          showDialog(
                            context: context,
                            builder: (context) => Dialog(
                              backgroundColor: Colors.transparent,
                              child: Stack(
                                children: [
                                  Center(
                                    child: Image.network(
                                      visitReport.photoUrls![index],
                                      fit: BoxFit.contain,
                                      errorBuilder:
                                          (context, error, stackTrace) {
                                            return Center(
                                              child: Icon(
                                                Icons.broken_image,
                                                color: Colors.white,
                                                size: 48,
                                              ),
                                            );
                                          },
                                    ),
                                  ),
                                  Positioned(
                                    top: 40,
                                    right: 20,
                                    child: IconButton(
                                      icon: const Icon(
                                        Icons.close,
                                        color: Colors.white,
                                        size: 28,
                                      ),
                                      onPressed: () =>
                                          Navigator.of(context).pop(),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          );
                        },
                        child: Container(
                          margin: const EdgeInsets.only(right: 12),
                          width: 140,
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(16),
                            border: Border.all(
                              color: isCheckInSelfie
                                  ? Colors.green.withValues(alpha:0.6)
                                  : colorScheme.outline.withValues(alpha:0.2),
                              width: isCheckInSelfie ? 2.5 : 1,
                            ),
                            boxShadow: [
                              BoxShadow(
                                color: Colors.black.withValues(alpha: 0.05),
                                blurRadius: 3,
                                offset: const Offset(0, 1),
                              ),
                            ],
                          ),
                          child: Stack(
                            children: [
                              ClipRRect(
                                borderRadius: BorderRadius.circular(16),
                                child: Image.network(
                                  visitReport.photoUrls![index],
                                  fit: BoxFit.cover,
                                  width: double.infinity,
                                  height: double.infinity,
                                  loadingBuilder:
                                      (context, child, loadingProgress) {
                                        if (loadingProgress == null) {
                                          return child;
                                        }
                                        return Center(
                                          child: CircularProgressIndicator(
                                            value:
                                                loadingProgress
                                                        .expectedTotalBytes !=
                                                    null
                                                ? loadingProgress
                                                          .cumulativeBytesLoaded /
                                                      loadingProgress
                                                          .expectedTotalBytes!
                                                : null,
                                          ),
                                        );
                                      },
                                  errorBuilder: (context, error, stackTrace) {
                                    return Container(
                                      color:
                                          colorScheme.surfaceContainerHighest,
                                      child: Center(
                                        child: Icon(
                                          Icons.broken_image,
                                          color: colorScheme.onSurface
                                              .withValues(alpha:0.3),
                                          size: 32,
                                        ),
                                      ),
                                    );
                                  },
                                ),
                              ),
                              if (isCheckInSelfie)
                                Positioned(
                                  top: 6,
                                  right: 6,
                                  child: Container(
                                    padding: const EdgeInsets.symmetric(
                                      horizontal: 8,
                                      vertical: 4,
                                    ),
                                    decoration: BoxDecoration(
                                      color: Colors.green.withValues(alpha:0.95),
                                      borderRadius: BorderRadius.circular(6),
                                      boxShadow: [
                                        BoxShadow(
                                          color: Colors.black.withValues(alpha: 0.1),
                                          blurRadius: 3,
                                          offset: const Offset(0, 1),
                                        ),
                                      ],
                                    ),
                                    child: Row(
                                      mainAxisSize: MainAxisSize.min,
                                      children: [
                                        const Icon(
                                          Icons.camera_alt,
                                          color: Colors.white,
                                          size: 12,
                                        ),
                                        const SizedBox(width: 4),
                                        Text(
                                          'Check-in',
                                          style: theme.textTheme.bodySmall
                                              ?.copyWith(
                                                color: Colors.white,
                                                fontSize: 10,
                                                fontWeight: FontWeight.bold,
                                              ),
                                        ),
                                      ],
                                    ),
                                  ),
                                ),
                              // Tap indicator
                              Positioned(
                                bottom: 8,
                                left: 8,
                                right: 8,
                                child: Container(
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: 8,
                                    vertical: 4,
                                  ),
                                  decoration: BoxDecoration(
                                    color: Colors.black.withValues(alpha:0.5),
                                    borderRadius: BorderRadius.circular(6),
                                  ),
                                  child: Row(
                                    mainAxisSize: MainAxisSize.min,
                                    children: [
                                      Icon(
                                        Icons.zoom_in,
                                        color: Colors.white,
                                        size: 12,
                                      ),
                                      const SizedBox(width: 4),
                                      Text(
                                        'Tap to view',
                                        style: theme.textTheme.bodySmall
                                            ?.copyWith(
                                              color: Colors.white,
                                              fontSize: 10,
                                            ),
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  ),
                ),
                const SizedBox(height: 16),
              ],
              // Action Buttons
              // Only show actions if status is 'draft' (can update, delete, check-in, check-out, submit)
              if (visitReport.status.toLowerCase() == 'draft') ...[
                // Check-in button (only if not checked in)
                if (visitReport.checkInTime == null)
                  FilledButton.icon(
                    onPressed: formState.isLoading ? null : _handleCheckIn,
                    icon: formState.isLoading
                        ? SizedBox(
                            width: 16,
                            height: 16,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              valueColor: AlwaysStoppedAnimation<Color>(
                                Colors.white,
                              ),
                            ),
                          )
                        : const Icon(Icons.login),
                    label: Text(l10n.checkIn),
                    style: FilledButton.styleFrom(
                      minimumSize: const Size(double.infinity, 48),
                      backgroundColor: Colors.green,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                  ),
                // Check-out button (only if checked in but not checked out)
                if (visitReport.checkInTime != null &&
                    visitReport.checkOutTime == null) ...[
                  const SizedBox(height: 12),
                  FilledButton.icon(
                    onPressed: formState.isLoading ? null : _handleCheckOut,
                    icon: formState.isLoading
                        ? SizedBox(
                            width: 16,
                            height: 16,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              valueColor: AlwaysStoppedAnimation<Color>(
                                Colors.white,
                              ),
                            ),
                          )
                        : const Icon(Icons.logout),
                    label: Text(l10n.checkOut),
                    style: FilledButton.styleFrom(
                      minimumSize: const Size(double.infinity, 48),
                      backgroundColor: Colors.blue,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                  ),
                ],
                // Submit button (only if checked in and checked out)
                if (visitReport.checkInTime != null &&
                    visitReport.checkOutTime != null) ...[
                  const SizedBox(height: 12),
                  FilledButton.icon(
                    onPressed: formState.isLoading ? null : _handleSubmit,
                    icon: formState.isLoading
                        ? SizedBox(
                            width: 16,
                            height: 16,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              valueColor: AlwaysStoppedAnimation<Color>(
                                Colors.white,
                              ),
                            ),
                          )
                        : const Icon(Icons.send),
                    label: Text(l10n.submitVisitReport),
                    style: FilledButton.styleFrom(
                      minimumSize: const Size(double.infinity, 48),
                      backgroundColor: Colors.orange,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                  ),
                ],
                // Update and Delete buttons (always available for draft status)
                const SizedBox(height: 12),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton.icon(
                        onPressed: formState.isLoading ? null : _handleUpdate,
                        icon: const Icon(Icons.edit),
                        label: Text(l10n.updateVisitReport),
                        style: OutlinedButton.styleFrom(
                          minimumSize: const Size(double.infinity, 48),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: OutlinedButton.icon(
                        onPressed: formState.isLoading ? null : _handleDelete,
                        icon: const Icon(Icons.delete_outline),
                        label: Text(l10n.delete),
                        style: OutlinedButton.styleFrom(
                          minimumSize: const Size(double.infinity, 48),
                          foregroundColor: Colors.red,
                          side: const BorderSide(color: Colors.red),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ],
              // Show status info for submitted/approved/rejected
              if (visitReport.status.toLowerCase() == 'submitted') ...[
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: colorScheme.secondaryContainer.withValues(alpha:0.3),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: colorScheme.secondary.withValues(alpha:0.3),
                    ),
                  ),
                  child: Row(
                    children: [
                      Icon(
                        Icons.pending_actions,
                        color: colorScheme.secondary,
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Text(
                          'Visit report submitted. Waiting for approval.',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: colorScheme.onSecondaryContainer,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
              if (visitReport.status.toLowerCase() == 'approved') ...[
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.green.withValues(alpha:0.1),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: Colors.green.withValues(alpha:0.3),
                    ),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.check_circle, color: Colors.green),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Text(
                          'Visit report approved.',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: Colors.green.shade700,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
              if (visitReport.status.toLowerCase() == 'rejected') ...[
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.red.withValues(alpha:0.1),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: Colors.red.withValues(alpha:0.3),
                    ),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          const Icon(Icons.cancel, color: Colors.red),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Text(
                              'Visit report rejected.',
                              style: theme.textTheme.bodyMedium?.copyWith(
                                color: Colors.red.shade700,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ),
                        ],
                      ),
                      if (visitReport.rejectionReason != null) ...[
                        const SizedBox(height: 8),
                        Text(
                          visitReport.rejectionReason!,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: Colors.red.shade700,
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
              ],
            ],
          ),
        ),
        loading: () => const LoadingWidget(),
        error: (error, stack) => ErrorStateWidget(
          message: error.toString().replaceFirst('Exception: ', ''),
          onRetry: () =>
              ref.invalidate(visitReportDetailProvider(widget.visitReportId)),
        ),
      ),
    );
  }
}

class _TypeBadge extends StatelessWidget {
  const _TypeBadge({
    required this.type,
    required this.theme,
    required this.colorScheme,
    required this.l10n,
  });

  final String type;
  final ThemeData theme;
  final ColorScheme colorScheme;
  final AppLocalizations l10n;

  @override
  Widget build(BuildContext context) {
    String typeText;
    Color backgroundColor;
    Color textColor;

    // Use type from API response
    switch (type.toLowerCase()) {
      case 'lead':
        typeText = l10n.lead;
        backgroundColor = Colors.purple.withValues(alpha:0.1);
        textColor = Colors.purple;
        break;
      case 'deal':
        typeText = l10n.deal;
        backgroundColor = Colors.orange.withValues(alpha:0.1);
        textColor = Colors.orange;
        break;
      case 'account':
      default:
        // Default to Account type
        typeText = l10n.accounts;
        backgroundColor = colorScheme.primary.withValues(alpha:0.1);
        textColor = colorScheme.primary;
        break;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        typeText.toUpperCase(),
        style: theme.textTheme.bodySmall?.copyWith(
          color: textColor,
          fontWeight: FontWeight.w600,
          fontSize: 10,
        ),
      ),
    );
  }
}

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({
    required this.status,
    required this.theme,
    required this.colorScheme,
  });

  final String status;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    Color backgroundColor;
    Color textColor;
    String displayText;

    switch (status.toLowerCase()) {
      case 'draft':
        backgroundColor = colorScheme.onSurface.withValues(alpha:0.1);
        textColor = colorScheme.onSurface.withValues(alpha:0.7);
        displayText = 'DRAFT';
        break;
      case 'in_progress':
        backgroundColor = Colors.orange.withValues(alpha:0.1);
        textColor = Colors.orange;
        displayText = 'IN PROGRESS';
        break;
      case 'completed':
        backgroundColor = colorScheme.primary.withValues(alpha:0.1);
        textColor = colorScheme.primary;
        displayText = 'COMPLETED';
        break;
      case 'approved':
        backgroundColor = Colors.green.withValues(alpha:0.1);
        textColor = Colors.green;
        displayText = 'APPROVED';
        break;
      case 'rejected':
        backgroundColor = colorScheme.error.withValues(alpha:0.1);
        textColor = colorScheme.error;
        displayText = 'REJECTED';
        break;
      default:
        backgroundColor = colorScheme.onSurface.withValues(alpha:0.1);
        textColor = colorScheme.onSurface.withValues(alpha:0.7);
        displayText = status.toUpperCase();
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        displayText,
        style: theme.textTheme.bodySmall?.copyWith(
          color: textColor,
          fontWeight: FontWeight.w600,
          fontSize: 12,
        ),
      ),
    );
  }
}

class _SectionTitle extends StatelessWidget {
  const _SectionTitle({
    required this.title,
    required this.theme,
    required this.colorScheme,
  });

  final String title;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    return Text(
      title,
      style: theme.textTheme.titleMedium?.copyWith(
        fontWeight: FontWeight.w600,
        color: colorScheme.onSurface,
      ),
    );
  }
}

class _InfoCard extends StatelessWidget {
  const _InfoCard({
    required this.children,
    required this.theme,
    required this.colorScheme,
  });

  final List<Widget> children;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: colorScheme.surface,
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 3,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      padding: const EdgeInsets.all(16),
      child: Column(
        children: children.asMap().entries.map((entry) {
          final index = entry.key;
          final child = entry.value;
          if (index == children.length - 1) {
            return child;
          }
          return Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: child,
          );
        }).toList(),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  const _InfoRow({
    required this.icon,
    required this.label,
    required this.value,
    required this.theme,
    required this.colorScheme,
  });

  final IconData icon;
  final String label;
  final String value;
  final ThemeData theme;
  final ColorScheme colorScheme;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Icon(icon, size: 20, color: colorScheme.onSurface.withValues(alpha:0.7)),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: theme.textTheme.bodySmall?.copyWith(
                  color: colorScheme.onSurface.withValues(alpha:0.7),
                ),
              ),
              const SizedBox(height: 4),
              Text(
                value,
                style: theme.textTheme.bodyMedium?.copyWith(
                  fontWeight: FontWeight.w500,
                  color: colorScheme.onSurface,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}
