import 'package:dio/dio.dart';
import '../../../core/network/pagination.dart';

import '../../../core/network/connectivity_service.dart';
import '../../../core/storage/offline_storage.dart';
import 'models/task.dart';

class TaskRepository {
  TaskRepository(this._dio, this._connectivity);

  final Dio _dio;
  final ConnectivityService _connectivity;

  Future<TaskListResponse> getTasks({
    int page = 1,
    int perPage = 20,
    String? search,
    String? status,
    String? priority,
    String? type,
    String? assignedTo,
    String? accountId,
    String? contactId,
    DateTime? dueDateFrom,
    DateTime? dueDateTo,
    bool forceRefresh = false,
    Function(TaskListResponse)? onBackgroundUpdate,
  }) async {
    // 1. Try to load from cache first (offline-first) - only for first page and no filters
    if (!forceRefresh &&
        page == 1 &&
        (search == null || search.isEmpty) &&
        status == null &&
        priority == null &&
        accountId == null &&
        contactId == null) {
      final cachedTasks = await OfflineStorage.getTasks();
      if (cachedTasks != null && cachedTasks.isNotEmpty) {
        try {
          final tasks = cachedTasks.map((json) => Task.fromJson(json)).toList();
          final cachedResponse = TaskListResponse(
            items: tasks,
            pagination: Pagination(
              page: 1,
              perPage: tasks.length,
              total: tasks.length,
              totalPages: 1,
            ),
          );

          // Trigger background refresh if online
          if (_connectivity.isOnline && !forceRefresh) {
            _fetchAndUpdateInBackground(
              page: page,
              perPage: perPage,
              search: search,
              status: status,
              priority: priority,
              type: type,
              assignedTo: assignedTo,
              accountId: accountId,
              contactId: contactId,
              dueDateFrom: dueDateFrom,
              dueDateTo: dueDateTo,
              onBackgroundUpdate: onBackgroundUpdate,
            );
          }

          return cachedResponse;
        } catch (e) {
          // If parsing fails, continue to API call
        }
      }
    }

    // 2. If online, fetch from API
    if (_connectivity.isOnline) {
      try {
        final queryParams = <String, dynamic>{
          'page': page,
          'per_page': perPage,
        };

        if (search != null && search.isNotEmpty) {
          queryParams['search'] = search;
        }
        if (status != null && status.isNotEmpty) {
          queryParams['status'] = status;
        }
        if (priority != null && priority.isNotEmpty) {
          queryParams['priority'] = priority;
        }
        if (type != null && type.isNotEmpty) {
          queryParams['type'] = type;
        }
        if (assignedTo != null && assignedTo.isNotEmpty) {
          queryParams['assigned_to'] = assignedTo;
        }
        if (accountId != null && accountId.isNotEmpty) {
          queryParams['account_id'] = accountId;
        }
        if (contactId != null && contactId.isNotEmpty) {
          queryParams['contact_id'] = contactId;
        }
        if (dueDateFrom != null) {
          queryParams['due_date_from'] = dueDateFrom.toIso8601String();
        }
        if (dueDateTo != null) {
          queryParams['due_date_to'] = dueDateTo.toIso8601String();
        }

        // Use mobile endpoint for tasks
        final response = await _dio.get(
          '/api/v1/mobile/tasks/my-tasks',
          queryParameters: queryParams,
        );

        if (response.data is Map<String, dynamic>) {
          final responseData = response.data as Map<String, dynamic>;
          if (responseData['success'] == true) {
            final taskListResponse = TaskListResponse.fromJson(responseData);

            // 3. Save to cache (only for first page and no filters)
            if (page == 1 &&
                (search == null || search.isEmpty) &&
                status == null &&
                priority == null &&
                accountId == null &&
                contactId == null) {
              final tasksJson = taskListResponse.items
                  .map((task) => task.toJson())
                  .toList();
              await OfflineStorage.saveTasks(tasksJson);
            }

            return taskListResponse;
          } else {
            throw Exception(
              responseData['error']?['message'] ?? 'Failed to fetch tasks',
            );
          }
        } else {
          throw Exception('Invalid response format');
        }
      } on DioException catch (e) {
        // If API fails, try to return cached data if available
        if (page == 1 &&
            (search == null || search.isEmpty) &&
            status == null &&
            priority == null &&
            accountId == null &&
            contactId == null) {
          final cachedTasks = await OfflineStorage.getTasks();
          if (cachedTasks != null && cachedTasks.isNotEmpty) {
            try {
              final tasks = cachedTasks
                  .map((json) => Task.fromJson(json))
                  .toList();
              return TaskListResponse(
                items: tasks,
                pagination: Pagination(
                  page: 1,
                  perPage: tasks.length,
                  total: tasks.length,
                  totalPages: 1,
                ),
              );
            } catch (_) {
              // Ignore parsing errors
            }
          }
        }

        if (e.response != null) {
          final errorData = e.response!.data;
          throw Exception(
            errorData['error']?['message'] ?? 'Failed to fetch tasks',
          );
        } else {
          throw Exception('Network error: ${e.message}');
        }
      } catch (e) {
        // If other error, try cached data
        if (page == 1 &&
            (search == null || search.isEmpty) &&
            status == null &&
            priority == null &&
            accountId == null &&
            contactId == null) {
          final cachedTasks = await OfflineStorage.getTasks();
          if (cachedTasks != null && cachedTasks.isNotEmpty) {
            try {
              final tasks = cachedTasks
                  .map((json) => Task.fromJson(json))
                  .toList();
              return TaskListResponse(
                items: tasks,
                pagination: Pagination(
                  page: 1,
                  perPage: tasks.length,
                  total: tasks.length,
                  totalPages: 1,
                ),
              );
            } catch (_) {
              // Ignore parsing errors
            }
          }
        }
        throw Exception('Failed to fetch tasks: $e');
      }
    }

    // 4. Offline: return cached data or throw error
    if (page == 1 &&
        (search == null || search.isEmpty) &&
        status == null &&
        priority == null &&
        accountId == null &&
        contactId == null) {
      final cachedTasks = await OfflineStorage.getTasks();
      if (cachedTasks != null && cachedTasks.isNotEmpty) {
        try {
          final tasks = cachedTasks.map((json) => Task.fromJson(json)).toList();
          return TaskListResponse(
            items: tasks,
            pagination: Pagination(
              page: 1,
              perPage: tasks.length,
              total: tasks.length,
              totalPages: 1,
            ),
          );
        } catch (e) {
          throw Exception('Failed to load cached tasks: $e');
        }
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Fetch tasks in background and update cache + UI
  Future<void> _fetchAndUpdateInBackground({
    required int page,
    required int perPage,
    String? search,
    String? status,
    String? priority,
    String? type,
    String? assignedTo,
    String? accountId,
    String? contactId,
    DateTime? dueDateFrom,
    DateTime? dueDateTo,
    Function(TaskListResponse)? onBackgroundUpdate,
  }) async {
    try {
      final queryParams = <String, dynamic>{'page': page, 'per_page': perPage};

      if (search != null && search.isNotEmpty) {
        queryParams['search'] = search;
      }
      if (status != null && status.isNotEmpty) {
        queryParams['status'] = status;
      }
      if (priority != null && priority.isNotEmpty) {
        queryParams['priority'] = priority;
      }
      if (type != null && type.isNotEmpty) {
        queryParams['type'] = type;
      }
      if (assignedTo != null && assignedTo.isNotEmpty) {
        queryParams['assigned_to'] = assignedTo;
      }
      if (accountId != null && accountId.isNotEmpty) {
        queryParams['account_id'] = accountId;
      }
      if (contactId != null && contactId.isNotEmpty) {
        queryParams['contact_id'] = contactId;
      }
      if (dueDateFrom != null) {
        queryParams['due_date_from'] = dueDateFrom.toIso8601String();
      }
      if (dueDateTo != null) {
        queryParams['due_date_to'] = dueDateTo.toIso8601String();
      }

      final response = await _dio.get(
        '/api/v1/mobile/tasks/my-tasks',
        queryParameters: queryParams,
      );

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          final taskListResponse = TaskListResponse.fromJson(responseData);

          // Save to cache
          if (page == 1 &&
              (search == null || search.isEmpty) &&
              status == null &&
              priority == null &&
              accountId == null &&
              contactId == null) {
            final tasksJson = taskListResponse.items
                .map((task) => task.toJson())
                .toList();
            await OfflineStorage.saveTasks(tasksJson);
          }

          // Notify UI to update with fresh data
          onBackgroundUpdate?.call(taskListResponse);
        }
      }
    } catch (e) {
      // Silently fail in background
    }
  }

  Future<Task> getTaskById(String id, {bool forceRefresh = false}) async {
    // Get cached task for fallback (only used when offline)
    final cachedTask = await OfflineStorage.getTaskDetail(id);

    // If online, always fetch from API (skip cache)
    if (_connectivity.isOnline) {
      try {
        final response = await _dio.get('/api/v1/tasks/$id');

        if (response.data is Map<String, dynamic>) {
          final responseData = response.data as Map<String, dynamic>;
          if (responseData['success'] == true) {
            final task = Task.fromJson(
              responseData['data'] as Map<String, dynamic>,
            );

            // Save to cache for offline use
            await OfflineStorage.saveTaskDetail(id, task.toJson());

            return task;
          } else {
            throw Exception(
              responseData['error']?['message'] ?? 'Failed to fetch task',
            );
          }
        } else {
          throw Exception('Invalid response format');
        }
      } on DioException catch (e) {
        // If API fails and we have cached data, use it as fallback
        if (cachedTask != null) {
          try {
            return Task.fromJson(cachedTask);
          } catch (_) {
            // Ignore parsing errors, continue to throw original error
          }
        }

        if (e.response != null) {
          final errorData = e.response!.data;
          throw Exception(
            errorData['error']?['message'] ?? 'Failed to fetch task',
          );
        } else {
          throw Exception('Network error: ${e.message}');
        }
      } catch (e) {
        // If other error and we have cached data, use it as fallback
        if (cachedTask != null) {
          try {
            return Task.fromJson(cachedTask);
          } catch (_) {
            // Ignore parsing errors
          }
        }
        throw Exception('Failed to fetch task: $e');
      }
    }

    // Offline: return cached data or throw error
    if (cachedTask != null) {
      try {
        return Task.fromJson(cachedTask);
      } catch (e) {
        throw Exception('Failed to load cached task: $e');
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  Future<Task> markInProgress(String id) async {
    try {
      final response = await _dio.post('/api/v1/tasks/$id/mark-in-progress');

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return Task.fromJson(responseData['data'] as Map<String, dynamic>);
        } else {
          throw Exception(
            responseData['error']?['message'] ??
                'Failed to mark task in progress',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to mark task in progress',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to mark task in progress: $e');
    }
  }

  Future<Task> completeTask(String id) async {
    try {
      final response = await _dio.post('/api/v1/tasks/$id/complete');

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return Task.fromJson(responseData['data'] as Map<String, dynamic>);
        } else {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to complete task',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to complete task',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to complete task: $e');
    }
  }

  Future<Reminder> createReminder({
    required String taskId,
    required DateTime remindAt,
    String reminderType = 'in_app',
    String? message,
  }) async {
    try {
      // Format remind_at in ISO8601 format with timezone offset
      // Example: "2024-01-20T10:00:00+07:00"
      final localDate = remindAt.toLocal();
      final offset = localDate.timeZoneOffset;

      final year = localDate.year.toString().padLeft(4, '0');
      final month = localDate.month.toString().padLeft(2, '0');
      final day = localDate.day.toString().padLeft(2, '0');
      final hour = localDate.hour.toString().padLeft(2, '0');
      final minute = localDate.minute.toString().padLeft(2, '0');
      final second = localDate.second.toString().padLeft(2, '0');

      final offsetHours = offset.inHours.abs().toString().padLeft(2, '0');
      final offsetMinutes = (offset.inMinutes.abs() % 60).toString().padLeft(
        2,
        '0',
      );
      final offsetSign = offset.isNegative ? '-' : '+';

      final remindAtFormatted =
          '$year-$month-${day}T$hour:$minute:$second$offsetSign$offsetHours:$offsetMinutes';

      final data = <String, dynamic>{
        'task_id': taskId,
        'remind_at': remindAtFormatted,
        'reminder_type': reminderType,
      };

      if (message != null && message.trim().isNotEmpty) {
        data['message'] = message.trim();
      }

      final response = await _dio.post('/api/v1/tasks/reminders', data: data);

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return Reminder.fromJson(
            responseData['data'] as Map<String, dynamic>,
          );
        } else {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to create reminder',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to create reminder',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to create reminder: $e');
    }
  }

  Future<void> deleteReminder(String id) async {
    try {
      final response = await _dio.delete('/api/v1/tasks/reminders/$id');

      // Handle response - DELETE might return 204 No Content or success response
      if (response.statusCode == 204 || response.statusCode == 200) {
        // Success - no need to check response.data
        return;
      }

      // If response has data, check for success flag
      if (response.data != null && response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == false) {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to delete reminder',
          );
        }
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to delete reminder',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to delete reminder: $e');
    }
  }
}
