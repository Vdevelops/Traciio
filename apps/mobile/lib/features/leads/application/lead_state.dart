import '../data/models/lead_model.dart';

class LeadListState {
  const LeadListState({
    this.leads = const [],
    this.isLoading = false,
    this.isLoadingMore = false,
    this.errorMessage,
    this.pagination,
    this.searchQuery = '',
    this.selectedStatus = '',
    this.selectedSource = '',
    this.selectedIndustry = '',
    this.selectedProvince = '',
  });

  final List<Lead> leads;
  final bool isLoading;
  final bool isLoadingMore;
  final String? errorMessage;
  final Pagination? pagination;
  final String searchQuery;
  final String selectedStatus;
  final String selectedSource;
  final String selectedIndustry;
  final String selectedProvince;

  bool get hasActiveFilters =>
      selectedStatus.isNotEmpty ||
      selectedSource.isNotEmpty ||
      selectedIndustry.isNotEmpty ||
      selectedProvince.isNotEmpty;

  String? get error => errorMessage;

  LeadListState copyWith({
    List<Lead>? leads,
    bool? isLoading,
    bool? isLoadingMore,
    String? errorMessage,
    Pagination? pagination,
    String? searchQuery,
    String? selectedStatus,
    String? selectedSource,
    String? selectedIndustry,
    String? selectedProvince,
    bool clearLeads = false,
  }) {
    return LeadListState(
      leads: clearLeads ? const [] : (leads ?? this.leads),
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      errorMessage: errorMessage,
      pagination: pagination ?? this.pagination,
      searchQuery: searchQuery ?? this.searchQuery,
      selectedStatus: selectedStatus ?? this.selectedStatus,
      selectedSource: selectedSource ?? this.selectedSource,
      selectedIndustry: selectedIndustry ?? this.selectedIndustry,
      selectedProvince: selectedProvince ?? this.selectedProvince,
    );
  }
}

class LeadFormState {
  const LeadFormState({
    this.isLoading = false,
    this.errorMessage,
    this.formData,
  });

  final bool isLoading;
  final String? errorMessage;
  final LeadFormData? formData;

  LeadFormState copyWith({
    bool? isLoading,
    String? errorMessage,
    LeadFormData? formData,
  }) {
    return LeadFormState(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
      formData: formData ?? this.formData,
    );
  }
}
