class Notification {
  final String id;
  final String userId;
  final String title;
  final String message;
  final NotificationType type;
  final bool isRead;
  final DateTime? readAt;
  final String? data; // JSON string
  final DateTime createdAt;
  final DateTime updatedAt;

  Notification({
    required this.id,
    required this.userId,
    required this.title,
    required this.message,
    required this.type,
    required this.isRead,
    this.readAt,
    this.data,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Notification.fromJson(Map<String, dynamic> json) {
    // Safely parse DateTime, handling null and invalid formats
    DateTime? parseDateTime(dynamic value) {
      if (value == null) return null;
      if (value is String) {
        try {
          return DateTime.parse(value);
        } catch (e) {
          return null;
        }
      }
      return null;
    }

    return Notification(
      id: json['id']?.toString() ?? '',
      userId: json['user_id']?.toString() ?? '',
      title: json['title']?.toString() ?? '',
      message: json['message']?.toString() ?? '',
      type: NotificationType.fromString(json['type']?.toString() ?? 'reminder'),
      isRead: json['is_read'] == true || json['is_read'] == 'true',
      readAt: parseDateTime(json['read_at']),
      data: json['data']?.toString(),
      createdAt: parseDateTime(json['created_at']) ?? DateTime.now(),
      updatedAt: parseDateTime(json['updated_at']) ?? DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'title': title,
      'message': message,
      'type': type.value,
      'is_read': isRead,
      'read_at': readAt?.toIso8601String(),
      'data': data,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

enum NotificationType {
  reminder('reminder'),
  task('task'),
  deal('deal'),
  activity('activity');

  const NotificationType(this.value);
  final String value;

  static NotificationType fromString(String value) {
    return NotificationType.values.firstWhere(
      (type) => type.value == value,
      orElse: () => NotificationType.reminder,
    );
  }
}

class NotificationListResponse {
  final List<Notification> items;
  final Pagination pagination;

  NotificationListResponse({
    required this.items,
    required this.pagination,
  });

  factory NotificationListResponse.fromJson(Map<String, dynamic> json) {
    // API response structure:
    // {
    //   "success": true,
    //   "data": Notification[],
    //   "meta": {
    //     "pagination": { ... }
    //   }
    // }
    
    final dynamic rawData = json['data'];
    List<Notification> items = [];

    // Parse notifications array
    if (rawData is List) {
      items = rawData
          .map((item) {
            if (item is Map<String, dynamic>) {
              return Notification.fromJson(item);
            }
            return null;
          })
          .whereType<Notification>()
          .toList();
    } else {
      // If data is not a list, return empty list
      items = [];
    }

    // Parse pagination from meta.pagination
    Pagination pagination;
    final meta = json['meta'] as Map<String, dynamic>?;
    if (meta != null && meta['pagination'] != null) {
      final paginationData = meta['pagination'];
      if (paginationData is Map<String, dynamic>) {
        pagination = Pagination.fromJson(paginationData);
      } else {
        // Fallback pagination
        pagination = Pagination(
          page: 1,
          perPage: items.length,
          total: items.length,
          totalPages: 1,
        );
      }
    } else {
      // Fallback pagination if meta is missing
      pagination = Pagination(
        page: 1,
        perPage: items.length,
        total: items.length,
        totalPages: 1,
      );
    }

    return NotificationListResponse(
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
    // Safely parse pagination data, handling both int and String types
    int parseToInt(dynamic value, int defaultValue) {
      if (value == null) return defaultValue;
      if (value is int) return value;
      if (value is String) {
        final parsed = int.tryParse(value);
        return parsed ?? defaultValue;
      }
      return defaultValue;
    }

    return Pagination(
      page: parseToInt(json['page'], 1),
      perPage: parseToInt(json['per_page'], 20),
      total: parseToInt(json['total'], 0),
      totalPages: parseToInt(json['total_pages'], 1),
    );
  }

  bool get hasNextPage => page < totalPages;
  bool get hasPrevPage => page > 1;
}

class UnreadCountResponse {
  final int unreadCount;

  UnreadCountResponse({required this.unreadCount});

  factory UnreadCountResponse.fromJson(Map<String, dynamic> json) {
    // Safely parse unread count, handling both int and String types
    int parseToInt(dynamic value, int defaultValue) {
      if (value == null) return defaultValue;
      if (value is int) return value;
      if (value is String) {
        final parsed = int.tryParse(value);
        return parsed ?? defaultValue;
      }
      return defaultValue;
    }

    final dynamic rawData = json['data'];
    if (rawData is Map<String, dynamic>) {
      return UnreadCountResponse(
        unreadCount: parseToInt(rawData['unread_count'], 0),
      );
    }
    return UnreadCountResponse(unreadCount: 0);
  }
}

