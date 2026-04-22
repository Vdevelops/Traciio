import '../../domain/models/google_calendar_event.dart';
import '../../domain/models/google_calendar_status.dart';
import '../../domain/repositories/google_calendar_repository.dart';
import '../datasources/google_calendar_remote_datasource.dart';

/// Implementation of Google Calendar repository
class GoogleCalendarRepositoryImpl implements GoogleCalendarRepository {
  final GoogleCalendarRemoteDataSource _dataSource;

  GoogleCalendarRepositoryImpl(this._dataSource);

  @override
  Future<GoogleCalendarStatus> getConnectionStatus() async {
    return await _dataSource.getConnectionStatus();
  }

  @override
  Future<String> getAuthUrl() async {
    return await _dataSource.getAuthUrl();
  }

  @override
  Future<void> disconnect() async {
    return await _dataSource.disconnect();
  }

  @override
  Future<GoogleCalendarEvent> syncSchedule(String scheduleId) async {
    return await _dataSource.syncSchedule(scheduleId);
  }

  @override
  Future<void> unsyncSchedule(String scheduleId) async {
    return await _dataSource.unsyncSchedule(scheduleId);
  }

  @override
  Future<void> exchangeCode(String code, String state) async {
    return await _dataSource.exchangeCode(code, state);
  }
}
