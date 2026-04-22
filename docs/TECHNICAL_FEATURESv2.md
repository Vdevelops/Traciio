# Technical Features Documentation v2.0

## 📋 Table of Contents

1. [Overview](#overview)
2. [New Features Summary](#new-features-summary)
3. [Tech Stack Additions](#tech-stack-additions)
4. [New Backend Modules](#new-backend-modules)
5. [New Frontend Modules](#new-frontend-modules)
6. [Enhanced Existing Modules](#enhanced-existing-modules)
7. [Infrastructure & Deployment](#infrastructure--deployment)
8. [Budget Breakdown](#budget-breakdown)
9. [Implementation Timeline](#implementation-timeline)

---

## Overview

**TECHNICAL_FEATURESv2.md** adalah dokumentasi untuk update fitur-fitur baru pada CRM Healthcare yang mencakup:

- **Performance Monitoring** - Leaderboard dan KPI tracking
- **Route Optimization** - Optimasi rute kunjungan untuk efisiensi
- **Advanced Scheduling** - Schedule assignment untuk task dan visit management
- **Geographic Mapping** - Area mapping dan location capture
- **Analytics Enhancement** - Top products, sales overview, dan time-based analytics
- **Server Installation** - Dokumentasi dan setup untuk instalasi server

**Budget Total:** Rp 450.000.000 (450 Juta Rupiah)

---

## New Features Summary

### 1. Leaderboard & KPI Monitoring
**Tujuan:** Memantau performa tim berdasarkan KPI untuk meningkatkan produktivitas dan motivasi.

**Fitur Utama:**
- Leaderboard berdasarkan berbagai metrik KPI (deals closed, revenue, visits, tasks completed)
- Ranking per periode (daily, weekly, monthly, yearly)
- KPI tracking per sales rep
- Achievement badges dan milestones
- Performance trends dan comparison

### 2. Route Optimization
**Tujuan:** Memberikan rute kunjungan yang efisien untuk menghemat waktu dan biaya perjalanan.

**Fitur Utama:**
- Multi-stop route optimization menggunakan algoritma TSP (Traveling Salesman Problem)
- Integration dengan Google Maps API / Mapbox untuk routing
- Optimasi berdasarkan jarak, waktu, atau prioritas
- Visualisasi rute di peta
- Export rute ke GPS/navigation apps
- Real-time traffic consideration

### 3. Schedule Assignment & Advanced Task Management
**Tujuan:** Memungkinkan manager untuk assign jadwal kunjungan dan task kepada tim sales dengan lebih terstruktur.

**Fitur Utama:**
- Schedule assignment untuk visit reports
- Bulk task assignment
- Calendar view untuk schedule management
- Recurring schedule support
- Schedule conflict detection
- Notification untuk assigned schedules
- Schedule approval workflow

### 4. Area Mapping & Location Capture
**Tujuan:** Capture dan mapping area saat visit untuk tracking coverage dan territory management.

**Fitur Utama:**
- GPS location capture saat check-in/check-out
- Area polygon drawing untuk territory mapping
- Map visualization dengan markers untuk visits
- Coverage area analysis
- Territory assignment per sales rep
- Heat map untuk visit frequency

### 5. Top Products Analytics
**Tujuan:** Visualisasi top 10 produk yang terjual untuk analisis penjualan.

**Fitur Utama:**
- Top 10 products chart (bar chart, pie chart)
- Filter berdasarkan periode (day, week, month, year)
- Product performance metrics (quantity, revenue, growth)
- Product comparison view
- Export to Excel/PDF

### 6. Sales Overview & Incentive Management
**Tujuan:** Overview performa top 10 sales dengan sistem insentif.

**Fitur Utama:**
- Top 10 sales rep leaderboard
- Sales performance metrics (revenue, deals, conversion rate)
- Incentive calculation based on performance
- Incentive tier system (bronze, silver, gold, platinum)
- Incentive history tracking
- Incentive payout management

### 7. Time-Based Sales Analytics
**Tujuan:** Analisis penjualan berdasarkan waktu (month, year, month-to-year).

**Fitur Utama:**
- Sales by month (line chart, bar chart)
- Sales by year (yearly comparison)
- Month-to-year comparison (YoY growth)
- Date range filtering
- Trend analysis
- Forecast projection

### 8. Server Installation & Deployment
**Tujuan:** Dokumentasi dan setup untuk instalasi server production.

**Fitur Utama:**
- Docker-based deployment
- Production environment setup
- Database migration scripts
- Environment configuration
- Monitoring and logging setup
- Backup and recovery procedures

---

## Tech Stack Additions

### Backend (API) - New Dependencies

**Route Optimization:**
- **Google Maps API** / **Mapbox API** - Geocoding and routing
- **OR-Tools (Go bindings)** - Optimization algorithms (TSP solver)
- **PostGIS extension** - Geographic data support (optional)

**Geographic & Mapping:**
- **PostGIS** - Spatial database extension for PostgreSQL
- **Geolib** - Distance and area calculations

**Analytics & Reporting:**
- **Chart.js Go bindings** - Server-side chart generation (optional)
- Enhanced **Excelize** usage** - Advanced Excel reports

**Scheduling:**
- **Cron parser** - Recurring schedule parsing
- **Time zone support** - Multi-timezone scheduling

### Frontend (Web) - New Dependencies

**Maps & Geographic:**
- **Leaflet** / **Mapbox GL JS** - Interactive maps
- **react-leaflet** - React bindings for Leaflet
- **@turf/turf** - Geospatial analysis library

**Charts & Visualization:**
- **Recharts** (existing) - Enhanced usage
- **Victory** (optional) - Additional chart types

**Calendar & Scheduling:**
- **react-big-calendar** / **@fullcalendar/react** - Calendar component
- **date-fns** (existing) - Enhanced usage

**Route Visualization:**
- **react-map-gl** - Mapbox React integration
- **@mapbox/mapbox-gl-directions** - Route directions

---

## New Backend Modules

### 16. Leaderboard Module (`leaderboard`)

**Features:**
- KPI calculation per sales rep
- Leaderboard ranking (daily, weekly, monthly, yearly)
- Performance metrics aggregation
- Achievement tracking
- Badge system

**Endpoints:**
- `GET /api/v1/leaderboard` - Get leaderboard (with filters: period, metric, limit)
- `GET /api/v1/leaderboard/kpi/:userId` - Get KPI for specific user
- `GET /api/v1/leaderboard/achievements/:userId` - Get achievements for user
- `GET /api/v1/leaderboard/trends/:userId` - Get performance trends

**Entities:**
- `LeaderboardEntry` - Leaderboard entry with ranking
- `KPIMetric` - KPI metric definition
- `Achievement` - Achievement/badge entity
- `PerformanceTrend` - Performance trend data

**KPI Metrics:**
- Total deals closed
- Total revenue generated
- Number of visits completed
- Tasks completed
- Conversion rate
- Average deal value
- Customer satisfaction score

**Database Schema:**
```sql
-- Leaderboard entries (computed view or materialized table)
CREATE TABLE leaderboard_entries (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    period_type VARCHAR(20) NOT NULL, -- daily, weekly, monthly, yearly
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    metric_type VARCHAR(50) NOT NULL, -- deals, revenue, visits, tasks
    metric_value DECIMAL(15,2) NOT NULL,
    rank INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, period_type, period_start, metric_type)
);

-- Achievements
CREATE TABLE achievements (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    achievement_type VARCHAR(50) NOT NULL, -- first_deal, top_seller, etc.
    achieved_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB
);

-- KPI settings
CREATE TABLE kpi_settings (
    id UUID PRIMARY KEY,
    metric_name VARCHAR(50) UNIQUE NOT NULL,
    weight DECIMAL(5,2) DEFAULT 1.0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

### 17. Route Optimization Module (`route_optimization`)

**Features:**
- Multi-stop route optimization
- Distance and time calculation
- Route generation with waypoints
- Traffic-aware routing
- Route export (GPX, JSON)
- Route history

**Endpoints:**
- `POST /api/v1/route-optimization/optimize` - Optimize route from waypoints
- `GET /api/v1/route-optimization/route/:id` - Get optimized route
- `POST /api/v1/route-optimization/calculate-distance` - Calculate distance between points
- `GET /api/v1/route-optimization/history` - Get route optimization history

**Entities:**
- `RouteOptimizationRequest` - Request with waypoints
- `OptimizedRoute` - Optimized route response
- `Waypoint` - Location waypoint
- `RouteSegment` - Route segment between waypoints

**Algorithm:**
- TSP (Traveling Salesman Problem) solver using OR-Tools
- Nearest Neighbor heuristic for quick solutions
- 2-opt improvement for better routes
- Consider traffic data if available

**Database Schema:**
```sql
CREATE TABLE optimized_routes (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    route_name VARCHAR(255),
    waypoints JSONB NOT NULL, -- Array of {lat, lng, address, account_id}
    optimized_order JSONB NOT NULL, -- Optimized waypoint order
    total_distance DECIMAL(10,2), -- in km
    total_duration INTEGER, -- in seconds
    route_polyline TEXT, -- Encoded polyline for map display
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_optimized_routes_user_id ON optimized_routes(user_id);
CREATE INDEX idx_optimized_routes_created_at ON optimized_routes(created_at);
```

**Integration:**
- Google Maps Directions API or Mapbox Directions API
- OR-Tools for optimization algorithm

---

### 18. Schedule Assignment Module (`schedule_assignment`)

**Features:**
- Schedule creation and assignment
- Bulk schedule assignment
- Recurring schedules
- Schedule conflict detection
- Calendar integration
- Schedule approval workflow

**Endpoints:**
- `POST /api/v1/schedules` - Create schedule
- `GET /api/v1/schedules` - List schedules (with filters)
- `GET /api/v1/schedules/:id` - Get schedule details
- `PUT /api/v1/schedules/:id` - Update schedule
- `DELETE /api/v1/schedules/:id` - Delete schedule
- `POST /api/v1/schedules/bulk-assign` - Bulk assign schedules
- `POST /api/v1/schedules/:id/assign` - Assign schedule to user
- `GET /api/v1/schedules/conflicts` - Check for schedule conflicts
- `POST /api/v1/schedules/:id/approve` - Approve schedule
- `POST /api/v1/schedules/:id/reject` - Reject schedule

**Entities:**
- `Schedule` - Schedule entity
- `ScheduleAssignment` - Assignment of schedule to user
- `ScheduleConflict` - Schedule conflict detection
- `RecurringSchedule` - Recurring schedule pattern

**Database Schema:**
```sql
CREATE TABLE schedules (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL, -- visit, task, meeting, other
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    assigned_to UUID, -- User ID
    assigned_by UUID NOT NULL, -- User who assigned
    account_id UUID, -- Optional: linked account
    contact_id UUID, -- Optional: linked contact
    deal_id UUID, -- Optional: linked deal
    visit_report_id UUID, -- Optional: linked visit report
    task_id UUID, -- Optional: linked task
    location JSONB, -- {lat, lng, address}
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected, completed, cancelled
    is_recurring BOOLEAN DEFAULT false,
    recurring_pattern JSONB, -- {type: daily/weekly/monthly, interval, end_date}
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE schedule_assignments (
    id UUID PRIMARY KEY,
    schedule_id UUID NOT NULL,
    user_id UUID NOT NULL,
    assigned_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'pending', -- pending, accepted, rejected
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
);

CREATE INDEX idx_schedules_assigned_to ON schedules(assigned_to);
CREATE INDEX idx_schedules_start_time ON schedules(start_time);
CREATE INDEX idx_schedules_status ON schedules(status);
```

**Recurring Pattern:**
```json
{
  "type": "daily|weekly|monthly",
  "interval": 1, // Every 1 day/week/month
  "days_of_week": [1, 3, 5], // For weekly: Monday, Wednesday, Friday
  "day_of_month": 15, // For monthly: day 15
  "end_date": "2024-12-31",
  "occurrences": 10 // Or limit by number of occurrences
}
```

---

### 19. Area Mapping Module (`area_mapping`)

**Features:**
- GPS location capture
- Area polygon drawing
- Territory mapping
- Coverage analysis
- Heat map generation
- Territory assignment

**Endpoints:**
- `POST /api/v1/area-mapping/capture` - Capture area from visit
- `POST /api/v1/area-mapping/territory` - Create territory polygon
- `GET /api/v1/area-mapping/territories` - List territories
- `GET /api/v1/area-mapping/coverage` - Get coverage analysis
- `GET /api/v1/area-mapping/heatmap` - Get heat map data
- `POST /api/v1/area-mapping/assign-territory` - Assign territory to user

**Entities:**
- `AreaCapture` - Captured area from visit
- `Territory` - Territory polygon definition
- `CoverageAnalysis` - Coverage analysis data
- `HeatMapData` - Heat map data points

**Database Schema:**
```sql
-- Using PostGIS extension for spatial data
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE area_captures (
    id UUID PRIMARY KEY,
    visit_report_id UUID NOT NULL,
    capture_type VARCHAR(20) NOT NULL, -- check_in, check_out, area
    location GEOGRAPHY(POINT, 4326) NOT NULL, -- GPS point
    address TEXT,
    accuracy DECIMAL(10,2), -- GPS accuracy in meters
    captured_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (visit_report_id) REFERENCES visit_reports(id) ON DELETE CASCADE
);

CREATE TABLE territories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    polygon GEOGRAPHY(POLYGON, 4326) NOT NULL, -- Territory polygon
    assigned_to UUID, -- User ID
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE coverage_analysis (
    id UUID PRIMARY KEY,
    territory_id UUID,
    user_id UUID,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    visit_count INTEGER DEFAULT 0,
    coverage_percent DECIMAL(5,2), -- Coverage percentage
    analyzed_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (territory_id) REFERENCES territories(id) ON DELETE SET NULL
);

-- Spatial indexes for performance
CREATE INDEX idx_area_captures_location ON area_captures USING GIST(location);
CREATE INDEX idx_territories_polygon ON territories USING GIST(polygon);
```

**Integration:**
- PostGIS for spatial queries
- Leaflet/Mapbox for map visualization
- Turf.js for client-side geospatial calculations

---

### 20. Product Analytics Module (`product_analytics`)

**Features:**
- Top products calculation
- Product performance metrics
- Product comparison
- Sales trend by product
- Product revenue analysis

**Endpoints:**
- `GET /api/v1/product-analytics/top-products` - Get top N products (default: 10)
- `GET /api/v1/product-analytics/product/:id/performance` - Get product performance
- `GET /api/v1/product-analytics/product-comparison` - Compare multiple products
- `GET /api/v1/product-analytics/product-trends/:id` - Get product sales trends

**Entities:**
- `TopProduct` - Top product entry
- `ProductPerformance` - Product performance metrics
- `ProductComparison` - Product comparison data
- `ProductTrend` - Product trend over time

**Database Schema:**
```sql
-- Product sales tracking (if not exists)
CREATE TABLE product_sales (
    id UUID PRIMARY KEY,
    deal_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    sold_at TIMESTAMP NOT NULL,
    sales_rep_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (deal_id) REFERENCES deals(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_sales_product_id ON product_sales(product_id);
CREATE INDEX idx_product_sales_sold_at ON product_sales(sold_at);
CREATE INDEX idx_product_sales_sales_rep_id ON product_sales(sales_rep_id);
```

**Analytics Queries:**
- Top products by quantity sold
- Top products by revenue
- Top products by growth rate
- Product performance by period

---

### 21. Sales Overview & Incentive Module (`sales_overview`, `incentive`)

**Features:**
- Top 10 sales rep calculation
- Sales performance metrics
- Incentive calculation
- Incentive tier system
- Incentive payout tracking

**Endpoints:**
- `GET /api/v1/sales-overview/top-sales` - Get top N sales reps (default: 10)
- `GET /api/v1/sales-overview/performance/:userId` - Get sales performance
- `GET /api/v1/incentives` - List incentives
- `GET /api/v1/incentives/:id` - Get incentive details
- `POST /api/v1/incentives/calculate` - Calculate incentive for period
- `POST /api/v1/incentives/payout` - Record incentive payout

**Entities:**
- `TopSalesRep` - Top sales rep entry
- `SalesPerformance` - Sales performance metrics
- `Incentive` - Incentive entity
- `IncentiveTier` - Incentive tier definition
- `IncentivePayout` - Incentive payout record

**Database Schema:**
```sql
CREATE TABLE incentive_tiers (
    id UUID PRIMARY KEY,
    tier_name VARCHAR(50) NOT NULL, -- bronze, silver, gold, platinum
    min_revenue DECIMAL(15,2) NOT NULL,
    max_revenue DECIMAL(15,2),
    incentive_percent DECIMAL(5,2) NOT NULL, -- Percentage of revenue
    bonus_amount DECIMAL(15,2) DEFAULT 0, -- Fixed bonus
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE incentives (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    period_type VARCHAR(20) NOT NULL, -- monthly, quarterly, yearly
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    revenue DECIMAL(15,2) NOT NULL,
    deals_closed INTEGER NOT NULL,
    tier_id UUID,
    incentive_amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, paid
    calculated_at TIMESTAMP DEFAULT NOW(),
    approved_at TIMESTAMP,
    paid_at TIMESTAMP,
    FOREIGN KEY (tier_id) REFERENCES incentive_tiers(id)
);

CREATE TABLE incentive_payouts (
    id UUID PRIMARY KEY,
    incentive_id UUID NOT NULL,
    payout_amount DECIMAL(15,2) NOT NULL,
    payout_method VARCHAR(50), -- bank_transfer, cash, etc.
    payout_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (incentive_id) REFERENCES incentives(id)
);

CREATE INDEX idx_incentives_user_id ON incentives(user_id);
CREATE INDEX idx_incentives_period ON incentives(period_start, period_end);
CREATE INDEX idx_incentives_status ON incentives(status);
```

**Incentive Calculation Logic:**
1. Calculate total revenue for period
2. Calculate deals closed
3. Determine tier based on revenue threshold
4. Calculate incentive: `(revenue * tier.incentive_percent / 100) + tier.bonus_amount`
5. Apply any additional bonuses (milestones, achievements)

---

### 22. Time-Based Sales Analytics Module (`sales_analytics`)

**Features:**
- Sales by month
- Sales by year
- Month-to-year comparison (YoY)
- Date range filtering
- Trend analysis
- Forecast projection

**Endpoints:**
- `GET /api/v1/sales-analytics/by-month` - Get sales by month
- `GET /api/v1/sales-analytics/by-year` - Get sales by year
- `GET /api/v1/sales-analytics/month-to-year` - Get month-to-year comparison
- `GET /api/v1/sales-analytics/trends` - Get sales trends
- `GET /api/v1/sales-analytics/forecast` - Get sales forecast

**Entities:**
- `MonthlySales` - Monthly sales data
- `YearlySales` - Yearly sales data
- `SalesTrend` - Sales trend data
- `SalesForecast` - Sales forecast projection

**Database Schema:**
```sql
-- Materialized view for monthly sales (for performance)
CREATE MATERIALIZED VIEW monthly_sales_summary AS
SELECT
    DATE_TRUNC('month', d.closed_at) AS month,
    COUNT(DISTINCT d.id) AS deals_count,
    COUNT(DISTINCT d.sales_rep_id) AS sales_reps_count,
    SUM(d.value) AS total_revenue,
    AVG(d.value) AS avg_deal_value,
    COUNT(DISTINCT CASE WHEN d.status = 'won' THEN d.id END) AS won_deals,
    SUM(CASE WHEN d.status = 'won' THEN d.value ELSE 0 END) AS won_revenue
FROM deals d
WHERE d.closed_at IS NOT NULL
GROUP BY DATE_TRUNC('month', d.closed_at);

CREATE UNIQUE INDEX idx_monthly_sales_summary_month ON monthly_sales_summary(month);

-- Refresh materialized view (run periodically)
-- REFRESH MATERIALIZED VIEW CONCURRENTLY monthly_sales_summary;

-- Similar for yearly
CREATE MATERIALIZED VIEW yearly_sales_summary AS
SELECT
    DATE_TRUNC('year', d.closed_at) AS year,
    COUNT(DISTINCT d.id) AS deals_count,
    COUNT(DISTINCT d.sales_rep_id) AS sales_reps_count,
    SUM(d.value) AS total_revenue,
    AVG(d.value) AS avg_deal_value,
    COUNT(DISTINCT CASE WHEN d.status = 'won' THEN d.id END) AS won_deals,
    SUM(CASE WHEN d.status = 'won' THEN d.value ELSE 0 END) AS won_revenue
FROM deals d
WHERE d.closed_at IS NOT NULL
GROUP BY DATE_TRUNC('year', d.closed_at);

CREATE UNIQUE INDEX idx_yearly_sales_summary_year ON yearly_sales_summary(year);
```

**Analytics Features:**
- Month-over-month growth
- Year-over-year growth
- Moving averages
- Seasonal patterns
- Forecast using linear regression or time series

---

## New Frontend Modules

### 14. Leaderboard Module (`leaderboard`)

**Components:**
- `LeaderboardView` - Main leaderboard view
- `LeaderboardTable` - Leaderboard table with rankings
- `KPICard` - KPI metric card
- `AchievementBadge` - Achievement badge component
- `PerformanceTrendChart` - Performance trend line chart
- `LeaderboardFilters` - Filter component (period, metric)

**Hooks:**
- `useLeaderboard` - Leaderboard data fetching
- `useKPI` - KPI metrics fetching
- `useAchievements` - Achievements fetching
- `usePerformanceTrends` - Performance trends fetching

**Features:**
- Real-time leaderboard updates
- Multiple metric views (deals, revenue, visits, tasks)
- Period filtering (daily, weekly, monthly, yearly)
- Achievement display
- Performance trend visualization
- Export to Excel/PDF

**Routes:**
- `/leaderboard` - Main leaderboard page
- `/leaderboard/kpi/:userId` - User KPI details
- `/leaderboard/achievements/:userId` - User achievements

---

### 15. Route Optimization Module (`route-optimization`)

**Components:**
- `RouteOptimizationView` - Main route optimization view
- `RouteMap` - Interactive map with route visualization
- `WaypointList` - List of waypoints with drag-and-drop reordering
- `RouteOptimizationForm` - Form to create optimization request
- `RouteDetails` - Route details (distance, duration, steps)
- `RouteHistory` - History of optimized routes
- `RouteExportDialog` - Export route dialog (GPX, JSON)

**Hooks:**
- `useRouteOptimization` - Route optimization mutation
- `useRouteHistory` - Route history fetching
- `useDistanceCalculation` - Distance calculation

**Features:**
- Interactive map with markers
- Drag-and-drop waypoint reordering
- Route visualization with polyline
- Distance and duration display
- Route export (GPX for GPS, JSON for API)
- Route history
- Integration with visit reports (auto-populate waypoints)

**Routes:**
- `/route-optimization` - Main route optimization page
- `/route-optimization/history` - Route history

**Map Integration:**
- Leaflet or Mapbox GL JS
- Marker clustering for multiple waypoints
- Route polyline visualization
- Directions display

---

### 16. Schedule Assignment Module (`schedule-assignment`)

**Components:**
- `ScheduleManagementView` - Main schedule management view
- `ScheduleCalendar` - Calendar view with schedules
- `ScheduleList` - List view of schedules
- `ScheduleForm` - Create/edit schedule form
- `ScheduleAssignmentDialog` - Assign schedule dialog
- `BulkAssignmentDialog` - Bulk assignment dialog
- `ScheduleConflictAlert` - Conflict detection alert
- `RecurringScheduleConfig` - Recurring schedule configuration

**Hooks:**
- `useSchedules` - Schedule data management
- `useScheduleAssignment` - Schedule assignment mutation
- `useScheduleConflicts` - Conflict detection
- `useBulkAssignment` - Bulk assignment mutation
- `useCalendarSchedules` - Calendar data fetching

**Features:**
- Calendar view (month, week, day)
- List view with filters
- Schedule creation with time picker
- Bulk assignment to multiple users
- Recurring schedule support
- Conflict detection and alerts
- Schedule approval workflow
- Integration with visit reports and tasks
- Notification for assigned schedules

**Routes:**
- `/schedules` - Main schedule management page
- `/schedules/calendar` - Calendar view
- `/schedules/:id` - Schedule details

**Calendar Library:**
- react-big-calendar or @fullcalendar/react
- Support for drag-and-drop scheduling
- Time slot selection
- Multiple view modes

---

### 17. Area Mapping Module (`area-mapping`)

**Components:**
- `AreaMappingView` - Main area mapping view
- `TerritoryMap` - Interactive map with territory polygons
- `AreaCaptureDialog` - Capture area dialog
- `TerritoryForm` - Create/edit territory form
- `CoverageAnalysis` - Coverage analysis component
- `HeatMapView` - Heat map visualization
- `TerritoryAssignmentDialog` - Assign territory dialog

**Hooks:**
- `useAreaCapture` - Area capture mutation
- `useTerritories` - Territory data management
- `useCoverageAnalysis` - Coverage analysis fetching
- `useHeatMapData` - Heat map data fetching
- `useTerritoryAssignment` - Territory assignment mutation

**Features:**
- Interactive map with drawing tools
- Polygon drawing for territories
- GPS location capture
- Territory visualization
- Coverage analysis (visit frequency, coverage percentage)
- Heat map for visit density
- Territory assignment to sales reps
- Integration with visit reports (auto-capture on check-in)

**Routes:**
- `/area-mapping` - Main area mapping page
- `/area-mapping/territories` - Territory management
- `/area-mapping/coverage` - Coverage analysis
- `/area-mapping/heatmap` - Heat map view

**Map Integration:**
- Leaflet with drawing plugin
- Turf.js for geospatial calculations
- PostGIS for server-side spatial queries

---

### 18. Product Analytics Module (`product-analytics`)

**Components:**
- `ProductAnalyticsView` - Main product analytics view
- `TopProductsChart` - Top 10 products chart (bar, pie)
- `ProductPerformanceCard` - Product performance card
- `ProductComparisonView` - Product comparison view
- `ProductTrendChart` - Product trend line chart
- `ProductAnalyticsFilters` - Filter component

**Hooks:**
- `useTopProducts` - Top products fetching
- `useProductPerformance` - Product performance fetching
- `useProductComparison` - Product comparison fetching
- `useProductTrends` - Product trends fetching

**Features:**
- Top 10 products visualization
- Multiple chart types (bar, pie, line)
- Period filtering (day, week, month, year)
- Product performance metrics
- Product comparison (side-by-side)
- Product trend analysis
- Export to Excel/PDF

**Routes:**
- `/product-analytics` - Main product analytics page
- `/product-analytics/top-products` - Top products view
- `/product-analytics/product/:id` - Product details

**Charts:**
- Recharts for bar, pie, line charts
- Responsive design
- Interactive tooltips
- Data export

---

### 19. Sales Overview Module (`sales-overview`)

**Components:**
- `SalesOverviewView` - Main sales overview view
- `TopSalesLeaderboard` - Top 10 sales rep leaderboard
- `SalesPerformanceCard` - Sales performance card
- `IncentiveCard` - Incentive card with tier badge
- `IncentiveHistory` - Incentive history table
- `IncentiveCalculationDialog` - Incentive calculation dialog
- `IncentivePayoutDialog` - Incentive payout dialog

**Hooks:**
- `useTopSales` - Top sales fetching
- `useSalesPerformance` - Sales performance fetching
- `useIncentives` - Incentive data management
- `useIncentiveCalculation` - Incentive calculation
- `useIncentivePayout` - Incentive payout mutation

**Features:**
- Top 10 sales rep leaderboard
- Sales performance metrics
- Incentive tier display (bronze, silver, gold, platinum)
- Incentive calculation
- Incentive payout management
- Incentive history tracking
- Period filtering

**Routes:**
- `/sales-overview` - Main sales overview page
- `/sales-overview/incentives` - Incentive management
- `/sales-overview/performance/:userId` - User performance

---

### 20. Time-Based Sales Analytics Module (`sales-analytics`)

**Components:**
- `SalesAnalyticsView` - Main sales analytics view
- `MonthlySalesChart` - Monthly sales chart
- `YearlySalesChart` - Yearly sales chart
- `MonthToYearComparison` - Month-to-year comparison chart
- `SalesTrendChart` - Sales trend line chart
- `SalesForecastChart` - Sales forecast chart
- `DateRangeFilter` - Date range filter component

**Hooks:**
- `useMonthlySales` - Monthly sales fetching
- `useYearlySales` - Yearly sales fetching
- `useMonthToYearComparison` - Month-to-year comparison fetching
- `useSalesTrends` - Sales trends fetching
- `useSalesForecast` - Sales forecast fetching

**Features:**
- Sales by month (line chart, bar chart)
- Sales by year (yearly comparison)
- Month-to-year comparison (YoY growth)
- Date range filtering
- Trend analysis
- Forecast projection
- Export to Excel/PDF

**Routes:**
- `/sales-analytics` - Main sales analytics page
- `/sales-analytics/monthly` - Monthly sales view
- `/sales-analytics/yearly` - Yearly sales view
- `/sales-analytics/comparison` - Comparison view

**Charts:**
- Recharts for time series charts
- Multiple chart types (line, bar, area)
- Interactive tooltips
- Zoom and pan
- Data export

---

## Enhanced Existing Modules

### Enhanced Dashboard Module

**New Components:**
- `LeaderboardWidget` - Leaderboard widget for dashboard
- `TopProductsWidget` - Top products widget
- `SalesOverviewWidget` - Sales overview widget
- `RouteOptimizationQuickAccess` - Quick access to route optimization

**New Endpoints:**
- Enhanced `GET /api/v1/dashboard/overview` - Include leaderboard, top products, sales overview

---

### Enhanced Task Management Module

**New Features:**
- Schedule assignment integration
- Calendar view for tasks
- Recurring task support
- Task scheduling with time slots

**New Components:**
- `TaskScheduleDialog` - Schedule task dialog
- `TaskCalendarView` - Calendar view for tasks
- `RecurringTaskConfig` - Recurring task configuration

---

### Enhanced Visit Report Module

**New Features:**
- Area capture on check-in/check-out
- Location mapping visualization
- Route optimization integration
- Schedule assignment integration

**New Components:**
- `VisitLocationMap` - Map showing visit location
- `AreaCaptureButton` - Capture area button
- `RouteOptimizationButton` - Optimize route button

**Enhanced Endpoints:**
- `POST /api/v1/visit-reports/:id/check-in` - Enhanced with area capture
- `POST /api/v1/visit-reports/:id/check-out` - Enhanced with area capture

---

## Infrastructure & Deployment

### Server Installation Requirements

**Hardware Requirements:**
- CPU: 4+ cores
- RAM: 8GB+ (16GB recommended)
- Storage: 100GB+ SSD
- Network: Stable internet connection

**Software Requirements:**
- Ubuntu 22.04 LTS / Debian 11+
- Docker 24.0+
- Docker Compose 2.20+
- PostgreSQL 15+
- Nginx (for reverse proxy)

**Environment Setup:**
```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install PostgreSQL (if not using Docker)
sudo apt update
sudo apt install postgresql postgresql-contrib postgis

# Install Nginx
sudo apt install nginx
```

**Docker Compose Configuration:**
```yaml
version: '3.8'

services:
  postgres:
    image: postgis/postgis:15-3.3
    environment:
      POSTGRES_DB: crm_healthcare
      POSTGRES_USER: crm_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  api:
    build: ./apps/api
    environment:
      DATABASE_URL: postgres://crm_user:${DB_PASSWORD}@postgres:5432/crm_healthcare?sslmode=disable
      JWT_SECRET: ${JWT_SECRET}
      # ... other env vars
    depends_on:
      - postgres
    ports:
      - "8080:8080"

  web:
    build: ./apps/web
    environment:
      NEXT_PUBLIC_API_URL: ${API_URL}
    depends_on:
      - api
    ports:
      - "3000:3000"

volumes:
  postgres_data:
```

**Nginx Configuration:**
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**SSL Setup (Let's Encrypt):**
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

**Backup Script:**
```bash
#!/bin/bash
# backup.sh
DATE=$(date +%Y%m%d_%H%M%S)
docker exec postgres pg_dump -U crm_user crm_healthcare > backup_$DATE.sql
# Upload to S3 or backup storage
```

**Monitoring:**
- Application logs: Docker logs
- Database monitoring: pgAdmin or similar
- Server monitoring: Prometheus + Grafana (optional)

---

## Budget Breakdown

**Total Budget: Rp 450.000.000 (450 Juta Rupiah)**

### 1. Development (Rp 280.000.000 - 62%)

**Backend Development:**
- Leaderboard & KPI Module: Rp 35.000.000
- Route Optimization Module: Rp 45.000.000
- Schedule Assignment Module: Rp 40.000.000
- Area Mapping Module: Rp 35.000.000
- Product Analytics Module: Rp 25.000.000
- Sales Overview & Incentive Module: Rp 35.000.000
- Time-Based Sales Analytics Module: Rp 30.000.000
- API Integration (Maps, Optimization): Rp 15.000.000
- Testing & QA: Rp 20.000.000

**Frontend Development:**
- Leaderboard UI: Rp 25.000.000
- Route Optimization UI: Rp 30.000.000
- Schedule Assignment UI: Rp 28.000.000
- Area Mapping UI: Rp 25.000.000
- Product Analytics UI: Rp 20.000.000
- Sales Overview UI: Rp 25.000.000
- Sales Analytics UI: Rp 22.000.000
- UI/UX Enhancement: Rp 15.000.000

### 2. Third-Party Services & APIs (Rp 50.000.000 - 11%)

- Google Maps API / Mapbox API (1 year): Rp 20.000.000
- OR-Tools License (if needed): Rp 5.000.000
- Additional Cloud Services: Rp 10.000.000
- SSL Certificates: Rp 2.000.000
- Domain & Hosting (1 year): Rp 8.000.000
- Monitoring Tools: Rp 5.000.000

### 3. Infrastructure & Server (Rp 60.000.000 - 13%)

- Server Hardware / Cloud Setup: Rp 30.000.000
- Database Setup & Optimization: Rp 10.000.000
- PostGIS Extension Setup: Rp 5.000.000
- Backup & Recovery Setup: Rp 8.000.000
- Security Hardening: Rp 7.000.000

### 4. Documentation & Training (Rp 30.000.000 - 7%)

- Technical Documentation: Rp 10.000.000
- User Manual: Rp 8.000.000
- Video Tutorials: Rp 7.000.000
- Training Sessions: Rp 5.000.000

### 5. Testing & Quality Assurance (Rp 20.000.000 - 4%)

- Integration Testing: Rp 8.000.000
- Performance Testing: Rp 6.000.000
- Security Testing: Rp 4.000.000
- User Acceptance Testing: Rp 2.000.000

### 6. Deployment & Go-Live (Rp 10.000.000 - 2%)

- Production Deployment: Rp 5.000.000
- Data Migration: Rp 3.000.000
- Go-Live Support: Rp 2.000.000

---

## Implementation Timeline

### Phase 1: Foundation (Weeks 1-4)
- Database schema design and migration
- PostGIS setup
- Basic API structure
- Authentication and authorization

### Phase 2: Core Modules (Weeks 5-10)
- Leaderboard & KPI Module
- Product Analytics Module
- Sales Overview & Incentive Module
- Time-Based Sales Analytics Module

### Phase 3: Advanced Features (Weeks 11-16)
- Route Optimization Module
- Schedule Assignment Module
- Area Mapping Module

### Phase 4: Integration & Enhancement (Weeks 17-20)
- Integration with existing modules
- UI/UX enhancement
- Performance optimization
- Testing and bug fixes

### Phase 5: Deployment (Weeks 21-22)
- Server installation
- Production deployment
- Data migration
- Go-live support

### Phase 6: Documentation & Training (Weeks 23-24)
- Technical documentation
- User manual
- Video tutorials
- Training sessions

**Total Duration: 24 weeks (6 months)**

---

## Summary

TECHNICAL_FEATURESv2.md mencakup 8 fitur utama baru:

1. **Leaderboard & KPI Monitoring** - Performance tracking dan motivasi tim
2. **Route Optimization** - Efisiensi rute kunjungan
3. **Schedule Assignment** - Advanced task dan visit scheduling
4. **Area Mapping** - Geographic tracking dan territory management
5. **Top Products Analytics** - Product performance visualization
6. **Sales Overview & Incentive** - Performance monitoring dengan sistem insentif
7. **Time-Based Sales Analytics** - Analisis penjualan berdasarkan waktu
8. **Server Installation** - Production deployment setup

**Total Budget:** Rp 450.000.000
**Timeline:** 24 weeks (6 months)
**New Backend Modules:** 7 modules
**New Frontend Modules:** 7 modules
**Enhanced Modules:** 3 modules (Dashboard, Task Management, Visit Report)

Sistem ini akan memberikan kemampuan analitik yang lebih kuat, optimasi operasional, dan tools manajemen yang lebih canggih untuk meningkatkan produktivitas tim sales.

