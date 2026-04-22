import '../data/models/optimized_route.dart';
import '../data/route_optimization_repository.dart';

class RouteListState {
  const RouteListState({
    this.routes = const [],
    this.isLoading = false,
    this.isLoadingMore = false,
    this.errorMessage,
    this.pagination,
    this.isOffline = false,
  });

  final List<OptimizedRoute> routes;
  final bool isLoading;
  final bool isLoadingMore;
  final String? errorMessage;
  final Pagination? pagination;
  final bool isOffline;

  RouteListState copyWith({
    List<OptimizedRoute>? routes,
    bool? isLoading,
    bool? isLoadingMore,
    String? errorMessage,
    Pagination? pagination,
    bool? isOffline,
    bool clearRoutes = false,
  }) {
    return RouteListState(
      routes: clearRoutes ? const [] : (routes ?? this.routes),
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      errorMessage: errorMessage,
      pagination: pagination ?? this.pagination,
      isOffline: isOffline ?? this.isOffline,
    );
  }
}

class RouteOptimizationState {
  const RouteOptimizationState({
    this.isLoading = false,
    this.errorMessage,
    this.isOptimizing = false,
    this.optimizedRoute,
  });

  final bool isLoading;
  final String? errorMessage;
  final bool isOptimizing;
  final OptimizedRoute? optimizedRoute;

  RouteOptimizationState copyWith({
    bool? isLoading,
    String? errorMessage,
    bool? isOptimizing,
    OptimizedRoute? optimizedRoute,
    bool clearOptimizedRoute = false,
  }) {
    return RouteOptimizationState(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
      isOptimizing: isOptimizing ?? this.isOptimizing,
      optimizedRoute: clearOptimizedRoute ? null : (optimizedRoute ?? this.optimizedRoute),
    );
  }
}

class DistanceCalculationState {
  const DistanceCalculationState({
    this.isCalculating = false,
    this.errorMessage,
    this.distance,
    this.duration,
    this.distanceFormatted,
    this.durationFormatted,
  });

  final bool isCalculating;
  final String? errorMessage;
  final double? distance;
  final int? duration;
  final String? distanceFormatted;
  final String? durationFormatted;

  DistanceCalculationState copyWith({
    bool? isCalculating,
    String? errorMessage,
    double? distance,
    int? duration,
    String? distanceFormatted,
    String? durationFormatted,
    bool clearResult = false,
  }) {
    return DistanceCalculationState(
      isCalculating: isCalculating ?? this.isCalculating,
      errorMessage: errorMessage,
      distance: clearResult ? null : (distance ?? this.distance),
      duration: clearResult ? null : (duration ?? this.duration),
      distanceFormatted: clearResult ? null : (distanceFormatted ?? this.distanceFormatted),
      durationFormatted: clearResult ? null : (durationFormatted ?? this.durationFormatted),
    );
  }
}





