class Deal {
  final String id;
  final String title;
  final String stageId;
  final PipelineStage? stage;
  final String? accountId;
  final AccountInfo? account;
  final String? contactId;
  final ContactInfo? contact;
  final int value;
  final int probability;
  final String status;
  final String? assignedTo;
  final DateTime? expectedCloseDate;
  final List<DealProductItem>? productItems;
  final String? notes;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Deal({
    required this.id,
    required this.title,
    required this.stageId,
    this.stage,
    this.accountId,
    this.account,
    this.contactId,
    this.contact,
    required this.value,
    required this.probability,
    required this.status,
    this.assignedTo,
    this.expectedCloseDate,
    this.productItems,
    this.notes,
    this.createdAt,
    this.updatedAt,
  });

  factory Deal.fromJson(Map<String, dynamic> json) {
    return Deal(
      id: json['id'] as String,
      title: json['title'] as String? ?? json['name'] as String? ?? '',
      stageId: json['stage_id'] as String,
      stage: json['stage'] != null
          ? PipelineStage.fromJson(json['stage'] as Map<String, dynamic>)
          : null,
      accountId: json['account_id'] as String?,
      account: json['account'] != null
          ? AccountInfo.fromJson(json['account'] as Map<String, dynamic>)
          : null,
      contactId: json['contact_id'] as String?,
      contact: json['contact'] != null
          ? ContactInfo.fromJson(json['contact'] as Map<String, dynamic>)
          : null,
      value: json['value'] as int? ?? json['amount'] as int? ?? 0,
      probability: json['probability'] as int? ?? 0,
      status: json['status'] as String? ?? 'open',
      assignedTo: json['assigned_to'] as String?,
      expectedCloseDate: json['expected_close_date'] != null
          ? DateTime.parse(json['expected_close_date'] as String)
          : null,
      productItems: json['product_items'] != null
          ? (json['product_items'] as List)
                .map((e) => DealProductItem.fromJson(e as Map<String, dynamic>))
                .toList()
          : null,
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
      'title': title,
      'stage_id': stageId,
      'stage': stage?.toJson(),
      'account_id': accountId,
      'account': account?.toJson(),
      'contact_id': contactId,
      'contact': contact?.toJson(),
      'value': value,
      'probability': probability,
      'status': status,
      'assigned_to': assignedTo,
      'expected_close_date': expectedCloseDate?.toIso8601String(),
      'product_items': productItems?.map((e) => e.toJson()).toList(),
      'notes': notes,
      'created_at': createdAt?.toIso8601String(),
      'updated_at': updatedAt?.toIso8601String(),
    };
  }
}

class PipelineStage {
  final String id;
  final String name;
  final int order;
  final int probability;
  final String? color;

  PipelineStage({
    required this.id,
    required this.name,
    required this.order,
    required this.probability,
    this.color,
  });

  factory PipelineStage.fromJson(Map<String, dynamic> json) {
    return PipelineStage(
      id: json['id'] as String,
      name: json['name'] as String,
      order: json['order'] as int? ?? 0,
      probability: json['probability'] as int? ?? 0,
      color: json['color'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'order': order,
      'probability': probability,
      'color': color,
    };
  }
}

class AccountInfo {
  final String id;
  final String name;

  AccountInfo({required this.id, required this.name});

  factory AccountInfo.fromJson(Map<String, dynamic> json) {
    return AccountInfo(id: json['id'] as String, name: json['name'] as String);
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name};
  }
}

class ContactInfo {
  final String id;
  final String name;

  ContactInfo({required this.id, required this.name});

  factory ContactInfo.fromJson(Map<String, dynamic> json) {
    return ContactInfo(id: json['id'] as String, name: json['name'] as String);
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name};
  }
}

class Product {
  final String id;
  final String name;
  final String sku;
  final int price;

  Product({
    required this.id,
    required this.name,
    required this.sku,
    required this.price,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['id'] as String,
      name: json['name'] as String,
      sku: json['sku'] as String? ?? '',
      price: json['price'] as int? ?? json['unit_price'] as int? ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name, 'sku': sku, 'price': price};
  }
}

class DealProductItem {
  final String id;
  final String productId;
  final String productName;
  final String productSku;
  final int unitPrice;
  final int quantity;
  final int discountAmount;
  final int subtotal;
  final String? notes;

  DealProductItem({
    required this.id,
    required this.productId,
    required this.productName,
    required this.productSku,
    required this.unitPrice,
    required this.quantity,
    required this.discountAmount,
    required this.subtotal,
    this.notes,
  });

  factory DealProductItem.fromJson(Map<String, dynamic> json) {
    return DealProductItem(
      id: json['id'] as String,
      productId: json['product_id'] as String,
      productName: json['product_name'] as String,
      productSku: json['product_sku'] as String,
      unitPrice: json['unit_price'] as int? ?? 0,
      quantity: json['quantity'] as int? ?? 0,
      discountAmount: json['discount_amount'] as int? ?? 0,
      subtotal: json['subtotal'] as int? ?? 0,
      notes: json['notes'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'product_id': productId,
      'product_name': productName,
      'product_sku': productSku,
      'unit_price': unitPrice,
      'quantity': quantity,
      'discount_amount': discountAmount,
      'subtotal': subtotal,
      'notes': notes,
    };
  }
}

class DealListResponse {
  final List<Deal> items;
  final Pagination pagination;

  DealListResponse({required this.items, required this.pagination});

  factory DealListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'];
    List<Deal> items = [];
    Pagination pagination;

    if (data is List) {
      items = data
          .map((e) => Deal.fromJson(e as Map<String, dynamic>))
          .toList();
      pagination = Pagination(
        page: 1,
        perPage: items.length,
        total: items.length,
        totalPages: 1,
      );
    } else if (data is Map<String, dynamic>) {
      items =
          (data['items'] as List<dynamic>?)
              ?.map((e) => Deal.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [];
      pagination = data['pagination'] != null
          ? Pagination.fromJson(data['pagination'] as Map<String, dynamic>)
          : Pagination(
              page: 1,
              perPage: 20,
              total: items.length,
              totalPages: 1,
            );
    } else {
      pagination = Pagination(page: 1, perPage: 20, total: 0, totalPages: 0);
    }

    return DealListResponse(items: items, pagination: pagination);
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
      totalPages: json['total_pages'] as int? ?? 1,
    );
  }

  bool get hasNextPage => page < totalPages;
  bool get hasPreviousPage => page > 1;
}
