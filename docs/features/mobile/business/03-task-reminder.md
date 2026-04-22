# Business - Task & Reminder Management

## CRM Healthcare Mobile App - Flutter

**Module**: Business Domain  
**Sprint**: Sprint 3  
**Version**: 1.0  
**Status**: ✅ **Completed**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Ringkasan Fitur](#ringkasan-fitur)
2. [Fitur Utama](#fitur-utama)
3. [Business Rules](#business-rules)
4. [Keputusan Teknis & Trade-offs](#keputusan-teknis--trade-offs)
5. [Struktur Folder](#struktur-folder)
6. [API Endpoints](#api-endpoints)
7. [Data Models](#data-models)
8. [Configuration](#configuration)
9. [Usage Examples](#usage-examples)
10. [Cara Test Manual](#cara-test-manual)
11. [Dependencies](#dependencies)
12. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Fitur **Task & Reminder Management** memungkinkan sales rep dan supervisor untuk membuat, mengelola, dan melacak tasks terkait dengan sales activities. Tasks dapat dikaitkan dengan accounts, contacts, visit reports, atau deals. Reminders memastikan user tidak melewatkan deadlines atau follow-ups penting.

### Goals

- **Task Creation**: Create tasks dengan due date, priority, dan assignment
- **Task Tracking**: Track status dan progress dari tasks
- **Reminders**: Notifikasi untuk upcoming dan overdue tasks
- **Organization**: Tasks terorganisir berdasarkan account, contact, atau project
- **Collaboration**: Task assignment dan delegation

---

## Fitur Utama

### 1. Task List

**Features**:

- List tasks dengan pagination
- Filter by status (pending, in_progress, completed)
- Filter by priority (low, medium, high)
- Filter by due date (today, this week, overdue)
- Sort by due date atau priority
- Search by title atau description
- Pull-to-refresh

**Views**:

- List View: Compact list dengan quick actions
- Calendar View: Tasks displayed di calendar
- Kanban View: Board dengan columns per status

### 2. Task Creation

**Form Fields**:

- Title (required)
- Description
- Due date & time
- Priority (low, medium, high)
- Status (default: pending)
- Assigned to (self atau other user)
- Related to (account, contact, deal, visit report)
- Tags/Labels
- Attachments (optional)

**Quick Create**:

- Create task from account detail
- Create task from visit report
- Create task from dashboard

### 3. Task Detail

**Information Displayed**:

- Title dan description
- Due date dengan countdown
- Status dan priority badges
- Assignee info dengan avatar
- Related entity (account/contact info)
- Creation dan update timestamps
- Comments/activity log
- Subtasks (jika ada)
- Reminders list

### 4. Task Actions

**Quick Actions**:

- ✅ Complete Task
- 🔄 Change Status
- 🔔 Set Reminder
- ✏️ Edit Task
- 🗑️ Delete Task
- 👤 Reassign Task

### 5. Reminders

**Types**:

- **Due Date Reminder**: Notifikasi saat mendekati due date
- **Overdue Reminder**: Daily reminder untuk overdue tasks
- **Custom Reminder**: User-defined reminder time

**Reminder Times**:

- 1 hour before due date
- 1 day before due date (untuk high priority)
- At due date
- Custom time

### 6. Task Completion

**Completion Flow**:

1. User tap "Complete"
2. Optional: Add completion notes
3. Task status berubah ke "completed"
4. Completion timestamp recorded
5. Related entities updated

---

## Business Rules

### 1. Task Creation Rules

**Required Fields**:

- Title
- Due date

**Optional Fields**:

- Description
- Priority (default: medium)
- Assignee (default: creator)
- Related entity

**Validation**:

- Due date tidak boleh di masa lalu (warning, bukan error)
- Title max 200 characters
- Description max 2000 characters

### 2. Assignment Rules

**Self Assignment**:

- Creator otomatis menjadi assignee
- Creator dapat reassign ke user lain

**Assignment to Others**:

- Hanya supervisor atau admin yang dapat assign ke user lain
- User yang di-assign menerima notifikasi
- Assigned user dapat accept atau decline (jika enabled)

**Team Assignment**:

- Task dapat di-assign ke team/department
- First person yang complete claim the task

### 3. Due Date Rules

**Overdue Tasks**:

- Task dengan due date < now() dianggap overdue
- Overdue tasks highlighted dengan warna merah
- Daily reminder untuk overdue tasks

**Due Time**:

- Jika due time tidak di-set, default ke 23:59
- Due date dengan time zone consideration

### 4. Status Rules

**Status Values**:

- **Pending**: Task belum dimulai
- **In Progress**: Task sedang dikerjakan
- **Completed**: Task selesai
- **Cancelled**: Task dibatalkan (soft delete)

**Status Transitions**:

```
Pending → In Progress
Pending → Completed
Pending → Cancelled
In Progress → Completed
In Progress → Cancelled
Completed → Pending (reopen)
```

### 5. Reminder Rules

**Default Reminders**:

- High priority: 1 day before + 1 hour before
- Medium priority: 1 day before
- Low priority: At due date

**Custom Reminders**:

- User dapat set custom reminder times
- Max 3 custom reminders per task
- Reminder times tidak boleh di masa lalu

**Reminder Delivery**:

- Push notification jika FCM available
- Local notification sebagai fallback
- In-app notification badge

### 6. Permission Rules

**Task Creator**:

- Full CRUD access
- Can reassign task
- Can delete task

**Assignee**:

- Can view task details
- Can update status
- Can add comments
- Cannot delete task

**Supervisor/Admin**:

- Can view all tasks
- Can edit any task
- Can reassign any task
- Can delete any task

---

## Keputusan Teknis & Trade-offs

### Mengapa Multiple Views (List, Calendar, Kanban)?

**Keputusan**: Support multiple view modes untuk tasks.

**Alasan**:

- **User Preference**: Different users prefer different views
- **Use Cases**: List untuk quick overview, Calendar untuk time-based planning, Kanban untuk workflow tracking
- **Flexibility**: Users dapat switch views based on their current need

**Trade-off**: Lebih complex UI dan maintenance. **Mitigasi**: Default ke List view, others sebagai optional.

### Mengapa Soft Delete vs Hard Delete?

**Keputusan**: Implement soft delete (status = cancelled).

**Alasan**:

- **Audit Trail**: Keep record dari cancelled tasks
- **Recovery**: Can restore cancelled tasks
- **Analytics**: Include cancelled tasks di productivity reports

**Trade-off**: Database akan tumbuh lebih besar. **Mitigasi**: Archive old cancelled tasks setelah 1 tahun.

### Mengapa Local Reminders sebagai Fallback?

**Keputusan**: Use local notifications jika FCM tidak available.

**Alasan**:

- **Reliability**: Ensure reminders always delivered
- **Offline Support**: Reminders work saat offline
- **Battery**: Local reminders lebih battery efficient

**Trade-off**: Dual reminder system lebih complex. **Mitigasi**: Abstract reminder logic di service layer.

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── tasks/
│       ├── data/
│       │   ├── models/
│       │   │   ├── task_model.dart           # Task entity
│       │   │   ├── reminder_model.dart       # Reminder entity
│       │   │   └── task_filter.dart          # Filter parameters
│       │   └── task_repository.dart          # API & cache
│       ├── application/
│       │   ├── task_list_provider.dart       # List state
│       │   ├── task_detail_provider.dart     # Detail state
│       │   ├── task_form_provider.dart       # Form state
│       │   ├── task_completion_provider.dart # Completion logic
│       │   └── reminder_provider.dart        # Reminder management
│       └── presentation/
│           ├── screens/
│           │   ├── task_list_screen.dart
│           │   ├── task_detail_screen.dart
│           │   ├── task_form_screen.dart
│           │   └── task_calendar_screen.dart
│           └── widgets/
│               ├── task_card.dart
│               ├── task_list_item.dart
│               ├── task_status_badge.dart
│               ├── priority_indicator.dart
│               ├── due_date_chip.dart
│               └── task_filter_sheet.dart
├── core/
│   └── services/
│       └── reminder_service.dart             # Reminder scheduling
```

---

## API Endpoints

### Task CRUD

#### GET /api/v1/tasks

List tasks dengan filter.

**Query Parameters**:

```
?page=1&limit=20&status=pending&priority=high&assigned_to=uuid&account_id=uuid&due_before=2025-01-20&search=keyword
```

**Response**:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "task-uuid",
        "title": "Follow up with RS Medika",
        "description": "Call Dr. Smith to discuss quotation",
        "status": "pending",
        "priority": "high",
        "due_date": "2025-01-20T14:00:00Z",
        "assigned_to": {
          "id": "user-uuid",
          "name": "John Doe",
          "avatar": "https://..."
        },
        "created_by": {
          "id": "user-uuid-2",
          "name": "Supervisor Name"
        },
        "related_to": {
          "type": "account",
          "id": "account-uuid",
          "name": "RS Medika Hospital"
        },
        "tags": ["follow-up", "urgent"],
        "reminders": [
          {
            "id": "reminder-uuid",
            "remind_at": "2025-01-20T13:00:00Z"
          }
        ],
        "completed_at": null,
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 3,
      "total_items": 50
    },
    "stats": {
      "total": 50,
      "pending": 20,
      "in_progress": 15,
      "completed": 15,
      "overdue": 5
    }
  }
}
```

#### GET /api/v1/tasks/:id

Get task detail dengan comments dan activity log.

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "task-uuid",
    "title": "Follow up with RS Medika",
    "description": "Call Dr. Smith to discuss quotation",
    "status": "pending",
    "priority": "high",
    "due_date": "2025-01-20T14:00:00Z",
    "assigned_to": {
      "id": "user-uuid",
      "name": "John Doe",
      "avatar": "https://..."
    },
    "created_by": {
      "id": "user-uuid-2",
      "name": "Supervisor Name"
    },
    "related_to": {
      "type": "account",
      "id": "account-uuid",
      "name": "RS Medika Hospital"
    },
    "tags": ["follow-up", "urgent"],
    "reminders": [
      {
        "id": "reminder-uuid",
        "remind_at": "2025-01-20T13:00:00Z",
        "is_triggered": false
      }
    ],
    "comments": [
      {
        "id": "comment-uuid",
        "user": {
          "id": "user-uuid",
          "name": "John Doe"
        },
        "content": "Called, left voicemail",
        "created_at": "2025-01-19T10:00:00Z"
      }
    ],
    "activity_log": [
      {
        "action": "created",
        "user": "Supervisor Name",
        "timestamp": "2025-01-15T10:00:00Z"
      },
      {
        "action": "assigned",
        "user": "Supervisor Name",
        "timestamp": "2025-01-15T10:00:00Z"
      }
    ],
    "completed_at": null,
    "completion_notes": null,
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

#### POST /api/v1/tasks

Create new task.

**Request**:

```json
{
  "title": "Follow up with RS Medika",
  "description": "Call Dr. Smith to discuss quotation",
  "priority": "high",
  "due_date": "2025-01-20T14:00:00Z",
  "assigned_to": "user-uuid",
  "related_type": "account",
  "related_id": "account-uuid",
  "tags": ["follow-up", "urgent"],
  "reminders": [{ "remind_at": "2025-01-20T13:00:00Z" }]
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "task-uuid",
    "message": "Task created successfully"
  }
}
```

#### PUT /api/v1/tasks/:id

Update task.

**Request**:

```json
{
  "title": "Updated title",
  "description": "Updated description",
  "status": "in_progress",
  "priority": "medium",
  "due_date": "2025-01-21T14:00:00Z"
}
```

#### POST /api/v1/tasks/:id/complete

Complete task.

**Request**:

```json
{
  "completion_notes": "Successfully discussed and agreed on terms"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "status": "completed",
    "completed_at": "2025-01-20T15:30:00Z",
    "message": "Task completed successfully"
  }
}
```

#### POST /api/v1/tasks/:id/assign

Reassign task.

**Request**:

```json
{
  "assigned_to": "new-user-uuid"
}
```

#### DELETE /api/v1/tasks/:id

Soft delete task (set status = cancelled).

### Reminder Endpoints

#### POST /api/v1/tasks/:id/reminders

Add reminder ke task.

**Request**:

```json
{
  "remind_at": "2025-01-20T13:00:00Z"
}
```

#### DELETE /api/v1/tasks/:id/reminders/:reminderId

Delete reminder.

---

## Data Models

### Task Model

```dart
@freezed
class Task with _$Task {
  const factory Task({
    required String id,
    required String title,
    String? description,
    @Default('pending') String status,
    @Default('medium') String priority,
    required DateTime dueDate,
    required User assignedTo,
    required User createdBy,
    RelatedEntity? relatedTo,
    @Default([]) List<String> tags,
    @Default([]) List<Reminder> reminders,
    @Default([]) List<Comment> comments,
    DateTime? completedAt,
    String? completionNotes,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _Task;

  factory Task.fromJson(Map<String, dynamic> json) =
      _$TaskFromJson(json);
}

enum TaskStatus {
  pending,
  inProgress,
  completed,
  cancelled;

  String get displayName {
    switch (this) {
      case TaskStatus.pending:
        return 'Pending';
      case TaskStatus.inProgress:
        return 'In Progress';
      case TaskStatus.completed:
        return 'Completed';
      case TaskStatus.cancelled:
        return 'Cancelled';
    }
  }

  Color get color {
    switch (this) {
      case TaskStatus.pending:
        return Colors.orange;
      case TaskStatus.inProgress:
        return Colors.blue;
      case TaskStatus.completed:
        return Colors.green;
      case TaskStatus.cancelled:
        return Colors.grey;
    }
  }
}

enum TaskPriority {
  low,
  medium,
  high;

  String get displayName {
    switch (this) {
      case TaskPriority.low:
        return 'Low';
      case TaskPriority.medium:
        return 'Medium';
      case TaskPriority.high:
        return 'High';
    }
  }

  Color get color {
    switch (this) {
      case TaskPriority.low:
        return Colors.grey;
      case TaskPriority.medium:
        return Colors.orange;
      case TaskPriority.high:
        return Colors.red;
    }
  }

  IconData get icon {
    switch (this) {
      case TaskPriority.low:
        return Icons.arrow_downward;
      case TaskPriority.medium:
        return Icons.remove;
      case TaskPriority.high:
        return Icons.arrow_upward;
    }
  }
}
```

### Reminder Model

```dart
@freezed
class Reminder with _$Reminder {
  const factory Reminder({
    required String id,
    required DateTime remindAt,
    @Default(false) bool isTriggered,
    DateTime? triggeredAt,
  }) = _Reminder;

  factory Reminder.fromJson(Map<String, dynamic> json) =
      _$ReminderFromJson(json);
}
```

### Related Entity Model

```dart
@freezed
class RelatedEntity with _$RelatedEntity {
  const factory RelatedEntity({
    required String type, // account, contact, deal, visit_report
    required String id,
    required String name,
    String? subtitle,
  }) = _RelatedEntity;

  factory RelatedEntity.fromJson(Map<String, dynamic> json) =
      _$RelatedEntityFromJson(json);

  factory RelatedEntity.fromAccount(Account account) {
    return RelatedEntity(
      type: 'account',
      id: account.id,
      name: account.name,
      subtitle: account.type,
    );
  }
}
```

### Task Filter Model

```dart
@freezed
class TaskFilter with _$TaskFilter {
  const factory TaskFilter({
    String? search,
    List<String>? statuses,
    List<String>? priorities,
    DateTime? dueBefore,
    DateTime? dueAfter,
    String? assignedTo,
    String? accountId,
    String? sortBy, // due_date, priority, created_at
    String? sortOrder, // asc, desc
    @Default(1) int page,
    @Default(20) int limit,
  }) = _TaskFilter;

  Map<String, dynamic> toQueryParameters() {
    return {
      if (search != null && search!.isNotEmpty) 'search': search,
      if (statuses != null && statuses!.isNotEmpty) 'status': statuses,
      if (priorities != null && priorities!.isNotEmpty) 'priority': priorities,
      if (dueBefore != null) 'due_before': dueBefore!.toIso8601String(),
      if (dueAfter != null) 'due_after': dueAfter!.toIso8601String(),
      if (assignedTo != null) 'assigned_to': assignedTo,
      if (accountId != null) 'account_id': accountId,
      if (sortBy != null) 'sort_by': sortBy,
      if (sortOrder != null) 'sort_order': sortOrder,
      'page': page,
      'limit': limit,
    };
  }
}
```

---

## Configuration

### Reminder Service

**File**: `core/services/reminder_service.dart`

```dart
class ReminderService {
  final LocalNotificationService _localNotifications;
  final ApiClient _apiClient;

  ReminderService(this._localNotifications, this._apiClient);

  Future<void> scheduleReminder(Task task, Reminder reminder) async {
    await _localNotifications.schedule(
      id: reminder.id.hashCode,
      title: 'Task Reminder',
      body: task.title,
      scheduledDate: reminder.remindAt,
      payload: jsonEncode({
        'type': 'task_reminder',
        'task_id': task.id,
      }),
    );
  }

  Future<void> scheduleDefaultReminders(Task task) async {
    final now = DateTime.now();
    final reminders = <Reminder>[];

    // High priority: 1 day + 1 hour before
    if (task.priority == 'high') {
      final oneDayBefore = task.dueDate.subtract(const Duration(days: 1));
      if (oneDayBefore.isAfter(now)) {
        reminders.add(Reminder(
          id: '${task.id}_1d',
          remindAt: oneDayBefore,
        ));
      }

      final oneHourBefore = task.dueDate.subtract(const Duration(hours: 1));
      if (oneHourBefore.isAfter(now)) {
        reminders.add(Reminder(
          id: '${task.id}_1h',
          remindAt: oneHourBefore,
        ));
      }
    }
    // Medium priority: 1 day before
    else if (task.priority == 'medium') {
      final oneDayBefore = task.dueDate.subtract(const Duration(days: 1));
      if (oneDayBefore.isAfter(now)) {
        reminders.add(Reminder(
          id: '${task.id}_1d',
          remindAt: oneDayBefore,
        ));
      }
    }
    // Low priority: at due date
    else {
      reminders.add(Reminder(
        id: '${task.id}_0d',
        remindAt: task.dueDate,
      ));
    }

    // Schedule all reminders
    for (final reminder in reminders) {
      await scheduleReminder(task, reminder);
    }

    // Save to backend
    for (final reminder in reminders) {
      await _apiClient.post(
        '/api/v1/tasks/${task.id}/reminders',
        data: {'remind_at': reminder.remindAt.toIso8601String()},
      );
    }
  }

  Future<void> cancelReminders(String taskId) async {
    // Cancel local notifications
    await _localNotifications.cancel(taskId.hashCode);
  }
}
```

---

## Usage Examples

### Task List dengan Filter

```dart
class TaskListScreen extends ConsumerWidget {
  const TaskListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(taskListProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Tasks'),
        actions: [
          IconButton(
            icon: const Icon(Icons.filter_list),
            onPressed: () => _showFilterSheet(context, ref),
          ),
          IconButton(
            icon: const Icon(Icons.calendar_today),
            onPressed: () => context.push(AppRoutes.taskCalendar),
          ),
        ],
      ),
      body: Column(
        children: [
          // Quick filters
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Row(
              children: [
                FilterChip(
                  label: const Text('Today'),
                  selected: ref.watch(taskFilterProvider).dueAfter?.isToday ?? false,
                  onSelected: (selected) {
                    ref.read(taskFilterProvider.notifier).setDueDate(
                      start: DateTime.now(),
                      end: DateTime.now(),
                    );
                  },
                ),
                const SizedBox(width: 8),
                FilterChip(
                  label: const Text('Overdue'),
                  selected: false,
                  onSelected: (selected) {
                    ref.read(taskFilterProvider.notifier).setOverdue();
                  },
                ),
                const SizedBox(width: 8),
                FilterChip(
                  label: const Text('High Priority'),
                  selected: ref.watch(taskFilterProvider).priorities?.contains('high') ?? false,
                  onSelected: (selected) {
                    ref.read(taskFilterProvider.notifier).setPriority('high');
                  },
                ),
              ],
            ),
          ),
          // Task list
          Expanded(
            child: state.when(
              loading: () => const TaskListSkeleton(),
              error: (error) => ErrorWidget(error: error),
              data: (result) => ListView.builder(
                itemCount: result.tasks.length,
                itemBuilder: (context, index) {
                  final task = result.tasks[index];
                  return TaskListItem(
                    task: task,
                    onTap: () => context.push(
                      AppRoutes.taskDetailPath(task.id),
                    ),
                    onComplete: () => _completeTask(context, ref, task),
                  );
                },
              ),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push(AppRoutes.taskCreate),
        child: const Icon(Icons.add),
      ),
    );
  }

  Future<void> _completeTask(BuildContext context, WidgetRef ref, Task task) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Complete Task'),
        content: const Text('Mark this task as completed?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Complete'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      await ref.read(taskCompletionProvider.notifier).completeTask(task.id);
    }
  }
}
```

### Task Form

```dart
class TaskFormScreen extends ConsumerStatefulWidget {
  final String? taskId;
  final String? prefillAccountId;

  const TaskFormScreen({
    super.key,
    this.taskId,
    this.prefillAccountId,
  });

  @override
  ConsumerState<TaskFormScreen> createState() => _TaskFormScreenState();
}

class _TaskFormScreenState extends ConsumerState<TaskFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();

  DateTime? _dueDate;
  TaskPriority _priority = TaskPriority.medium;
  String? _assignedTo;
  RelatedEntity? _relatedTo;

  @override
  void initState() {
    super.initState();
    if (widget.prefillAccountId != null) {
      // Prefill account
      _loadAccount(widget.prefillAccountId!);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.taskId != null ? 'Edit Task' : 'Create Task'),
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Title
            TextFormField(
              controller: _titleController,
              decoration: const InputDecoration(
                labelText: 'Title *',
                hintText: 'What needs to be done?',
              ),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Title is required';
                }
                return null;
              },
            ),

            const SizedBox(height: 16),

            // Due Date
            ListTile(
              leading: const Icon(Icons.calendar_today),
              title: const Text('Due Date *'),
              subtitle: Text(
                _dueDate != null
                    ? DateFormat('MMM dd, yyyy - HH:mm').format(_dueDate!)
                    : 'Select due date',
              ),
              onTap: () => _pickDueDate(context),
            ),

            const SizedBox(height: 16),

            // Priority
            SegmentedButton<TaskPriority>(
              segments: TaskPriority.values.map((priority) {
                return ButtonSegment(
                  value: priority,
                  label: Text(priority.displayName),
                  icon: Icon(priority.icon),
                );
              }).toList(),
              selected: {_priority},
              onSelectionChanged: (selected) {
                setState(() {
                  _priority = selected.first;
                });
              },
            ),

            const SizedBox(height: 16),

            // Assigned To
            ListTile(
              leading: const Icon(Icons.person),
              title: const Text('Assigned To'),
              subtitle: Text(_assignedTo ?? 'Assign to me'),
              onTap: () => _pickAssignee(context),
            ),

            const SizedBox(height: 16),

            // Related To
            ListTile(
              leading: const Icon(Icons.link),
              title: const Text('Related To'),
              subtitle: Text(_relatedTo?.name ?? 'Select related account/contact'),
              onTap: () => _pickRelatedEntity(context),
            ),

            const SizedBox(height: 16),

            // Description
            TextFormField(
              controller: _descriptionController,
              maxLines: 5,
              decoration: const InputDecoration(
                labelText: 'Description',
                hintText: 'Add more details...',
                alignLabelWithHint: true,
              ),
            ),

            const SizedBox(height: 24),

            // Submit Button
            FilledButton(
              onPressed: _submit,
              child: const Text('Save Task'),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _pickDueDate(BuildContext context) async {
    final date = await showDatePicker(
      context: context,
      initialDate: _dueDate ?? DateTime.now(),
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 365)),
    );

    if (date != null) {
      final time = await showTimePicker(
        context: context,
        initialTime: TimeOfDay.now(),
      );

      if (time != null) {
        setState(() {
          _dueDate = DateTime(
            date.year,
            date.month,
            date.day,
            time.hour,
            time.minute,
          );
        });
      }
    }
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    if (_dueDate == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please select due date')),
      );
      return;
    }

    final task = Task(
      id: widget.taskId ?? '',
      title: _titleController.text,
      description: _descriptionController.text.isNotEmpty
          ? _descriptionController.text
          : null,
      priority: _priority.name,
      dueDate: _dueDate!,
      assignedTo: User(id: _assignedTo ?? 'self', name: ''),
      createdBy: User(id: 'self', name: ''),
      relatedTo: _relatedTo,
    );

    if (widget.taskId != null) {
      await ref.read(taskFormProvider.notifier).updateTask(task);
    } else {
      await ref.read(taskFormProvider.notifier).createTask(task);
    }

    if (mounted) {
      Navigator.pop(context);
    }
  }
}
```

---

## Cara Test Manual

### Test Task Creation

1. **Create Task**:
   - Tap FAB (+)
   - Fill title dan due date
   - Select priority
   - Save
   - Verifikasi: Task muncul di list

2. **Task dengan Reminder**:
   - Create high priority task
   - Verifikasi: Reminders scheduled
   - Wait untuk reminder time
   - Verifikasi: Notification muncul

3. **Quick Create dari Account**:
   - Buka account detail
   - Tap "Add Task"
   - Verifikasi: Account pre-selected

### Test Task Completion

1. **Complete Task**:
   - Tap task di list
   - Tap "Complete"
   - Verifikasi: Status berubah ke completed
   - Verifikasi: Completion timestamp recorded

2. **Overdue Tasks**:
   - Create task dengan due date di masa lalu
   - Verifikasi: Overdue indicator muncul
   - Verifikasi: Warna merah/urgent styling

### Test Filters

1. **Filter by Status**:
   - Apply filter "Pending"
   - Verifikasi: Hanya pending tasks yang ditampilkan

2. **Filter by Priority**:
   - Apply filter "High"
   - Verifikasi: Hanya high priority tasks

3. **Filter by Date**:
   - Filter "Today"
   - Verifikasi: Hanya tasks due today

---

## Dependencies

### Internal

- `core/services/reminder_service.dart` - Reminder scheduling
- `core/network/api_client.dart` - API calls
- `features/accounts/data/account_repository.dart` - Account selection

### External

- `table_calendar: ^3.0.0` - Calendar view (optional)
- `intl: ^0.18.0` - Date formatting
- `flutter_riverpod: ^2.4.0` - State management
- `freezed: ^2.4.0` - Immutable models

---

## Notes & Improvements

### Known Limitations

1. **No Recurring Tasks**: Belum support recurring tasks (daily, weekly, monthly).

2. **Limited Collaboration**: No real-time collaboration features.

3. **No Time Tracking**: Tidak ada tracking waktu yang dihabiskan untuk task.

4. **No Dependencies**: Tasks tidak dapat memiliki dependencies satu sama lain.

### Future Improvements

1. **Recurring Tasks**: Support untuk repeating tasks

2. **Time Tracking**: Track actual time spent on tasks

3. **Task Dependencies**: Set task dependencies (task A harus selesai sebelum task B)

4. **Gantt Chart**: Visual timeline untuk complex projects

5. **Team Collaboration**: Real-time comments dan updates

6. **Task Templates**: Pre-defined templates untuk common tasks

7. **AI Suggestions**: AI-powered task prioritization dan scheduling

---

**Document Status**: Active  
**Last Updated**: January 2025  
**Maintained By**: Dev3 (Mobile Development Team)
