import { z } from "zod";

export const analyzeVisitReportSchema = z.object({
  visit_report_id: z.string().uuid("Invalid visit report ID"),
});

export type AnalyzeVisitReportFormData = z.infer<
  typeof analyzeVisitReportSchema
>;

export const chatMessageSchema = z.object({
  role: z.enum(["user", "assistant"]),
  content: z.string(),
});

export const domainSchema = z.enum([
  "route_optimization",
  "sales",
  "inventory",
  "customers",
  "analytics",
  "management",
  "general",
]);

export const chatSchema = z.object({
  message: z.string().min(1, "Message is required"),
  context: z.string().uuid().optional(),
  context_type: z
    .enum(["visit_report", "deal", "contact", "account", "lead"])
    .optional(),
  conversation_history: z.array(chatMessageSchema).optional(),
  domain: domainSchema.optional(),
});

export type ChatFormData = z.infer<typeof chatSchema>;

export const aiDataPrivacySchema = z.object({
  // Sales domain
  allow_leads: z.boolean(),
  allow_deals: z.boolean(),
  allow_visit_reports: z.boolean(),
  allow_activities: z.boolean(),
  allow_tasks: z.boolean(),
  allow_schedule: z.boolean(),
  allow_pipelines: z.boolean(),
  // Customer domain
  allow_accounts: z.boolean(),
  allow_contacts: z.boolean(),
  // Inventory domain
  allow_products: z.boolean(),
  // Analytics domain
  allow_sales_performance: z.boolean(),
  allow_product_analysis: z.boolean(),
  allow_reports: z.boolean(),
  // Management domain
  allow_users: z.boolean(),
  allow_roles: z.boolean(),
  allow_groups: z.boolean(),
  allow_brick_management: z.boolean(),
  allow_target: z.boolean(),
  // Route Optimization domain
  allow_route_optimization: z.boolean(),
});

export const aiSettingsSchema = z.object({
  enabled: z.boolean(),
  data_privacy: aiDataPrivacySchema,
});

export type AISettingsFormData = z.infer<typeof aiSettingsSchema>;

