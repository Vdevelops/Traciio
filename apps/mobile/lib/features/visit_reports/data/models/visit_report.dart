class VisitReport {
  final String id;
  final String? accountId; // Made optional to support lead-only visit reports
  final AccountInfo? account;
  final String? contactId;
  final ContactInfo? contact;
  final String? dealId; // Added for deal support
  final DealInfo? deal; // Added for deal info
  final String? leadId; // Added for lead support
  final LeadInfo? lead; // Added for lead info
  final String type; // "account", "deal", or "lead"
  final String? salesRepId;
  final String visitDate;
  final String? purpose;
  final String? notes;
  final String status;
  final Location? checkInLocation;
  final DateTime? checkInTime;
  final Location? checkOutLocation;
  final DateTime? checkOutTime;
  final List<String>? photoUrls;
  final String? outcome; // Added for submit (positive, neutral, negative, very_positive)
  final String? nextSteps; // Added for submit
  final String? rejectionReason; // Added for rejection
  final String? approvedBy; // Added for approval
  final DateTime? approvedAt; // Added for approval
  final DateTime? createdAt;
  final DateTime? updatedAt;

  VisitReport({
    required this.id,
    this.accountId,
    this.account,
    this.contactId,
    this.contact,
    this.dealId,
    this.deal,
    this.leadId,
    this.lead,
    required this.type,
    this.salesRepId,
    required this.visitDate,
    this.purpose,
    this.notes,
    required this.status,
    this.checkInLocation,
    this.checkInTime,
    this.checkOutLocation,
    this.checkOutTime,
    this.photoUrls,
    this.outcome,
    this.nextSteps,
    this.rejectionReason,
    this.approvedBy,
    this.approvedAt,
    this.createdAt,
    this.updatedAt,
  });

  factory VisitReport.fromJson(Map<String, dynamic> json) {
    return VisitReport(
      id: json['id'] as String,
      accountId: json['account_id'] as String?,
      account: json['account'] != null
          ? AccountInfo.fromJson(json['account'] as Map<String, dynamic>)
          : null,
      contactId: json['contact_id'] as String?,
      contact: json['contact'] != null
          ? ContactInfo.fromJson(json['contact'] as Map<String, dynamic>)
          : null,
      dealId: json['deal_id'] as String?,
      deal: json['deal'] != null
          ? DealInfo.fromJson(json['deal'] as Map<String, dynamic>)
          : null,
      leadId: json['lead_id'] as String?,
      lead: json['lead'] != null
          ? LeadInfo.fromJson(json['lead'] as Map<String, dynamic>)
          : null,
      type: json['type'] as String? ?? 'account', // Default to 'account' if not provided
      salesRepId: json['sales_rep_id'] as String?,
      visitDate: json['visit_date'] as String,
      purpose: json['purpose'] as String?,
      notes: json['notes'] as String?,
      status: json['status'] as String? ?? 'draft',
      checkInLocation: json['check_in_location'] != null
          ? Location.fromJson(
              json['check_in_location'] as Map<String, dynamic>,
            )
          : null,
      checkInTime: json['check_in_time'] != null
          ? DateTime.parse(json['check_in_time'] as String)
          : null,
      checkOutLocation: json['check_out_location'] != null
          ? Location.fromJson(
              json['check_out_location'] as Map<String, dynamic>,
            )
          : null,
      checkOutTime: json['check_out_time'] != null
          ? DateTime.parse(json['check_out_time'] as String)
          : null,
      photoUrls: json['photos'] != null
          ? (json['photos'] as List<dynamic>)
              .map((e) => e as String)
              .toList()
          : null,
      outcome: json['outcome'] as String?,
      nextSteps: json['next_steps'] as String?,
      rejectionReason: json['rejection_reason'] as String?,
      approvedBy: json['approved_by'] as String?,
      approvedAt: json['approved_at'] != null
          ? DateTime.parse(json['approved_at'] as String)
          : null,
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
      'contact_id': contactId,
      'contact': contact?.toJson(),
      'deal_id': dealId,
      'deal': deal?.toJson(),
      'lead_id': leadId,
      'lead': lead?.toJson(),
      'type': type,
      'sales_rep_id': salesRepId,
      'visit_date': visitDate,
      'purpose': purpose,
      'notes': notes,
      'status': status,
      'check_in_location': checkInLocation?.toJson(),
      'check_in_time': checkInTime?.toIso8601String(),
      'check_out_location': checkOutLocation?.toJson(),
      'check_out_time': checkOutTime?.toIso8601String(),
      'photos': photoUrls,
      'outcome': outcome,
      'next_steps': nextSteps,
      'rejection_reason': rejectionReason,
      'approved_by': approvedBy,
      'approved_at': approvedAt?.toIso8601String(),
      'created_at': createdAt?.toIso8601String(),
      'updated_at': updatedAt?.toIso8601String(),
    };
  }
}

class DealInfo {
  final String id;
  final String title;
  final String? accountId;
  final AccountInfo? account;

  DealInfo({
    required this.id,
    required this.title,
    this.accountId,
    this.account,
  });

  factory DealInfo.fromJson(Map<String, dynamic> json) {
    return DealInfo(
      id: json['id'] as String,
      title: json['title'] as String,
      accountId: json['account_id'] as String?,
      account: json['account'] != null
          ? AccountInfo.fromJson(json['account'] as Map<String, dynamic>)
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'account_id': accountId,
      'account': account?.toJson(),
    };
  }
}

class LeadInfo {
  final String id;
  final String firstName;
  final String? lastName;
  final String? companyName;

  LeadInfo({
    required this.id,
    required this.firstName,
    this.lastName,
    this.companyName,
  });

  factory LeadInfo.fromJson(Map<String, dynamic> json) {
    return LeadInfo(
      id: json['id'] as String,
      firstName: json['first_name'] as String,
      lastName: json['last_name'] as String?,
      companyName: json['company_name'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'first_name': firstName,
      'last_name': lastName,
      'company_name': companyName,
    };
  }
}

class AccountInfo {
  final String id;
  final String name;
  final String? city;

  AccountInfo({
    required this.id,
    required this.name,
    this.city,
  });

  factory AccountInfo.fromJson(Map<String, dynamic> json) {
    return AccountInfo(
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

class ContactInfo {
  final String id;
  final String name;
  final String? position;

  ContactInfo({
    required this.id,
    required this.name,
    this.position,
  });

  factory ContactInfo.fromJson(Map<String, dynamic> json) {
    return ContactInfo(
      id: json['id'] as String,
      name: json['name'] as String,
      position: json['position'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'position': position,
    };
  }
}

class Location {
  final double latitude;
  final double longitude;
  final String? address;

  Location({
    required this.latitude,
    required this.longitude,
    this.address,
  });

  factory Location.fromJson(Map<String, dynamic> json) {
    return Location(
      latitude: (json['latitude'] as num).toDouble(),
      longitude: (json['longitude'] as num).toDouble(),
      address: json['address'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'latitude': latitude,
      'longitude': longitude,
      'address': address,
    };
  }
}

class VisitReportListResponse {
  final List<VisitReport> items;
  final Pagination pagination;

  VisitReportListResponse({
    required this.items,
    required this.pagination,
  });

  factory VisitReportListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'];
    
    // Handle different response formats
    List<VisitReport> items;
    Pagination pagination;

    if (data is List) {
      // Format: { success: true, data: [...] }
      items = data
          .map((item) => VisitReport.fromJson(item as Map<String, dynamic>))
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
              ?.map((item) => VisitReport.fromJson(item as Map<String, dynamic>))
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

    return VisitReportListResponse(
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

