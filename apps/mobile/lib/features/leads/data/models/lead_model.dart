class Lead {
  final String id;
  final String firstName;
  final String? lastName;
  final String? companyName;
  final String? title;
  final String? email;
  final String? phone;
  final String? address;
  final String? notes;
  final String? leadSourceId;
  final String? leadStatusId;
  final String? industry;
  final String? province;
  final String? leadSource;
  final LeadStatus? leadStatus;
  final DateTime createdAt;
  final DateTime updatedAt;

  Lead({
    required this.id,
    required this.firstName,
    this.lastName,
    this.companyName,
    this.title,
    this.email,
    this.phone,
    this.address,
    this.notes,
    this.leadSourceId,
    this.leadStatusId,
    this.industry,
    this.province,
    this.leadSource,
    this.leadStatus,
    required this.createdAt,
    required this.updatedAt,
  });

  String get name => '$firstName ${lastName ?? ""}'.trim();
  String get company => companyName ?? '';
  String get status => leadStatus?.code ?? leadStatusId ?? 'new';
  String get source => leadSource ?? leadSourceId ?? '';

  factory Lead.fromJson(Map<String, dynamic> json) {
    return Lead(
      id: json['id'] as String,
      firstName: json['first_name'] as String,
      lastName: json['last_name'] as String?,
      companyName: json['company_name'] as String?,
      title: json['job_title'] as String?,
      email: json['email'] as String?,
      phone: json['phone'] as String?,
      address: json['address'] as String?,
      notes: json['notes'] as String?,
      leadSourceId: json['lead_source_id'] as String?,
      leadStatusId:
          json['lead_status_id'] as String? ?? json['lead_status'] as String?,
      industry: json['industry'] as String?,
      province: json['province'] as String?,
      leadSource: json['lead_source'] as String?,
      leadStatus: json['lead_status_ref'] != null
          ? LeadStatus.fromJson(json['lead_status_ref'] as Map<String, dynamic>)
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'first_name': firstName,
      'last_name': lastName,
      'company_name': companyName,
      'job_title': title,
      'email': email,
      'phone': phone,
      'address': address,
      'notes': notes,
      'lead_source_id': leadSourceId,
      'lead_status_id': leadStatusId,
      'industry': industry,
      'province': province,
      'lead_source': leadSource,
      'lead_status_ref': leadStatus?.toJson(), // ✅ Tambahkan ini
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class LeadStatus {
  final String id;
  final String name;
  final String code;
  final String color;

  LeadStatus({
    required this.id,
    required this.name,
    required this.code,
    required this.color,
  });

  factory LeadStatus.fromJson(Map<String, dynamic> json) {
    return LeadStatus(
      id: json['id'] as String,
      name: json['name'] as String,
      code: json['code'] as String,
      color: json['color'] as String? ?? '#808080',
    );
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name, 'code': code, 'color': color};
  }
}

class LeadSource {
  final String id;
  final String name;
  final String code;

  LeadSource({required this.id, required this.name, required this.code});

  factory LeadSource.fromJson(Map<String, dynamic> json) {
    return LeadSource(
      id: json['id'] as String,
      name: json['name'] as String,
      code: json['code'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name, 'code': code};
  }
}

class LeadListResponse {
  final List<Lead> items;
  final Pagination pagination;

  LeadListResponse({required this.items, required this.pagination});

  factory LeadListResponse.fromJson(Map<String, dynamic> json) {
    final rawData = json['data'];
    final metaData = json['meta'];

    List<Lead> items = [];
    Pagination pagination;

    if (rawData is List) {
      items = rawData
          .map((item) => Lead.fromJson(item as Map<String, dynamic>))
          .toList();

      if (metaData != null && metaData is Map<String, dynamic>) {
        final paginationData = metaData['pagination'] as Map<String, dynamic>?;
        if (paginationData != null) {
          pagination = Pagination.fromJson(paginationData);
        } else {
          pagination = Pagination(
            page: 1,
            perPage: items.length,
            total: items.length,
            totalPages: 1,
          );
        }
      } else {
        pagination = Pagination(
          page: 1,
          perPage: items.length,
          total: items.length,
          totalPages: 1,
        );
      }
    } else {
      pagination = Pagination(page: 1, perPage: 0, total: 0, totalPages: 1);
    }

    return LeadListResponse(items: items, pagination: pagination);
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
      perPage: json['per_page'] as int? ?? 10,
      total: json['total'] as int? ?? 0,
      totalPages: json['total_pages'] as int? ?? 1,
    );
  }
}

class LeadFormData {
  final List<FormOption> leadSources;
  final List<FormOption> leadStatuses;
  final List<String> industries;
  final List<String> provinces;

  LeadFormData({
    required this.leadSources,
    required this.leadStatuses,
    required this.industries,
    required this.provinces,
  });

  factory LeadFormData.fromJson(Map<String, dynamic> json) {
    return LeadFormData(
      leadSources: (json['lead_sources'] as List? ?? [])
          .map((e) => FormOption.fromJson(e as Map<String, dynamic>))
          .toList(),
      leadStatuses: (json['lead_statuses'] as List? ?? [])
          .map((e) => FormOption.fromJson(e as Map<String, dynamic>))
          .toList(),
      industries: List<String>.from(json['industries'] ?? []),
      provinces: List<String>.from(json['provinces'] ?? []),
    );
  }
}

class FormOption {
  final String? id;
  final String value;
  final String label;

  FormOption({this.id, required this.value, required this.label});

  factory FormOption.fromJson(Map<String, dynamic> json) {
    return FormOption(
      id: json['id'] as String?,
      value:
          json['code'] as String? ??
          json['value'] as String? ??
          json['id'] as String? ??
          '',
      label: json['label'] as String? ?? json['name'] as String? ?? '',
    );
  }
}
