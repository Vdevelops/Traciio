import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../../core/network/api_client.dart';
import '../data/datasources/google_calendar_remote_datasource.dart';
import '../data/repositories/google_calendar_repository_impl.dart';
import '../domain/models/google_calendar_event.dart';
import '../domain/models/google_calendar_status.dart';
import '../domain/repositories/google_calendar_repository.dart';

/// Provider for Google Calendar repository
final googleCalendarRepositoryProvider = Provider<GoogleCalendarRepository>((
  ref,
) {
  final dio = ApiClient.dio;
  final dataSource = GoogleCalendarRemoteDataSource(dio);
  return GoogleCalendarRepositoryImpl(dataSource);
});

/// Provider for Google Calendar connection status
final googleCalendarStatusProvider = FutureProvider<GoogleCalendarStatus>((
  ref,
) async {
  final repository = ref.watch(googleCalendarRepositoryProvider);
  return await repository.getConnectionStatus();
});

/// Notifier for Google Calendar operations
class GoogleCalendarNotifier
    extends Notifier<AsyncValue<GoogleCalendarStatus?>> {
  GoogleCalendarRepository get _repository =>
      ref.read(googleCalendarRepositoryProvider);

  @override
  AsyncValue<GoogleCalendarStatus?> build() {
    // Load status saat initialization
    Future.microtask(() => _loadStatus());
    return const AsyncValue.loading();
  }

  /// Load connection status
  Future<void> _loadStatus() async {
    state = const AsyncValue.loading();
    try {
      final status = await _repository.getConnectionStatus();
      state = AsyncValue.data(status);
    } catch (e, stackTrace) {
      state = AsyncValue.error(e, stackTrace);
    }
  }

  /// Refresh connection status
  Future<void> refreshStatus() async {
    await _loadStatus();
  }

  /// Connect to Google Calendar (opens external browser)
  Future<void> connect() async {
    try {
      final authUrl = await _repository.getAuthUrl();

      // Open external browser for OAuth
      final uri = Uri.parse(authUrl);
      if (await canLaunchUrl(uri)) {
        await launchUrl(uri, mode: LaunchMode.externalApplication);
      } else {
        throw Exception('Could not launch OAuth URL');
      }
    } catch (e, stackTrace) {
      state = AsyncValue.error(e, stackTrace);
    }
  }

  /// Disconnect from Google Calendar
  Future<void> disconnect() async {
    state = const AsyncValue.loading();
    try {
      await _repository.disconnect();
      await refreshStatus();
    } catch (e, stackTrace) {
      state = AsyncValue.error(e, stackTrace);
    }
  }

  /// Sync schedule to Google Calendar
  Future<GoogleCalendarEvent?> syncSchedule(String scheduleId) async {
    try {
      return await _repository.syncSchedule(scheduleId);
    } catch (e) {
      // Return null on error, caller should handle
      return null;
    }
  }

  /// Unsync schedule from Google Calendar
  Future<void> unsyncSchedule(String scheduleId) async {
    try {
      await _repository.unsyncSchedule(scheduleId);
    } catch (e) {
      // Error handled silently
    }
  }

  /// Exchange authorization code for token (Option 2 - Direct Deep Link)
  Future<bool> exchangeCode(String code, String oauthState) async {
    state = const AsyncValue.loading();
    try {
      await _repository.exchangeCode(code, oauthState);
      await refreshStatus();
      return true;
    } catch (e, stackTrace) {
      state = AsyncValue.error(e, stackTrace);
      return false;
    }
  }
}

/// Provider for Google Calendar notifier
final googleCalendarNotifierProvider =
    NotifierProvider<GoogleCalendarNotifier, AsyncValue<GoogleCalendarStatus?>>(
      GoogleCalendarNotifier.new,
    );
