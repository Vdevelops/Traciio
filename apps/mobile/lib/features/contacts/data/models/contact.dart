class Contact {
  final String id;
  final String accountId;
  final AccountContact? account;
  final String name;
  final String? roleId;
  final Role? role;
  final String? phone;
  final String? email;
  final String? position;
  final String? notes;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Contact({
    required this.id,
    required this.accountId,
    this.account,
    required this.name,
    this.roleId,
    this.role,
    this.phone,
    this.email,
    this.position,
    this.notes,
    this.createdAt,
    this.updatedAt,
  });

  factory Contact.fromJson(Map<String, dynamic> json) {
    return Contact(
      id: json['id'] as String,
      accountId: json['account_id'] as String,
      account: json['account'] != null
          ? AccountContact.fromJson(json['account'] as Map<String, dynamic>)
          : null,
      name: json['name'] as String,
      roleId: json['role_id'] as String?,
      role: json['role'] != null
          ? Role.fromJson(json['role'] as Map<String, dynamic>)
          : null,
      phone: json['phone'] as String?,
      email: json['email'] as String?,
      position: json['position'] as String?,
      notes: json['notes'] as String?,
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
      'account_id': accountId,
      'account': account?.toJson(),
      'name': name,
      'role_id': roleId,
      'role': role?.toJson(),
      'phone': phone,
      'email': email,
      'position': position,
      'notes': notes,
      'created_at': createdAt?.toIso8601String(),
      'updated_at': updatedAt?.toIso8601String(),
    };
  }
}

class Role {
  final String id;
  final String name;
  final String? code;

  Role({
    required this.id,
    required this.name,
    this.code,
  });

  factory Role.fromJson(Map<String, dynamic> json) {
    return Role(
      id: json['id'] as String,
      name: json['name'] as String,
      code: json['code'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'code': code,
    };
  }
}

// Simplified Account model for Contact (nested account data)
class AccountContact {
  final String id;
  final String name;
  final String? city;

  AccountContact({
    required this.id,
    required this.name,
    this.city,
  });

  factory AccountContact.fromJson(Map<String, dynamic> json) {
    return AccountContact(
      id: json['id'] as String,
      name: json['name'] as String,
      city: json['city'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'city': city,
    };
  }
}

class ContactListResponse {
  final List<Contact> items;
  final Pagination pagination;

  ContactListResponse({
    required this.items,
    required this.pagination,
  });

  factory ContactListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'];
    
    // Handle different response formats
    List<Contact> items;
    Pagination pagination;

    if (data is List) {
      // Format: { success: true, data: [...] }
      items = data
          .map((item) => Contact.fromJson(item as Map<String, dynamic>))
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
              ?.map((item) => Contact.fromJson(item as Map<String, dynamic>))
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

    return ContactListResponse(
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

