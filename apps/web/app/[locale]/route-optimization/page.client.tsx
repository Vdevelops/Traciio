"use client";

import { useRef } from "react";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { 
  RouteOptimizationManagement, 
  type RouteOptimizationManagementRef 
} from "@/features/route-optimization/components/route-optimization-management";

function RouteOptimizationPageContent() {
  const managementRef = useRef<RouteOptimizationManagementRef>(null);

  return (
    <div className="absolute inset-0 overflow-hidden">
      <RouteOptimizationManagement ref={managementRef} />
    </div>
  );
}

export default function RouteOptimizationPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="route-optimization.view">
        <RouteOptimizationPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
