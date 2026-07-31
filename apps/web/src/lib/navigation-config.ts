import { type MenuIcon } from "@/lib/menu-icons";

export interface NavigationItem {
  id: string; // Unique ID for keying
  name: string; // Translation key or label
  href: string;
  icon: MenuIcon; // Key for getMenuIcon
  permission?: string; // Required permission code (e.g., "dashboard.view")
  roles?: string[]; // Optional role allow-list in addition to permission
  children?: NavigationItem[];
}

export interface NavigationGroup {
  label: string; // Group label (e.g., "Main", "CRM")
  items: NavigationItem[];
}

export const NAVIGATION_CONFIG: NavigationGroup[] = [
  {
    label: "Main",
    items: [
      {
        id: "dashboard",
        name: "Dashboard",
        href: "/dashboard",
        icon: "layout-dashboard",
        permission: "dashboard.view",
      },
      {
        id: "route-optimization",
        name: "Route Optimization",
        href: "/route-optimization",
        icon: "map",
        permission: "route-optimization.view",
      },
    ],
  },
  {
    label: "Sales",
    items: [
      {
        id: "leads",
        name: "Leads",
        href: "/leads",
        icon: "users",
        permission: "leads.view",
      },
      {
        id: "pipeline",
        name: "Pipeline",
        href: "/pipeline",
        icon: "kanban",
        permission: "pipeline.view",
      },
      {
        id: "schedules",
        name: "Schedules",
        href: "/schedules",
        icon: "calendar",
        permission: "schedules.view",
      },
      {
        id: "visit-reports",
        name: "Visits",
        href: "/visit-reports",
        icon: "file-text",
        permission: "visit-reports.view",
      },
      {
        id: "tasks",
        name: "Tasks",
        href: "/tasks",
        icon: "check-square",
        permission: "tasks.view",
      },
    ],
  },
  {
    label: "Inventory",
    items: [
      {
        id: "products",
        name: "Products",
        href: "/products",
        icon: "package",
        permission: "products.view",
      },
    ],
  },
  {
    label: "Customers",
    items: [
      {
        id: "accounts",
        name: "Accounts",
        href: "/accounts",
        icon: "building",
        permission: "accounts.view",
      },
    ],
  },
  {
    label: "Analytics",
    items: [
      {
        id: "sales-overview",
        name: "Sales Performance",
        href: "/sales-overview",
        icon: "bar-chart-3",
        permission: "sales-overview.view",
      },
      {
        id: "kpi",
        name: "KPI",
        href: "/kpi",
        icon: "leaderboard",
        // Visible only to sales roles; backend enforces final authorization
        roles: ["sales_rep", "sales_manager"],
      },
      {
        id: "product-analytics",
        name: "Product Analytics",
        href: "/product-analytics",
        icon: "pie-chart",
        permission: "product-analytics.view",
      },

      {
        id: "reports",
        name: "Reports",
        href: "/reports",
        icon: "file-bar-chart",
        permission: "reports.view",
      },
    ],
  },
  {
    label: "Management",
    items: [
      {
        id: "users",
        name: "Users",
        href: "/master-data/users",
        icon: "users-2",
        permission: "users.view",
      },
      {
        id: "roles",
        name: "Roles",
        href: "/master-data/users?tab=roles", // Assuming roles are a tab in users or separate
        icon: "shield",
        permission: "users.roles",
      },
      {
        id: "groups",
        name: "Groups",
        href: "/master-data/groups",
        icon: "users-round",
        permission: "groups.view",
      },
      {
        id: "bricks",
        name: "Bricks",
        href: "/master-data/bricks",
        icon: "map-pin",
        permission: "bricks.view",
      },
      {
        id: "targets",
        name: "Targets",
        href: "/master-data/monthly-targets",
        icon: "target",
        permission: "monthly-targets.view",
      },
    ],
  },
  {
    label: "AI",
    items: [
      {
        id: "ai-chatbot",
        name: "Chatbot",
        href: "/ai-chatbot",
        icon: "bot",
        permission: "ai-chatbot.view",
      },
    ],
  },
];
