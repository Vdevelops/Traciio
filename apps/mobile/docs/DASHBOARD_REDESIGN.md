# Dokumentasi Rancangan Dashboard Mobile - Redesign untuk Sales

## 📋 Daftar Isi

1. [Overview](#overview)
2. [UI/UX Design](#uiux-design)
3. [API Design](#api-design)
4. [Data Models](#data-models)
5. [Implementation Plan](#implementation-plan)
6. [Testing Strategy](#testing-strategy)

---

## Overview

### Tujuan
Membuat dashboard mobile yang lebih sederhana dan mudah digunakan untuk user sales yang sering berpergian ke client. Dashboard akan fokus pada 3 layout utama:

1. **Quick Stats** - Target user dan progress
2. **Visits** - Visit yang sedang dikerjakan atau sudah dilakukan
3. **Upcoming Tasks** - Tasks yang di-assign ke user

### Prinsip Desain
- **Simplicity**: Fokus pada informasi yang paling penting
- **Mobile-First**: Optimized untuk penggunaan di mobile saat di lapangan
- **Quick Access**: Informasi penting dapat diakses dengan cepat
- **Offline Support**: Dapat melihat data cached saat offline

---

## UI/UX Design

### Layout Struktur

```
┌─────────────────────────────────────┐
│         Dashboard Header            │
│  (Logo, Notifications, Profile)    │
├─────────────────────────────────────┤
│      Quick Stats Section            │
│  ┌──────────┐  ┌──────────┐       │
│  │  Target  │  │ Progress │       │
│  └──────────┘  └──────────┘       │
├─────────────────────────────────────┤
│      Visits Section                 │
│  ┌─────────────────────────────┐   │
│  │  [Tab: Active | Completed]  │   │
│  ├─────────────────────────────┤   │
│  │  ← [Card] [Card] [Card] →  │   │
│  │  (Horizontal Scrollable)    │   │
│  │  (Max 5 items)              │   │
│  └─────────────────────────────┘   │
├─────────────────────────────────────┤
│      Upcoming Tasks Section         │
│  ┌─────────────────────────────┐   │
│  │  Tasks (Max 3) [View All →]│   │
│  ├─────────────────────────────┤   │
│  │  - Task Card 1              │   │
│  │  - Task Card 2              │   │
│  │  - Task Card 3              │   │
│  │  (Fixed List, No Scroll)     │   │
│  └─────────────────────────────┘   │
└─────────────────────────────────────┘
```

### 1. Quick Stats Section

#### Komponen
- **Target Card**: Menampilkan target bulanan user
  - Target amount (Rp format)
  - Achieved amount (Rp format)
  - Progress bar dengan persentase
  - Progress indicator (color-coded: green >80%, yellow 50-80%, red <50%)

- **Quick Metrics** (Optional - bisa di-expand):
  - Total Visits Today
  - Completed Visits
  - Pending Visits
  - Total Tasks

#### Design Specs
```
┌─────────────────────────────────────┐
│  🎯 Target Progress                  │
│  ─────────────────────────────────  │
│  Rp 15.000.000 / Rp 20.000.000     │
│  ▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░ 75%          │
│  ─────────────────────────────────  │
│  [View Details →]                  │
└─────────────────────────────────────┘
```

#### Behavior
- Tap untuk expand detail target
- Swipe untuk refresh data
- Auto-refresh setiap 5 menit saat app aktif

### 2. Visits Section

#### Tab Structure
- **Active Tab**: Visit yang sedang dikerjakan (status: draft, in-progress)
- **Completed Tab**: Visit yang sudah selesai (status: approved, completed)

#### Visit Card Design
```
┌─────────────────────┐
│  🏢 Account Name    │
│  📍 Address         │
│  📅 Date | ⏰ Time  │
│  🏷️  Status Badge   │
│  [Action Button]   │
└─────────────────────┘
```

**Card Dimensions**:
- Width: 280dp (fixed)
- Height: 180dp (fixed)
- Margin: 12dp between cards

#### Card States
- **Draft**: Gray badge, "Start Visit" button
- **In Progress**: Orange badge, "Check Out" button
- **Completed**: Green badge, "View Details" button
- **Pending Approval**: Yellow badge, "View Details" button

#### Behavior
- **Horizontal Scrollable**: Cards di-scroll secara horizontal (left-right)
- **Maximum Items**: Maksimal 5 items per tab (untuk performance)
- **Pagination Indicator**: Tampilkan dots indicator jika ada lebih dari 5 items
- **Tap card**: Navigate ke detail visit
- **Swipe left/right**: Scroll ke card berikutnya/sebelumnya
- **Pull to refresh**: Refresh data saat pull down pada section header
- **View All Button**: Jika ada lebih dari 5 items, tampilkan "View All" button di kanan header untuk navigate ke full visits list

### 3. Upcoming Tasks Section

#### Section Header
```
┌─────────────────────────────────────┐
│  📋 Upcoming Tasks    [View All →] │
└─────────────────────────────────────┘
```

#### Task Card Design
```
┌─────────────────────────────────────┐
│  ⚠️  [Priority Indicator]           │
│  📋 Task Title                      │
│  📅 Due Date | ⏰ Time              │
│  👤 Assigned By                     │
│  [Action: Complete/Start]          │
└─────────────────────────────────────┘
```

#### Priority Indicators
- **High**: Red dot + "High Priority" badge
- **Medium**: Orange dot + "Medium Priority" badge
- **Low**: Blue dot + "Low Priority" badge
- **Overdue**: Red border + "Overdue" badge

#### Behavior
- **Fixed List**: Tidak scrollable, maksimal 3 items ditampilkan
- **Sort by**: Due date (ascending), Priority (high first)
- **View All Button**: Tampilkan di header section untuk navigate ke full tasks list screen
- **Tap card**: Navigate ke detail task
- **Empty State**: Jika tidak ada tasks, tampilkan empty state dengan message "No upcoming tasks"
- **Overflow Handling**: Jika ada lebih dari 3 tasks, hanya tampilkan 3 teratas (sorted by priority & due date)

### Color Scheme

#### Primary Colors
- **Primary**: Orange (#FF6B35) - untuk CTA buttons
- **Success**: Green (#4CAF50) - untuk completed status
- **Warning**: Orange (#FF9800) - untuk pending/overdue
- **Error**: Red (#F44336) - untuk error states
- **Info**: Blue (#2196F3) - untuk info badges

#### Background Colors
- **Surface**: White (light mode) / Dark Gray (dark mode)
- **Surface Container**: Light Gray (light mode) / Darker Gray (dark mode)
- **Card Background**: White with subtle shadow

### Typography

- **Header**: Roboto Bold, 20sp
- **Title**: Roboto Medium, 16sp
- **Body**: Roboto Regular, 14sp
- **Caption**: Roboto Regular, 12sp

### Spacing

- **Card Padding**: 16dp
- **Section Spacing**: 24dp
- **Item Spacing**: 12dp
- **Screen Padding**: 16dp

### Animations

- **Page Entrance**: Fade in + Slide up (300ms)
- **Card Entrance**: Stagger animation (100ms delay per card)
- **Pull to Refresh**: Circular progress indicator
- **Loading States**: Skeleton loaders

---

## API Design

### Endpoint Overview

#### 1. Mobile Dashboard Overview
```
GET /api/v1/mobile/dashboard/overview
```

**Description**: Get simplified dashboard data untuk mobile sales user

**Query Parameters**:
- `period` (optional): `today`, `week`, `month` (default: `today`)
- `start_date` (optional): ISO 8601 date string
- `end_date` (optional): ISO 8601 date string

**Response**:
```json
{
  "success": true,
  "data": {
    "target": {
      "target_amount": 20000000,
      "target_amount_formatted": "Rp 20.000.000",
      "achieved_amount": 15000000,
      "achieved_amount_formatted": "Rp 15.000.000",
      "progress_percent": 75.0,
      "period": "2024-01",
      "brick_name": "Jakarta Selatan"
    },
    "visit_summary": {
      "total_today": 5,
      "active": 2,
      "completed": 3,
      "pending_approval": 1
    },
    "task_summary": {
      "total": 8,
      "today": 3,
      "overdue": 1,
      "upcoming": 4
    }
  },
  "error": null
}
```

#### 2. Mobile Visits List
```
GET /api/v1/mobile/dashboard/visits
```

**Description**: Get visits list untuk logged-in user (active dan completed)

**Query Parameters**:
- `status` (optional): `active`, `completed`, `all` (default: `all`)
- `limit` (optional): Items per page (default: 5, max: 5 untuk dashboard)
- `date` (optional): Filter by date (ISO 8601)

**Note**: Untuk dashboard, limit maksimal adalah 5 items. Jika perlu lebih banyak, gunakan endpoint `/api/v1/mobile/visit-reports/my-visit-reports` dengan pagination.

**Response**:
```json
{
  "success": true,
  "data": {
    "visits": [
      {
        "id": "visit-123",
        "account_id": "acc-456",
        "account_name": "RS Pondok Indah",
        "account_address": "Jl. Metro Duta Kav. UE",
        "visit_date": "2024-01-15",
        "visit_time": "09:00",
        "status": "in_progress",
        "check_in_time": "2024-01-15T09:00:00Z",
        "check_in_location": {
          "latitude": -6.2654,
          "longitude": 106.7833,
          "address": "Jl. Metro Duta Kav. UE"
        },
        "check_out_time": null,
        "check_out_location": null,
        "created_at": "2024-01-15T08:30:00Z",
        "updated_at": "2024-01-15T09:00:00Z"
      }
    ],
    "total": 25,
    "has_more": true
  },
  "error": null
}
```

#### 3. Mobile Upcoming Tasks
```
GET /api/v1/mobile/dashboard/tasks
```

**Description**: Get upcoming tasks untuk logged-in user

**Query Parameters**:
- `status` (optional): `pending`, `in_progress`, `completed`, `all` (default: `pending`)
- `filter` (optional): `today`, `week`, `overdue`, `all` (default: `all`)
- `limit` (optional): Items per page (default: 3, max: 3 untuk dashboard)

**Note**: Untuk dashboard, limit maksimal adalah 3 items. Jika perlu lebih banyak, gunakan endpoint `/api/v1/mobile/tasks/my-tasks` dengan pagination.

**Response**:
```json
{
  "success": true,
  "data": {
    "tasks": [
      {
        "id": "task-789",
        "title": "Follow up dengan RS Pondok Indah",
        "description": "Follow up proposal yang sudah dikirim",
        "due_date": "2024-01-16",
        "due_time": "14:00",
        "priority": "high",
        "status": "pending",
        "assigned_by": {
          "id": "user-123",
          "name": "Manager Sales"
        },
        "created_at": "2024-01-15T10:00:00Z",
        "is_overdue": false
      }
    ],
    "total": 15,
    "has_more": true
  },
  "error": null
}
```

### Error Responses

Semua endpoint mengikuti standard error response format:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message",
    "details": {}
  }
}
```

### Error Codes

- `DASHBOARD_FETCH_FAILED`: Gagal mengambil data dashboard
- `VISITS_FETCH_FAILED`: Gagal mengambil data visits
- `TASKS_FETCH_FAILED`: Gagal mengambil data tasks
- `TARGET_NOT_FOUND`: Target tidak ditemukan untuk user
- `INVALID_PERIOD`: Period parameter tidak valid
- `UNAUTHORIZED`: User tidak terautentikasi

---

## Data Models

### Mobile Dashboard Overview Model

```dart
class MobileDashboardOverview {
  final TargetSummary target;
  final VisitSummary visitSummary;
  final TaskSummary taskSummary;

  MobileDashboardOverview({
    required this.target,
    required this.visitSummary,
    required this.taskSummary,
  });

  factory MobileDashboardOverview.fromJson(Map<String, dynamic> json) {
    return MobileDashboardOverview(
      target: TargetSummary.fromJson(json['target']),
      visitSummary: VisitSummary.fromJson(json['visit_summary']),
      taskSummary: TaskSummary.fromJson(json['task_summary']),
    );
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
      targetAmount: json['target_amount'] as int,
      targetAmountFormatted: json['target_amount_formatted'] as String,
      achievedAmount: json['achieved_amount'] as int,
      achievedAmountFormatted: json['achieved_amount_formatted'] as String,
      progressPercent: (json['progress_percent'] as num).toDouble(),
      period: json['period'] as String,
      brickName: json['brick_name'] as String?,
    );
  }
}

class VisitSummary {
  final int totalToday;
  final int active;
  final int completed;
  final int pendingApproval;

  VisitSummary({
    required this.totalToday,
    required this.active,
    required this.completed,
    required this.pendingApproval,
  });

  factory VisitSummary.fromJson(Map<String, dynamic> json) {
    return VisitSummary(
      totalToday: json['total_today'] as int,
      active: json['active'] as int,
      completed: json['completed'] as int,
      pendingApproval: json['pending_approval'] as int,
    );
  }
}

class TaskSummary {
  final int total;
  final int today;
  final int overdue;
  final int upcoming;

  TaskSummary({
    required this.total,
    required this.today,
    required this.overdue,
    required this.upcoming,
  });

  factory TaskSummary.fromJson(Map<String, dynamic> json) {
    return TaskSummary(
      total: json['total'] as int,
      today: json['today'] as int,
      overdue: json['overdue'] as int,
      upcoming: json['upcoming'] as int,
    );
  }
}
```

### Mobile Visit Model

```dart
class MobileVisit {
  final String id;
  final String accountId;
  final String accountName;
  final String? accountAddress;
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
    required this.accountId,
    required this.accountName,
    this.accountAddress,
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
      accountId: json['account_id'] as String,
      accountName: json['account_name'] as String,
      accountAddress: json['account_address'] as String?,
      visitDate: json['visit_date'] as String,
      visitTime: json['visit_time'] as String?,
      status: json['status'] as String,
      checkInTime: json['check_in_time'] != null
          ? DateTime.parse(json['check_in_time'] as String)
          : null,
      checkInLocation: json['check_in_location'] != null
          ? VisitLocation.fromJson(json['check_in_location'] as Map<String, dynamic>)
          : null,
      checkOutTime: json['check_out_time'] != null
          ? DateTime.parse(json['check_out_time'] as String)
          : null,
      checkOutLocation: json['check_out_location'] != null
          ? VisitLocation.fromJson(json['check_out_location'] as Map<String, dynamic>)
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
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
}
```

### Mobile Task Model

```dart
class MobileTask {
  final String id;
  final String title;
  final String? description;
  final String? dueDate;
  final String? dueTime;
  final String priority;
  final String status;
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
      assignedBy: json['assigned_by'] != null
          ? TaskAssignee.fromJson(json['assigned_by'] as Map<String, dynamic>)
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
      isOverdue: json['is_overdue'] as bool? ?? false,
    );
  }
}

class TaskAssignee {
  final String id;
  final String name;

  TaskAssignee({
    required this.id,
    required this.name,
  });

  factory TaskAssignee.fromJson(Map<String, dynamic> json) {
    return TaskAssignee(
      id: json['id'] as String,
      name: json['name'] as String,
    );
  }
}
```

---

## Implementation Plan

### Phase 1: Backend API Development (Week 1)

#### 1.1 Create Mobile Dashboard Routes
- [ ] Add routes di `apps/api/internal/api/routes/dashboard_routes.go`
- [ ] Create mobile-specific route group: `/mobile/dashboard`

#### 1.2 Create Mobile Dashboard Handler
- [ ] Create `GetMobileOverview` handler
- [ ] Create `GetMobileVisits` handler
- [ ] Create `GetMobileTasks` handler
- [ ] Implement error handling sesuai standard

#### 1.3 Create Mobile Dashboard Service
- [ ] Create `GetMobileOverview` service method
  - Get user target dari monthly_targets table
  - Calculate progress
  - Get visit summary stats
  - Get task summary stats
- [ ] Create `GetMobileVisits` service method
  - Filter visits by user_id
  - Support status filter (active/completed)
  - Limit maksimal 5 items (untuk horizontal scroll)
  - Return `has_more` flag jika ada lebih dari 5 items
- [ ] Create `GetMobileTasks` service method
  - Filter tasks by assigned_to user_id
  - Support status dan filter (today/week/overdue)
  - Limit maksimal 3 items (untuk fixed list)
  - Return `has_more` flag jika ada lebih dari 3 items

#### 1.4 Update Repository (if needed)
- [ ] Check if existing repositories support mobile queries
- [ ] Add methods if needed:
  - `GetUserTarget(userID, period)`
  - `GetUserVisits(userID, filters, limit)` - limit maksimal 5
  - `GetUserTasks(userID, filters, limit)` - limit maksimal 3

#### 1.5 Testing
- [ ] Unit tests untuk service methods
- [ ] Integration tests untuk API endpoints
- [ ] Test dengan Postman collection

### Phase 2: Frontend Data Layer (Week 1-2)

#### 2.1 Create Data Models
- [ ] Create `MobileDashboardOverview` model
- [ ] Create `TargetSummary` model
- [ ] Create `VisitSummary` model
- [ ] Create `TaskSummary` model
- [ ] Create `MobileVisit` model
- [ ] Create `MobileTask` model
- [ ] Create `VisitLocation` model
- [ ] Create `TaskAssignee` model

#### 2.2 Create Repository
- [ ] Create `MobileDashboardRepository`
  - `getOverview(period)`
  - `getVisits(status, limit)` - limit maksimal 5
  - `getTasks(status, filter, limit)` - limit maksimal 3
- [ ] Implement offline caching
- [ ] Implement error handling

#### 2.3 Create Provider/State Management
- [ ] Create `MobileDashboardProvider` (Riverpod)
- [ ] Create `MobileDashboardState`
- [ ] Implement loading states
- [ ] Implement error states
- [ ] Implement refresh logic

### Phase 3: UI Components (Week 2)

#### 3.1 Create Quick Stats Widget
- [ ] Create `TargetProgressCard` widget
- [ ] Create `QuickMetricsCard` widget (optional)
- [ ] Implement progress bar dengan color coding
- [ ] Implement tap to expand detail

#### 3.2 Create Visits Section Widget
- [ ] Create `VisitsSection` widget
- [ ] Create `VisitTabBar` widget (Active/Completed)
- [ ] Create `VisitCard` widget (fixed width: 280dp, height: 180dp)
- [ ] Implement horizontal scrollable list (ListView horizontal)
- [ ] Implement status badges
- [ ] Implement action buttons
- [ ] Implement pull to refresh
- [ ] Implement pagination indicator (dots) jika `has_more = true`
- [ ] Implement "View All" button di header jika `has_more = true`
- [ ] Limit maksimal 5 items per tab

#### 3.3 Create Tasks Section Widget
- [ ] Create `TasksSection` widget dengan header "Upcoming Tasks" + "View All" button
- [ ] Create `TaskCard` widget
- [ ] Implement priority indicators
- [ ] Implement overdue detection
- [ ] Implement fixed list (tidak scrollable, maksimal 3 items)
- [ ] Implement "View All" button untuk navigate ke TaskListScreen
- [ ] Implement empty state jika tidak ada tasks
- [ ] Limit maksimal 3 items (sorted by priority & due date)

#### 3.4 Create Dashboard Screen
- [ ] Create `MobileDashboardScreen` widget
- [ ] Integrate semua sections
- [ ] Implement pull to refresh untuk seluruh screen
- [ ] Implement loading states
- [ ] Implement error states
- [ ] Implement empty states

### Phase 4: Integration & Polish (Week 2-3)

#### 4.1 Navigation Integration
- [ ] Update routing untuk dashboard screen baru
- [ ] Update bottom navigation
- [ ] Test navigation flow

#### 4.2 Offline Support
- [ ] Test offline caching
- [ ] Test sync saat online kembali
- [ ] Handle offline error messages

#### 4.3 Performance Optimization
- [ ] Implement lazy loading untuk lists
- [ ] Optimize image loading
- [ ] Implement limit maksimal untuk visits (5) dan tasks (3)
- [ ] Test dengan large datasets

#### 4.4 Localization
- [ ] Add localization keys untuk semua strings
- [ ] Test dengan multiple languages
- [ ] Ensure date/time formatting sesuai locale

#### 4.5 Dark Mode Support
- [ ] Test semua components di dark mode
- [ ] Ensure color contrast sesuai accessibility
- [ ] Fix any dark mode issues

### Phase 5: Testing & Documentation (Week 3)

#### 5.1 Unit Testing
- [ ] Test data models
- [ ] Test repository methods
- [ ] Test provider logic

#### 5.2 Widget Testing
- [ ] Test individual widgets
- [ ] Test widget interactions
- [ ] Test error states

#### 5.3 Integration Testing
- [ ] Test full user flow
- [ ] Test API integration
- [ ] Test offline scenarios

#### 5.4 Documentation
- [ ] Update README dengan new dashboard info
- [ ] Document API endpoints
- [ ] Create user guide (optional)

---

## Testing Strategy

### Unit Tests

#### Backend
- Test service methods dengan mock repositories
- Test handler methods dengan mock services
- Test error handling scenarios

#### Frontend
- Test data model parsing
- Test repository methods
- Test provider state management

### Integration Tests

#### Backend
- Test API endpoints dengan real database
- Test authentication/authorization
- Test limit maksimal (5 untuk visits, 3 untuk tasks)
- Test filtering
- Test `has_more` flag

#### Frontend
- Test API integration
- Test offline caching
- Test state management flow

### E2E Tests

- Test complete user flow:
  1. User opens dashboard
  2. Dashboard loads data
  3. User views target progress
  4. User switches visit tab (Active/Completed)
  5. User scrolls visits horizontally
  6. User taps "View All" untuk visits (jika has_more = true)
  7. User views tasks (maksimal 3 items)
  8. User taps "View All" untuk navigate ke tasks list
  9. User completes task
  10. Dashboard refreshes

### Performance Tests

- Test dengan large datasets (1000+ visits, 500+ tasks)
- Test performance dengan limit maksimal
- Test offline sync performance
- Test memory usage

### Accessibility Tests

- Test screen reader compatibility
- Test color contrast
- Test touch target sizes
- Test keyboard navigation (if applicable)

---

## Migration Strategy

### From Old Dashboard

1. **Keep Old Dashboard**: Jangan hapus old dashboard dulu, biarkan sebagai fallback
2. **Feature Flag**: Gunakan feature flag untuk switch antara old dan new dashboard
3. **Gradual Rollout**: 
   - Week 1: Internal testing
   - Week 2: Beta testing dengan beberapa users
   - Week 3: Full rollout
4. **Monitor**: Monitor error rates dan user feedback
5. **Deprecate**: Setelah stable, deprecate old dashboard

### Data Migration

- Tidak ada data migration yang diperlukan
- API baru menggunakan data yang sama
- Hanya format response yang berbeda

---

## Success Metrics

### Performance Metrics
- Dashboard load time < 2 seconds
- API response time < 500ms
- Offline data availability > 95%

### User Experience Metrics
- User engagement dengan dashboard
- Task completion rate
- Visit completion rate
- User satisfaction score

### Technical Metrics
- Error rate < 1%
- Crash rate < 0.1%
- API success rate > 99%

---

## Future Enhancements

### Phase 2 Features (Future)
1. **Quick Actions**: Swipe actions untuk quick complete/start
2. **Notifications**: Push notifications untuk overdue tasks
3. **Analytics**: Track user behavior untuk improve UX
4. **Customization**: Allow users to customize dashboard layout
5. **Widgets**: Add more optional widgets (revenue, deals, etc.)

---

## Notes

- Semua API endpoints harus mengikuti standard response format
- Semua error harus menggunakan standard error codes
- Offline support adalah critical untuk sales users di lapangan
- Performance optimization sangat penting untuk mobile experience
- Accessibility harus dipertimbangkan dari awal

---

**Document Version**: 1.0  
**Last Updated**: 2024-01-15  
**Author**: Development Team
