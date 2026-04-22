"use client";

import { useState, useEffect, useRef } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import {
  createBrickSchema,
  updateBrickSchema,
  type CreateBrickFormData,
  type UpdateBrickFormData,
} from "../schemas/brick.schema";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Map, FileText } from "lucide-react";
import type { Brick } from "../types";
import { useUsers } from "@/features/master-data/user-management/hooks/useUsers";
import { useBricks } from "../hooks/useBricks";
import { BrickMapSelector } from "./brick-map-selector";

interface BrickFormProps {
  readonly brick?: Brick;
  readonly onSubmit: (
    data: CreateBrickFormData | UpdateBrickFormData
  ) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
  /** Pre-fill regency when creating from map click */
  readonly prefillRegency?: string;
  /** Pre-fill province when creating from map click */
  readonly prefillProvince?: string;
}

type BrickFormData = CreateBrickFormData | UpdateBrickFormData;

// Type for form errors that includes province and regency
type BrickFormErrors = {
  province?: { message?: string };
  regency?: { message?: string };
  name?: { message?: string };
  code?: { message?: string };
  description?: { message?: string };
  manager_id?: { message?: string };
  status?: { message?: string };
};

export function BrickForm({
  brick,
  onSubmit,
  onCancel,
  isLoading,
  prefillRegency,
  prefillProvince,
}: BrickFormProps) {
  const isEdit = !!brick;
  const hasPrefill = !!(prefillRegency && prefillProvince);
  const t = useTranslations("brickManagement.form");
  // Default to form mode for edit, map mode for create (unless prefilled from map click)
  const [mapMode, setMapMode] = useState(!isEdit && !hasPrefill);

  // Get users for manager selection (filter by sales_manager role if needed)
  // Only fetch active users to reduce data
  const { data: usersData, isLoading: isLoadingUsers } = useUsers({ 
    per_page: 100,
    status: "active", // Only get active users
  });
  
  // Filter managers: sales_manager or admin roles
  const managers = usersData?.data?.filter((user) => {
    const roleCode = user.role?.code?.toLowerCase();
    return roleCode === "sales_manager" || roleCode === "admin";
  }) ?? [];

  // Get existing bricks for map visualization (exclude current brick if editing)
  const { data: bricksData } = useBricks({ per_page: 100 });
  const existingBricks = bricksData?.data
    ?.filter((b) => !isEdit || b.id !== brick?.id) // Exclude current brick when editing
    .map((b) => ({
      regency: b.regency,
      province: b.province,
    })) ?? [];

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<BrickFormData>({
    resolver: zodResolver(isEdit ? updateBrickSchema : createBrickSchema),
    defaultValues: brick
      ? {
          name: brick.name,
          code: brick.code,
          description: brick.description || "",
          province: brick.province,
          regency: brick.regency,
          manager_id: brick.manager_id,
          status: brick.status,
        }
      : {
          status: "active",
          province: prefillProvince || "",
          regency: prefillRegency || "",
          name: hasPrefill ? `${prefillRegency}, ${prefillProvince}` : "",
          code: "",
        },
  });

  const managerId = watch("manager_id");
  const regencyValue = watch("regency");
  const provinceValue = watch("province");

  // Track previous province to detect changes (not initial load)
  const prevProvinceRef = useRef<string | undefined>(provinceValue);
  const isUpdatingFromMapRef = useRef(false); // Flag to track if update is from map selector
  
  // Reset regency when province changes (to ensure consistency)
  // BUT: Don't reset if update is from map selector (map selector handles both province and regency together)
  useEffect(() => {
    // Skip if update is from map selector
    if (isUpdatingFromMapRef.current) {
      isUpdatingFromMapRef.current = false; // Reset flag
      prevProvinceRef.current = provinceValue; // Update ref
      return;
    }
    
    // Only reset regency if province actually changed (not initial load)
    // Check if we have a previous value and it's different from current
    if (prevProvinceRef.current && prevProvinceRef.current !== provinceValue && regencyValue) {
      // Province changed manually (not from map), reset regency to ensure consistency
      setValue("regency", "", { shouldValidate: false });
    }
    // Update previous province value
    prevProvinceRef.current = provinceValue;
  }, [provinceValue, regencyValue, setValue]);

  const handleRegencySelectFromMap = (regency: { name: string; province: string; district?: string }) => {
    // Set flag to prevent useEffect from resetting regency
    isUpdatingFromMapRef.current = true;
    
    // Update prevProvinceRef first to prevent useEffect from resetting regency
    prevProvinceRef.current = regency.province;
    
    // Always update regency and province when selecting from map (both create and edit mode)
    // Use shouldDirty: true to mark field as changed, shouldTouch: true to mark as touched
    // Update province first, then regency to ensure consistency
    setValue("province", regency.province, { shouldValidate: true, shouldDirty: true, shouldTouch: true });
    setValue("regency", regency.name, { shouldValidate: true, shouldDirty: true, shouldTouch: true });
    
    // Update name when selecting from map (both create and edit mode)
    setValue("name", `${regency.name}, ${regency.province}`, { shouldValidate: false, shouldDirty: true, shouldTouch: true });
  };

  const handleFormSubmit = async (data: BrickFormData) => {
    await onSubmit(data);
  };

  // Show map in both create and edit mode
  const showMapTab = true;

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      {showMapTab && (
        <Tabs value={mapMode ? "map" : "form"} onValueChange={(v) => setMapMode(v === "map")}>
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="map" className="flex items-center gap-2">
              <Map className="h-4 w-4" />
              Map Selection
            </TabsTrigger>
            <TabsTrigger value="form" className="flex items-center gap-2">
              <FileText className="h-4 w-4" />
              Manual Input
            </TabsTrigger>
          </TabsList>
          <TabsContent value="map" className="space-y-4">
            {mapMode && (
              <BrickMapSelector
                selectedRegency={
                  regencyValue && provinceValue
                    ? { name: regencyValue, province: provinceValue }
                    : undefined
                }
                onRegencySelect={handleRegencySelectFromMap}
                existingBricks={existingBricks}
              />
            )}
          </TabsContent>
          <TabsContent value="form" className="space-y-4">
            {/* Form fields will continue below */}
          </TabsContent>
        </Tabs>
      )}

      <Field orientation="vertical">
        <FieldLabel>{t("nameLabel")}</FieldLabel>
        <Input
          {...register("name")}
          placeholder={t("namePlaceholder")}
          disabled={isLoading}
        />
        {errors.name && <FieldError>{errors.name.message}</FieldError>}
      </Field>

{/* Code field only shown in edit mode — auto-generated by server on create */}
      {isEdit && (
        <Field orientation="vertical">
          <FieldLabel>{t("codeLabel")}</FieldLabel>
          <Input
            {...register("code")}
            placeholder={t("codePlaceholder")}
            disabled={isLoading}
          />
          {errors.code && <FieldError>{errors.code.message}</FieldError>}
        </Field>
      )}

      {(!showMapTab || !mapMode) && (
        <>
          <Field orientation="vertical">
            <FieldLabel>{t("provinceLabel")}</FieldLabel>
            <Input
              {...register("province")}
              placeholder={t("provincePlaceholder")}
              disabled={isLoading || (isEdit && !mapMode)}
              readOnly={isEdit && mapMode}
            />
            {(errors as BrickFormErrors).province && <FieldError>{(errors as BrickFormErrors).province?.message}</FieldError>}
          </Field>

          <Field orientation="vertical">
            <FieldLabel>{t("regencyLabel")}</FieldLabel>
            <Input
              {...register("regency")}
              placeholder={t("regencyPlaceholder")}
              disabled={isLoading || (isEdit && !mapMode)}
              readOnly={isEdit && mapMode}
            />
            {(errors as BrickFormErrors).regency && <FieldError>{(errors as BrickFormErrors).regency?.message}</FieldError>}
          </Field>
        </>
      )}
      
      {/* Show province and regency values when in map mode (even in edit mode) */}
      {showMapTab && mapMode && (
        <div className="space-y-2 p-4 bg-muted/50 rounded-lg border border-border/50">
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium text-muted-foreground">{t("provinceLabel")}:</span>
            <span className="text-sm font-medium">{provinceValue || "-"}</span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium text-muted-foreground">{t("regencyLabel")}:</span>
            <span className="text-sm font-medium">{regencyValue || "-"}</span>
          </div>
        </div>
      )}

      <Field orientation="vertical">
        <FieldLabel>{t("descriptionLabel")}</FieldLabel>
        <Textarea
          {...register("description")}
          placeholder={t("descriptionPlaceholder")}
          disabled={isLoading}
          rows={3}
        />
        {errors.description && (
          <FieldError>{errors.description.message}</FieldError>
        )}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("managerLabel")}</FieldLabel>
        <Select
          value={managerId || "none"}
          onValueChange={(value) => setValue("manager_id", value === "none" ? undefined : value)}
          disabled={isLoading || isLoadingUsers}
        >
          <SelectTrigger>
            <SelectValue placeholder={isLoadingUsers ? "Loading managers..." : t("managerPlaceholder")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="none">{t("noManager")}</SelectItem>
            {(() => {
              if (isLoadingUsers) {
                return <div className="px-2 py-1.5 text-sm text-muted-foreground">Loading...</div>;
              }
              if (managers.length === 0) {
                return <div className="px-2 py-1.5 text-sm text-muted-foreground">No managers available</div>;
              }
              return managers.map((manager) => (
                <SelectItem key={manager.id} value={manager.id}>
                  {manager.name} ({manager.email})
                </SelectItem>
              ));
            })()}
          </SelectContent>
        </Select>
        {errors.manager_id && (
          <FieldError>{errors.manager_id.message}</FieldError>
        )}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("statusLabel")}</FieldLabel>
        <Select
          value={watch("status") || "active"}
          onValueChange={(value) =>
            setValue("status", value as "active" | "inactive")
          }
          disabled={isLoading}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="active">{t("active")}</SelectItem>
            <SelectItem value="inactive">{t("inactive")}</SelectItem>
          </SelectContent>
        </Select>
        {errors.status && <FieldError>{errors.status.message}</FieldError>}
      </Field>

      <div className="flex justify-end gap-3 pt-4">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isLoading}
        >
          {t("cancel")}
        </Button>
        <Button type="submit" disabled={isLoading}>
          {(() => {
            if (isLoading) return t("saving");
            return isEdit ? t("update") : t("create");
          })()}
        </Button>
      </div>
    </form>
  );
}

