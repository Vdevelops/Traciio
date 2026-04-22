import '../data/models/task.dart';
import '../../../core/network/pagination.dart';

class TaskListState {
  const TaskListState({
    this.tasks = const [],
    this.isLoading = false,
    this.isLoadingMore = false,
    this.errorMessage,
    this.pagination,
    this.searchQuery = '',
    this.selectedStatus,
    this.selectedPriority,
    this.selectedType,
    this.dueDateFrom,
    this.dueDateTo,
    this.isOffline = false,
  });

  final List<Task> tasks;
  final bool isLoading;
  final bool isLoadingMore;
  final String? errorMessage;
  final Pagination? pagination;
  final String searchQuery;
  final String? selectedStatus;
  final String? selectedPriority;
  final String? selectedType;
  final DateTime? dueDateFrom;
  final DateTime? dueDateTo;
  final bool isOffline;

  TaskListState copyWith({
    List<Task>? tasks,
    bool? isLoading,
    bool? isLoadingMore,
    String? errorMessage,
    Pagination? pagination,
    String? searchQuery,
    String? selectedStatus,
    String? selectedPriority,
    String? selectedType,
    DateTime? dueDateFrom,
    DateTime? dueDateTo,
    bool? isOffline,
    bool clearTasks = false,
  }) {
    return TaskListState(
      tasks: clearTasks ? const [] : (tasks ?? this.tasks),
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      errorMessage: errorMessage,
      pagination: pagination ?? this.pagination,
      searchQuery: searchQuery ?? this.searchQuery,
      selectedStatus: selectedStatus,
      selectedPriority: selectedPriority,
      selectedType: selectedType,
      dueDateFrom: dueDateFrom,
      dueDateTo: dueDateTo,
      isOffline: isOffline ?? this.isOffline,
    );
  }
}

class TaskFormState {
  const TaskFormState({this.isLoading = false, this.errorMessage});

  final bool isLoading;
  final String? errorMessage;

  TaskFormState copyWith({bool? isLoading, String? errorMessage}) {
    return TaskFormState(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
    );
  }
}
