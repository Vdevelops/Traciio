import '../data/models/deal_model.dart';

class PipelineState {
  const PipelineState({
    this.stages = const [],
    this.deals = const [],
    this.isLoading = false,
    this.isLoadingMore = false,
    this.errorMessage,
    this.pagination,
    this.searchQuery = '',
    this.selectedStageId,
    this.isOffline = false,
    this.formData = const {},
    this.products = const [],
  });

  final List<PipelineStage> stages;
  final List<Deal> deals;
  final bool isLoading;
  final bool isLoadingMore;
  final String? errorMessage;
  final Pagination? pagination;
  final String searchQuery;
  final String? selectedStageId;
  final bool isOffline;
  final Map<String, dynamic> formData;
  final List<Product> products;

  PipelineState copyWith({
    List<PipelineStage>? stages,
    List<Deal>? deals,
    bool? isLoading,
    bool? isLoadingMore,
    String? errorMessage,
    Pagination? pagination,
    String? searchQuery,
    String? selectedStageId,
    bool? isOffline,
    Map<String, dynamic>? formData,
    List<Product>? products,
    bool clearDeals = false,
  }) {
    return PipelineState(
      stages: stages ?? this.stages,
      deals: clearDeals ? const [] : (deals ?? this.deals),
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      errorMessage: errorMessage,
      pagination: pagination ?? this.pagination,
      searchQuery: searchQuery ?? this.searchQuery,
      selectedStageId: selectedStageId ?? this.selectedStageId,
      isOffline: isOffline ?? this.isOffline,
      formData: formData ?? this.formData,
      products: products ?? this.products,
    );
  }
}
