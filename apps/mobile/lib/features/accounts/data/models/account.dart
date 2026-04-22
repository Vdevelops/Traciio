class Account {
  final String id;
  final String name;
  final String? categoryId;
  final Category? category;
  final String? address;
  final String? city;
  final String? province;
  final String? phone;
  final String? email;
  final String status;
  final String? assignedTo;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Account({
    required this.id,
    required this.name,
    this.categoryId,
    this.category,
    this.address,
    this.city,
    this.province,
    this.phone,
    this.email,
    required this.status,
    this.assignedTo,
    this.createdAt,
    this.updatedAt,
  });

  factory Account.fromJson(Map<String, dynamic> json) {
    return Account(
      id: json['id'] as String,
      name: json['name'] as String,
      categoryId: json['category_id'] as String?,
      category: json['category'] != null
          ? Category.fromJson(json['category'] as Map<String, dynamic>)
          : null,
      address: json['address'] as String?,
      city: json['city'] as String?,
      province: json['province'] as String?,
      phone: json['phone'] as String?,
      email: json['email'] as String?,
      status: json['status'] as String? ?? 'active',
      assignedTo: json['assigned_to'] as String?,
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'] as String)
          : null,
      updatedAt: json['updated_at'] != null
          ? DateTime.parse(json['updated_at'] as String)
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'category_id': categoryId,
      'category': category?.toJson(),
      'address': address,
      'city': city,
      'province': province,
      'phone': phone,
      'email': email,
      'status': status,
      'assigned_to': assignedTo,
      'created_at': createdAt?.toIso8601String(),
      'updated_at': updatedAt?.toIso8601String(),
    };
  }
}

class Category {
  final String id;
  final String name;
  final String? code;
  final String? badgeColor;

  Category({
    required this.id,
    required this.name,
    this.code,
    this.badgeColor,
  });

  factory Category.fromJson(Map<String, dynamic> json) {
    return Category(
      id: json['id'] as String,
      name: json['name'] as String,
      code: json['code'] as String?,
      badgeColor: json['badge_color'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'code': code,
      'badge_color': badgeColor,
    };
  }
}

class AccountListResponse {
  final List<Account> items;
  final Pagination pagination;

  AccountListResponse({
    required this.items,
    required this.pagination,
  });

  factory AccountListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'];
    
    // Handle different response formats
    List<Account> items;
    Pagination pagination;

    if (data is List) {
      // Format: { success: true, data: [...] }
      items = data
          .map((item) => Account.fromJson(item as Map<String, dynamic>))
          .toList();
      // Create default pagination if data is a list
      pagination = Pagination(
        page: json['page'] as int? ?? 1,
        perPage: json['per_page'] as int? ?? data.length,
        total: json['total'] as int? ?? data.length,
        totalPages: json['total_pages'] as int? ?? 1,
      );
    } else if (data is Map<String, dynamic>) {
      // Format: { success: true, data: { items: [...], pagination: {...} } }
      items = (data['items'] as List<dynamic>?)
              ?.map((item) => Account.fromJson(item as Map<String, dynamic>))
              .toList() ??
          [];
      pagination = data['pagination'] != null
          ? Pagination.fromJson(data['pagination'] as Map<String, dynamic>)
          : Pagination(
              page: json['page'] as int? ?? 1,
              perPage: json['per_page'] as int? ?? 20,
              total: json['total'] as int? ?? items.length,
              totalPages: json['total_pages'] as int? ?? 1,
            );
    } else {
      items = [];
      pagination = Pagination(
        page: 1,
        perPage: 20,
        total: 0,
        totalPages: 0,
      );
    }

    return AccountListResponse(
      items: items,
      pagination: pagination,
    );
  }
}

class Pagination {
  final int page;
  final int perPage;
  final int total;
  final int totalPages;

  Pagination({
    required this.page,
    required this.perPage,
    required this.total,
    required this.totalPages,
  });

  factory Pagination.fromJson(Map<String, dynamic> json) {
    return Pagination(
      page: json['page'] as int? ?? 1,
      perPage: json['per_page'] as int? ?? 20,
      total: json['total'] as int? ?? 0,
      totalPages: json['total_pages'] as int? ?? 0,
    );
  }

  bool get hasNextPage => page < totalPages;
  bool get hasPreviousPage => page > 1;
}

