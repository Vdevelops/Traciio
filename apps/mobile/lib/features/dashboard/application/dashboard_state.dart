import '../data/models/dashboard.dart';

class DashboardState {
  final MobileDashboardOverview? overview;
  final List<MobileVisit>? visits;
  final List<MobileTask>? tasks;

  final bool isLoading;
  final bool isLoadingOverview;
  final bool isLoadingVisits;
  final bool isLoadingTasks;

  final String? errorMessage;
  final String selectedPeriod;
  final bool isOffline;

  // Filter for visits tab
  final String visitStatusFilter; // 'draft', 'submitted', 'approved', 'rejected'

  DashboardState({
    this.overview,
    this.visits,
    this.tasks,
    this.isLoading = false,
    this.isLoadingOverview = false,
    this.isLoadingVisits = false,
    this.isLoadingTasks = false,
    this.errorMessage,
    this.selectedPeriod = 'today',
    this.isOffline = false,
    this.visitStatusFilter = 'draft',
  });

  DashboardState copyWith({
    MobileDashboardOverview? overview,
    List<MobileVisit>? visits,
    List<MobileTask>? tasks,
    bool? isLoading,
    bool? isLoadingOverview,
    bool? isLoadingVisits,
    bool? isLoadingTasks,
    String? errorMessage,
    String? selectedPeriod,
    bool? isOffline,
    String? visitStatusFilter,
  }) {
    return DashboardState(
      overview: overview ?? this.overview,
      visits: visits ?? this.visits,
      tasks: tasks ?? this.tasks,
      isLoading: isLoading ?? this.isLoading,
      isLoadingOverview: isLoadingOverview ?? this.isLoadingOverview,
      isLoadingVisits: isLoadingVisits ?? this.isLoadingVisits,
      isLoadingTasks: isLoadingTasks ?? this.isLoadingTasks,
      errorMessage: errorMessage ?? this.errorMessage,
      selectedPeriod: selectedPeriod ?? this.selectedPeriod,
      isOffline: isOffline ?? this.isOffline,
      visitStatusFilter: visitStatusFilter ?? this.visitStatusFilter,
    );
  }
}
