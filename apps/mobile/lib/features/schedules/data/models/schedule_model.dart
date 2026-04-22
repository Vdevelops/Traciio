import '../../../../core/network/pagination.dart';
import '../../../tasks/data/models/task.dart';

class ScheduleModel {
  final String id;
  final String title;
  final String? description;
  final DateTime scheduledAt;
  final int? reminderMinutesBefore;
  final String status;
  final String? location;
  final String? activityType;
  final String? visitReportId;
  final String? taskId;
  final Task? task;
  final String createdBy;
  final DateTime createdAt;
  final DateTime updatedAt;
  final bool syncToCalendar;

  ScheduleModel({
    required this.id,
    required this.title,
    this.description,
    required this.scheduledAt,
    this.reminderMinutesBefore,
    required this.status,
    this.location,
    this.activityType,
    this.visitReportId,
    this.taskId,
    this.task,
    required this.createdBy,
    required this.createdAt,
    required this.updatedAt,
    this.syncToCalendar = false,
  });

  factory ScheduleModel.fromJson(Map<String, dynamic> json) {
    return ScheduleModel(
      id: (json['id'] ?? '').toString(),
      title: (json['title'] ?? '').toString(),
      description: json['description']?.toString(),
      scheduledAt: json['scheduled_at'] != null
          ? DateTime.tryParse(json['scheduled_at'].toString()) ?? DateTime.now()
          : DateTime.now(),
      reminderMinutesBefore: json['reminder_minutes_before'] != null
          ? int.tryParse(json['reminder_minutes_before'].toString())
          : null,
      status: (json['status'] ?? 'pending').toString(),
      location: json['location']?.toString(),
      activityType: json['activity_type']?.toString(),
      visitReportId: json['visit_report_id']?.toString(),
      taskId: json['task_id']?.toString(),
      task: json['task'] != null && json['task'] is Map<String, dynamic>
          ? Task.fromJson(json['task'] as Map<String, dynamic>)
          : null,
      createdBy: (json['created_by'] ?? json['user_id'] ?? '').toString(),
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at'].toString()) ?? DateTime.now()
          : DateTime.now(),
      updatedAt: json['updated_at'] != null
          ? DateTime.tryParse(json['updated_at'].toString()) ?? DateTime.now()
          : DateTime.now(),
      syncToCalendar:
          json['sync_to_calendar'] == true || json['sync_to_calendar'] == 1,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'scheduled_at': scheduledAt.toIso8601String(),
      'reminder_minutes_before': reminderMinutesBefore,
      'status': status,
      'task_id': taskId,
      'created_by': createdBy,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  ScheduleModel copyWith({
    String? id,
    String? title,
    String? description,
    DateTime? scheduledAt,
    int? reminderMinutesBefore,
    String? status,
    String? location,
    String? activityType,
    String? visitReportId,
    String? taskId,
    Task? task,
    String? createdBy,
    DateTime? createdAt,
    DateTime? updatedAt,
    bool? syncToCalendar,
  }) {
    return ScheduleModel(
      id: id ?? this.id,
      title: title ?? this.title,
      description: description ?? this.description,
      scheduledAt: scheduledAt ?? this.scheduledAt,
      reminderMinutesBefore:
          reminderMinutesBefore ?? this.reminderMinutesBefore,
      status: status ?? this.status,
      location: location ?? this.location,
      activityType: activityType ?? this.activityType,
      visitReportId: visitReportId ?? this.visitReportId,
      taskId: taskId ?? this.taskId,
      task: task ?? this.task,
      createdBy: createdBy ?? this.createdBy,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      syncToCalendar: syncToCalendar ?? this.syncToCalendar,
    );
  }
}

class ScheduleRequest {
  final String title;
  final String? description;
  final DateTime scheduledAt;
  final int? reminderMinutesBefore;
  final String status;
  final String? taskId;
  final String? location;
  final String? activityType;

  ScheduleRequest({
    required this.title,
    this.description,
    required this.scheduledAt,
    this.reminderMinutesBefore,
    this.status = 'pending',
    this.taskId,
    this.location,
    this.activityType,
  });

  Map<String, dynamic> toJson() {
    return {
      'title': title,
      'description': description,
      'scheduled_at': scheduledAt.toIso8601String(),
      'reminder_minutes_before': reminderMinutesBefore,
      'status': status,
      'location': location,
      'activity_type': activityType,
      'task_id': taskId,
    };
  }
}

class ScheduleListResponse {
  final List<ScheduleModel> items;
  final Pagination pagination;

  ScheduleListResponse({required this.items, required this.pagination});

  factory ScheduleListResponse.fromJson(Map<String, dynamic> json) {
    final dynamic rawData = json['data'];
    final dynamic metaData = json['meta'];

    List<ScheduleModel> items = [];
    Pagination pagination;

    if (rawData is List) {
      items = rawData
          .where((item) => item != null && item is Map<String, dynamic>)
          .map((item) => ScheduleModel.fromJson(item as Map<String, dynamic>))
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
      // If data is null or not a List, return empty list instead of throwing
      pagination = Pagination(page: 1, perPage: 0, total: 0, totalPages: 1);
    }

    return ScheduleListResponse(items: items, pagination: pagination);
  }
}
