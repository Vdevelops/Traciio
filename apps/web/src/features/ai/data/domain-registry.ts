import type { AIDomain, AIDomainModule, AIDataPrivacySettings } from "../types";

/**
 * Modular domain registry for the AI system.
 *
 * Each domain defines:
 * - Which privacy keys it needs
 * - Keywords for intent detection
 * - Supported CRUD entities
 * - A compact prompt fragment injected alongside the core system prompt
 *
 * This keeps token usage minimal: only the matched domain prompt is sent,
 * not the entire CRM schema.
 */
export const DOMAIN_MODULES: Record<AIDomain, AIDomainModule> = {
  route_optimization: {
    id: "route_optimization",
    label: "Route Optimization",
    description: "Optimize travel routes for field sales visits",
    requiredPrivacyKeys: ["allow_route_optimization", "allow_accounts", "allow_schedule"],
    intentKeywords: [
      "rute", "route", "optimize", "optimasi", "perjalanan", "travel",
      "jarak", "distance", "navigasi", "navigation", "kunjungan optimal",
      "visit route", "shortest path", "jalur", "efisien", "efficient",
      "lokasi", "location", "maps", "peta", "arah", "direction",
    ],
    supportedEntities: ["route", "schedule", "account"],
    promptFragment: `You are working in the Route Optimization module.
Capabilities: Create optimized visit routes, calculate distances, suggest visit sequences.
Entities: Routes, Schedules, Account locations.
When user requests route optimization:
1. Request or auto-detect user's current location.
2. If location access fails, immediately respond with the error detail and stop (do NOT consume more tokens).
3. Calculate optimal route considering distance, priority, and time windows.
4. Present route as numbered steps with estimated travel times.
CRUD: Create routes, Read existing schedules/locations, Update route plans.`,
  },

  sales: {
    id: "sales",
    label: "Sales",
    description: "Leads, Pipeline, Schedules, Visits, Tasks",
    requiredPrivacyKeys: [
      "allow_leads", "allow_deals", "allow_visit_reports",
      "allow_activities", "allow_tasks", "allow_schedule", "allow_pipelines",
    ],
    intentKeywords: [
      "lead", "leads", "pipeline", "deal", "deals", "opportunity",
      "visit", "kunjungan", "task", "tugas", "schedule", "jadwal",
      "follow up", "follow-up", "konversi", "conversion", "prospect",
      "prospek", "sales", "penjualan", "closed won", "closed lost",
      "stage", "tahap", "aktivitas", "activity", "report", "laporan kunjungan",
    ],
    supportedEntities: ["lead", "deal", "visit_report", "task", "schedule", "activity"],
    promptFragment: `You are working in the Sales module.
Entities: Leads, Pipeline/Deals, Schedules, Visit Reports, Tasks, Activities.
Capabilities:
- Leads: List, filter by status (new/contacted/qualified/converted/lost), analyze conversion rates, score leads.
- Pipeline: Show deals by stage, calculate values, forecast revenue.
- Schedules: View upcoming visits, suggest optimal scheduling.
- Visit Reports: Analyze reports, identify patterns, track approval status.
- Tasks: List pending/completed tasks, prioritize by urgency.
CRUD: Create leads/tasks/schedules, Read all sales data, Update lead status/deal stage.`,
  },

  inventory: {
    id: "inventory",
    label: "Inventory",
    description: "Products catalog and inventory management",
    requiredPrivacyKeys: ["allow_products"],
    intentKeywords: [
      "product", "produk", "inventory", "inventaris", "stok", "stock",
      "obat", "medicine", "pharmaceutical", "farmasi", "harga", "price",
      "katalog", "catalog", "barang", "item",
    ],
    supportedEntities: ["product"],
    promptFragment: `You are working in the Inventory module.
Entities: Products.
Capabilities:
- List products with filtering (category, price range, availability).
- Analyze product performance and demand trends.
- Identify top-selling and underperforming products.
CRUD: Read product data, analyze product metrics.`,
  },

  customers: {
    id: "customers",
    label: "Customers",
    description: "Accounts and contact management",
    requiredPrivacyKeys: ["allow_accounts", "allow_contacts"],
    intentKeywords: [
      "account", "akun", "customer", "pelanggan", "contact", "kontak",
      "rumah sakit", "hospital", "klinik", "clinic", "apotek", "pharmacy",
      "dokter", "doctor", "rs", "rsud", "faskes", "facility",
    ],
    supportedEntities: ["account", "contact"],
    promptFragment: `You are working in the Customers module.
Entities: Accounts (hospitals, clinics, pharmacies), Contacts (doctors, pharmacists, staff).
Capabilities:
- List and filter accounts by type/category/region.
- Show contacts linked to accounts with roles and relationships.
- Analyze account activity and engagement levels.
CRUD: Read account/contact data, analyze customer relationships.`,
  },

  analytics: {
    id: "analytics",
    label: "Analytics",
    description: "Sales performance, product analytics, and reports",
    requiredPrivacyKeys: ["allow_sales_performance", "allow_product_analysis", "allow_reports"],
    intentKeywords: [
      "analytics", "analitik", "performance", "performa", "report",
      "laporan", "statistik", "statistic", "dashboard", "chart",
      "grafik", "graph", "trend", "tren", "kpi", "metric",
      "revenue", "pendapatan", "target", "pencapaian", "achievement",
      "forecast", "prediksi", "prediction", "growth", "pertumbuhan",
    ],
    supportedEntities: ["sales_performance", "product_analytics", "report"],
    promptFragment: `You are working in the Analytics module.
Entities: Sales Performance, Product Analytics, Reports.
Capabilities:
- Sales Performance: Revenue trends, conversion rates, rep productivity, target vs actual.
- Product Analytics: Product demand analysis, category performance, pricing insights.
- Reports: Generate summary reports, export data insights.
Calculations must use real data only. Never fabricate numbers.
Present analytics with tables and clear formatting.`,
  },

  management: {
    id: "management",
    label: "Management",
    description: "Users, Roles, Groups, Bricks, Targets",
    requiredPrivacyKeys: [
      "allow_users", "allow_roles", "allow_groups",
      "allow_brick_management", "allow_target",
    ],
    intentKeywords: [
      "user", "pengguna", "role", "peran", "group", "grup",
      "brick", "wilayah", "territory",
      "target", "sasaran", "permission", "izin", "hak akses",
      "admin", "management", "manajemen", "assign", "struktur",
    ],
    supportedEntities: ["user", "role", "group", "brick", "target"],
    promptFragment: `You are working in the Management module.
Entities: Users, Roles, Groups, Bricks (territories), Targets.
Capabilities:
- Users: List users, view assignments, check activity.
- Roles: View role definitions and permissions.
- Groups: Show group structures and memberships.
- Bricks: Territory management, area assignments, coverage analysis.
- Targets: View sales targets, track progress, compare actuals.
CRUD: Read management data, analyze organizational structure.`,
  },

  general: {
    id: "general",
    label: "General",
    description: "General CRM assistance and multi-domain queries",
    requiredPrivacyKeys: [],
    intentKeywords: [],
    supportedEntities: [],
    promptFragment: `You are a general CRM assistant.
Help the user with any CRM-related question. If the query relates to a specific domain
(Sales, Inventory, Customers, Analytics, Management, Route Optimization),
identify the domain and respond accordingly.
For ambiguous queries, ask the user to clarify which module they need help with.`,
  },
};

/**
 * All available domain IDs excluding "general"
 */
export const SPECIFIC_DOMAINS: AIDomain[] = [
  "route_optimization",
  "sales",
  "inventory",
  "customers",
  "analytics",
  "management",
];

/**
 * Check if the required privacy keys for a domain are all enabled.
 */
export function isDomainAccessible(
  domain: AIDomain,
  privacy: AIDataPrivacySettings,
): boolean {
  const domainDef = DOMAIN_MODULES[domain];
  if (domainDef.requiredPrivacyKeys.length === 0) return true;

  return domainDef.requiredPrivacyKeys.some((key) => privacy[key]);
}

/**
 * Get the list of accessible domains based on current privacy settings.
 */
export function getAccessibleDomains(
  privacy: AIDataPrivacySettings,
): AIDomain[] {
  return SPECIFIC_DOMAINS.filter((d) => isDomainAccessible(d, privacy));
}

/**
 * Privacy settings grouped by domain category for the settings UI.
 */
export interface PrivacyGroup {
  domain: AIDomain;
  label: string;
  icon: string;
  keys: { key: keyof AIDataPrivacySettings; label: string; description: string }[];
}

export const PRIVACY_GROUPS: PrivacyGroup[] = [
  {
    domain: "route_optimization",
    label: "Route Optimization",
    icon: "MapPin",
    keys: [
      {
        key: "allow_route_optimization",
        label: "Route Optimization",
        description: "Allow AI to optimize visit routes and calculate distances",
      },
    ],
  },
  {
    domain: "sales",
    label: "Sales",
    icon: "TrendingUp",
    keys: [
      { key: "allow_leads", label: "Leads", description: "Allow AI to access lead data" },
      { key: "allow_pipelines", label: "Pipeline", description: "Allow AI to access pipeline data" },
      { key: "allow_deals", label: "Deals", description: "Allow AI to access deal/opportunity data" },
      { key: "allow_schedule", label: "Schedules", description: "Allow AI to access schedule data" },
      { key: "allow_visit_reports", label: "Visit Reports", description: "Allow AI to access visit report data" },
      { key: "allow_tasks", label: "Tasks", description: "Allow AI to access task data" },
      { key: "allow_activities", label: "Activities", description: "Allow AI to access activity data" },
    ],
  },
  {
    domain: "inventory",
    label: "Inventory",
    icon: "Package",
    keys: [
      { key: "allow_products", label: "Products", description: "Allow AI to access product catalog data" },
    ],
  },
  {
    domain: "customers",
    label: "Customers",
    icon: "Users",
    keys: [
      { key: "allow_accounts", label: "Accounts", description: "Allow AI to access account/customer data" },
      { key: "allow_contacts", label: "Contacts", description: "Allow AI to access contact data" },
    ],
  },
  {
    domain: "analytics",
    label: "Analytics",
    icon: "BarChart3",
    keys: [
      { key: "allow_sales_performance", label: "Sales Performance", description: "Allow AI to access sales performance analytics" },
      { key: "allow_product_analysis", label: "Product Analytics", description: "Allow AI to access product analysis data" },
      { key: "allow_reports", label: "Reports", description: "Allow AI to generate and access reports" },
    ],
  },
  {
    domain: "management",
    label: "Management",
    icon: "Settings",
    keys: [
      { key: "allow_users", label: "Users", description: "Allow AI to access user data" },
      { key: "allow_roles", label: "Roles", description: "Allow AI to access role and permission data" },
      { key: "allow_groups", label: "Groups", description: "Allow AI to access group structure data" },
      { key: "allow_brick_management", label: "Bricks", description: "Allow AI to access territory/brick data" },
      { key: "allow_target", label: "Targets", description: "Allow AI to access sales target data" },
    ],
  },
];
