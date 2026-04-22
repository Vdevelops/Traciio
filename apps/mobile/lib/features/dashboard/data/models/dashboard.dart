class MobileDashboardOverview {
  final TargetSummary target;

  MobileDashboardOverview({
    required this.target,
  });

  factory MobileDashboardOverview.fromJson(Map<String, dynamic> json) {
    return MobileDashboardOverview(
      target: TargetSummary.fromJson(json['target']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'target': target.toJson(),
    };
  }
}

class TargetSummary {
  final int targetAmount;
  final String targetAmountFormatted;
  final int achievedAmount;
  final String achievedAmountFormatted;
  final double progressPercent;
  final String period;
  final String? brickName;

  TargetSummary({
    required this.targetAmount,
    required this.targetAmountFormatted,
    required this.achievedAmount,
    required this.achievedAmountFormatted,
    required this.progressPercent,
    required this.period,
    this.brickName,
  });

  factory TargetSummary.fromJson(Map<String, dynamic> json) {
    return TargetSummary(
      targetAmount: (json['target_amount'] as num?)?.toInt() ?? 0,
      targetAmountFormatted:
          json['target_amount_formatted'] as String? ?? 'RP 0',
      achievedAmount: (json['achieved_amount'] as num?)?.toInt() ?? 0,
      achievedAmountFormatted:
          json['achieved_amount_formatted'] as String? ?? 'Rp 0',
      progressPercent: (json['progress_percent'] as num?)?.toDouble() ?? 0.0,
      period: json['period'] as String? ?? '',
      brickName: json['brick_name'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'target_amount': targetAmount,
      'target_amount_formatted': targetAmountFormatted,
      'achieved_amount': achievedAmount,
      'achieved_amount_formatted': achievedAmountFormatted,
      'progress_percent': progressPercent,
      'period': period,
      'brick_name': brickName,
    };
  }
}


class MobileVisit {
  final String id;
  final String type; // "account", "deal", or "lead"
  final String purpose;
  final String? accountId;
  final String? accountName;
  final String? accountAddress;
  final String? contactId;
  final String? contactName;
  final String? dealId;
  final String? dealTitle;
  final String? leadId;
  final String? leadName;
  final String visitDate;
  final String? visitTime;
  final String status;
  final DateTime? checkInTime;
  final VisitLocation? checkInLocation;
  final DateTime? checkOutTime;
  final VisitLocation? checkOutLocation;
  final DateTime createdAt;
  final DateTime updatedAt;

  MobileVisit({
    required this.id,
    required this.type,
    required this.purpose,
    this.accountId,
    this.accountName,
    this.accountAddress,
    this.contactId,
    this.contactName,
    this.dealId,
    this.dealTitle,
    this.leadId,
    this.leadName,
    required this.visitDate,
    this.visitTime,
    required this.status,
    this.checkInTime,
    this.checkInLocation,
    this.checkOutTime,
    this.checkOutLocation,
    required this.createdAt,
    required this.updatedAt,
  });

  factory MobileVisit.fromJson(Map<String, dynamic> json) {
    return MobileVisit(
      id: json['id'] as String,
      type: json['type'] as String? ?? 'account',
      purpose: json['purpose'] as String? ?? '',
      accountId: json['account_id'] as String?,
      accountName: json['account_name'] as String?,
      accountAddress: json['account_address'] as String?,
      contactId: json['contact_id'] as String?,
      contactName: json['contact_name'] as String?,
      dealId: json['deal_id'] as String?,
      dealTitle: json['deal_title'] as String?,
      leadId: json['lead_id'] as String?,
      leadName: json['lead_name'] as String?,
      visitDate: json['visit_date'] as String,
      visitTime: json['visit_time'] as String?,
      status: json['status'] as String,
      checkInTime: json['check_in_time'] != null
          ? DateTime.parse(json['check_in_time'] as String)
          : null,
      checkInLocation: json['check_in_location'] != null
          ? VisitLocation.fromJson(
              json['check_in_location'] as Map<String, dynamic>,
            )
          : null,
      checkOutTime: json['check_out_time'] != null
          ? DateTime.parse(json['check_out_time'] as String)
          : null,
      checkOutLocation: json['check_out_location'] != null
          ? VisitLocation.fromJson(
              json['check_out_location'] as Map<String, dynamic>,
            )
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'type': type,
      'purpose': purpose,
      'account_id': accountId,
      'account_name': accountName,
      'account_address': accountAddress,
      'contact_id': contactId,
      'contact_name': contactName,
      'deal_id': dealId,
      'deal_title': dealTitle,
      'lead_id': leadId,
      'lead_name': leadName,
      'visit_date': visitDate,
      'visit_time': visitTime,
      'status': status,
      'check_in_time': checkInTime?.toIso8601String(),
      'check_in_location': checkInLocation?.toJson(),
      'check_out_time': checkOutTime?.toIso8601String(),
      'check_out_location': checkOutLocation?.toJson(),
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class VisitLocation {
  final double latitude;
  final double longitude;
  final String? address;

  VisitLocation({
    required this.latitude,
    required this.longitude,
    this.address,
  });

  factory VisitLocation.fromJson(Map<String, dynamic> json) {
    return VisitLocation(
      latitude: (json['latitude'] as num).toDouble(),
      longitude: (json['longitude'] as num).toDouble(),
      address: json['address'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {'latitude': latitude, 'longitude': longitude, 'address': address};
  }
}

class MobileTask {
  final String id;
  final String title;
  final String? description;
  final String? dueDate;
  final String? dueTime;
  final String priority;
  final String status;
  final String type; // general, call, email, meeting, follow_up
  final TaskAssignee? assignedBy;
  final DateTime createdAt;
  final bool isOverdue;

  MobileTask({
    required this.id,
    required this.title,
    this.description,
    this.dueDate,
    this.dueTime,
    required this.priority,
    required this.status,
    this.type = 'general', // Default to general if not provided
    this.assignedBy,
    required this.createdAt,
    required this.isOverdue,
  });

  factory MobileTask.fromJson(Map<String, dynamic> json) {
    return MobileTask(
      id: json['id'] as String,
      title: json['title'] as String,
      description: json['description'] as String?,
      dueDate: json['due_date'] as String?,
      dueTime: json['due_time'] as String?,
      priority: json['priority'] as String,
      status: json['status'] as String,
      type: json['type'] as String? ?? 'general', // Default to general if not provided
      assignedBy: json['assigned_by'] != null
          ? TaskAssignee.fromJson(json['assigned_by'] as Map<String, dynamic>)
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
      isOverdue: json['is_overdue'] as bool? ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'due_date': dueDate,
      'due_time': dueTime,
      'priority': priority,
      'status': status,
      'type': type,
      'assigned_by': assignedBy?.toJson(),
      'created_at': createdAt.toIso8601String(),
      'is_overdue': isOverdue,
    };
  }
}

class TaskAssignee {
  final String id;
  final String name;

  TaskAssignee({required this.id, required this.name});

  factory TaskAssignee.fromJson(Map<String, dynamic> json) {
    return TaskAssignee(id: json['id'] as String, name: json['name'] as String);
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name};
  }
}
