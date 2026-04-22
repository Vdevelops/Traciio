# Sprint Planning - Developer 4 (Fullstack Developer)

## CRM Healthcare/Pharmaceutical Platform - Analytics & Performance

**Developer**: Fullstack Developer (Go Backend + Next.js Frontend)  
**Role**: Develop modul-modul Analytics & Performance secara fullstack (backend + frontend)  
**Versi**: 2.0  
**Status**: Active  
**Last Updated**: 2025-01-15

---

## 📋 Overview

Developer 4 bertanggung jawab untuk:

- **Fullstack Development**: Develop modul-modul Analytics & Performance secara lengkap (backend API + frontend)
- **Backend**: Go (Gin) APIs untuk modul yang ditugaskan
- **Frontend**: Next.js 16 frontend untuk modul yang ditugaskan
- **Database**: Design dan implement database schema untuk modul yang ditugaskan
- **Postman Collection**: Update Postman collection untuk modul yang ditugaskan

**Modul yang ditugaskan ke Dev4**:

1. Leaderboard & KPI Monitoring (Fullstack) - General ranking, KPI tracking, achievements
2. Sales Performance & Incentive Management (Fullstack) - Performance detail per user, incentive calculation & management
3. Time-Based Sales Analytics (Fullstack) - Sales by month/year, trends, forecast
4. Top Products Analytics (Fullstack) - Top products visualization, product performance

**Parallel Development Strategy**:

- ✅ **TIDAK bergantung ke Dev5** - bisa dikerjakan paralel
- ✅ Setiap modul dikerjakan fullstack sampai selesai
- ✅ **Hackathon mode** - tidak ada unit test
- ✅ Manual testing saja
- ✅ Menggunakan data dari modul yang sudah ada (Deals, Visit Reports, Tasks, Products)

---

## 🎯 Sprint Details

### Sprint 1: Leaderboard & KPI Monitoring (Fullstack)

**Goal**: Implement Leaderboard & KPI Monitoring secara fullstack (backend + frontend)

**Backend Tasks**:

- [x] Create leaderboard model dan migration
- [x] Create achievement model dan migration
- [x] Create kpi_settings model dan migration
- [x] Create leaderboard repository interface dan implementation
- [x] Create achievement repository interface dan implementation
- [x] Create leaderboard service
- [x] Create achievement service
- [x] Implement leaderboard list API (`GET /api/v1/leaderboard`)
  - Support filters: period (daily, weekly, monthly, yearly), metric (deals, revenue, visits, tasks), limit
  - Calculate rankings based on selected metric
- [x] Implement KPI API (`GET /api/v1/leaderboard/kpi/:userId`)
  - Calculate all KPI metrics for user
  - Support period filtering
- [x] Implement achievements API (`GET /api/v1/leaderboard/achievements/:userId`)
- [x] Implement performance trends API (`GET /api/v1/leaderboard/trends/:userId`)
- [x] Add KPI calculation logic:
  - Total deals closed (from deals table where status='won' and assigned_to=userId)
  - Total revenue generated (sum of deal values where status='won')
  - Number of visits completed (from visit_reports where sales_rep_id=userId and status='approved')
  - Tasks completed (from tasks where assigned_to=userId and status='completed')
  - Conversion rate (won deals / total deals)
  - Average deal value
- [x] Add achievement detection logic (first deal, top seller milestone, etc.)
- [x] Add pagination support
- [x] Add validation
- [x] Add leaderboard and achievement seeders

**Frontend Tasks**:

- [x] Create leaderboard types (`types/leaderboard.d.ts`)
- [x] Create achievement types (`types/achievement.d.ts`)
- [x] Create KPI types (`types/kpi.d.ts`)
- [x] Create leaderboard service (`leaderboardService`)
- [x] Create achievement service (`achievementService`)
- [x] Create leaderboard list page (`/leaderboard`)
- [x] Create leaderboard table component (`LeaderboardTable`)
- [x] Create KPI card component (`KPICard`)
- [x] Create achievement badge component (`AchievementBadge`)
- [x] Create performance trend chart component (`PerformanceTrendChart`)
- [x] Create leaderboard filters component (`LeaderboardFilters`)
- [x] Add period filter (daily, weekly, monthly, yearly)
- [x] Add metric filter (deals, revenue, visits, tasks)
- [x] Add ranking display with badges (1st, 2nd, 3rd)
- [x] Add KPI detail view for user
- [x] Add achievement display
- [x] Add performance trend visualization (line chart)
- [ ] Add export functionality (Excel/PDF) - Optional, can be added later

**Postman Collection**:

- [x] Add leaderboard APIs ke Postman collection (Web section)

**Menu & Permissions**:

- [x] Add Leaderboard menu to menu seeder
- [x] Add Leaderboard permissions to permission seeder

**Acceptance Criteria**:

- ✅ Leaderboard APIs bekerja dengan baik
- ✅ KPI calculation accurate
- ✅ Achievement system bekerja
- ✅ Performance trends ditampilkan dengan benar
- ✅ Frontend terintegrasi dengan backend APIs
- ✅ Filters bekerja optimal
- ✅ Charts dan visualizations ditampilkan dengan benar
- ✅ Postman collection updated

**Testing** (Manual testing):

- Test leaderboard calculation (backend + frontend)
- Test KPI metrics accuracy
- Test achievement detection
- Test performance trends
- Test filters and period selection

**Estimated Time**: 6-7 days

**Database Schema**:

```sql
-- Leaderboard entries (materialized view for performance)
CREATE TABLE leaderboard_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    period_type VARCHAR(20) NOT NULL, -- daily, weekly, monthly, yearly
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    metric_type VARCHAR(50) NOT NULL, -- deals, revenue, visits, tasks
    metric_value DECIMAL(15,2) NOT NULL,
    rank INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, period_type, period_start, metric_type),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_leaderboard_entries_period ON leaderboard_entries(period_type, period_start, metric_type);
CREATE INDEX idx_leaderboard_entries_user_id ON leaderboard_entries(user_id);

-- Achievements
CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    achievement_type VARCHAR(50) NOT NULL, -- first_deal, top_seller, milestone_100, etc.
    achievement_name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(100), -- icon name or URL
    achieved_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB, -- Additional data (e.g., deal_id, revenue_amount)
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_achievements_user_id ON achievements(user_id);
CREATE INDEX idx_achievements_type ON achievements(achievement_type);

-- KPI settings
CREATE TABLE kpi_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    metric_name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    weight DECIMAL(5,2) DEFAULT 1.0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**Relasi dengan Table Existing**:

- `leaderboard_entries.user_id` → `users.id`
- KPI calculation menggunakan:
  - `deals` table (status='won', assigned_to=userId)
  - `visit_reports` table (status='approved', sales_rep_id=userId)
  - `tasks` table (status='completed', assigned_to=userId)

---

### Sprint 2: Sales Performance & Incentive Management (Fullstack)

**Goal**: Implement Sales Performance Detail & Incentive Management secara fullstack (backend + frontend)

**Note**: Top Sales ranking sudah ada di Leaderboard (bisa filter `metric=revenue, limit=10`). Sprint ini fokus pada **Sales Performance Detail per User** dan **Incentive Management**.

**Backend Tasks**:

- [x] Create incentive_tiers model dan migration
- [x] Create incentives model dan migration
- [x] Create incentive_payouts model dan migration
- [x] Create incentive repository interface dan implementation
- [x] Create sales overview repository interface dan implementation
- [x] Create incentive service
- [x] Create sales overview service
- [x] Implement sales performance detail API (`GET /api/v1/sales-overview/performance/:userId`)
  - Get detailed performance metrics for specific user
  - Support period filtering
  - Calculate: revenue breakdown, deals breakdown, conversion rate, avg deal value, visit stats, task completion
  - Return detailed metrics (not ranking)
- [x] Implement sales rep detail page API (`GET /api/v1/sales-overview/sales-rep/:userId`)
  - Get comprehensive detail for sales rep detail page
  - Include: statistics, activities, check-in locations
  - Support period filtering
- [x] Implement sales rep activities API (`GET /api/v1/sales-overview/sales-rep/:userId/activities`)
  - Get activities timeline for sales rep
  - Include: visit reports, deals, tasks
  - Support pagination and date filtering
- [x] Implement sales rep check-in locations API (`GET /api/v1/sales-overview/sales-rep/:userId/check-in-locations`)
  - Get all check-in locations for sales rep
  - Ordered by visit number (1, 2, 3, ...)
  - Include: visit number, location (lat, lng, address), visit date, account name
  - Support period filtering
- [x] Implement sales performance list API (`GET /api/v1/sales-overview/performance`)
  - List all sales reps with performance summary (for management view)
  - Support filters: period, search by name
  - Not a ranking, but a management overview
  - Make sales rep name clickable (link to detail page)
- [x] Implement incentive tiers API (`GET /api/v1/incentives/tiers`)
- [x] Implement incentive tiers management API (`POST /api/v1/incentives/tiers`, `PUT /api/v1/incentives/tiers/:id`, `DELETE /api/v1/incentives/tiers/:id`)
- [x] Implement incentive list API (`GET /api/v1/incentives`)
  - Support filters: user_id, period, status
- [x] Implement incentive detail API (`GET /api/v1/incentives/:id`)
- [x] Implement incentive calculation API (`POST /api/v1/incentives/calculate`)
  - Calculate incentive based on revenue, deals, tier
  - Support period (monthly, quarterly, yearly)
  - Auto-determine tier based on revenue threshold
- [x] Implement incentive payout API (`POST /api/v1/incentives/:id/payout`)
- [x] Implement incentive approval API (`POST /api/v1/incentives/:id/approve`)
- [x] Add incentive calculation logic:
  - Determine tier based on revenue threshold
  - Calculate: (revenue * tier.incentive_percent / 100) + tier.bonus_amount
  - Apply milestone bonuses
- [x] Add pagination support
- [x] Add validation
- [x] Add incentive tiers seeder
- [x] Create handlers for all APIs
- [x] Register routes for all APIs
- [x] Integrate services and handlers in main.go

**Frontend Tasks**:

- [x] Create sales overview types (`types/sales-overview.d.ts`)
- [x] Create incentive types (`types/incentive.d.ts`)
- [x] Create sales overview service (`salesOverviewService`)
- [x] Create incentive service (`incentiveService`)
- [x] Create sales overview page (`/sales-overview`)
- [x] Create sales performance list component (`SalesPerformanceTable`)
  - Management overview (all sales reps)
  - Not ranking, but summary table
  - Make sales rep name clickable → navigate to detail page
- [ ] Create sales rep detail page (`/sales-overview/sales-rep/:userId`)
  - Sales rep information header
  - Statistics cards section:
    - Total revenue, deals closed, visits completed, tasks completed
    - Conversion rate, avg deal value
    - Period comparison (vs previous period)
  - Activities timeline section:
    - Timeline/list of activities (visit reports, deals, tasks)
    - Filter by activity type
    - Support pagination
  - Check-in locations map section:
    - Interactive map with numbered markers (1, 2, 3, ...)
    - List of check-in locations with visit number
    - Click marker → show visit details
    - Support period filtering
- [ ] Create sales performance detail view (`SalesPerformanceDetail`)
  - Detailed metrics per user (not ranking)
  - Revenue breakdown chart
  - Deals breakdown chart
  - Conversion rate analysis
  - Visit statistics
  - Task completion stats
- [ ] Create sales performance card component (`SalesPerformanceCard`)
- [ ] Create sales rep statistics component (`SalesRepStatistics`)
  - Statistics cards with metrics
- [ ] Create sales rep activities component (`SalesRepActivities`)
  - Activities timeline/list
  - Activity type filter
- [ ] Create sales rep check-in map component (`SalesRepCheckInMap`)
  - Map with numbered markers (1, 2, 3, ...)
  - Check-in locations list
  - Integration with Leaflet/Mapbox
- [x] Create incentive management page (`/incentives`)
- [x] Create incentive list component (`IncentiveTable`)
- [ ] Create incentive card component (`IncentiveCard`)
- [x] Create incentive tier badge component (`IncentiveTierBadge` - integrated in table)
- [ ] Create incentive tier management component (`IncentiveTierManagement`)
- [ ] Create incentive history table component (`IncentiveHistory`)
- [ ] Create incentive calculation dialog (`IncentiveCalculationDialog`)
- [ ] Create incentive payout dialog (`IncentivePayoutDialog`)
- [x] Add tier visualization (bronze, silver, gold, platinum) - in table
- [ ] Add incentive calculation UI
- [ ] Add incentive payout management
- [x] Add period filtering - basic implemented
- [ ] Add export functionality (Excel/PDF) - Optional, can be added later

**Postman Collection**:

- [x] Add sales overview APIs ke Postman collection (Web section)
- [x] Add incentive APIs ke Postman collection (Web section)

**Menu & Permissions**:

- [x] Add Sales Performance menu to menu seeder
- [x] Add Incentive Management menu to menu seeder
- [x] Add Sales Performance permissions to permission seeder
- [x] Add Incentive Management permissions to permission seeder

**Acceptance Criteria**:

- ✅ Sales performance APIs bekerja dengan baik
- ✅ Sales rep detail page APIs bekerja dengan baik
- ✅ Sales rep activities API bekerja dengan baik
- ✅ Sales rep check-in locations API bekerja dengan baik (ordered by visit number)
- ✅ Sales performance detail accurate (per user, bukan ranking)
- ✅ Sales rep detail page menampilkan:
  - ✅ Statistics cards dengan metrics lengkap
  - ✅ Activities timeline/list dengan filter
  - ✅ Check-in locations map dengan numbered markers (1, 2, 3, ...)
- ✅ Incentive calculation accurate
- ✅ Incentive tier system bekerja
- ✅ Incentive tier management bekerja
- ✅ Incentive payout management bekerja
- ✅ Frontend terintegrasi dengan backend APIs
- ✅ Charts dan visualizations ditampilkan dengan benar
- ✅ Map integration bekerja dengan baik
- ✅ Postman collection updated

**Testing** (Manual testing):

- Test sales performance detail (backend + frontend)
- Test sales rep detail page (backend + frontend)
- Test sales rep activities API
- Test sales rep check-in locations API (verify ordering 1, 2, 3, ...)
- Test check-in locations map display
- Test activities timeline display
- Test statistics calculation accuracy
- Test incentive calculation
- Test tier assignment
- Test incentive tier management
- Test incentive payout flow
- Test incentive approval workflow

**Estimated Time**: 7-8 days (ditambah 1 hari untuk sales rep detail page dengan map integration dan check-in locations numbering)

**Database Schema**:

```sql
-- Incentive tiers
CREATE TABLE incentive_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tier_name VARCHAR(50) NOT NULL, -- bronze, silver, gold, platinum
    display_name VARCHAR(255) NOT NULL,
    min_revenue DECIMAL(15,2) NOT NULL,
    max_revenue DECIMAL(15,2), -- NULL for highest tier
    incentive_percent DECIMAL(5,2) NOT NULL, -- Percentage of revenue
    bonus_amount DECIMAL(15,2) DEFAULT 0, -- Fixed bonus
    icon VARCHAR(100), -- icon name or URL
    color VARCHAR(50), -- badge color
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_incentive_tiers_revenue ON incentive_tiers(min_revenue, max_revenue);

-- Incentives
CREATE TABLE incentives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    period_type VARCHAR(20) NOT NULL, -- monthly, quarterly, yearly
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    revenue DECIMAL(15,2) NOT NULL,
    deals_closed INTEGER NOT NULL,
    conversion_rate DECIMAL(5,2), -- Percentage
    tier_id UUID,
    incentive_amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, paid
    calculated_at TIMESTAMP DEFAULT NOW(),
    approved_at TIMESTAMP,
    approved_by UUID,
    paid_at TIMESTAMP,
    notes TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tier_id) REFERENCES incentive_tiers(id),
    FOREIGN KEY (approved_by) REFERENCES users(id)
);

CREATE INDEX idx_incentives_user_id ON incentives(user_id);
CREATE INDEX idx_incentives_period ON incentives(period_start, period_end);
CREATE INDEX idx_incentives_status ON incentives(status);

-- Incentive payouts
CREATE TABLE incentive_payouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incentive_id UUID NOT NULL,
    payout_amount DECIMAL(15,2) NOT NULL,
    payout_method VARCHAR(50), -- bank_transfer, cash, etc.
    payout_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID,
    FOREIGN KEY (incentive_id) REFERENCES incentives(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX idx_incentive_payouts_incentive_id ON incentive_payouts(incentive_id);
```

**Relasi dengan Table Existing**:

- `incentives.user_id` → `users.id`
- `incentives.approved_by` → `users.id`
- `incentive_payouts.created_by` → `users.id`
- Sales Performance calculation menggunakan:
  - `deals` table (status='won', assigned_to=userId, closed_at between period_start and period_end)
  - `visit_reports` table (status='approved', sales_rep_id=userId)
  - `tasks` table (status='completed', assigned_to=userId)
- Sales Rep Activities menggunakan:
  - `visit_reports` table (sales_rep_id=userId) - ordered by visit_date, created_at
  - `deals` table (assigned_to=userId) - ordered by created_at
  - `tasks` table (assigned_to=userId) - ordered by created_at
- Sales Rep Check-in Locations menggunakan:
  - `visit_reports` table (sales_rep_id=userId, check_in_location IS NOT NULL)
  - Ordered by visit_date ASC, check_in_time ASC (untuk numbering 1, 2, 3, ...)
  - Include visit number (ROW_NUMBER() berdasarkan urutan check-in)
- Incentive calculation menggunakan:
  - `deals` table (status='won', assigned_to=userId, closed_at between period_start and period_end)

**Catatan Penting**:
- **Leaderboard** (Sprint 1): Untuk ranking/gamification, bisa filter `metric=revenue, limit=10` untuk top sales
- **Sales Performance** (Sprint 2): Untuk management analysis, detail performance per user (bukan ranking), fokus pada breakdown metrics dan incentive management

**API Response Structure untuk Check-in Locations**:

```json
{
  "success": true,
  "data": {
    "sales_rep": {
      "id": "user-uuid",
      "name": "John Doe",
      "email": "john@example.com"
    },
    "check_in_locations": [
      {
        "visit_number": 1,
        "visit_report_id": "visit-uuid-1",
        "visit_date": "2024-01-15",
        "check_in_time": "2024-01-15T09:00:00+07:00",
        "location": {
          "latitude": -6.2088,
          "longitude": 106.8456,
          "address": "Jl. Sudirman No. 123, Jakarta"
        },
        "account": {
          "id": "account-uuid",
          "name": "RS Cipto Mangunkusumo"
        },
        "purpose": "Product presentation"
      },
      {
        "visit_number": 2,
        "visit_report_id": "visit-uuid-2",
        "visit_date": "2024-01-16",
        "check_in_time": "2024-01-16T10:30:00+07:00",
        "location": {
          "latitude": -6.2146,
          "longitude": 106.8451,
          "address": "Jl. Thamrin No. 45, Jakarta"
        },
        "account": {
          "id": "account-uuid-2",
          "name": "RS Pondok Indah"
        },
        "purpose": "Follow-up meeting"
      }
      // ... dst (ordered by visit_date ASC, check_in_time ASC)
    ],
    "total_visits": 2345,
    "period": {
      "start": "2024-01-01",
      "end": "2024-12-31"
    }
  }
}
```

**Frontend Map Implementation**:
- Map markers harus dinomori 1, 2, 3, ... sesuai `visit_number`
- Click marker → show visit details popup
- List di sidebar dengan visit number dan account name
- Click list item → pan map ke marker tersebut

---

### Sprint 3: Time-Based Sales Analytics (Fullstack)

**Goal**: Implement Time-Based Sales Analytics secara fullstack (backend + frontend)

**Backend Tasks**:

- [ ] Create sales_analytics service
- [ ] Implement sales by month API (`GET /api/v1/sales-analytics/by-month`)
  - Support date range filtering
  - Group by month
  - Calculate: deals_count, revenue, avg_deal_value, won_deals, won_revenue
- [ ] Implement sales by year API (`GET /api/v1/sales-analytics/by-year`)
  - Group by year
  - Calculate yearly metrics
- [ ] Implement month-to-year comparison API (`GET /api/v1/sales-analytics/month-to-year`)
  - Compare current month with same month previous year
  - Calculate YoY growth percentage
- [ ] Implement sales trends API (`GET /api/v1/sales-analytics/trends`)
  - Calculate moving averages
  - Identify seasonal patterns
- [ ] Implement sales forecast API (`GET /api/v1/sales-analytics/forecast`)
  - Simple linear regression or time series forecast
  - Project next N months
- [ ] Add materialized views for performance:
  - `monthly_sales_summary`
  - `yearly_sales_summary`
- [ ] Add refresh logic for materialized views
- [ ] Add date range validation
- [ ] Add pagination support

**Frontend Tasks**:

- [ ] Create sales analytics types (`types/sales-analytics.d.ts`)
- [ ] Create sales analytics service (`salesAnalyticsService`)
- [ ] Create sales analytics page (`/sales-analytics`)
- [ ] Create monthly sales chart component (`MonthlySalesChart`)
- [ ] Create yearly sales chart component (`YearlySalesChart`)
- [ ] Create month-to-year comparison component (`MonthToYearComparison`)
- [ ] Create sales trend chart component (`SalesTrendChart`)
- [ ] Create sales forecast chart component (`SalesForecastChart`)
- [ ] Create date range filter component (`DateRangeFilter`)
- [ ] Add chart types: line, bar, area
- [ ] Add interactive tooltips
- [ ] Add zoom and pan functionality
- [ ] Add export functionality (Excel/PDF)
- [ ] Add period comparison UI

**Postman Collection**:

- [ ] Add sales analytics APIs ke Postman collection (Web section)

**Menu & Permissions**:

- [ ] Add Sales Analytics menu to menu seeder
- [ ] Add Sales Analytics permissions to permission seeder

**Acceptance Criteria**:

- ✅ Sales analytics APIs bekerja dengan baik
- ✅ Monthly/yearly aggregation accurate
- ✅ Month-to-year comparison accurate
- ✅ Trend analysis accurate
- ✅ Forecast projection reasonable
- ✅ Frontend terintegrasi dengan backend APIs
- ✅ Charts ditampilkan dengan benar
- ✅ Date range filtering bekerja
- ✅ Postman collection updated

**Testing** (Manual testing):

- Test monthly sales aggregation (backend + frontend)
- Test yearly sales aggregation
- Test month-to-year comparison
- Test trend calculation
- Test forecast projection
- Test date range filtering

**Estimated Time**: 5-6 days

**Database Schema**:

```sql
-- Materialized view for monthly sales (for performance)
CREATE MATERIALIZED VIEW monthly_sales_summary AS
SELECT
    DATE_TRUNC('month', d.actual_close_date) AS month,
    COUNT(DISTINCT d.id) AS deals_count,
    COUNT(DISTINCT d.assigned_to) AS sales_reps_count,
    SUM(d.value) AS total_revenue,
    AVG(d.value) AS avg_deal_value,
    COUNT(DISTINCT CASE WHEN d.status = 'won' THEN d.id END) AS won_deals,
    SUM(CASE WHEN d.status = 'won' THEN d.value ELSE 0 END) AS won_revenue
FROM deals d
WHERE d.actual_close_date IS NOT NULL
GROUP BY DATE_TRUNC('month', d.actual_close_date);

CREATE UNIQUE INDEX idx_monthly_sales_summary_month ON monthly_sales_summary(month);

-- Materialized view for yearly sales
CREATE MATERIALIZED VIEW yearly_sales_summary AS
SELECT
    DATE_TRUNC('year', d.actual_close_date) AS year,
    COUNT(DISTINCT d.id) AS deals_count,
    COUNT(DISTINCT d.assigned_to) AS sales_reps_count,
    SUM(d.value) AS total_revenue,
    AVG(d.value) AS avg_deal_value,
    COUNT(DISTINCT CASE WHEN d.status = 'won' THEN d.id END) AS won_deals,
    SUM(CASE WHEN d.status = 'won' THEN d.value ELSE 0 END) AS won_revenue
FROM deals d
WHERE d.actual_close_date IS NOT NULL
GROUP BY DATE_TRUNC('year', d.actual_close_date);

CREATE UNIQUE INDEX idx_yearly_sales_summary_year ON yearly_sales_summary(year);

-- Refresh materialized views (run periodically via cron or worker)
-- REFRESH MATERIALIZED VIEW CONCURRENTLY monthly_sales_summary;
-- REFRESH MATERIALIZED VIEW CONCURRENTLY yearly_sales_summary;
```

**Relasi dengan Table Existing**:

- Analytics menggunakan `deals` table
- Filter berdasarkan `actual_close_date` untuk closed deals
- Filter berdasarkan `assigned_to` untuk per sales rep

---

### Sprint 4: Top Products Analytics (Fullstack)

**Goal**: Implement Top Products Analytics secara fullstack (backend + frontend)

**Backend Tasks**:

- [x] Create product_sales model dan migration (if not exists)
- [x] Create product_analytics service
- [x] Implement top products API (`GET /api/v1/product-analytics/top-products`)
  - Get top N products (default: 10)
  - Support filters: period, metric (quantity, revenue, growth)
  - Calculate: quantity_sold, total_revenue, avg_price, growth_rate
- [x] Implement product performance API (`GET /api/v1/product-analytics/product/:id/performance`)
- [x] Implement product comparison API (`GET /api/v1/product-analytics/product-comparison`)
  - Compare multiple products side-by-side
- [x] Implement product trends API (`GET /api/v1/product-analytics/product-trends/:id`)
  - Product sales trend over time
- [x] Add product sales tracking (if not exists):
  - Link deals to products
  - Track quantity and revenue per product
- [x] Add aggregation logic:
  - Top products by quantity sold
  - Top products by revenue
  - Top products by growth rate
- [x] Add pagination support
- [x] Add validation

**Frontend Tasks**:

- [x] Create product analytics types (`types/product-analytics.d.ts`)
- [x] Create product analytics service (`productAnalyticsService`)
- [x] Create product analytics page (`/product-analytics`)
- [x] Create top products chart component (`TopProductsChart`)
  - Bar chart for horizontal display
- [x] Create product performance card component (`ProductPerformanceCard`)
- [x] Create product comparison view component (`ProductComparisonView`)
- [x] Create product trend chart component (`ProductTrendChart`)
  - Line chart with dual Y-axis (quantity & revenue)
- [x] Create product analytics filters component (`ProductAnalyticsFilters`)
- [x] Add period filtering (day, week, month, year)
- [x] Add metric selection (quantity, revenue, growth)
- [x] Add date range filtering for trends
- [x] Add group by selection for trends
- [x] Add i18n translations (en.json, id.json)
- [ ] Add export functionality (Excel/PDF) - Optional, can be added later

**Postman Collection**:

- [x] Add product analytics APIs ke Postman collection (Web section)

**Menu & Permissions**:

- [x] Add Product Analytics menu to menu seeder
- [x] Add Product Analytics permissions to permission seeder

**Acceptance Criteria**:

- ✅ Product analytics APIs bekerja dengan baik
- ✅ Top products calculation accurate
- ✅ Product performance metrics accurate
- ✅ Product comparison bekerja
- ✅ Product trends ditampilkan dengan benar
- ✅ Frontend terintegrasi dengan backend APIs
- ✅ Charts ditampilkan dengan benar
- ✅ Filters bekerja optimal
- ✅ Postman collection updated

**Testing** (Manual testing):

- Test top products calculation (backend + frontend)
- Test product performance metrics
- Test product comparison
- Test product trends
- Test filters and period selection

**Estimated Time**: 5-6 days

**Database Schema**:

```sql
-- Product sales tracking (if not exists)
CREATE TABLE IF NOT EXISTS product_sales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deal_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    sold_at TIMESTAMP NOT NULL,
    sales_rep_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (deal_id) REFERENCES deals(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (sales_rep_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_sales_product_id ON product_sales(product_id);
CREATE INDEX idx_product_sales_sold_at ON product_sales(sold_at);
CREATE INDEX idx_product_sales_sales_rep_id ON product_sales(sales_rep_id);
CREATE INDEX idx_product_sales_deal_id ON product_sales(deal_id);
```

**Relasi dengan Table Existing**:

- `product_sales.deal_id` → `deals.id`
- `product_sales.product_id` → `products.id`
- `product_sales.sales_rep_id` → `users.id`
- Analytics menggunakan `product_sales` table dengan join ke `products` dan `deals`

---

## 📊 Sprint Summary

| Sprint   | Goal                            | Duration | Status       |
| -------- | ------------------------------- | -------- | ------------ |
| Sprint 1 | Leaderboard & KPI (Fullstack)   | 6-7 days | ✅ Completed |
| Sprint 2 | Sales Performance & Incentive    | 7-8 days | ✅ Completed |
| Sprint 3 | Time-Based Sales Analytics      | 5-6 days | ⏳ Pending   |
| Sprint 4 | Top Products Analytics          | 5-6 days | ✅ Completed |

**Total Estimated Time**: 23-27 days (3.3-3.9 weeks)

---

## 🔗 Coordination dengan Dev5

### Modul yang dikerjakan Dev5 (untuk referensi):

- Route Optimization (Fullstack)
- Schedule Assignment (Fullstack)
- Area Mapping (Fullstack)

### Integration Points:

- Leaderboard: General ranking system (termasuk top sales dengan filter `metric=revenue, limit=10`)
- Sales Performance: Detail performance per user (bukan ranking), untuk management analysis
- Sales Performance bisa menggunakan data dari Schedule Assignment (completed schedules count)
- Route Optimization efficiency metrics bisa digunakan untuk Sales Performance (future enhancement)
- Area Mapping bisa digunakan untuk territory-based analytics (future enhancement)

### Coordination:

- [ ] Week 1: Coordinate API contract untuk integration points (if needed)
- [ ] Week 2: Mid-sprint review - check integration points
- [ ] Week 3: Pre-integration review
- [ ] Week 4: Final integration testing (if needed)

---

## 📝 Notes

1. **Fullstack Development**: Setiap modul dikerjakan fullstack sampai selesai
2. **No Dependencies**: Tidak bergantung ke Dev5, bisa dikerjakan paralel
3. **Hackathon Mode**: Tidak ada unit test, manual testing saja
4. **Code Review**: Lakukan code review sebelum merge
5. **Documentation**: Update documentation setelah setiap sprint
6. **Postman Collection**: Update Postman collection untuk setiap modul
7. **API Standards**: Follow `/docs/api-standart/api-response-standards.md` dan `/docs/api-standart/api-error-codes.md`
8. **Frontend Standards**: Follow `.cursor/rules/standart.mdc` untuk folder structure, types, hooks, services, components

---

**Dokumen ini akan diupdate sesuai dengan progress development.**

