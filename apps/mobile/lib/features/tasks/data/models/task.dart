import '../../../../core/network/pagination.dart';

class Task {
  final String id;
  final String title;
  final String? description;
  final String type;
  final String priority;
  final String status;
  final DateTime? dueDate;
  final DateTime? completedAt;
  final String? assignedTo;
  final String? assignedFrom;
  final String? accountId;
  final String? contactId;
  final String? dealId;
  final AccountInfo? account;
  final ContactInfo? contact;
  final DealInfo? deal;
  final AssignedUser? assignedUser;
  final AssignedUser? assignedFromUser;
  final DateTime createdAt;
  final DateTime updatedAt;
  final List<Reminder> reminders;

  Task({
    required this.id,
    required this.title,
    this.description,
    required this.type,
    required this.priority,
    required this.status,
    this.dueDate,
    this.completedAt,
    this.assignedTo,
    this.assignedFrom,
    this.accountId,
    this.contactId,
    this.dealId,
    this.account,
    this.contact,
    this.deal,
    this.assignedUser,
    this.assignedFromUser,
    required this.createdAt,
    required this.updatedAt,
    this.reminders = const [],
  });

  factory Task.fromJson(Map<String, dynamic> json) {
    return Task(
      id: (json['id'] ?? '').toString(),
      title: (json['title'] ?? '').toString(),
      description: json['description'] as String?,
      type: (json['type'] ?? 'general').toString(),
      priority: (json['priority'] ?? 'medium').toString(),
      status: (json['status'] ?? 'pending').toString(),
      dueDate: json['due_date'] != null
          ? DateTime.tryParse(json['due_date'].toString())
          : null,
      completedAt: json['completed_at'] != null
          ? DateTime.tryParse(json['completed_at'].toString())
          : null,
      assignedTo: json['assigned_to']?.toString(),
      assignedFrom: json['assigned_from']?.toString(),
      accountId: json['account_id']?.toString(),
      contactId: json['contact_id']?.toString(),
      dealId: json['deal_id']?.toString(),
      account: json['account'] != null
          ? AccountInfo.fromJson(json['account'] as Map<String, dynamic>)
          : null,
      contact: json['contact'] != null
          ? ContactInfo.fromJson(json['contact'] as Map<String, dynamic>)
          : null,
      deal: json['deal'] != null
          ? DealInfo.fromJson(json['deal'] as Map<String, dynamic>)
          : null,
      assignedUser: json['assigned_user'] != null
          ? AssignedUser.fromJson(json['assigned_user'] as Map<String, dynamic>)
          : null,
      assignedFromUser: json['assigned_from_user'] != null
          ? AssignedUser.fromJson(
              json['assigned_from_user'] as Map<String, dynamic>,
            )
          : null,
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at'].toString()) ?? DateTime.now()
          : DateTime.now(),
      updatedAt: json['updated_at'] != null
          ? DateTime.tryParse(json['updated_at'].toString()) ?? DateTime.now()
          : DateTime.now(),
      reminders: json['reminders'] != null && json['reminders'] is List
          ? (json['reminders'] as List<dynamic>)
                .map((e) => Reminder.fromJson(e as Map<String, dynamic>))
                .toList()
          : [],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'type': type,
      'priority': priority,
      'status': status,
      'due_date': dueDate?.toIso8601String(),
      'completed_at': completedAt?.toIso8601String(),
      'assigned_to': assignedTo,
      'assigned_from': assignedFrom,
      'account_id': accountId,
      'contact_id': contactId,
      'deal_id': dealId,
      'account': account?.toJson(),
      'contact': contact?.toJson(),
      'deal': deal?.toJson(),
      'assigned_user': assignedUser?.toJson(),
      'assigned_from_user': assignedFromUser?.toJson(),
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
      'reminders': reminders.map((r) => r.toJson()).toList(),
    };
  }

  bool get isOverdue {
    if (dueDate == null || status == 'completed' || status == 'cancelled') {
      return false;
    }
    return dueDate!.isBefore(DateTime.now());
  }

  bool get isDueToday {
    if (dueDate == null) return false;
    final now = DateTime.now();
    return dueDate!.year == now.year &&
        dueDate!.month == now.month &&
        dueDate!.day == now.day;
  }
}

class AccountInfo {
  final String id;
  final String name;

  AccountInfo({required this.id, required this.name});

  factory AccountInfo.fromJson(Map<String, dynamic> json) {
    return AccountInfo(
      id: (json['id'] ?? '').toString(),
      name: (json['name'] ?? '').toString(),
    );
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
    return ContactInfo(
      id: (json['id'] ?? '').toString(),
      name: (json['name'] ?? '').toString(),
    );
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name};
  }
}

class DealInfo {
  final String id;
  final String title;

  DealInfo({required this.id, required this.title});

  factory DealInfo.fromJson(Map<String, dynamic> json) {
    return DealInfo(
      id: (json['id'] ?? '').toString(),
      title: (json['title'] ?? '').toString(),
    );
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'title': title};
  }
}

class AssignedUser {
  final String id;
  final String name;
  final String? email;

  AssignedUser({required this.id, required this.name, this.email});

  factory AssignedUser.fromJson(Map<String, dynamic> json) {
    return AssignedUser(
      id: (json['id'] ?? '').toString(),
      name: (json['name'] ?? '').toString(),
      email: json['email']?.toString(),
    );
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name, 'email': email};
  }
}

class Reminder {
  final String id;
  final String taskId;
  final DateTime remindAt;
  final String reminderType;
  final String? message;
  final bool isSent;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Reminder({
    required this.id,
    required this.taskId,
    required this.remindAt,
    required this.reminderType,
    this.message,
    this.isSent = false,
    this.createdAt,
    this.updatedAt,
  });

  factory Reminder.fromJson(Map<String, dynamic> json) {
    return Reminder(
      id: (json['id'] ?? '').toString(),
      taskId: (json['task_id'] ?? '').toString(),
      remindAt: json['remind_at'] != null
          ? DateTime.tryParse(json['remind_at'].toString()) ?? DateTime.now()
          : DateTime.now(),
      reminderType: (json['reminder_type'] ?? 'in_app').toString(),
      message: json['message']?.toString(),
      isSent: json['is_sent'] as bool? ?? false,
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at'].toString())
          : null,
      updatedAt: json['updated_at'] != null
          ? DateTime.tryParse(json['updated_at'].toString())
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'task_id': taskId,
      'remind_at': remindAt.toIso8601String(),
      'reminder_type': reminderType,
      'message': message,
    };
  }
}

class TaskListResponse {
  final List<Task> items;
  final Pagination pagination;

  TaskListResponse({required this.items, required this.pagination});

  factory TaskListResponse.fromJson(Map<String, dynamic> json) {
    final dynamic rawData = json['data'];
    final dynamic metaData = json['meta'];

    List<Task> items;
    Pagination pagination;

    // Mobile API returns data as array directly with meta.pagination
    if (rawData is List) {
      items = rawData
          .map((item) => Task.fromJson(item as Map<String, dynamic>))
          .toList();

      // Extract pagination from meta
      if (metaData != null && metaData is Map<String, dynamic>) {
        final paginationData = metaData['pagination'] as Map<String, dynamic>?;
        if (paginationData != null) {
          pagination = Pagination.fromJson(paginationData);
        } else {
          // Fallback if pagination not in meta
          pagination = Pagination(
            page: 1,
            perPage: items.length,
            total: items.length,
            totalPages: 1,
            hasNext: false,
            hasPrev: false,
          );
        }
      } else {
        // Fallback if no meta
        pagination = Pagination(
          page: 1,
          perPage: items.length,
          total: items.length,
          totalPages: 1,
          hasNext: false,
          hasPrev: false,
        );
      }
    } else if (rawData is Map<String, dynamic>) {
      // Web API format (backward compatibility)
      items =
          (rawData['items'] as List<dynamic>?)
              ?.map((item) => Task.fromJson(item as Map<String, dynamic>))
              .toList() ??
          [];
      pagination = Pagination.fromJson(
        rawData['pagination'] as Map<String, dynamic>,
      );
    } else {
      throw Exception('Invalid data format for TaskListResponse');
    }

    return TaskListResponse(items: items, pagination: pagination);
  }
}
