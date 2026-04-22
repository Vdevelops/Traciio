class Pagination {
  final int page;
  final int perPage;
  final int total;
  final int totalPages;
  final bool hasNext;
  final bool hasPrev;

  Pagination({
    required this.page,
    required this.perPage,
    required this.total,
    required this.totalPages,
    this.hasNext = false,
    this.hasPrev = false,
  });

  factory Pagination.fromJson(Map<String, dynamic> json) {
    return Pagination(
      page: json['page'] as int? ?? 1,
      perPage: json['per_page'] as int? ?? 20,
      total: json['total'] as int? ?? 0,
      totalPages: json['total_pages'] as int? ?? 1,
      hasNext:
          json['has_next'] as bool? ??
          (json['page'] != null &&
              json['total_pages'] != null &&
              (json['page'] as int) < (json['total_pages'] as int)),
      hasPrev:
          json['has_prev'] as bool? ??
          (json['page'] != null && (json['page'] as int) > 1),
    );
  }

  bool get hasNextPage => hasNext || page < totalPages;
  bool get hasPreviousPage => hasPrev || page > 1;

  Map<String, dynamic> toJson() {
    return {
      'page': page,
      'per_page': perPage,
      'total': total,
      'total_pages': totalPages,
      'has_next': hasNext,
      'has_prev': hasPrev,
    };
  }
}
