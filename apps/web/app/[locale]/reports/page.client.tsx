"use client";

import { motion } from "framer-motion";
import { useTranslations } from "next-intl";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { ReportGenerator } from "@/features/reports/components/report-generator";

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
    },
  },
};

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.4 },
  },
};

function ReportsPageContent() {
  const t = useTranslations("reports");

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="space-y-6"
    >
      <motion.div variants={itemVariants}>
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-medium tracking-tight">{t("title")}</h1>
            <p className="text-muted-foreground mt-1">{t("description")}</p>
          </div>
        </div>
      </motion.div>

      <motion.div variants={itemVariants}>
        <ReportGenerator />
      </motion.div>
    </motion.div>
  );
}

export default function ReportsPageClient() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="reports.view">
        <ReportsPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
