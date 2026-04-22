"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { 
  X, 
  Eye, 
  Trash2, 
  Route as RouteIcon, 
  Calendar, 
  Clock, 
  Navigation, 
  Search,
  MapPin
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { useRouteList, useDeleteRoute } from "../hooks/useRouteOptimization";
import { format } from "date-fns";
import { toast } from "sonner";
import type { OptimizedRoute } from "../types";

interface RouteHistoryDialogProps {
  readonly isOpen: boolean;
  readonly onClose: () => void;
  readonly userId?: string;
  readonly onSelectRoute: (route: OptimizedRoute) => void;
}

export function RouteHistoryDialog({
  isOpen,
  onClose,
  userId,
  onSelectRoute,
}: RouteHistoryDialogProps) {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [deletingRouteId, setDeletingRouteId] = useState<string | null>(null);

  const { data, isLoading, isFetching } = useRouteList({
    page,
    per_page: 20,
    user_id: userId,
  });

  const deleteRoute = useDeleteRoute();

  const routes = data?.data ?? [];
  const pagination = data?.meta?.pagination;

  const filteredRoutes = routes.filter((route) => {
    if (!search) return true;
    const searchLower = search.toLowerCase();
    return (
      route.route_name?.toLowerCase().includes(searchLower) ||
      route.id.toLowerCase().includes(searchLower)
    );
  });

  const handleDelete = async () => {
    if (deletingRouteId) {
      try {
        await deleteRoute.mutateAsync(deletingRouteId);
        toast.success("Route deleted successfully");
        setDeletingRouteId(null);
      } catch (error) {
        toast.error("Failed to delete route");
        throw error;
      }
    }
  };

  const handleViewRoute = (route: OptimizedRoute) => {
    onSelectRoute(route);
    onClose();
  };

  return (
    <>
      <AnimatePresence>
        {isOpen && (
          <>
            {/* Backdrop */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.2 }}
              className="fixed inset-0 bg-black/30 backdrop-blur-sm z-1000"
              onClick={onClose}
            />

            {/* Dialog */}
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 20 }}
              transition={{ type: "spring", damping: 25, stiffness: 300 }}
              className="fixed inset-x-4 top-[5%] md:inset-auto md:left-1/2 md:top-1/2 md:-translate-x-1/2 md:-translate-y-1/2 w-auto md:w-full md:max-w-2xl max-h-[90vh] bg-card/95 backdrop-blur-xl rounded-2xl shadow-2xl z-1001 border border-border/50 overflow-hidden"
            >
              {/* Header */}
              <div className="sticky top-0 bg-card/95 backdrop-blur-md border-b border-border/50 p-5 z-10">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-linear-to-br from-primary to-primary/80 flex items-center justify-center shadow-lg">
                      <Calendar className="w-5 h-5 text-white" />
                    </div>
                    <div>
                      <h2 className="font-medium text-lg text-foreground">Route History</h2>
                      <p className="text-xs text-muted-foreground">
                        {pagination?.total ?? 0} saved routes
                      </p>
                    </div>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={onClose}
                    className="h-9 w-9 rounded-lg hover:bg-muted"
                  >
                    <X className="h-5 w-5" />
                  </Button>
                </div>

                {/* Search */}
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    placeholder="Search routes..."
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    className="pl-10 bg-muted/50 border-border/50"
                  />
                </div>
              </div>

              {/* Content */}
              <ScrollArea className="h-[60vh] md:h-[50vh]">
                <div className="p-4 space-y-3">
                  {isLoading ? (
                    // Loading skeleton
                    Array.from({ length: 5 }).map((_, i) => (
                      <div key={i} className="rounded-xl border border-border/50 p-4">
                        <div className="flex items-start gap-3">
                          <Skeleton className="w-10 h-10 rounded-lg" />
                          <div className="flex-1 space-y-2">
                            <Skeleton className="h-5 w-32" />
                            <Skeleton className="h-4 w-48" />
                          </div>
                        </div>
                      </div>
                    ))
                  ) : filteredRoutes.length === 0 ? (
                    // Empty state
                    <div className="text-center py-12">
                      <div className="w-16 h-16 mx-auto mb-4 rounded-2xl bg-muted/50 flex items-center justify-center">
                        <RouteIcon className="w-8 h-8 text-muted-foreground" />
                      </div>
                      <h3 className="font-medium text-foreground mb-1">No routes found</h3>
                      <p className="text-sm text-muted-foreground">
                        {search ? "Try a different search term" : "Create your first optimized route"}
                      </p>
                    </div>
                  ) : (
                    // Routes list
                    filteredRoutes.map((route, index) => (
                      <motion.div
                        key={route.id}
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.03 }}
                        className="group relative rounded-xl border border-border/50 bg-card hover:bg-muted/30 hover:border-primary/30 p-4"
                      >
                        <div className="flex items-start gap-3">
                          {/* Icon */}
                          <div className="w-10 h-10 rounded-lg bg-linear-to-br from-primary/20 to-primary/10 flex items-center justify-center shrink-0">
                            <RouteIcon className="w-5 h-5 text-primary" />
                          </div>

                          {/* Content */}
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2 mb-1">
                              <h4 className="font-medium text-foreground truncate">
                                {route.route_name || `Route ${route.id.slice(0, 8)}`}
                              </h4>
                              <Badge variant="secondary" className="text-xs">
                                <MapPin className="w-3 h-3 mr-1" />
                                {route.waypoints.length} stops
                              </Badge>
                            </div>

                            <div className="flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                              <span className="flex items-center gap-1">
                                <Navigation className="w-3 h-3" />
                                {route.total_distance_formatted || "N/A"}
                              </span>
                              <span className="flex items-center gap-1">
                                <Clock className="w-3 h-3" />
                                {route.total_duration_formatted || "N/A"}
                              </span>
                              <span className="flex items-center gap-1">
                                <Calendar className="w-3 h-3" />
                                {format(new Date(route.created_at), "MMM d, yyyy")}
                              </span>
                            </div>
                          </div>

                          {/* Actions */}
                          <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleViewRoute(route)}
                              className="h-8 w-8 text-primary hover:text-primary hover:bg-primary/10 cursor-pointer"
                            >
                              <Eye className="h-4 w-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => setDeletingRouteId(route.id)}
                              className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10 cursor-pointer"
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </div>
                      </motion.div>
                    ))
                  )}
                </div>

                {/* Pagination */}
                {pagination && pagination.total_pages > 1 && (
                  <div className="sticky bottom-0 bg-card/95 backdrop-blur-md border-t border-border/50 p-4">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-muted-foreground">
                        Page {page} of {pagination.total_pages}
                      </span>
                      <div className="flex gap-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => setPage((p) => Math.max(1, p - 1))}
                          disabled={!pagination.has_prev || isFetching}
                        >
                          Previous
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => setPage((p) => p + 1)}
                          disabled={!pagination.has_next || isFetching}
                        >
                          Next
                        </Button>
                      </div>
                    </div>
                  </div>
                )}
              </ScrollArea>
            </motion.div>
          </>
        )}
      </AnimatePresence>

      <DeleteDialog
        open={!!deletingRouteId}
        onOpenChange={(open) => !open && setDeletingRouteId(null)}
        onConfirm={handleDelete}
        title="Delete Route"
        description="Are you sure you want to delete this route? This action cannot be undone."
        isLoading={deleteRoute.isPending}
      />
    </>
  );
}
