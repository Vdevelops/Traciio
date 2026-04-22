import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../../core/network/connectivity_service.dart';
import '../../../core/sync/auto_sync_manager.dart';
import '../data/models/schedule_model.dart';
import '../data/schedule_repository.dart';
import 'schedule_state.dart';

final scheduleRepositoryProvider = Provider<ScheduleRepository>((ref) {
  return ScheduleRepository(
    ApiClient.dio,
    ref.watch(connectivityServiceProvider),
  );
});

final scheduleListProvider =
    NotifierProvider<ScheduleListNotifier, ScheduleListState>(
      ScheduleListNotifier.new,
    );

final scheduleDetailProvider =
    NotifierProvider.family<
      ScheduleDetailNotifier,
      ScheduleDetailState,
      String
    >(ScheduleDetailNotifier.new);

final scheduleFormProvider =
    NotifierProvider<ScheduleFormNotifier, ScheduleFormState>(
      ScheduleFormNotifier.new,
    );

class ScheduleListNotifier extends Notifier<ScheduleListState>
    with WidgetsBindingObserver {
  late final ScheduleRepository _repository;

  @override
  ScheduleListState build() {
    _repository = ref.read(scheduleRepositoryProvider);

    // Listen to app lifecycle changes for auto-sync
    WidgetsBinding.instance.addObserver(this);
    ref.onDispose(() {
      WidgetsBinding.instance.removeObserver(this);
    });

    // Register with AutoSyncManager for centralized sync
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref
          .read(autoSyncManagerProvider.notifier)
          .registerFeature('schedules', () => silentRefresh());
    });

    return const ScheduleListState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);
    if (state == AppLifecycleState.resumed) {
      // Auto sync when app comes back to foreground
      ref.read(autoSyncManagerProvider.notifier).syncFeature('schedules');
    }
  }

  Future<void> loadSchedules({bool forceRefresh = false}) async {
    if (state.isLoading && !forceRefresh) return;

    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final response = await _repository.getSchedules(
        page: 1,
        search: state.search,
        status: state.selectedStatus,
        scheduledFrom: state.scheduledFrom,
        scheduledTo: state.scheduledTo,
        forceRefresh: forceRefresh,
        onBackgroundUpdate: (freshData) {
          // Update UI when fresh data arrives from background fetch
          _updateStateWithFreshData(freshData);
        },
      );

      state = state.copyWith(
        isLoading: false,
        schedules: response.items,
        page: 1,
        hasNextPage: response.pagination.hasNextPage,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, errorMessage: e.toString());
    }
  }

  /// Silent refresh (dipanggil oleh AutoSyncManager)
  Future<void> silentRefresh() async {
    try {
      final response = await _repository.getSchedules(
        page: 1,
        search: state.search,
        status: state.selectedStatus,
        scheduledFrom: state.scheduledFrom,
        scheduledTo: state.scheduledTo,
        forceRefresh: true,
      );

      state = state.copyWith(
        schedules: response.items,
        page: 1,
        hasNextPage: response.pagination.hasNextPage,
      );
    } catch (e) {
      // Silent fail
    }
  }

  /// Update state with fresh data from background fetch
  void _updateStateWithFreshData(ScheduleListResponse freshData) {
    state = state.copyWith(
      schedules: freshData.items,
      page: 1,
      hasNextPage: freshData.pagination.hasNextPage,
    );
  }

  Future<void> loadMore() async {
    if (state.isLoading || !state.hasNextPage) return;

    state = state.copyWith(isLoading: true);

    try {
      final nextPage = state.page + 1;
      final response = await _repository.getSchedules(
        page: nextPage,
        search: state.search,
        status: state.selectedStatus,
        scheduledFrom: state.scheduledFrom,
        scheduledTo: state.scheduledTo,
      );

      state = state.copyWith(
        isLoading: false,
        schedules: [...state.schedules, ...response.items],
        page: nextPage,
        hasNextPage: response.pagination.hasNextPage,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, errorMessage: e.toString());
    }
  }

  void updateSearchQuery(String search) {
    state = state.copyWith(search: search);
    loadSchedules(forceRefresh: true);
  }

  void setSearch(String search) {
    state = state.copyWith(search: search);
    loadSchedules(forceRefresh: true);
  }

  void setStatus(String? status) {
    state = state.copyWith(selectedStatus: status);
    loadSchedules(forceRefresh: true);
  }

  void setDateRange(DateTime? from, DateTime? to) {
    state = state.copyWith(scheduledFrom: from, scheduledTo: to);
    loadSchedules(forceRefresh: true);
  }

  void clearFilters() {
    state = const ScheduleListState();
    loadSchedules(forceRefresh: true);
  }
}

class ScheduleDetailNotifier extends Notifier<ScheduleDetailState> {
  ScheduleDetailNotifier(this.id);
  final String id;
  late final ScheduleRepository _repository;

  @override
  ScheduleDetailState build() {
    _repository = ref.read(scheduleRepositoryProvider);
    Future.microtask(() => loadSchedule());
    return const ScheduleDetailState();
  }

  Future<void> loadSchedule() async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final schedule = await _repository.getScheduleById(id);
      state = state.copyWith(isLoading: false, schedule: schedule);
    } catch (e) {
      state = state.copyWith(isLoading: false, errorMessage: e.toString());
    }
  }

  Future<bool> deleteSchedule() async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      await _repository.deleteSchedule(id);
      state = state.copyWith(isLoading: false);
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, errorMessage: e.toString());
      return false;
    }
  }

  Future<void> updateStatus(String status) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final updated = await _repository.updateStatus(id, status);
      state = state.copyWith(isLoading: false, schedule: updated);
    } catch (e) {
      state = state.copyWith(isLoading: false, errorMessage: e.toString());
    }
  }

  Future<bool> syncToGoogleCalendar() async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final updated = await _repository.syncToGoogleCalendar(id);
      state = state.copyWith(isLoading: false, schedule: updated);
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, errorMessage: e.toString());
      return false;
    }
  }

  Future<bool> unsyncFromGoogleCalendar() async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final updated = await _repository.unsyncFromGoogleCalendar(id);
      state = state.copyWith(isLoading: false, schedule: updated);
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, errorMessage: e.toString());
      return false;
    }
  }
}

class ScheduleFormNotifier extends Notifier<ScheduleFormState> {
  late final ScheduleRepository _repository;

  @override
  ScheduleFormState build() {
    _repository = ref.read(scheduleRepositoryProvider);
    return const ScheduleFormState();
  }

  Future<ScheduleModel?> loadInitialData(String id) async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      final schedule = await _repository.getScheduleById(id);
      state = state.copyWith(isLoading: false, initialSchedule: schedule);
      return schedule;
    } catch (e) {
      state = state.copyWith(isLoading: false, errorMessage: e.toString());
      return null;
    }
  }

  Future<bool> submit(Map<String, dynamic> data) async {
    state = state.copyWith(
      isLoading: true,
      errorMessage: null,
      isSuccess: false,
    );

    try {
      final request = ScheduleRequest(
        title: data['title'] as String,
        description: data['description'] as String?,
        scheduledAt: DateTime.parse(data['scheduled_at'] as String),
        reminderMinutesBefore: data['reminder_minutes_before'] as int?,
        status: data['status'] as String? ?? 'pending',
        taskId: data['task_id'] as String?,
        location: data['location'] as String?,
        activityType: data['activity_type'] as String?,
      );

      if (state.initialSchedule != null) {
        await _repository.updateSchedule(state.initialSchedule!.id, request);
      } else {
        await _repository.createSchedule(request);
      }
      state = state.copyWith(isLoading: false, isSuccess: true);
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, errorMessage: e.toString());
      return false;
    }
  }
}
