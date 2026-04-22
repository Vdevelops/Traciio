"use client";

import { motion } from "framer-motion";
import { useTranslations } from "next-intl";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { ProductManagement } from "@/features/sales-crm/product-management/components/product-management";

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

function ProductsPageContent() {
  const t = useTranslations("productManagement.page");

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
        <ProductManagement />
      </motion.div>
    </motion.div>
  );
}

export default function ProductsPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="products.view">
        <ProductsPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
