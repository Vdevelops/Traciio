"use client";

import { useState } from "react";
import { Search, MapPin, Calendar, Building2, Check } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Checkbox } from "@/components/ui/checkbox";
import { Badge } from "@/components/ui/badge";
import { useAccounts } from "@/features/sales-crm/account-management/hooks/useAccounts";
import { useVisitReports } from "@/features/sales-crm/visit-report/hooks/useVisitReports";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { toast } from "sonner";
import type { Waypoint } from "../types";
import type { Account } from "@/features/sales-crm/account-management/types";
import type { VisitReport } from "@/features/sales-crm/visit-report/types";

interface WaypointSelectorDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onSelect: (waypoints: Waypoint[]) => void;
}

export function WaypointSelectorDialog({
  open,
  onOpenChange,
  onSelect,
}: WaypointSelectorDialogProps) {
  const [activeTab, setActiveTab] = useState<"accounts" | "visits">("accounts");
  const [selectedAccounts, setSelectedAccounts] = useState<Set<string>>(new Set());
  const [selectedVisitReports, setSelectedVisitReports] = useState<Set<string>>(new Set());
  const [searchAccounts, setSearchAccounts] = useState("");
  const [searchVisits, setSearchVisits] = useState("");

  const user = useAuthStore((state) => state.user);

  // Fetch accounts 
  const { data: accountsData, isLoading: accountsLoading } = useAccounts({
    per_page: 100,
    search: searchAccounts,
    status: "active",
  });
  const accounts = accountsData?.data ?? [];

  // Fetch visit reports with location data
  // Note: Fetching all visit reports (not filtered by user) so reps can reuse historical routes
  const { data: visitsData, isLoading: visitsLoading } = useVisitReports({
    per_page: 100,
    search: searchVisits,
    // sales_rep_id: user?.id, // Removed: Allow viewing all visit reports with location data
    status: "approved", // Only show approved visits for route planning
  });
  const visitReports = visitsData?.data?.filter(
    (vr) => vr.check_in_location || vr.check_out_location
  ) ?? [];

  const handleToggleAccount = (accountId: string) => {
    const newSelected = new Set(selectedAccounts);
    if (newSelected.has(accountId)) {
      newSelected.delete(accountId);
    } else {
      newSelected.add(accountId);
    }
    setSelectedAccounts(newSelected);
  };

  const handleToggleVisitReport = (visitReportId: string) => {
    const newSelected = new Set(selectedVisitReports);
    if (newSelected.has(visitReportId)) {
      newSelected.delete(visitReportId);
    } else {
      newSelected.add(visitReportId);
    }
    setSelectedVisitReports(newSelected);
  };

  const handleConfirm = async () => {
    const waypoints: Waypoint[] = [];

    // Add selected visit reports with location (these have valid coordinates)
    visitReports
      .filter((vr) => selectedVisitReports.has(vr.id))
      .forEach((visitReport) => {
        // Prefer check_in_location, fallback to check_out_location
        const location = visitReport.check_in_location || visitReport.check_out_location;
        if (location && location.latitude && location.longitude) {
          waypoints.push({
            lat: location.latitude,
            lng: location.longitude,
            address: location.address || `${visitReport.account?.name || "Visit Location"}`,
            account_id: visitReport.account_id,
            visit_report_id: visitReport.id,
            account: visitReport.account
              ? {
                  id: visitReport.account.id,
                  name: visitReport.account.name,
                }
              : undefined,
          });
        }
      });

    // Add selected accounts - prioritize visit report locations, then stored account coords
    if (selectedAccounts.size > 0) {
      const selectedAccountsList = accounts.filter((acc) => selectedAccounts.has(acc.id));
      
      for (const account of selectedAccountsList) {
        // First, try to find visit report with location for this account
        const accountVisitReport = visitReports.find(
          (vr) => vr.account_id === account.id && (vr.check_in_location || vr.check_out_location)
        );

        if (accountVisitReport) {
          // Use location from visit report (more accurate, no geocoding needed)
          const location = accountVisitReport.check_in_location || accountVisitReport.check_out_location;
          if (location && location.latitude && location.longitude) {
            waypoints.push({
              lat: location.latitude,
              lng: location.longitude,
              address: location.address || account.address || account.name,
              account_id: account.id,
              account_name: account.name,
              visit_report_id: accountVisitReport.id,
              account: {
                id: account.id,
                name: account.name,
              },
            });
            continue; // Skip geocoding, we already have location
          }
        }

        if (typeof account.latitude === "number" && typeof account.longitude === "number") {
          waypoints.push({
            lat: account.latitude,
            lng: account.longitude,
            address: account.address || account.name,
            account_id: account.id,
            account_name: account.name,
            account: {
              id: account.id,
              name: account.name,
            },
          });
          continue;
        }

        toast.error(
          `Account "${account.name}" has no saved coordinates and no visit report location. Please add latitude/longitude to the account or pick from Visit Reports.`
        );
      }
    }

    if (waypoints.length > 0) {
      onSelect(waypoints);
      setSelectedAccounts(new Set());
      setSelectedVisitReports(new Set());
      onOpenChange(false);
    } else {
      toast.error("No waypoints with valid locations found");
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-3xl max-h-[80vh] overflow-hidden flex flex-col">
        <DialogHeader>
          <DialogTitle>Select Waypoints</DialogTitle>
        </DialogHeader>

        <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as "accounts" | "visits")} className="flex-1 flex flex-col overflow-hidden">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="accounts">
              <Building2 className="h-4 w-4 mr-2" />
              Accounts ({selectedAccounts.size} selected)
            </TabsTrigger>
            <TabsTrigger value="visits">
              <Calendar className="h-4 w-4 mr-2" />
              Visit Reports ({selectedVisitReports.size} selected)
            </TabsTrigger>
          </TabsList>

          <TabsContent value="accounts" className="flex-1 flex flex-col overflow-hidden mt-4">
            <div className="space-y-4 flex-1 overflow-hidden flex flex-col">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search accounts..."
                  value={searchAccounts}
                  onChange={(e) => setSearchAccounts(e.target.value)}
                  className="pl-10"
                />
              </div>

              <div className="flex-1 overflow-y-auto space-y-2">
                {accountsLoading ? (
                  <div className="text-center py-8 text-muted-foreground">Loading accounts...</div>
                ) : accounts.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No accounts found
                  </div>
                ) : (
                  accounts.map((account) => {
                    const hasAddress = !!(account.address || account.city || account.province);
                    const isSelected = selectedAccounts.has(account.id);
                    // Check if account has visit report with location
                    const hasVisitReportLocation = visitReports.some(
                      (vr) => vr.account_id === account.id && (vr.check_in_location || vr.check_out_location)
                    );
                    const hasSavedCoordinates = typeof account.latitude === "number" && typeof account.longitude === "number";
                    const canSelect = hasSavedCoordinates || hasVisitReportLocation || hasAddress;
                    
                    return (
                      <div
                        key={account.id}
                        className={`flex items-start gap-3 p-3 border rounded-lg hover:bg-muted/50 cursor-pointer transition-colors ${
                          !canSelect ? "opacity-50 cursor-not-allowed" : ""
                        }`}
                        onClick={() => canSelect && handleToggleAccount(account.id)}
                      >
                        <Checkbox
                          checked={isSelected}
                          disabled={!canSelect}
                          onCheckedChange={() => canSelect && handleToggleAccount(account.id)}
                        />
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <Building2 className="h-4 w-4 text-muted-foreground" />
                            <span className="font-medium">{account.name}</span>
                            {isSelected && (
                              <Badge variant="secondary" className="ml-2 text-xs">
                                Selected
                              </Badge>
                            )}
                            {hasVisitReportLocation && (
                              <Badge variant="outline" className="ml-2 text-xs text-green-600 border-green-600">
                                <MapPin className="h-3 w-3 mr-1" />
                                Has GPS Location
                              </Badge>
                            )}
                          </div>
                          {hasVisitReportLocation ? (
                            <div className="flex items-center gap-2 mt-1 text-sm text-green-600">
                              <MapPin className="h-3 w-3" />
                              <span className="text-xs font-medium">
                                Location available from visit report (will use GPS coordinates)
                              </span>
                            </div>
                          ) : hasAddress ? (
                            <>
                              {account.address && (
                                <div className="flex items-center gap-2 mt-1 text-sm text-muted-foreground">
                                  <MapPin className="h-3 w-3" />
                                  <span>{account.address}</span>
                                </div>
                              )}
                              {account.city && account.province && (
                                <p className="text-xs text-muted-foreground mt-1">
                                  {account.city}, {account.province}
                                </p>
                              )}
                            </>
                          ) : (
                            <div className="flex items-center gap-2 mt-1 text-sm text-muted-foreground">
                              <MapPin className="h-3 w-3" />
                              <span className="text-xs italic">No address or visit report location available</span>
                            </div>
                          )}
                        </div>
                      </div>
                    );
                  })
                )}
              </div>
            </div>
          </TabsContent>

          <TabsContent value="visits" className="flex-1 flex flex-col overflow-hidden mt-4">
            <div className="space-y-4 flex-1 overflow-hidden flex flex-col">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search visit reports..."
                  value={searchVisits}
                  onChange={(e) => setSearchVisits(e.target.value)}
                  className="pl-10"
                />
              </div>

              <div className="flex-1 overflow-y-auto space-y-2">
                {visitsLoading ? (
                  <div className="text-center py-8 text-muted-foreground">Loading visit reports...</div>
                ) : visitReports.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No visit reports with location data found
                  </div>
                ) : (
                  visitReports.map((visitReport) => {
                    const location = visitReport.check_in_location || visitReport.check_out_location;
                    return (
                      <div
                        key={visitReport.id}
                        className="flex items-start gap-3 p-3 border rounded-lg hover:bg-muted/50 cursor-pointer"
                        onClick={() => handleToggleVisitReport(visitReport.id)}
                      >
                        <Checkbox
                          checked={selectedVisitReports.has(visitReport.id)}
                          onCheckedChange={() => handleToggleVisitReport(visitReport.id)}
                        />
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <Calendar className="h-4 w-4 text-muted-foreground" />
                            <span className="font-medium">
                              {visitReport.account?.name || "Visit Report"}
                            </span>
                            <Badge variant="outline" className="ml-2">
                              {visitReport.status}
                            </Badge>
                          </div>
                          {location && (
                            <div className="flex items-center gap-2 mt-1 text-sm text-muted-foreground">
                              <MapPin className="h-3 w-3" />
                              <span>
                                {location.address || `${location.latitude}, ${location.longitude}`}
                              </span>
                              <Badge variant="secondary" className="ml-2 text-xs">
                                {location.latitude.toFixed(6)}, {location.longitude.toFixed(6)}
                              </Badge>
                            </div>
                          )}
                          <p className="text-xs text-muted-foreground mt-1">
                            {new Date(visitReport.visit_date).toLocaleDateString()}
                          </p>
                        </div>
                      </div>
                    );
                  })
                )}
              </div>
            </div>
          </TabsContent>
        </Tabs>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={handleConfirm}
            disabled={selectedAccounts.size === 0 && selectedVisitReports.size === 0}
          >
            Add {selectedAccounts.size + selectedVisitReports.size} Waypoint
            {selectedAccounts.size + selectedVisitReports.size !== 1 ? "s" : ""}
            {selectedAccounts.size > 0 && selectedVisitReports.size > 0
              ? " (Accounts + Visits)"
              : selectedAccounts.size > 0
                ? " from Accounts"
                : " from Visit Reports"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

