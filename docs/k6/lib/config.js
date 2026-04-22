// K6 Config - CRM Healthcare Load Testing
// Backend API: https://api.gilabs.id (Go + Gin)
// Frontend: https://crm-demo.gilabs.id (Next.js)

export const BASE_URL = __ENV.BASE_URL || "https://api.gilabs.id";
export const API_V1 = `${BASE_URL}/api/v1`;

// Default credentials (from seeders/auth_seeder.go)
// All passwords: admin123
export const USERS = {
  admin: { email: "admin@example.com", password: "admin123" },
  salesManager: { email: "salesmanager@example.com", password: "admin123" },
  sales: { email: "sales@example.com", password: "admin123" },
};

// Reusable HTTP params with Bearer token
export function authHeaders(token) {
  return {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    timeout: "30s",
  };
}

// Reusable HTTP params without auth
export function jsonHeaders() {
  return {
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    timeout: "30s",
  };
}

// =============================================
// Verified API Endpoints (from routes/*.go & main.go)
// =============================================
export const ENDPOINTS = {
  // Public (root level)
  health: `${BASE_URL}/health`,
  healthCache: `${BASE_URL}/health/cache`,
  healthCacheMetrics: `${BASE_URL}/health/cache/metrics`,
  healthCircuitBreakers: `${BASE_URL}/health/circuit-breakers`,
  healthRuntime: `${BASE_URL}/health/runtime`,
  ping: `${BASE_URL}/ping`,

  // Auth (public - /api/v1/auth/*)
  login: `${API_V1}/auth/login`,
  refresh: `${API_V1}/auth/refresh`,
  logout: `${API_V1}/auth/logout`,

  // Users (authenticated - /api/v1/users/*)
  myProfile: `${API_V1}/users/me/settings-summary`,
  users: `${API_V1}/users`,
  user: (id) => `${API_V1}/users/${id}`,

  // Leads (authenticated - /api/v1/leads/*)
  leads: `${API_V1}/leads`,
  lead: (id) => `${API_V1}/leads/${id}`,
  leadConvert: (id) => `${API_V1}/leads/${id}/convert`,
  leadCreateAccount: (id) => `${API_V1}/leads/${id}/create-account`,
  leadFormData: `${API_V1}/leads/form-data`,
  leadAnalytics: `${API_V1}/leads/analytics`,

  // Accounts (authenticated - /api/v1/accounts/*) — NOTE: not "companies"
  accounts: `${API_V1}/accounts`,
  account: (id) => `${API_V1}/accounts/${id}`,
  accountsMap: `${API_V1}/accounts/map`,

  // Contacts (authenticated - /api/v1/contacts/*)
  contacts: `${API_V1}/contacts`,
  contact: (id) => `${API_V1}/contacts/${id}`,

  // Pipeline & Deals (authenticated - /api/v1/pipelines/* & /api/v1/deals/*)
  pipelines: `${API_V1}/pipelines`,
  pipeline: (id) => `${API_V1}/pipelines/${id}`,
  pipelineSummary: `${API_V1}/pipelines/summary`,
  pipelineForecast: `${API_V1}/pipelines/forecast`,
  deals: `${API_V1}/deals`,
  deal: (id) => `${API_V1}/deals/${id}`,
  dealsByStage: `${API_V1}/deals/by-stage`,

  // Visit Reports (authenticated - /api/v1/visit-reports/*)
  visitReports: `${API_V1}/visit-reports`,
  visitReport: (id) => `${API_V1}/visit-reports/${id}`,

  // Activities (authenticated - /api/v1/activities/*)
  activities: `${API_V1}/activities`,
  activity: (id) => `${API_V1}/activities/${id}`,
  activityTimeline: `${API_V1}/activities/timeline`,

  // Tasks & Reminders (authenticated - /api/v1/tasks/*)
  tasks: `${API_V1}/tasks`,
  task: (id) => `${API_V1}/tasks/${id}`,
  taskReminders: `${API_V1}/tasks/reminders`,

  // Schedules (authenticated - /api/v1/schedules/*)
  schedules: `${API_V1}/schedules`,
  schedule: (id) => `${API_V1}/schedules/${id}`,

  // Products (authenticated - /api/v1/products/*)
  products: `${API_V1}/products`,
  product: (id) => `${API_V1}/products/${id}`,

  // Dashboard (authenticated - /api/v1/dashboard/*)
  dashboardOverview: `${API_V1}/dashboard/overview`,
  dashboardVisits: `${API_V1}/dashboard/visits`,
  dashboardActivityTrends: `${API_V1}/dashboard/activity-trends`,
  dashboardPipeline: `${API_V1}/dashboard/pipeline`,
  dashboardTopAccounts: `${API_V1}/dashboard/top-accounts`,
  dashboardTopSalesRep: `${API_V1}/dashboard/top-sales-rep`,
  dashboardRecentActivities: `${API_V1}/dashboard/recent-activities`,

  // Reports (authenticated - /api/v1/reports/*)
  reportVisitReports: `${API_V1}/reports/visit-reports`,
  reportPipeline: `${API_V1}/reports/pipeline`,
  reportSalesPerformance: `${API_V1}/reports/sales-performance`,
  reportAccountActivity: `${API_V1}/reports/account-activity`,

  // Sales Overview (authenticated - /api/v1/sales-overview/*)
  salesOverviewPerformance: `${API_V1}/sales-overview/performance`,
  salesOverviewMonthly: `${API_V1}/sales-overview/monthly-overview`,

  // Notifications (authenticated - /api/v1/notifications/*)
  notifications: `${API_V1}/notifications`,
  notificationsUnreadCount: `${API_V1}/notifications/unread-count`,

  // Master Data — Reference data endpoints
  roles: `${API_V1}/roles`,
  groups: `${API_V1}/groups`,
  divisions: `${API_V1}/divisions`,
  categories: `${API_V1}/categories`,
  contactRoles: `${API_V1}/contact-roles`,
  industries: `${API_V1}/industries`,
  leadSources: `${API_V1}/lead-sources`,
  leadStatuses: `${API_V1}/lead-statuses`,
  monthlyTargets: `${API_V1}/monthly-targets`,
  bricks: `${API_V1}/bricks`,
};
