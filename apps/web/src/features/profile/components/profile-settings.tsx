"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import React, { useState } from "react";
import { Mail } from "lucide-react";
import { useProfile, useUpdateProfile, useChangePassword } from "../hooks/useProfile";
import { useSalesRepCheckInLocations } from "@/features/sales-overview/hooks/useSalesRepCheckInLocations";
import { updateProfileSchema, changePasswordSchema, type UpdateProfileFormData, type ChangePasswordFormData } from "../schemas/profile.schema";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { formatEmailToMailto } from "@/lib/utils";
import { UserStatistics } from "./UserStatistics";
import { UserDetailTabs } from "./UserDetailTabs";
import type { DateRange } from "react-day-picker";

// Get date range for current week
function getWeekDateRange(): DateRange {
  const now = new Date();
  const dayOfWeek = now.getDay();
  const diffToMonday = dayOfWeek === 0 ? -6 : 1 - dayOfWeek;
  const monday = new Date(now);
  monday.setDate(now.getDate() + diffToMonday);
  monday.setHours(0, 0, 0, 0);

  const sunday = new Date(monday);
  sunday.setDate(monday.getDate() + 6);
  sunday.setHours(23, 59, 59, 999);

  return {
    from: monday,
    to: sunday,
  };
}

function formatDateForAPI(date: Date | undefined): string {
  if (!date) return "";
  // Use local date components instead of toISOString() to avoid timezone conversion
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

export function ProfileSettings() {
  const { user } = useAuthStore();
  const [dateRange, setDateRange] = useState<DateRange | undefined>(getWeekDateRange());
  const [activeTab, setActiveTab] = useState("overview");

  const { data: profileData, isLoading: isLoadingProfile } = useProfile({
    start_date: formatDateForAPI(dateRange?.from),
    end_date: formatDateForAPI(dateRange?.to),
  });
  const updateProfile = useUpdateProfile();
  const changePassword = useChangePassword();
  const t = useTranslations("profile");

  const profile = profileData?.data?.user || user;
  const stats = profileData?.data?.stats;

  // Fetch check-in locations with date range
  const checkInLocations = useSalesRepCheckInLocations(user?.id ?? "", {
    start_date: formatDateForAPI(dateRange?.from),
    end_date: formatDateForAPI(dateRange?.to),
    page: 1,
    per_page: 50,
  });

  // Profile form
  const {
    register: registerProfile,
    handleSubmit: handleSubmitProfile,
    reset: resetProfile,
    formState: { errors: profileErrors },
  } = useForm<UpdateProfileFormData>({
    resolver: zodResolver(updateProfileSchema),
    defaultValues: {
      name: profile?.name || "",
    },
  });

  // Update form when profile data loads
  React.useEffect(() => {
    if (profile?.name) {
      resetProfile({ name: profile.name });
    }
  }, [profile?.name, resetProfile]);

  // Password form
  const {
    register: registerPassword,
    handleSubmit: handleSubmitPassword,
    reset: resetPassword,
    formState: { errors: passwordErrors },
  } = useForm<ChangePasswordFormData>({
    resolver: zodResolver(changePasswordSchema),
  });

  const handleUpdateProfile = async (data: UpdateProfileFormData) => {
    await updateProfile.mutateAsync(data);
  };

  const handleChangePassword = async (data: ChangePasswordFormData) => {
    await changePassword.mutateAsync(data);
    resetPassword();
  };

  const fallbackAvatarUrl = "/avatar-placeholder.svg";
  const currentAvatarUrl =
    profile?.avatar_url && profile.avatar_url.trim() !== ""
      ? profile.avatar_url
      : fallbackAvatarUrl;

  if (isLoadingProfile) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-64 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="overview" className="cursor-pointer">{t("tabs.overview")}</TabsTrigger>
          <TabsTrigger value="password" className="cursor-pointer">{t("tabs.password")}</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="mt-6">
          <div className="grid gap-6 lg:grid-cols-3">
            {/* Left Sidebar - User Profile (1 column) */}
            <div className="lg:col-span-1 space-y-6">
              {/* User Card */}
              <Card>
                <CardContent className="pt-6">
                  <div className="flex flex-col items-center text-center space-y-4">
                    <Avatar className="h-24 w-24">
                      <AvatarImage src={currentAvatarUrl} alt={profile?.name || "User"} />
                    </Avatar>
                    <div>
                      <div className="flex items-center justify-center gap-2">
                        <h3 className="text-lg font-medium">{profile?.name || "User"}</h3>
                      </div>
                      <p className="text-sm text-muted-foreground mt-1">
                        {(() => {
                          if (typeof profile?.role === "object" && profile.role?.name) {
                            return profile.role.name;
                          }
                          if (typeof profile?.role === "string") {
                            return profile.role;
                          }
                          return "User";
                        })()}
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Contact Information */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-base">{t("contact.title")}</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex items-center gap-3 text-sm">
                    <Mail className="h-4 w-4 text-muted-foreground" />
                    <a href={formatEmailToMailto(profile?.email)} className="text-muted-foreground hover:text-primary hover:underline cursor-pointer min-w-0">{profile?.email || "-"}</a>
                  </div>
                </CardContent>
              </Card>

            </div>

            {/* Right Content - Performance Stats & Tabs (2 columns) */}
            <div className="lg:col-span-2 space-y-6">
              {/* Date Range Filter Header */}
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold">{t("title")}</h3>
                </div>
                <DateRangePicker dateRange={dateRange} onDateChange={setDateRange} />
              </div>

              {/* Performance Stats */}
              <UserStatistics statistics={stats} />

              {/* Activity Details Tabs */}
              <Card>
                <CardHeader>
                  <CardTitle>Activity Details</CardTitle>
                  <CardDescription>
                    {formatDateForAPI(dateRange?.from)} - {formatDateForAPI(dateRange?.to)}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <UserDetailTabs
                    userId={user?.id ?? ""}
                    startDate={formatDateForAPI(dateRange?.from)}
                    endDate={formatDateForAPI(dateRange?.to)}
                    checkInLocationsProps={{
                      locations: checkInLocations.locations,
                      isLoading: checkInLocations.isLoading,
                      totalVisits: checkInLocations.totalVisits,
                      page: checkInLocations.page,
                      perPage: checkInLocations.perPage,
                      onPageChange: checkInLocations.setPage,
                      onPerPageChange: checkInLocations.setPerPage,
                    }}
                  />
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="password" className="mt-6">
          <div className="grid gap-6 lg:grid-cols-3">
            {/* Empty left column for layout consistency */}
            <div className="hidden lg:block lg:col-span-1" />
            
            {/* Password Change Form */}
            <div className="lg:col-span-2">
              <Card>
                <CardHeader>
                  <CardTitle>{t("password.title")}</CardTitle>
                  <CardDescription>{t("password.description")}</CardDescription>
                </CardHeader>
                <CardContent>
                  <form onSubmit={handleSubmitPassword(handleChangePassword)} className="space-y-6">
                    <Field>
                      <FieldLabel>{t("password.currentPasswordLabel")}</FieldLabel>
                      <Input
                        {...registerPassword("current_password")}
                        type="password"
                        placeholder={t("password.currentPasswordPlaceholder")}
                      />
                      {passwordErrors.current_password && (
                        <FieldError>{passwordErrors.current_password.message}</FieldError>
                      )}
                    </Field>

                    <Field>
                      <FieldLabel>{t("password.newPasswordLabel")}</FieldLabel>
                      <Input
                        {...registerPassword("password")}
                        type="password"
                        placeholder={t("password.newPasswordPlaceholder")}
                      />
                      {passwordErrors.password && (
                        <FieldError>{passwordErrors.password.message}</FieldError>
                      )}
                    </Field>

                    <Field>
                      <FieldLabel>{t("password.confirmPasswordLabel")}</FieldLabel>
                      <Input
                        {...registerPassword("confirm_password")}
                        type="password"
                        placeholder={t("password.confirmPasswordPlaceholder")}
                      />
                      {passwordErrors.confirm_password && (
                        <FieldError>{passwordErrors.confirm_password.message}</FieldError>
                      )}
                    </Field>

                    <div className="flex gap-4">
                      <Button
                        type="submit"
                        disabled={changePassword.isPending}
                        className="cursor-pointer"
                      >
                        {changePassword.isPending ? t("password.changing") : t("password.change")}
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => resetPassword()}
                        className="cursor-pointer"
                      >
                        {t("password.cancel")}
                      </Button>
                    </div>
                  </form>
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
