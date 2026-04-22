"use client";

import { useState, useCallback, useMemo } from "react";
import { toast } from "sonner";
import { motion, AnimatePresence } from "framer-motion";
import {
  Grid3X3,
  Plus,
  Search,
  Filter,
  List,
  X,
  ChevronRight,
  MapPin,
  Users,
  Loader2,
  BarChart3,
  Edit,
  Trash2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Drawer } from "@/components/ui/drawer";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useBricks, useCreateBrick, useUpdateBrick, useDeleteBrick } from "../hooks/useBricks";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import type { Brick } from "../types";
import type { CreateBrickFormData, UpdateBrickFormData } from "../schemas/brick.schema";
import { BrickForm } from "./brick-form";
import { useBrickPeriodQueryParams } from "../hooks/useBrickPeriodQueryParams";
import { BrickPeriodFilter } from "./brick-period-filter";

// Static import for BrickMapFull (it is already dynamic inside)
import { BrickMapFull } from "./brick-map-full";

export function BrickManagement() {
  const [listSearchInput, setListSearchInput] = useState("");
  const [status, setStatus] = useState<string>("");
  const [isListOpen, setIsListOpen] = useState(false);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [editingBrick, setEditingBrick] = useState<Brick | null>(null);
  const [selectedBrick, setSelectedBrick] = useState<Brick | null>(null);
  const [deletingBrickId, setDeletingBrickId] = useState<string | null>(null);
  // Pre-fill data from map click on unassigned area
  const [mapCreateData, setMapCreateData] = useState<{ regency: string; province: string } | null>(null);
  const [colorBy, setColorBy] = useState<"revenue" | "achievement">("revenue");

  const t = useTranslations("brickManagement");
  const router = useRouter();

  const { mode, periodStart, periodEnd } = useBrickPeriodQueryParams();

  const periodQueryString = useMemo(() => {
    const params = new URLSearchParams();
    params.set("period_mode", mode);
    params.set("period_start", periodStart);
    params.set("period_end", periodEnd);
    return params.toString();
  }, [mode, periodStart, periodEnd]);

  // Fetch bricks
  const { data: bricksData, isLoading } = useBricks({
    per_page: 100,
    status: status || undefined,
    search: undefined,
  });

  const createBrick = useCreateBrick();
  const updateBrick = useUpdateBrick();
  const deleteBrick = useDeleteBrick();

  const bricks = bricksData?.data ?? [];

  // Filter bricks for drawer list (client-side filtering)
  const filteredListBricks = useMemo(() => {
    if (!listSearchInput.trim()) return bricks;

    const query = listSearchInput.toLowerCase().trim();
    return bricks.filter(
      (brick) =>
        brick.name.toLowerCase().includes(query) ||
        brick.code.toLowerCase().includes(query) ||
        brick.province.toLowerCase().includes(query) ||
        brick.regency.toLowerCase().includes(query) ||
        brick.manager?.name?.toLowerCase().includes(query)
    );
  }, [bricks, listSearchInput]);

  // Handle brick selection from map
  const handleBrickSelect = useCallback((brick: Brick) => {
    setSelectedBrick(brick);
  }, []);

  // Handle create brick from map (clicking unassigned area)
  const handleCreateBrickFromMap = useCallback((regency: string, province: string) => {
    setMapCreateData({ regency, province });
    setIsFormOpen(true);
  }, []);

  // Handle view details
  const handleViewDetails = useCallback((brick: Brick) => {
    router.push(`/master-data/bricks/${brick.id}/dashboard?${periodQueryString}`);
  }, [router, periodQueryString]);

  // Handle edit
  const handleEdit = useCallback((brick: Brick) => {
    setEditingBrick(brick);
    setSelectedBrick(null);
  }, []);

  // Handle create
  const handleCreate = async (data: CreateBrickFormData | UpdateBrickFormData) => {
    try {
      const result = await createBrick.mutateAsync(data as CreateBrickFormData);
      setIsFormOpen(false);
      setMapCreateData(null);
      toast.success(t("form.createSuccess"));
      // Auto-select the newly created brick so it appears highlighted in the map and detail card updates immediately
      if (result?.data) {
        setSelectedBrick(result.data);
      }
    } catch {
      // Error handled by api-client
    }
  };

  // Handle update
  const handleUpdate = async (data: CreateBrickFormData | UpdateBrickFormData) => {
    if (editingBrick) {
      try {
        await updateBrick.mutateAsync({ id: editingBrick.id, data: data as UpdateBrickFormData });
        setEditingBrick(null);
        toast.success(t("form.updateSuccess"));
      } catch {
        // Error handled by api-client
      }
    }
  };

  // Handle delete
  const handleDeleteConfirm = async () => {
    if (deletingBrickId) {
      try {
        await deleteBrick.mutateAsync(deletingBrickId);
        setDeletingBrickId(null);
        if (selectedBrick?.id === deletingBrickId) {
          setSelectedBrick(null);
        }
        toast.success(t("form.deleteSuccess"));
      } catch {
        // Error handled by api-client
      }
    }
  };

  return (
    <div className="relative w-full h-full">
      {/* Full-screen Map */}
      <div className="absolute inset-0">
        <BrickMapFull
          bricks={bricks}
          searchQuery=""
          selectedBrickId={selectedBrick?.id}
          onBrickClick={handleBrickSelect}
          onCreateBrickFromMap={handleCreateBrickFromMap}
          isLoading={isLoading}
          colorBy={colorBy}
        />
      </div>

      {/* Floating Controls - Top Left */}
      <div className="absolute top-6 left-6 z-30 flex flex-col gap-3">
        {/* Logo/Title */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="bg-card/95 backdrop-blur-md rounded-2xl shadow-2xl border border-border/30 p-4 flex items-center gap-3"
        >
          <div className="w-10 h-10 rounded-xl bg-primary flex items-center justify-center shadow-lg">
            <Grid3X3 className="w-5 h-5 text-primary-foreground" />
          </div>
          <div>
            <h1 className="font-medium text-foreground">{t("page.title")}</h1>
            <p className="text-xs text-muted-foreground">
              {bricks.length > 0 ? `${bricks.length} ${t("list.bricks")}` : t("page.description")}
            </p>
          </div>
        </motion.div>

        {/* Period Filter (match Sales Performance) */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.15 }}
          className="bg-card/95 backdrop-blur-md rounded-xl shadow-lg border border-border/30 p-3"
        >
          <BrickPeriodFilter />
        </motion.div>


        {/* Action Buttons */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="flex flex-col gap-2"
        >
          {/* Create New Brick Button */}
          <Button
            onClick={() => {
              setMapCreateData(null);
              setIsFormOpen(true);
            }}
            className="bg-primary hover:bg-primary/90 text-primary-foreground shadow-lg rounded-xl h-12 px-5 font-medium cursor-pointer"
          >
            <Plus className="w-5 h-5 mr-2" />
            {t("list.addBrick")}
          </Button>

          {/* Brick List Button */}
          <Button
            variant="secondary"
            onClick={() => setIsListOpen(true)}
            className="bg-card/95 backdrop-blur-md hover:bg-card shadow-lg rounded-xl h-12 px-5 font-medium border-0 cursor-pointer"
          >
            <List className="w-5 h-5 mr-2" />
            {t("list.viewList")}
            {bricks.length > 0 && (
              <span className="ml-2 px-2 py-0.5 bg-primary/10 text-primary rounded-full text-xs font-medium">
                {bricks.length}
              </span>
            )}
          </Button>

          {/* Filter Button */}
          <Button
            variant="secondary"
            onClick={() => setIsFilterOpen(!isFilterOpen)}
            className="bg-card/95 backdrop-blur-md hover:bg-card shadow-lg rounded-xl h-12 px-5 font-medium border-0 cursor-pointer"
          >
            <Filter className="w-5 h-5 mr-2" />
            {t("list.filter")}
            {status && (
              <Badge variant="secondary" className="ml-2 text-xs">
                {status}
              </Badge>
            )}
          </Button>
        </motion.div>

        {/* Filter Dropdown */}
        <AnimatePresence>
          {isFilterOpen && (
            <motion.div
              initial={{ opacity: 0, y: -10, scale: 0.95 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -10, scale: 0.95 }}
              className="bg-card/95 backdrop-blur-md rounded-xl shadow-lg border border-border/30 p-4 space-y-4"
            >
              <div>
                <label className="text-sm font-medium text-foreground mb-2 block">
                  {t("list.filterStatus")}
                </label>
                <Select
                  value={status || "all"}
                  onValueChange={(value) => {
                    setStatus(value === "all" ? "" : value);
                  }}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("list.filterStatus")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">{t("list.allStatuses")}</SelectItem>
                    <SelectItem value="active">{t("list.active")}</SelectItem>
                    <SelectItem value="inactive">{t("list.inactive")}</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div>
                <span className="text-sm font-medium text-foreground mb-2 block">
                  Color By
                </span>
                <Select
                  value={colorBy}
                  onValueChange={(value) => setColorBy(value as "revenue" | "achievement")}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="revenue">Revenue</SelectItem>
                    <SelectItem value="achievement">Achievement</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* Selected Brick Info - Bottom Left */}
      <AnimatePresence>
        {selectedBrick && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 20 }}
            className="absolute bottom-6 left-6 z-30"
          >
            <div className="bg-card/95 backdrop-blur-md rounded-2xl shadow-2xl border border-border/30 p-5 min-w-[300px] max-w-[400px]">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center">
                    <Grid3X3 className="w-6 h-6 text-primary" />
                  </div>
                  <div>
                    <h3 className="font-medium text-foreground">{selectedBrick.name}</h3>
                    <p className="text-xs text-muted-foreground font-mono">{selectedBrick.code}</p>
                  </div>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setSelectedBrick(null)}
                  className="h-8 w-8 cursor-pointer"
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>

              <div className="space-y-2 mb-4">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <MapPin className="w-4 h-4" />
                  <span>{selectedBrick.regency}, {selectedBrick.province}</span>
                </div>
                {selectedBrick.manager && (
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Users className="w-4 h-4" />
                    <span>{selectedBrick.manager.name}</span>
                  </div>
                )}
                <div className="flex items-center gap-2">
                  <Badge variant={selectedBrick.status === "active" ? "default" : "secondary"}>
                    {selectedBrick.status === "active" ? t("list.active") : t("list.inactive")}
                  </Badge>
                </div>
              </div>

              <div className="flex gap-2">
                <Button
                  onClick={() => handleViewDetails(selectedBrick)}
                  className="flex-1 cursor-pointer"
                >
                  <BarChart3 className="w-4 h-4 mr-2" />
                  {t("list.viewDetails")}
                </Button>
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => handleEdit(selectedBrick)}
                  className="cursor-pointer"
                >
                  <Edit className="w-4 h-4" />
                </Button>
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => setDeletingBrickId(selectedBrick.id)}
                  className="text-destructive hover:text-destructive cursor-pointer"
                >
                  <Trash2 className="w-4 h-4" />
                </Button>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Loading State */}
      <AnimatePresence>
        {isLoading && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="absolute inset-0 flex items-center justify-center z-20 bg-background/50 backdrop-blur-sm"
          >
            <div className="flex flex-col items-center gap-4">
              <div className="relative w-16 h-16">
                <div className="absolute inset-0 border-4 border-primary/20 rounded-full" />
                <div className="absolute inset-0 border-4 border-primary border-t-transparent rounded-full animate-spin" />
                <Grid3X3 className="absolute inset-0 m-auto w-6 h-6 text-primary/60" />
              </div>
              <p className="text-sm font-medium text-muted-foreground">Loading bricks...</p>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Brick List Drawer - z-index lebih tinggi dari legend (yang z-1000) */}
      <Drawer
        open={isListOpen}
        onOpenChange={(open) => {
          setIsListOpen(open);
          if (!open) {
            setListSearchInput("");
          }
        }}
        side="right"
        title={t("list.brickList")}
        description={`${bricks.length} ${t("list.bricks")}`}
        defaultWidth={480}
        minWidth={400}
        maxWidth={600}
      >
        <div className="space-y-4">
          {/* Search in drawer */}
          <div className="sticky top-0 bg-background pb-2 z-10">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder={t("list.searchPlaceholder")}
              value={listSearchInput}
              onChange={(e) => setListSearchInput(e.target.value)}
              className="pl-9 h-10"
            />
            {listSearchInput && (
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setListSearchInput("")}
                className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8 cursor-pointer"
              >
                <X className="h-4 w-4" />
              </Button>
            )}
          </div>

          {/* Filtered results count */}
          {listSearchInput && (
            <p className="text-xs text-muted-foreground">
              {filteredListBricks.length} of {bricks.length} {t("list.bricks")} found
            </p>
          )}

          {/* Brick items */}
          {filteredListBricks.map((brick) => (
            <motion.div
              key={brick.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              onClick={() => {
                setSelectedBrick(brick);
                setIsListOpen(false);
                setListSearchInput("");
              }}
              className="p-4 rounded-xl border border-border/50 hover:border-primary/30 hover:bg-muted/50 cursor-pointer"
            >
              <div className="flex items-start gap-3">
                <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center shrink-0">
                  <Grid3X3 className="w-5 h-5 text-primary" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <span className="font-medium text-foreground">{brick.name}</span>
                    <Badge variant={brick.status === "active" ? "default" : "secondary"} className="text-xs">
                      {brick.status}
                    </Badge>
                  </div>
                  <p className="text-xs text-muted-foreground font-mono mb-1">{brick.code}</p>
                  <p className="text-sm text-muted-foreground truncate">
                    {brick.regency}, {brick.province}
                  </p>
                  {brick.manager && (
                    <p className="text-xs text-muted-foreground mt-1 flex items-center gap-1">
                      <Users className="w-3 h-3" />
                      {brick.manager.name}
                    </p>
                  )}
                </div>
                <ChevronRight className="w-5 h-5 text-muted-foreground" />
              </div>
            </motion.div>
          ))}

          {filteredListBricks.length === 0 && !isLoading && (
            <div className="text-center py-12">
              <Grid3X3 className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
              <p className="text-muted-foreground">
                {listSearchInput ? t("list.noSearchResults") : t("list.noBricks")}
              </p>
            </div>
          )}
        </div>
      </Drawer>

      {/* Create Dialog */}
      <Dialog open={isFormOpen} onOpenChange={(open) => {
        setIsFormOpen(open);
        if (!open) setMapCreateData(null);
      }}>
        <DialogContent className="sm:max-w-[900px] max-w-[95vw] max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{t("form.createTitle")}</DialogTitle>
            <DialogDescription>{t("form.createDescription")}</DialogDescription>
          </DialogHeader>
          <BrickForm
            onSubmit={handleCreate}
            onCancel={() => {
              setIsFormOpen(false);
              setMapCreateData(null);
            }}
            isLoading={createBrick.isPending}
            prefillRegency={mapCreateData?.regency}
            prefillProvince={mapCreateData?.province}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={!!editingBrick} onOpenChange={(open) => !open && setEditingBrick(null)}>
        <DialogContent className="sm:max-w-[900px] max-w-[95vw] max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{t("form.editTitle")}</DialogTitle>
            <DialogDescription>{t("form.editDescription")}</DialogDescription>
          </DialogHeader>
          {editingBrick && (
            <BrickForm
              brick={editingBrick}
              onSubmit={handleUpdate}
              onCancel={() => setEditingBrick(null)}
              isLoading={updateBrick.isPending}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingBrickId}
        onOpenChange={(open) => !open && setDeletingBrickId(null)}
        onConfirm={handleDeleteConfirm}
        title={t("list.deleteTitle")}
        description={t("list.deleteDescription")}
        isLoading={deleteBrick.isPending}
      />
    </div>
  );
}
