"use client";

import { useState, useCallback, useRef, useEffect, useMemo } from "react";
import dynamic from "next/dynamic";
import { useRouter, usePathname } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import {
  Building2,
  Plus,
  Search,
  MapPin,
  Phone,
  ChevronRight,
  ChevronLeft,
  Loader2,
  Eye,
  Edit,
  Trash2,
  Tag,
  UserCircle,
  Filter,
  X,
  Navigation,
  Layers,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { toBadgeVariant } from "@/lib/badge-variant";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { AccountForm } from "./account-form";
import { AccountDetailModal } from "./account-detail-modal";
import { useAccountsByBBox, useCreateAccount, useUpdateAccount, useDeleteAccount, useAccount } from "../hooks/useAccounts";
import { useCategories } from "../hooks/useCategories";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { useDebounce } from "@/hooks/use-debounce";
import { toast } from "sonner";
import { useTheme } from "next-themes";
import { SmartTileLayer, TILE_SOURCES, LIGHT_FALLBACK_CHAIN, DARK_FALLBACK_CHAIN, type TileSource } from "@/components/ui/smart-tile-layer";
import type { Account } from "../types";
import type { CreateAccountFormData, UpdateAccountFormData } from "../schemas/account.schema";
import { useTranslations } from "next-intl";
import { cn } from "@/lib/utils";
import { useIsMobile } from "@/hooks/use-mobile";

// Dynamically import Leaflet components (client-side only)
const MapContainer = dynamic(
  () => import("react-leaflet").then((mod) => mod.MapContainer),
  { ssr: false }
);
const ZoomControl = dynamic(
  () => import("react-leaflet").then((mod) => mod.ZoomControl),
  { ssr: false }
);

// ============== MAP STYLES ==============
type MapStyle = "light" | "dark" | "satellite" | "streets";

const mapStyles: Record<MapStyle, { name: string; source: TileSource }> = {
  light: { name: "Light", source: TILE_SOURCES.cartoLight },
  dark: { name: "Dark", source: TILE_SOURCES.cartoDark },
  satellite: { name: "Satellite", source: TILE_SOURCES.esriSatellite },
  streets: { name: "Streets", source: TILE_SOURCES.openStreetMap },
};

// ============== TYPES ==============
interface BBoxParams {
  north: number;
  south: number;
  east: number;
  west: number;
  search?: string;
  status?: string;
  category_id?: string;
}

// ============== MAP EVENT HANDLER ==============
const MapEventHandler = dynamic(
  () => import("react-leaflet").then((mod) => {
    const { useMapEvents, useMap } = mod;
    return function MapEventHandlerInner({
      onBoundsChange,
    }: {
      onBoundsChange: (bounds: { north: number; south: number; east: number; west: number }) => void;
    }) {
      const map = useMap();

      // Fire initial bounds immediately on mount so the first query triggers
      // without requiring the user to pan/zoom
      useEffect(() => {
        const b = map.getBounds();
        onBoundsChange({
          north: b.getNorth(),
          south: b.getSouth(),
          east: b.getEast(),
          west: b.getWest(),
        });
      // eslint-disable-next-line react-hooks/exhaustive-deps
      }, []);

      useMapEvents({
        moveend: (e) => {
          const b = e.target.getBounds();
          onBoundsChange({
            north: b.getNorth(),
            south: b.getSouth(),
            east: b.getEast(),
            west: b.getWest(),
          });
        },
        zoomend: (e) => {
          const b = e.target.getBounds();
          onBoundsChange({
            north: b.getNorth(),
            south: b.getSouth(),
            east: b.getEast(),
            west: b.getWest(),
          });
        },
      });
      return null;
    };
  }),
  { ssr: false }
);

// Cluster component for grouping markers
// Cluster icon factory - returns a function that creates cluster div icons
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function createClusterIconFactory(LeafletLib: any) {
  const dimensionMap: Record<string, number> = { small: 36, medium: 44, large: 52 };
  const fontSizeMap: Record<string, number> = { small: 12, medium: 14, large: 16 };

  return (clusterObj: { getChildCount: () => number }) => {
    const count = clusterObj.getChildCount();
    let size = "small";
    if (count >= 100) size = "large";
    else if (count >= 10) size = "medium";

    const dim = dimensionMap[size];
    const fontSize = fontSizeMap[size];

    return LeafletLib.divIcon({
      className: "custom-cluster-icon",
      html: `<div style="
        background: linear-gradient(135deg, oklch(0.63 0.19 250) 0%, oklch(0.55 0.16 250) 100%);
        color: white;
        width: ${dim}px;
        height: ${dim}px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 700;
        font-size: ${fontSize}px;
        border: 3px solid white;
        box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4), 0 2px 4px rgba(0,0,0,0.2);
      ">${count}</div>`,
      iconSize: [dim, dim],
      iconAnchor: [dim / 2, dim / 2],
    });
  };
}

// Populates a cluster layer with account markers
function addAccountMarkers(
  accounts: Account[],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  cluster: any,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  LeafletLib: any,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  createIcon: (account: Account) => any,
  onAccountClick: (account: Account) => void,
) {
  for (const account of accounts) {
    if (account.latitude == null || account.longitude == null) continue;
    const icon = createIcon(account);
    if (!icon) continue;

    const marker = LeafletLib.marker([account.latitude, account.longitude], { icon });
    marker.on("click", () => onAccountClick(account));

    const categoryName = account.category?.name || "-";
    marker.bindPopup(`
      <div style="min-width: 200px; font-family: system-ui, sans-serif;">
        <div style="font-weight: 600; font-size: 14px; margin-bottom: 4px;">${account.name}</div>
        <div style="font-size: 12px; color: #6b7280; margin-bottom: 2px;">${categoryName}</div>
        ${account.address ? `<div style="font-size: 11px; color: #9ca3af;">${account.address}</div>` : ""}
        ${account.city ? `<div style="font-size: 11px; color: #9ca3af;">${account.city}</div>` : ""}
      </div>
    `, { maxWidth: 280, closeButton: true, autoPan: true });

    cluster.addLayer(marker);
  }
}

// Sets up cluster layer on the map, returns cleanup function
function setupClusterLayer(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  LeafletLib: any,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  map: any,
  accounts: Account[],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  createIcon: (account: Account) => any,
  onAccountClick: (account: Account) => void,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  clusterRef: React.RefObject<any>,
): (() => void) | undefined {
  if (!LeafletLib || !map) return undefined;

  if (clusterRef.current) {
    map.removeLayer(clusterRef.current);
    clusterRef.current = null;
  }

  // markerClusterGroup is added to L by the leaflet.markercluster import
  const mcg = LeafletLib.markerClusterGroup;
  if (!mcg) return undefined;

  const cluster = mcg({
    maxClusterRadius: 50,
    spiderfyOnMaxZoom: true,
    showCoverageOnHover: false,
    zoomToBoundsOnClick: true,
    disableClusteringAtZoom: 16,
    chunkedLoading: true,
    iconCreateFunction: createClusterIconFactory(LeafletLib),
  });

  addAccountMarkers(accounts, cluster, LeafletLib, createIcon, onAccountClick);

  map.addLayer(cluster);
  clusterRef.current = cluster;

  return () => {
    if (clusterRef.current) {
      map.removeLayer(clusterRef.current);
      clusterRef.current = null;
    }
  };
}

const MarkerClusterDynamic = dynamic(
  () => Promise.all([
    import("react-leaflet"),
    import("leaflet"),
    import("leaflet.markercluster"),
  ]).then(([reactLeaflet, leaflet]) => {
    const { useMap } = reactLeaflet;
    const L = leaflet.default;

    return function MarkerClusterInner({
      accounts,
      onAccountClick,
      createIcon,
    }: {
      accounts: Account[];
      onAccountClick: (account: Account) => void;
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      createIcon: (account: Account) => any;
    }) {
      const map = useMap();
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const clusterRef = useRef<any>(null);

      useEffect(() => {
        return setupClusterLayer(L, map, accounts, createIcon, onAccountClick, clusterRef);
      }, [accounts, map, onAccountClick, createIcon]);

      return null;
    };
  }),
  { ssr: false }
);

// ============== ACCOUNT LIST CONTENT ==============
interface AccountListContentProps {
  readonly accounts: Account[];
  readonly isBboxLoading: boolean;
  readonly hasActiveFilters: boolean;
  readonly selectedAccount: Account | null;
  readonly onAccountSelect: (account: Account) => void;
  readonly onViewAccount: (id: string) => void;
  readonly onEditAccount: (id: string) => void;
  readonly onDeleteAccount: (id: string) => void;
  readonly hasEditPermission: boolean;
  readonly hasDeletePermission: boolean;
}

function AccountListContent({
  accounts,
  isBboxLoading,
  hasActiveFilters,
  selectedAccount,
  onAccountSelect,
  onViewAccount,
  onEditAccount,
  onDeleteAccount,
  hasEditPermission,
  hasDeletePermission,
}: AccountListContentProps) {
  if (isBboxLoading && !accounts.length) {
    return (
      <div className="p-4 space-y-3">
        {Array.from({ length: 6 }, (_, i) => (
          <div key={i} className="flex items-start gap-3 p-3">
            <Skeleton className="w-10 h-10 rounded-lg shrink-0" />
            <div className="flex-1 space-y-2">
              <Skeleton className="h-4 w-3/4" />
              <Skeleton className="h-3 w-1/2" />
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (accounts.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 px-6 text-center">
        <div className="w-16 h-16 rounded-2xl bg-muted/50 flex items-center justify-center mb-4">
          <MapPin className="w-8 h-8 text-muted-foreground/50" />
        </div>
        <h3 className="text-sm font-medium text-foreground mb-1">No accounts in view</h3>
        <p className="text-xs text-muted-foreground">
          {hasActiveFilters
            ? "Try adjusting your filters or zoom out"
            : "Pan or zoom the map to find accounts"}
        </p>
      </div>
    );
  }

  return (
    <div className="divide-y divide-border/30">
      {accounts.map((account) => (
        <AccountSidebarItem
          key={account.id}
          account={account}
          isSelected={selectedAccount?.id === account.id}
          onClick={() => onAccountSelect(account)}
          onView={() => onViewAccount(account.id)}
          onEdit={() => onEditAccount(account.id)}
          onDelete={() => onDeleteAccount(account.id)}
          hasEditPermission={hasEditPermission}
          hasDeletePermission={hasDeletePermission}
        />
      ))}
    </div>
  );
}

// ============== MAIN COMPONENT ==============
export function AccountMapManagement() {
  const [leafletLoaded, setLeafletLoaded] = useState(false);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [L, setL] = useState<any>(null);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const mapRef = useRef<any>(null);
  const { theme, resolvedTheme } = useTheme();
  const [mapStyle, setMapStyle] = useState<MapStyle>("light");
  const router = useRouter();
  const pathname = usePathname();
  const isMobile = useIsMobile();
  const t = useTranslations("accountManagement");

  // Sidebar state
  const [isSidebarOpen, setIsSidebarOpen] = useState(!isMobile);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 400);
  const [statusFilter, setStatusFilter] = useState<string>("");
  const [categoryFilter, setCategoryFilter] = useState<string>("");
  const [showFilters, setShowFilters] = useState(false);

  // BBOX state
  const [bboxParams, setBboxParams] = useState<BBoxParams | null>(null);
  const debouncedBbox = useDebounce(bboxParams, 300);

  // Account state
  const [selectedAccount, setSelectedAccount] = useState<Account | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [viewingAccountId, setViewingAccountId] = useState<string | null>(null);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingAccountId, setEditingAccountId] = useState<string | null>(null);
  const [deletingAccountId, setDeletingAccountId] = useState<string | null>(null);

  // Permissions
  const hasCreatePermission = useHasPermission("accounts.create");
  const hasEditPermission = useHasPermission("accounts.update");
  const hasDeletePermission = useHasPermission("accounts.delete");

  // Data hooks
  const { data: categoriesData } = useCategories();
  const categories = categoriesData?.data || [];

  const queryParams = useMemo(() => {
    if (!debouncedBbox) return null;
    return {
      ...debouncedBbox,
      search: debouncedSearch || undefined,
      status: statusFilter || undefined,
      category_id: categoryFilter || undefined,
    };
  }, [debouncedBbox, debouncedSearch, statusFilter, categoryFilter]);

  const { data: bboxData, isLoading: isBboxLoading, isFetching } = useAccountsByBBox(queryParams);
  const accounts = bboxData?.data || [];

  const { data: editingAccountData } = useAccount(editingAccountId || "");
  const createAccount = useCreateAccount();
  const updateAccount = useUpdateAccount();
  const deleteAccount = useDeleteAccount();

  // Auto-select map style based on theme
  useEffect(() => {
    const currentTheme = resolvedTheme || theme;
    setMapStyle(currentTheme === "dark" ? "dark" : "light");
  }, [theme, resolvedTheme]);

  // Load Leaflet on client side
  useEffect(() => {
    if (globalThis.window === undefined) return;

    const link = document.createElement("link");
    link.rel = "stylesheet";
    link.href = "https://unpkg.com/leaflet@1.9.4/dist/leaflet.css";
    link.integrity = "sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY=";
    link.crossOrigin = "";
    document.head.appendChild(link);

    // Load markercluster CSS
    const clusterCss = document.createElement("link");
    clusterCss.rel = "stylesheet";
    clusterCss.href = "https://unpkg.com/leaflet.markercluster@1.5.3/dist/MarkerCluster.css";
    document.head.appendChild(clusterCss);

    const clusterDefaultCss = document.createElement("link");
    clusterDefaultCss.rel = "stylesheet";
    clusterDefaultCss.href = "https://unpkg.com/leaflet.markercluster@1.5.3/dist/MarkerCluster.Default.css";
    document.head.appendChild(clusterDefaultCss);

    import("leaflet").then(async (leaflet) => {
      const Leaflet = leaflet.default;
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      delete (Leaflet.Icon.Default.prototype as any)._getIconUrl;
      Leaflet.Icon.Default.mergeOptions({
        iconRetinaUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
        iconUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
        shadowUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
      });

      // Try loading markercluster plugin
      try {
        await import("leaflet.markercluster");
      } catch {
        // Plugin not available; clustering will be degraded
      }

      setL(Leaflet);
      setLeafletLoaded(true);
    });
  }, []);

  // Handle map bounds change (viewport)
  const handleBoundsChange = useCallback(
    (bounds: { north: number; south: number; east: number; west: number }) => {
      setBboxParams(bounds);
    },
    []
  );

  // Handle initial map load
  const handleMapReady = useCallback(() => {
    const map = mapRef.current;
    if (!map) return;
    const b = map.getBounds();
    if (b) {
      handleBoundsChange({
        north: b.getNorth(),
        south: b.getSouth(),
        east: b.getEast(),
        west: b.getWest(),
      });
    }
  }, [handleBoundsChange]);

  useEffect(() => {
    if (leafletLoaded && mapRef.current) {
      // Small delay to let the map initialize bounds
      const timer = setTimeout(handleMapReady, 200);
      return () => clearTimeout(timer);
    }
  }, [leafletLoaded, handleMapReady]);

  // Fly to account on map
  const flyToAccount = useCallback((account: Account) => {
    const map = mapRef.current;
    if (!map || account.latitude == null || account.longitude == null) return;
    try {
      map.flyTo([account.latitude, account.longitude], 17, { animate: true, duration: 0.6 });
    } catch {
      map.setView([account.latitude, account.longitude], 17);
    }
  }, []);

  // Handle account click from sidebar or marker
  const handleAccountSelect = useCallback((account: Account) => {
    setSelectedAccount(account);
    flyToAccount(account);
    if (isMobile) setIsSidebarOpen(false);
  }, [flyToAccount, isMobile]);

  const handleViewAccount = useCallback((accountId: string) => {
    setViewingAccountId(accountId);
    setIsDetailModalOpen(true);
  }, []);

  // Create account icon
  const createAccountIcon = useCallback((account: Account) => {
    if (!L) return null;
    const isSelected = selectedAccount?.id === account.id;
    const categoryColor = getCategoryColor(account.category?.code);

    return L.divIcon({
      className: "custom-div-icon",
      html: `
        <div style="
          background: ${categoryColor};
          color: white;
          width: ${isSelected ? 40 : 32}px;
          height: ${isSelected ? 40 : 32}px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          border: ${isSelected ? "4px" : "3px"} solid white;
          box-shadow: 0 4px 12px ${isSelected ? "rgba(0,0,0,0.4)" : "rgba(0,0,0,0.2)"}, 0 2px 4px rgba(0,0,0,0.1);
          transition: all 0.2s ease;
          ${isSelected ? "transform: scale(1.1);" : ""}
        ">
          <svg xmlns="http://www.w3.org/2000/svg" width="${isSelected ? 18 : 14}" height="${isSelected ? 18 : 14}" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M6 22V4a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v18Z"/>
            <path d="M6 12H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h2"/>
            <path d="M18 9h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-2"/>
            <path d="M10 6h4"/>
            <path d="M10 10h4"/>
            <path d="M10 14h4"/>
            <path d="M10 18h4"/>
          </svg>
        </div>
      `,
      iconSize: [isSelected ? 40 : 32, isSelected ? 40 : 32],
      iconAnchor: [isSelected ? 20 : 16, isSelected ? 20 : 16],
      popupAnchor: [0, isSelected ? -20 : -16],
    });
  }, [L, selectedAccount]);

  // CRUD handlers
  const handleCreate = useCallback(async (data: CreateAccountFormData | UpdateAccountFormData) => {
    try {
      await createAccount.mutateAsync(data as CreateAccountFormData);
      setIsCreateDialogOpen(false);
      toast.success("Account created successfully");
    } catch {
      // Error handled by interceptor
    }
  }, [createAccount]);

  const handleUpdate = useCallback(async (data: UpdateAccountFormData) => {
    if (!editingAccountId) return;
    try {
      await updateAccount.mutateAsync({ id: editingAccountId, data });
      setEditingAccountId(null);
      toast.success("Account updated successfully");
    } catch {
      // Error handled by interceptor
    }
  }, [editingAccountId, updateAccount]);

  const handleDeleteConfirm = useCallback(async () => {
    if (!deletingAccountId) return;
    try {
      await deleteAccount.mutateAsync(deletingAccountId);
      toast.success("Account deleted successfully");
      setDeletingAccountId(null);
      if (selectedAccount?.id === deletingAccountId) {
        setSelectedAccount(null);
      }
    } catch {
      // Error handled by interceptor
    }
  }, [deletingAccountId, deleteAccount, selectedAccount]);

  const clearFilters = useCallback(() => {
    setSearch("");
    setStatusFilter("");
    setCategoryFilter("");
  }, []);

  const hasActiveFilters = !!debouncedSearch || !!statusFilter || !!categoryFilter;

  const isClient = globalThis.window !== undefined;
  const currentStyle = mapStyles[mapStyle];

  // Loading state
  if (!isClient || !leafletLoaded || !L) {
    return (
      <div className="relative w-full h-full">
        <div className="absolute inset-0 bg-linear-to-br from-blue-50/50 to-indigo-50/50 dark:from-gray-900/50 dark:to-gray-800/50 flex items-center justify-center">
          <div className="text-center space-y-3">
            <div className="relative w-16 h-16 mx-auto">
              <div className="absolute inset-0 border-4 border-primary/20 rounded-full" />
              <div className="absolute inset-0 border-4 border-primary border-t-transparent rounded-full animate-spin" />
              <Navigation className="absolute inset-0 m-auto w-6 h-6 text-primary/60" />
            </div>
            <p className="text-sm text-muted-foreground font-medium">Loading map...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="relative w-full h-full">
      {/* Full-screen Map */}
      <div className="absolute inset-0">
        <MapContainer
          center={[-2.548926, 118.014863]}
          zoom={5}
          style={{ height: "100%", width: "100%" }}
          scrollWheelZoom
          className="z-0"
          zoomControl={false}
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          ref={mapRef as any}
          preferCanvas
          wheelPxPerZoomLevel={120}
          zoomSnap={0.5}
          zoomDelta={0.5}
        >
          <MapEventHandler onBoundsChange={handleBoundsChange} />
          <ZoomControl position="bottomright" />
          <SmartTileLayer
            key={`smart-tile-${mapStyle}`}
            source={currentStyle.source}
            fallbackSources={mapStyle === "dark" ? DARK_FALLBACK_CHAIN : LIGHT_FALLBACK_CHAIN}
            maxRetries={2}
            retryDelay={200}
            priorityMode="viewport"
          />

          {/* Render markers with clustering */}
          <MarkerClusterDynamic
            accounts={accounts}
            onAccountClick={handleAccountSelect}
            createIcon={createAccountIcon}
          />
        </MapContainer>
      </div>

      {/* Map Style Selector - Top Right */}
      <div className="absolute top-4 right-4 z-30">
        <div className="flex flex-col gap-2">
          <Button
            variant="secondary"
            size="icon"
            className="h-10 w-10 rounded-lg shadow-lg bg-card/95 backdrop-blur-sm border-0 hover:bg-card cursor-pointer"
            onClick={() => {
              const styles: MapStyle[] = ["light", "dark", "satellite", "streets"];
              const idx = styles.indexOf(mapStyle);
              setMapStyle(styles[(idx + 1) % styles.length]);
            }}
          >
            <Layers className="h-5 w-5" />
          </Button>
        </div>
      </div>

      {/* Loading indicator - Top Center */}
      <AnimatePresence>
        {isFetching && (
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="absolute top-4 left-1/2 -translate-x-1/2 z-30"
          >
            <div className="bg-card/95 backdrop-blur-md rounded-full shadow-lg border border-border/30 px-4 py-2 flex items-center gap-2">
              <Loader2 className="w-4 h-4 animate-spin text-primary" />
              <span className="text-xs font-medium text-muted-foreground">Loading accounts...</span>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Account count badge - Bottom Left */}
      <div className="absolute bottom-6 left-6 z-30">
        <div className="bg-card/95 backdrop-blur-md rounded-xl shadow-lg border border-border/30 px-4 py-2.5 flex items-center gap-2">
          <Building2 className="w-4 h-4 text-primary" />
          <span className="text-sm font-medium text-foreground">{accounts.length}</span>
          <span className="text-xs text-muted-foreground">on map</span>
          {hasActiveFilters && (
            <Badge variant="secondary" className="text-xs ml-1">filtered</Badge>
          )}
        </div>
      </div>

      {/* Sidebar Toggle Button */}
      {!isSidebarOpen && (
        <div className="absolute top-6 left-6 z-30">
          <Button
            onClick={() => setIsSidebarOpen(true)}
            className="bg-card/95 backdrop-blur-md hover:bg-card shadow-2xl rounded-2xl h-14 px-5 font-medium border border-border/30 text-foreground cursor-pointer"
            variant="ghost"
          >
            <Building2 className="w-5 h-5 mr-2 text-primary" />
            <span>Accounts</span>
            <Badge variant="secondary" className="ml-2 text-xs">{accounts.length}</Badge>
            <ChevronRight className="w-4 h-4 ml-2" />
          </Button>
        </div>
      )}

      {/* Left Sidebar */}
      {isSidebarOpen && (
        <div
          className={cn(
            "absolute top-0 left-0 bottom-0 z-30 flex flex-col bg-card/98 backdrop-blur-xl border-r border-border/30 shadow-2xl",
            isMobile ? "w-full" : "w-[380px]"
          )}
        >
            {/* Sidebar Header */}
            <div className="shrink-0 p-4 border-b border-border/30">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-2.5">
                  <div className="w-9 h-9 rounded-xl bg-linear-to-br from-primary to-primary/80 flex items-center justify-center shadow-lg">
                    <Building2 className="w-4.5 h-4.5 text-white" />
                  </div>
                  <div>
                    <h1 className="font-medium text-foreground text-sm">{t("tabs.accounts")}</h1>
                    <p className="text-xs text-muted-foreground">
                      {accounts.length} on map
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-1">
                  {hasCreatePermission && (
                    <Button
                      size="sm"
                      onClick={() => setIsCreateDialogOpen(true)}
                      className="h-8 px-3 rounded-lg text-xs font-medium cursor-pointer"
                    >
                      <Plus className="w-3.5 h-3.5 mr-1" />
                      Add
                    </Button>
                  )}
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => setIsSidebarOpen(false)}
                    className="h-8 w-8 cursor-pointer"
                  >
                    <ChevronLeft className="w-4 h-4" />
                  </Button>
                </div>
              </div>

              {/* Search */}
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input
                  placeholder="Search accounts..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-9 h-9 text-sm rounded-lg bg-muted/50 border-0 focus-visible:ring-1"
                />
                {search && (
                  <button
                    onClick={() => setSearch("")}
                    className="absolute right-3 top-1/2 -translate-y-1/2 cursor-pointer"
                  >
                    <X className="w-3.5 h-3.5 text-muted-foreground" />
                  </button>
                )}
              </div>

              {/* Filter Toggle */}
              <div className="flex items-center gap-2 mt-2">
                <Button
                  variant={showFilters ? "secondary" : "ghost"}
                  size="sm"
                  onClick={() => setShowFilters(!showFilters)}
                  className="h-7 px-2.5 text-xs cursor-pointer"
                >
                  <Filter className="w-3 h-3 mr-1" />
                  Filters
                  {hasActiveFilters && (
                    <span className="ml-1 w-1.5 h-1.5 rounded-full bg-primary" />
                  )}
                </Button>
                {hasActiveFilters && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={clearFilters}
                    className="h-7 px-2 text-xs text-muted-foreground cursor-pointer"
                  >
                    Clear all
                  </Button>
                )}
              </div>

              {/* Filters Panel */}
              {showFilters && (
                <div className="overflow-hidden">
                  <div className="grid grid-cols-2 gap-2 mt-2">
                      <Select value={statusFilter} onValueChange={setStatusFilter}>
                        <SelectTrigger className="h-8 text-xs">
                          <SelectValue placeholder="All Status" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value=" ">All Status</SelectItem>
                          <SelectItem value="active">Active</SelectItem>
                          <SelectItem value="inactive">Inactive</SelectItem>
                        </SelectContent>
                      </Select>
                      <Select value={categoryFilter} onValueChange={setCategoryFilter}>
                        <SelectTrigger className="h-8 text-xs">
                          <SelectValue placeholder="All Categories" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value=" ">All Categories</SelectItem>
                          {categories
                            .filter((c) => c.status === "active")
                            .map((cat) => (
                              <SelectItem key={cat.id} value={cat.id}>
                                {cat.name}
                              </SelectItem>
                            ))}
                        </SelectContent>
                      </Select>
                    </div>
                </div>
              )}
            </div>

            {/* Account List */}
            <div className="flex-1 overflow-y-auto">
              <AccountListContent
                accounts={accounts}
                isBboxLoading={isBboxLoading}
                hasActiveFilters={hasActiveFilters}
                selectedAccount={selectedAccount}
                onAccountSelect={handleAccountSelect}
                onViewAccount={handleViewAccount}
                onEditAccount={setEditingAccountId}
                onDeleteAccount={setDeletingAccountId}
                hasEditPermission={hasEditPermission}
                hasDeletePermission={hasDeletePermission}
              />
            </div>

            {/* Sidebar Footer - Sub-module navigation */}
            <div className="shrink-0 border-t border-border/30 p-3">
              <div className="flex gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  className="flex-1 h-8 text-xs cursor-pointer"
                  onClick={() => {
                    router.push(pathname + "?tab=categories");
                  }}
                >
                  <Tag className="w-3.5 h-3.5 mr-1.5" />
                  {t("tabs.categories")}
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="flex-1 h-8 text-xs cursor-pointer"
                  onClick={() => {
                    router.push(pathname + "?tab=contact-roles");
                  }}
                >
                  <UserCircle className="w-3.5 h-3.5 mr-1.5" />
                  {t("tabs.contactRoles")}
                </Button>
              </div>
            </div>
          </div>
        )}

      {/* Create Account Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className={cn("sm:max-w-[600px] max-h-[90vh] overflow-y-auto", isMobile && "mx-2 max-w-[calc(100vw-1rem)]")}>
          <DialogHeader>
            <DialogTitle>{t("accountList.createTitle")}</DialogTitle>
          </DialogHeader>
          <AccountForm
            onSubmit={handleCreate}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createAccount.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Account Dialog */}
      {editingAccountId && editingAccountData?.data && (
        <Dialog open={!!editingAccountId} onOpenChange={(open) => !open && setEditingAccountId(null)}>
          <DialogContent className={cn("sm:max-w-[600px] max-h-[90vh] overflow-y-auto", isMobile && "mx-2 max-w-[calc(100vw-1rem)]")}>
            <DialogHeader>
              <DialogTitle>{t("accountList.editTitle")}</DialogTitle>
            </DialogHeader>
            <AccountForm
              account={editingAccountData.data}
              onSubmit={handleUpdate}
              onCancel={() => setEditingAccountId(null)}
              isLoading={updateAccount.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Delete Account Dialog */}
      <DeleteDialog
        open={!!deletingAccountId}
        onOpenChange={(open) => !open && setDeletingAccountId(null)}
        onConfirm={handleDeleteConfirm}
        title={t("accountList.deleteTitle")}
        description={
          deletingAccountId
            ? t("accountList.deleteDescriptionWithName", {
                name: accounts.find((a) => a.id === deletingAccountId)?.name || "this account",
              })
            : t("accountList.deleteDescription")
        }
        itemName={t("accountList.deleteTitle")}
        isLoading={deleteAccount.isPending}
      />

      {/* Account Detail Modal */}
      <AccountDetailModal
        accountId={viewingAccountId}
        open={isDetailModalOpen}
        onOpenChange={(open) => {
          setIsDetailModalOpen(open);
          if (!open) setViewingAccountId(null);
        }}
      />

      {/* Custom CSS */}
      <style jsx global>{`
        .custom-popup .leaflet-popup-content-wrapper {
          border-radius: 12px;
          box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
          padding: 0;
          background: var(--card) !important;
          border: 2px solid var(--border);
        }
        .custom-popup .leaflet-popup-content {
          margin: 0;
          padding: 0;
          color: var(--card-foreground);
        }
        .custom-popup .leaflet-popup-tip {
          background: var(--card);
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
          border: 1px solid var(--border);
        }
        .custom-popup .leaflet-popup-close-button {
          color: var(--muted-foreground) !important;
          font-size: 20px !important;
          font-weight: bold !important;
        }
        .custom-div-icon {
          background: transparent !important;
          border: none !important;
        }
        .custom-cluster-icon {
          background: transparent !important;
          border: none !important;
        }
      `}</style>
    </div>
  );
}

// ============== SIDEBAR ITEM ==============
interface AccountSidebarItemProps {
  readonly account: Account;
  readonly isSelected: boolean;
  readonly onClick: () => void;
  readonly onView: () => void;
  readonly onEdit: () => void;
  readonly onDelete: () => void;
  readonly hasEditPermission: boolean;
  readonly hasDeletePermission: boolean;
}

function AccountSidebarItem({
  account,
  isSelected,
  onClick,
  onView,
  onEdit,
  onDelete,
  hasEditPermission,
  hasDeletePermission,
}: AccountSidebarItemProps) {
  const hasCoords = account.latitude != null && account.longitude != null;

  return (
    <div
      className={cn(
        "group flex items-start gap-3 px-4 py-3 transition-all duration-150 hover:bg-muted/50",
        isSelected && "bg-primary/5 border-l-2 border-l-primary"
      )}
    >
      {/* Clickable left section (icon + content) */}
      <button
        type="button"
        className="flex items-start gap-3 flex-1 min-w-0 text-left cursor-pointer bg-transparent border-0 p-0"
        onClick={onClick}
      >
        {/* Icon */}
        <div
          className={cn(
            "w-10 h-10 rounded-xl flex items-center justify-center shrink-0 transition-colors",
            isSelected ? "bg-primary text-white" : "bg-muted/80 text-muted-foreground"
          )}
        >
          <Building2 className="w-5 h-5" />
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <h3 className="text-sm font-medium text-foreground truncate">{account.name}</h3>
            {!hasCoords && (
              <span className="shrink-0 w-2 h-2 rounded-full bg-amber-400" title="No coordinates" />
            )}
          </div>

          <div className="flex items-center gap-1.5 mt-0.5 flex-wrap">
            {account.category && (
              <Badge
                variant={toBadgeVariant(account.category.badge_color, "secondary")}
                className="text-[10px] font-normal h-4 px-1.5"
              >
                {account.category.name}
              </Badge>
            )}
            <Badge
              variant={account.status === "active" ? "active" : "inactive"}
              className="text-[10px] h-4 px-1.5"
            >
              {account.status}
            </Badge>
          </div>

          {account.address && (
            <p className="text-xs text-muted-foreground mt-1 truncate">{account.address}</p>
          )}
          {account.city && !account.address && (
            <p className="text-xs text-muted-foreground mt-1">{account.city}</p>
          )}

          {account.phone && (
            <div className="flex items-center gap-1 mt-1">
              <Phone className="w-3 h-3 text-muted-foreground" />
              <span className="text-xs text-primary">{account.phone}</span>
            </div>
          )}
        </div>
      </button>

      {/* Actions - siblings of the button, not children */}
      <div className="flex items-center gap-0.5 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
        <Button
          variant="ghost"
          size="icon"
          className="h-7 w-7 cursor-pointer"
          onClick={onView}
          title="View"
        >
          <Eye className="h-3.5 w-3.5" />
        </Button>
        {hasEditPermission && (
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 cursor-pointer"
            onClick={onEdit}
            title="Edit"
          >
            <Edit className="h-3.5 w-3.5" />
          </Button>
        )}
        {hasDeletePermission && (
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 text-destructive hover:text-destructive cursor-pointer"
            onClick={onDelete}
            title="Delete"
          >
            <Trash2 className="h-3.5 w-3.5" />
          </Button>
        )}
      </div>
    </div>
  );
}

// ============== HELPERS ==============
function getCategoryColor(code?: string): string {
  const colors: Record<string, string> = {
    hospital: "linear-gradient(135deg, #ef4444 0%, #dc2626 100%)",
    clinic: "linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)",
    pharmacy: "linear-gradient(135deg, #10b981 0%, #059669 100%)",
    puskesmas: "linear-gradient(135deg, #f59e0b 0%, #d97706 100%)",
    lab: "linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%)",
  };
  return colors[code?.toLowerCase() || ""] || "linear-gradient(135deg, oklch(0.73 0.19 55) 0%, oklch(0.68 0.16 55) 100%)";
}
