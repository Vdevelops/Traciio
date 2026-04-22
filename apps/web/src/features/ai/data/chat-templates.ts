import type { AIDomain } from "../types";

export interface ChatTemplate {
  id: string;
  name: string;
  category: string;
  content: string;
  description?: string;
  /** Domain this template belongs to for filtering */
  domain?: AIDomain;
}

export const chatTemplates: ChatTemplate[] = [
  // ──────────────────────────────────────────
  // Route Optimization
  // ──────────────────────────────────────────
  {
    id: "route-1",
    name: "Optimasi rute kunjungan hari ini",
    category: "Route Optimization",
    content: "Buatkan rute optimal untuk kunjungan sales hari ini berdasarkan jadwal yang ada",
    description: "Optimize visit route based on today's schedule",
    domain: "route_optimization",
  },
  {
    id: "route-2",
    name: "Rute ke beberapa lokasi",
    category: "Route Optimization",
    content: "Buatkan saya rute yang optimize untuk pergi ke RSUD Jakarta, Klinik Sehat, dan Apotek Farma",
    description: "Multi-stop route optimization with shortest path",
    domain: "route_optimization",
  },
  {
    id: "route-3",
    name: "Rute terdekat dari lokasi saya",
    category: "Route Optimization",
    content: "Dari lokasi saya sekarang, akun mana yang paling dekat untuk dikunjungi?",
    description: "Find nearest accounts from current location",
    domain: "route_optimization",
  },

  // ──────────────────────────────────────────
  // Sales - Leads
  // ──────────────────────────────────────────
  {
    id: "sales-lead-1",
    name: "Tampilkan semua leads",
    category: "Sales - Leads",
    content: "Tampilkan semua leads yang ada di sistem",
    description: "List all leads with status breakdown",
    domain: "sales",
  },
  {
    id: "sales-lead-2",
    name: "Leads yang qualified",
    category: "Sales - Leads",
    content: "Tampilkan semua leads yang sudah qualified",
    description: "Filter leads with qualified status",
    domain: "sales",
  },
  {
    id: "sales-lead-3",
    name: "Lead conversion rate",
    category: "Sales - Leads",
    content: "Berapa conversion rate dari qualified leads ke converted?",
    description: "Calculate lead conversion metrics",
    domain: "sales",
  },
  {
    id: "sales-lead-4",
    name: "Leads per source",
    category: "Sales - Leads",
    content: "Breakdown leads berdasarkan lead source (website, referral, cold_call, event, dll)",
    description: "Group leads by source with counts",
    domain: "sales",
  },
  {
    id: "sales-lead-5",
    name: "Leads yang perlu di-follow up",
    category: "Sales - Leads",
    content: "Leads mana yang perlu segera di-follow up? (leads yang sudah contacted tapi belum qualified)",
    description: "Identify stale leads needing follow-up",
    domain: "sales",
  },
  {
    id: "sales-lead-6",
    name: "Leads dengan score tinggi",
    category: "Sales - Leads",
    content: "Tampilkan leads dengan lead score tinggi (>= 70)",
    description: "High-priority leads by score",
    domain: "sales",
  },

  // ──────────────────────────────────────────
  // Sales - Pipeline
  // ──────────────────────────────────────────
  {
    id: "sales-pipeline-1",
    name: "Jumlah deals di pipeline",
    category: "Sales - Pipeline",
    content: "Berapa banyak deals yang ada di pipeline?",
    description: "Pipeline summary with deal counts per stage",
    domain: "sales",
  },
  {
    id: "sales-pipeline-2",
    name: "Breakdown deals per stage",
    category: "Sales - Pipeline",
    content: "Berikan breakdown deals per stage dengan total value dan count",
    description: "Pipeline breakdown with values",
    domain: "sales",
  },
  {
    id: "sales-pipeline-3",
    name: "Deals perlu follow-up",
    category: "Sales - Pipeline",
    content: "Deals mana yang perlu segera di-follow up berdasarkan last interaction date?",
    description: "Identify stale deals needing attention",
    domain: "sales",
  },
  {
    id: "sales-pipeline-4",
    name: "Deals tanpa visit report",
    category: "Sales - Pipeline",
    content: "Deals mana yang belum ada visit report dalam 30 hari terakhir?",
    description: "Gaps between deals and visit reports",
    domain: "sales",
  },

  // ──────────────────────────────────────────
  // Sales - Schedules & Visits
  // ──────────────────────────────────────────
  {
    id: "sales-schedule-1",
    name: "Jadwal kunjungan minggu ini",
    category: "Sales - Schedules",
    content: "Tampilkan jadwal kunjungan saya untuk minggu ini",
    description: "Upcoming visit schedule",
    domain: "sales",
  },
  {
    id: "sales-visit-1",
    name: "Visit reports yang approved",
    category: "Sales - Visits",
    content: "Tampilkan semua visit reports yang sudah approved",
    description: "Approved visit reports",
    domain: "sales",
  },
  {
    id: "sales-visit-2",
    name: "Trend visit reports",
    category: "Sales - Visits",
    content: "Bagaimana trend visit reports dalam 30 hari terakhir? Apakah meningkat atau menurun?",
    description: "Visit report trend analysis",
    domain: "sales",
  },

  // ──────────────────────────────────────────
  // Sales - Tasks
  // ──────────────────────────────────────────
  {
    id: "sales-task-1",
    name: "Tasks pending",
    category: "Sales - Tasks",
    content: "Tampilkan semua tasks yang masih pending",
    description: "Pending tasks list",
    domain: "sales",
  },
  {
    id: "sales-task-2",
    name: "Tasks overdue",
    category: "Sales - Tasks",
    content: "Tasks mana yang sudah overdue dan perlu segera diselesaikan?",
    description: "Overdue tasks identification",
    domain: "sales",
  },

  // ──────────────────────────────────────────
  // Inventory - Products
  // ──────────────────────────────────────────
  {
    id: "inv-1",
    name: "Tampilkan semua produk",
    category: "Inventory",
    content: "Tampilkan semua produk yang ada di katalog",
    description: "Full product catalog listing",
    domain: "inventory",
  },
  {
    id: "inv-2",
    name: "Produk terlaris",
    category: "Inventory",
    content: "Produk mana yang paling banyak terjual bulan ini?",
    description: "Top-selling products analysis",
    domain: "inventory",
  },
  {
    id: "inv-3",
    name: "Analisis harga produk",
    category: "Inventory",
    content: "Tampilkan analisis harga produk per kategori",
    description: "Product pricing analysis by category",
    domain: "inventory",
  },

  // ──────────────────────────────────────────
  // Customers - Accounts
  // ──────────────────────────────────────────
  {
    id: "cust-1",
    name: "Tampilkan semua akun",
    category: "Customers",
    content: "Tampilkan semua akun yang ada di sistem",
    description: "List all accounts (hospitals, clinics, pharmacies)",
    domain: "customers",
  },
  {
    id: "cust-2",
    name: "Akun paling aktif",
    category: "Customers",
    content: "Berdasarkan visit reports yang sudah approved, akun mana yang paling aktif?",
    description: "Most active accounts by visit count",
    domain: "customers",
  },
  {
    id: "cust-3",
    name: "Akun tanpa contact",
    category: "Customers",
    content: "Akun mana yang belum punya contact sama sekali?",
    description: "Accounts without contacts",
    domain: "customers",
  },
  {
    id: "cust-4",
    name: "Kontak RSUD Jakarta",
    category: "Customers",
    content: "Siapa saja kontak yang terhubung dengan RSUD Jakarta?",
    description: "Contacts linked to a specific account",
    domain: "customers",
  },
  {
    id: "cust-5",
    name: "History interaksi contact",
    category: "Customers",
    content: "Contact Dr. Budi di RSUD Jakarta, apa history interaksi dan deals yang terkait?",
    description: "Contact interaction history with deals",
    domain: "customers",
  },

  // ──────────────────────────────────────────
  // Analytics
  // ──────────────────────────────────────────
  {
    id: "analytics-1",
    name: "Ringkasan dashboard",
    category: "Analytics",
    content: "Buat ringkasan lengkap: total akun, total deals, total visit reports, total revenue dari closed won",
    description: "Dashboard summary with KPIs",
    domain: "analytics",
  },
  {
    id: "analytics-2",
    name: "Performance summary",
    category: "Analytics",
    content: "Ringkasan performance: total deals, total value, conversion rate, average deal size",
    description: "Sales performance KPIs",
    domain: "analytics",
  },
  {
    id: "analytics-3",
    name: "Forecast revenue bulan depan",
    category: "Analytics",
    content: "Berdasarkan data deals yang ada, berapa forecast revenue untuk bulan depan?",
    description: "Revenue forecast from pipeline data",
    domain: "analytics",
  },
  {
    id: "analytics-4",
    name: "Sales rep paling produktif",
    category: "Analytics",
    content: "Sales rep mana yang paling produktif berdasarkan jumlah visit reports?",
    description: "Rep productivity ranking",
    domain: "analytics",
  },
  {
    id: "analytics-5",
    name: "Rata-rata nilai deal per kategori",
    category: "Analytics",
    content: "Berapa rata-rata nilai deal per kategori akun (RS, Klinik, Apotek)?",
    description: "Average deal value by account type",
    domain: "analytics",
  },
  {
    id: "analytics-6",
    name: "Analisis produk terlaris",
    category: "Analytics",
    content: "Analisis produk mana yang memiliki performa penjualan terbaik dalam 3 bulan terakhir",
    description: "Product performance analysis",
    domain: "analytics",
  },

  // ──────────────────────────────────────────
  // Management
  // ──────────────────────────────────────────
  {
    id: "mgmt-1",
    name: "Daftar semua users",
    category: "Management",
    content: "Tampilkan daftar semua users beserta role dan divisi mereka",
    description: "User listing with roles and divisions",
    domain: "management",
  },
  {
    id: "mgmt-2",
    name: "Struktur organisasi",
    category: "Management",
    content: "Tampilkan struktur organisasi: divisi, grup, dan jumlah anggota per grup",
    description: "Organizational structure overview",
    domain: "management",
  },
  {
    id: "mgmt-3",
    name: "Target vs Aktual",
    category: "Management",
    content: "Bandingkan target sales dengan pencapaian aktual per divisi",
    description: "Target vs actual comparison by division",
    domain: "management",
  },
  {
    id: "mgmt-4",
    name: "Coverage area brick",
    category: "Management",
    content: "Tampilkan coverage area brick dan sales rep yang bertanggung jawab",
    description: "Territory coverage with assigned reps",
    domain: "management",
  },
  {
    id: "mgmt-5",
    name: "Roles dan permissions",
    category: "Management",
    content: "Tampilkan semua roles yang ada beserta permissions-nya",
    description: "Role definitions and permissions",
    domain: "management",
  },

  // ──────────────────────────────────────────
  // Cross-Domain / Strategy
  // ──────────────────────────────────────────
  {
    id: "strategy-1",
    name: "Prioritas akun untuk kunjungan",
    category: "Sales Strategy",
    content: "Berdasarkan data deals dan visit reports, akun mana yang harus diprioritaskan untuk kunjungan berikutnya?",
    description: "Priority scoring based on multiple factors",
    domain: "sales",
  },
  {
    id: "strategy-2",
    name: "Strategi meningkatkan conversion",
    category: "Sales Strategy",
    content: "Untuk meningkatkan conversion rate, rekomendasikan strategi berdasarkan analisis deals yang won vs lost",
    description: "Win/loss analysis with recommendations",
    domain: "analytics",
  },
  {
    id: "strategy-3",
    name: "Analisis lengkap akun",
    category: "Sales Strategy",
    content: "Untuk akun RSUD Jakarta, tampilkan semua visit reports, deals, dan contacts yang terkait",
    description: "Multi-entity analysis for one account",
    domain: "customers",
  },
];

export const templateCategories = [
  "All",
  "Route Optimization",
  "Sales - Leads",
  "Sales - Pipeline",
  "Sales - Schedules",
  "Sales - Visits",
  "Sales - Tasks",
  "Inventory",
  "Customers",
  "Analytics",
  "Management",
  "Sales Strategy",
] as const;

/**
 * Get templates filtered by category name.
 */
export function getTemplatesByCategory(category: string): ChatTemplate[] {
  if (category === "All") {
    return chatTemplates;
  }
  return chatTemplates.filter((template) => template.category === category);
}

/**
 * Get templates filtered by AI domain.
 * Returns only templates relevant to the given domain for token-efficient suggestions.
 */
export function getTemplatesByDomain(domain: AIDomain): ChatTemplate[] {
  if (domain === "general") {
    return chatTemplates;
  }
  return chatTemplates.filter((template) => template.domain === domain);
}

