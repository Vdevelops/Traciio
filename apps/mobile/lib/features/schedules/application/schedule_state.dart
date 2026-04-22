import '../data/models/schedule_model.dart';

class ScheduleListState {
  final bool isLoading;
  final String? errorMessage;
  final List<ScheduleModel> schedules;
  final int page;
  final bool hasNextPage;
  final String search;
  final String? selectedStatus;
  final DateTime? scheduledFrom;
  final DateTime? scheduledTo;

  const ScheduleListState({
    this.isLoading = false,
    this.errorMessage,
    this.schedules = const [],
    this.page = 1,
    this.hasNextPage = false,
    this.search = '',
    this.selectedStatus,
    this.scheduledFrom,
    this.scheduledTo,
  });

  ScheduleListState copyWith({
    bool? isLoading,
    String? errorMessage,
    List<ScheduleModel>? schedules,
    int? page,
    bool? hasNextPage,
    String? search,
    String? selectedStatus,
    DateTime? scheduledFrom,
    DateTime? scheduledTo,
  }) {
    return ScheduleListState(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
      schedules: schedules ?? this.schedules,
      page: page ?? this.page,
      hasNextPage: hasNextPage ?? this.hasNextPage,
      search: search ?? this.search,
      selectedStatus: selectedStatus ?? this.selectedStatus,
      scheduledFrom: scheduledFrom ?? this.scheduledFrom,
      scheduledTo: scheduledTo ?? this.scheduledTo,
    );
  }

  bool get hasActiveFilters =>
      search.isNotEmpty ||
      selectedStatus != null ||
      scheduledFrom != null ||
      scheduledTo != null;
}

class ScheduleDetailState {
  final bool isLoading;
  final String? errorMessage;
  final ScheduleModel? schedule;

  const ScheduleDetailState({
    this.isLoading = false,
    this.errorMessage,
    this.schedule,
  });

  ScheduleDetailState copyWith({
    bool? isLoading,
    String? errorMessage,
    ScheduleModel? schedule,
  }) {
    return ScheduleDetailState(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
      schedule: schedule ?? this.schedule,
    );
  }
}

class ScheduleFormState {
  final bool isLoading;
  final String? errorMessage;
  final bool isSuccess;
  final ScheduleModel? initialSchedule;

  const ScheduleFormState({
    this.isLoading = false,
    this.errorMessage,
    this.isSuccess = false,
    this.initialSchedule,
  });

  ScheduleFormState copyWith({
    bool? isLoading,
    String? errorMessage,
    bool? isSuccess,
    ScheduleModel? initialSchedule,
  }) {
    return ScheduleFormState(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
      isSuccess: isSuccess ?? this.isSuccess,
      initialSchedule: initialSchedule ?? this.initialSchedule,
    );
  }
}
