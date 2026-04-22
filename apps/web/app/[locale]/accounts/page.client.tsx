"use client";

import { useSearchParams } from "next/navigation";
import { motion } from "framer-motion";
import { useTranslations } from "next-intl";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { AccountMapManagement } from "@/features/sales-crm/account-management/components/account-map-management";
import { AccountManagement } from "@/features/sales-crm/account-management/components/account-management";

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

function AccountsPageContent() {
  const searchParams = useSearchParams();
  const tab = searchParams.get("tab");
  const t = useTranslations("accountManagement.page");

  // Show traditional tabbed layout for categories/contact-roles sub-tabs
  if (tab === "categories" || tab === "contact-roles") {
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
              <h1 className="text-3xl font-medium tracking-tight">
                {t("title")}
              </h1>
              <p className="text-muted-foreground mt-1">
                {t("description")}
              </p>
            </div>
          </div>
        </motion.div>

        <motion.div variants={itemVariants}>
          <AccountManagement />
        </motion.div>
      </motion.div>
    );
  }

  // Default: Full-screen map view for accounts
  return (
    <div className="absolute inset-0 overflow-hidden">
      <AccountMapManagement />
    </div>
  );
}

export default function AccountsPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="accounts.view">
        <AccountsPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
