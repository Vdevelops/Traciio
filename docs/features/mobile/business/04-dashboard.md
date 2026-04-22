# Business - Dashboard

## CRM Healthcare Mobile App - Flutter

**Module**: Business Domain  
**Sprint**: Sprint 4  
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

Fitur **Dashboard** menyediakan overview singkat mengenai performance dan activities sales rep. Dashboard menampilkan key metrics seperti jumlah visits, tasks, accounts, dan pipeline summary. Data di-fetch dari backend secara real-time dengan caching untuk performa optimal.

### Goals

- **Overview**: Quick glance at key metrics dan performance
- **Performance Tracking**: Track daily, weekly, monthly performance
- **Actionable Insights**: Highlight items yang memerlukan attention
- **Quick Navigation**: Shortcuts ke features penting
- **Motivation**: Visual representation dari achievements

---

## Fitur Utama

### 1. Overview Stats

**Key Metrics Cards**:

- 📊 **Total Visits**: Jumlah visit reports periode ini
- 📋 **Pending Tasks**: Tasks yang belum selesai
- 🏢 **Active Accounts**: Jumlah accounts yang dikelola
- 💰 **Revenue**: Pipeline value atau actual revenue

**Display Format**:

- Current value dengan comparison (vs previous period)
- Trend indicator (↑ ↓ →)
- Percentage change

### 2. Pipeline Summary

**Visual Pipeline**:

- Stages: Lead → Qualified → Proposal → Negotiation → Closed Won
- Deal count dan value per stage
- Progress bar atau funnel chart
- Conversion rates antar stages

### 3. Recent Activities

**Activity Feed**:

- Recent visit reports
- Recent task completions
- Recent account additions
- Time-stamped dengan relative time ("2 hours ago")

### 4. Upcoming Tasks

**Priority Tasks**:

- Tasks due today
- Overdue tasks
- High priority tasks
- Quick complete action

### 5. Period Selector

**Time Periods**:

- Today
- This Week
- This Month
- This Quarter
- This Year

**Period Comparison**:

- Compare current period vs previous period
- Growth percentage

### 6. Quick Actions

**Floating Action Menu**:

- ➕ Create Visit Report
- ➕ Create Task
- ➕ Quick Account Search

---

## Business Rules

### 1. Data Visibility Rules

**Sales Rep View**:

- Hanya melihat data milik mereka sendiri
- Visits: Visit reports yang mereka create
- Tasks: Tasks yang di-assign ke mereka
- Accounts: Accounts yang di-assign ke mereka
- Pipeline: Deals yang mereka own

**Supervisor View**:

- Melihat aggregate data team-nya
- Filter by team member
- Compare performance antar team members

### 2. Data Freshness Rules

**Real-time Data**:

- Dashboard data di-fetch setiap kali screen opened
- Background refresh setiap 5 menit
- Pull-to-refresh untuk manual refresh

**Caching**:

- Cache data selama 5 menit
- Show cached data sementara fetch baru
- Offline mode: show last cached data

### 3. Metric Calculation Rules

**Total Visits**:

- Count visit reports dengan status approved
- Filter by selected period
- Include check-in time dalam period

**Pending Tasks**:

- Count tasks dengan status pending atau in_progress
- Include overdue tasks
- Filter by assigned_to = current user

**Active Accounts**:

- Count unique accounts dengan activity dalam period
- Atau count accounts yang di-assign ke user

**Revenue**:

- Sum dari deals dengan status closed_won
- Deal value dari closed date dalam period
- Atau pipeline value dari active deals

### 4. Period Comparison Rules

**Current Period**: Selected period (Today, Week, Month, etc.)
**Previous Period**: Same duration sebelumnya
**Comparison**: Percentage change = ((Current - Previous) / Previous) × 100

**Contoh**:

- Today: Jan 20 vs Jan 19
- This Week: Week 3 vs Week 2
- This Month: Jan vs Dec

### 5. Pipeline Stage Rules

**Standard Stages**:

1. **Lead**: Initial contact atau inquiry
2. **Qualified**: Qualified opportunity
3. **Proposal**: Proposal submitted
4. **Negotiation**: Under negotiation
5. **Closed Won**: Deal closed successfully
6. **Closed Lost**: Deal lost

**Calculation**:

- Deal count per stage
- Total value per stage
- Average deal size
- Conversion rate ke stage berikutnya

---

## Keputusan Teknis & Trade-offs

### Mengapa Mobile-Specific Endpoints?

**Keputusan**: Menggunakan mobile-specific endpoints (`/api/v1/dashboard/mobile/*`).

**Alasan**:

- **Performance**: Optimized response size untuk mobile
- **Relevance**: Hanya data yang diperlukan mobile app
- **Battery**: Reduce processing time dan data usage

**Trade-off**: Additional backend work untuk create mobile endpoints. **Mitigasi**: Reuse existing data services dengan mobile-specific DTOs.

### Mengapa Multiple Small Endpoints vs One Big Endpoint?

**Keputusan**: Separate endpoints untuk overview, visits, dan tasks.

**Alasan**:

- **Parallel Loading**: Multiple endpoints dapat di-fetch parallel
- **Partial Failure**: Jika satu endpoint gagal, others tetap berfungsi
- **Caching**: Granular caching per section
- **Progressive Loading**: Show data saat available, tidak perlu wait untuk all

**Trade-off**: Multiple network requests. **Mitigasi**: Parallel fetching dengan Future.wait.

### Mengapa Simplified Charts?

**Keputusan**: Use simple bar charts dan progress bars, bukan complex visualizations.

**Alasan**:

- **Mobile Screen**: Limited space untuk complex charts
- **Performance**: Simple charts render lebih cepat
- **Clarity**: Easier to understand pada glance
- **Library Size**: Smaller charting library

**Trade-off**: Less detailed analytics. **Mitigasi**: Full analytics available di Reports screen.

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── dashboard/
│       ├── data/
│       │   ├── models/
│       │   │   ├── dashboard_overview_model.dart
│       │   │   ├── dashboard_stats_model.dart
│       │   │   └── pipeline_summary_model.dart
│       │   └── dashboard_repository.dart
│       ├── application/
│       │   ├── dashboard_overview_provider.dart
│       │   ├── dashboard_stats_provider.dart
│       │   ├── dashboard_visits_provider.dart
│       │   ├── dashboard_tasks_provider.dart
│       │   └── dashboard_period_provider.dart
│       └── presentation/
│           ├── screens/
│           │   └── dashboard_screen.dart
│           └── widgets/
│               ├── dashboard_overview_section.dart
│               ├── stat_card.dart
│               ├── pipeline_summary_widget.dart
│               ├── recent_activities_list.dart
│               ├── upcoming_tasks_section.dart
│               ├── period_selector.dart
│               └── quick_actions_menu.dart
├── core/
│   └── widgets/
│       ├── charts/
│       │   └── simple_bar_chart.dart
│       └── skeletons/
│           └── dashboard_skeleton.dart
```

---

## API Endpoints

### Mobile Dashboard Endpoints

#### GET /api/v1/dashboard/mobile/overview

Get dashboard overview dengan key metrics.

**Query Parameters**:

```
?period=today&start_date=2025-01-20&end_date=2025-01-20
```

**Response**:

```json
{
  "success": true,
  "data": {
    "period": "today",
    "date_range": {
      "start": "2025-01-20T00:00:00Z",
      "end": "2025-01-20T23:59:59Z"
    },
    "stats": {
      "visits": {
        "current": 5,
        "previous": 3,
        "change_percent": 66.7,
        "trend": "up"
      },
      "tasks": {
        "current": 12,
        "previous": 15,
        "change_percent": -20.0,
        "trend": "down"
      },
      "accounts": {
        "current": 28,
        "previous": 28,
        "change_percent": 0.0,
        "trend": "neutral"
      },
      "revenue": {
        "current": 15000000,
        "previous": 12000000,
        "change_percent": 25.0,
        "trend": "up",
        "currency": "IDR"
      }
    },
    "summary": {
      "pending_tasks": 8,
      "overdue_tasks": 2,
      "upcoming_visits": 3,
      "deals_in_pipeline": 5
    }
  },
  "timestamp": "2025-01-20T10:30:45+07:00"
}
```

#### GET /api/v1/dashboard/mobile/pipeline

Get pipeline summary.

**Response**:

```json
{
  "success": true,
  "data": {
    "total_deals": 15,
    "total_value": 75000000,
    "currency": "IDR",
    "stages": [
      {
        "stage": "lead",
        "name": "Lead",
        "count": 5,
        "value": 25000000,
        "color": "#9E9E9E"
      },
      {
        "stage": "qualified",
        "name": "Qualified",
        "count": 4,
        "value": 20000000,
        "color": "#2196F3"
      },
      {
        "stage": "proposal",
        "name": "Proposal",
        "count": 3,
        "value": 15000000,
        "color": "#FF9800"
      },
      {
        "stage": "negotiation",
        "name": "Negotiation",
        "count": 2,
        "value": 10000000,
        "color": "#9C27B0"
      },
      {
        "stage": "closed_won",
        "name": "Closed Won",
        "count": 1,
        "value": 5000000,
        "color": "#4CAF50"
      }
    ],
    "conversion_rates": {
      "lead_to_qualified": 80.0,
      "qualified_to_proposal": 75.0,
      "proposal_to_negotiation": 66.7,
      "negotiation_to_closed": 50.0
    }
  }
}
```

#### GET /api/v1/dashboard/mobile/visits

Get recent visits.

**Query Parameters**:

```
?limit=5
```

**Response**:

```json
{
  "success": true,
  "data": {
    "visits": [
      {
        "id": "vr-uuid",
        "account_name": "RS Medika Hospital",
        "account_id": "account-uuid",
        "visit_date": "2025-01-20T09:00:00Z",
        "status": "approved",
        "has_photos": true
      }
    ]
  }
}
```

#### GET /api/v1/dashboard/mobile/tasks

Get upcoming tasks.

**Query Parameters**:

```
?limit=5&status=pending,in_progress&sort=due_date&order=asc
```

**Response**:

```json
{
  "success": true,
  "data": {
    "tasks": [
      {
        "id": "task-uuid",
        "title": "Follow up with Dr. Smith",
        "due_date": "2025-01-20T14:00:00Z",
        "priority": "high",
        "status": "pending",
        "related_account": "RS Medika Hospital"
      }
    ]
  }
}
```

---

## Data Models

### Dashboard Overview Model

```dart
@freezed
class DashboardOverview with _$DashboardOverview {
  const factory DashboardOverview({
    required String period,
    required DateRange dateRange,
    required DashboardStats stats,
    required DashboardSummary summary,
  }) = _DashboardOverview;

  factory DashboardOverview.fromJson(Map<String, dynamic> json) =
      _$DashboardOverviewFromJson(json);
}

@freezed
class DateRange with _$DateRange {
  const factory DateRange({
    required DateTime start,
    required DateTime end,
  }) = _DateRange;

  factory DateRange.fromJson(Map<String, dynamic> json) =
      _$DateRangeFromJson(json);
}

@freezed
class DashboardStats with _$DashboardStats {
  const factory DashboardStats({
    required StatMetric visits,
    required StatMetric tasks,
    required StatMetric accounts,
    required StatMetric revenue,
  }) = _DashboardStats;

  factory DashboardStats.fromJson(Map<String, dynamic> json) =
      _$DashboardStatsFromJson(json);
}

@freezed
class StatMetric with _$StatMetric {
  const factory StatMetric({
    required num current,
    required num previous,
    required double changePercent,
    required String trend,
    String? currency,
  }) = _StatMetric;

  factory StatMetric.fromJson(Map<String, dynamic> json) =
      _$StatMetricFromJson(json);

  String get formattedChange {
    final prefix = changePercent >= 0 ? '+' : '';
    return '$prefix${changePercent.toStringAsFixed(1)}%';
  }

  bool get isPositive => changePercent >= 0;
}

@freezed
class DashboardSummary with _$DashboardSummary {
  const factory DashboardSummary({
    required int pendingTasks,
    required int overdueTasks,
    required int upcomingVisits,
    required int dealsInPipeline,
  }) = _DashboardSummary;

  factory DashboardSummary.fromJson(Map<String, dynamic> json) =
      _$DashboardSummaryFromJson(json);
}
```

### Pipeline Summary Model

```dart
@freezed
class PipelineSummary with _$PipelineSummary {
  const factory PipelineSummary({
    required int totalDeals,
    required double totalValue,
    required String currency,
    required List<PipelineStage> stages,
    required Map<String, double> conversionRates,
  }) = _PipelineSummary;

  factory PipelineSummary.fromJson(Map<String, dynamic> json) =
      _$PipelineSummaryFromJson(json);
}

@freezed
class PipelineStage with _$PipelineStage {
  const factory PipelineStage({
    required String stage,
    required String name,
    required int count,
    required double value,
    required String color,
  }) = _PipelineStage;

  factory PipelineStage.fromJson(Map<String, dynamic> json) =
      _$PipelineStageFromJson(json);

  double get percentageOfTotal => count / (parent?.totalDeals ?? 1) * 100;
}
```

---

## Configuration

### Dashboard Repository

**File**: `features/dashboard/data/dashboard_repository.dart`

```dart
class DashboardRepository {
  final ApiClient _apiClient;
  final ConnectivityService _connectivity;
  final CacheManager _cache;

  DashboardRepository(
    this._apiClient,
    this._connectivity,
    this._cache,
  );

  Future<DashboardOverview> getOverview({
    required String period,
    DateTime? startDate,
    DateTime? endDate,
    bool forceRefresh = false,
  }) async {
    final cacheKey = 'dashboard_overview_$period';

    // Check cache
    if (!forceRefresh) {
      final cached = await _cache.get<DashboardOverview>(cacheKey);
      if (cached != null) {
        // Return cached, refresh in background
        if (_connectivity.isOnline) {
          _refreshInBackground(cacheKey, period, startDate, endDate);
        }
        return cached;
      }
    }

    // Fetch dari API
    if (_connectivity.isOnline) {
      final response = await _apiClient.get(
        '/api/v1/dashboard/mobile/overview',
        queryParameters: {
          'period': period,
          if (startDate != null)
            'start_date': startDate.toIso8601String(),
          if (endDate != null)
            'end_date': endDate.toIso8601String(),
        },
      );

      final overview = DashboardOverview.fromJson(response.data['data']);

      // Cache untuk 5 menit
      await _cache.set(cacheKey, overview, ttl: const Duration(minutes: 5));

      return overview;
    }

    throw Exception('No internet connection');
  }

  Future<PipelineSummary> getPipelineSummary() async {
    final response = await _apiClient.get(
      '/api/v1/dashboard/mobile/pipeline',
    );
    return PipelineSummary.fromJson(response.data['data']);
  }

  Future<List<DashboardVisit>> getRecentVisits({int limit = 5}) async {
    final response = await _apiClient.get(
      '/api/v1/dashboard/mobile/visits',
      queryParameters: {'limit': limit},
    );
    return (response.data['data']['visits'] as List)
        .map((json) => DashboardVisit.fromJson(json))
        .toList();
  }

  Future<List<DashboardTask>> getUpcomingTasks({int limit = 5}) async {
    final response = await _apiClient.get(
      '/api/v1/dashboard/mobile/tasks',
      queryParameters: {'limit': limit},
    );
    return (response.data['data']['tasks'] as List)
        .map((json) => DashboardTask.fromJson(json))
        .toList();
  }
}
```

---

## Usage Examples

### Dashboard Screen

```dart
class DashboardScreen extends ConsumerWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final period = ref.watch(dashboardPeriodProvider);
    final overviewAsync = ref.watch(dashboardOverviewProvider(period));

    return Scaffold(
      body: RefreshIndicator(
        onRefresh: () => ref.refresh(dashboardOverviewProvider(period).future),
        child: CustomScrollView(
          slivers: [
            // App Bar dengan period selector
            SliverAppBar(
              floating: true,
              title: const Text('Dashboard'),
              actions: [
                PeriodSelector(
                  selectedPeriod: period,
                  onPeriodChanged: (newPeriod) {
                    ref.read(dashboardPeriodProvider.notifier).state = newPeriod;
                  },
                ),
              ],
            ),

            // Dashboard Content
            SliverPadding(
              padding: const EdgeInsets.all(16),
              sliver: overviewAsync.when(
                loading: () => const SliverToBoxAdapter(
                  child: DashboardSkeleton(),
                ),
                error: (error, _) => SliverToBoxAdapter(
                  child: ErrorWidget(
                    error: error.toString(),
                    onRetry: () => ref.refresh(
                      dashboardOverviewProvider(period).future,
                    ),
                  ),
                ),
                data: (overview) => SliverList(
                  delegate: SliverChildListDelegate([
                    // Stats Grid
                    StatsGrid(stats: overview.stats),

                    const SizedBox(height: 24),

                    // Summary Chips
                    SummaryChips(summary: overview.summary),

                    const SizedBox(height: 24),

                    // Pipeline Summary
                    const PipelineSummaryWidget(),

                    const SizedBox(height: 24),

                    // Upcoming Tasks
                    const UpcomingTasksSection(),

                    const SizedBox(height: 24),

                    // Recent Visits
                    const RecentVisitsSection(),
                  ]),
                ),
              ),
            ),
          ],
        ),
      ),
      floatingActionButton: QuickActionsMenu(
        onCreateVisit: () => context.push(AppRoutes.visitReportCreate),
        onCreateTask: () => context.push(AppRoutes.taskCreate),
      ),
    );
  }
}
```

### Stat Card

```dart
class StatCard extends StatelessWidget {
  final String title;
  final StatMetric metric;
  final IconData icon;
  final Color color;

  const StatCard({
    super.key,
    required this.title,
    required this.metric,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(icon, color: color),
                const SizedBox(width: 8),
                Text(title, style: Theme.of(context).textTheme.titleSmall),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              _formatValue(metric.current, metric.currency),
              style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 4),
            Row(
              children: [
                Icon(
                  metric.isPositive ? Icons.arrow_upward : Icons.arrow_downward,
                  size: 16,
                  color: metric.isPositive ? Colors.green : Colors.red,
                ),
                const SizedBox(width: 4),
                Text(
                  metric.formattedChange,
                  style: TextStyle(
                    color: metric.isPositive ? Colors.green : Colors.red,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(width: 4),
                Text(
                  'vs last period',
                  style: Theme.of(context).textTheme.bodySmall,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  String _formatValue(num value, String? currency) {
    if (currency != null) {
      return NumberFormat.currency(
        symbol: currency == 'IDR' ? 'Rp ' : '\$',
        decimalDigits: 0,
      ).format(value);
    }
    return NumberFormat.compact().format(value);
  }
}
```

---

## Cara Test Manual

### Test Dashboard Loading

1. **Initial Load**:
   - Open dashboard
   - Verifikasi: Skeleton loading muncul
   - Verifikasi: Data loaded setelah beberapa detik
   - Verifikasi: All sections populated

2. **Period Switch**:
   - Tap period selector
   - Select "This Week"
   - Verifikasi: Data refresh untuk week view
   - Verifikasi: Stats change sesuai period

3. **Pull-to-Refresh**:
   - Pull down
   - Verifikasi: Refresh indicator muncul
   - Verifikasi: Data refreshed

### Test Offline Mode

1. **Offline Dashboard**:
   - Matikan internet
   - Open dashboard
   - Verifikasi: Cached data ditampilkan
   - Verifikasi: Offline indicator muncul

2. **Background Refresh**:
   - Online, load dashboard
   - Matikan internet
   - Tunggu beberapa menit
   - Turn on internet
   - Verifikasi: Auto-refresh terjadi

---

## Dependencies

### Internal

- `features/accounts/data/account_repository.dart` - Account data
- `features/tasks/data/task_repository.dart` - Task data
- `core/network/api_client.dart` - API calls
- `core/cache/cache_manager.dart` - Data caching

### External

- `fl_chart: ^0.66.0` - Charts (optional)
- `intl: ^0.18.0` - Number formatting
- `flutter_riverpod: ^2.4.0` - State management
- `skeletonizer: ^1.0.0` - Skeleton loading

---

## Notes & Improvements

### Known Limitations

1. **Static Data**: Dashboard data static, tidak ada real-time updates.

2. **Limited Customization**: User tidak dapat customize which metrics ditampilkan.

3. **No Drill-down**: Tidak bisa tap metric untuk melihat detail.

4. **No Goals**: Tidak ada goals/targets untuk comparison.

### Future Improvements

1. **Real-time Updates**: WebSocket atau polling untuk real-time dashboard updates

2. **Customizable Dashboard**: User dapat choose which widgets to display

3. **Drill-down Navigation**: Tap metric untuk navigate ke detail screen

4. **Goals & Targets**: Set daily/weekly/monthly targets dan track progress

5. **Performance Charts**: Historical performance charts

6. **Leaderboard**: Team performance ranking

7. **Predictive Analytics**: AI-powered sales predictions

---

**Document Status**: Active  
**Last Updated**: January 2025  
**Maintained By**: Dev3 (Mobile Development Team)
