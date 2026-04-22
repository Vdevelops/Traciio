import { useState, useMemo } from "react";
import { formatCurrency } from "@/lib/utils";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useMonthlyTargets, useBulkSetMonthlyTarget } from "../hooks/useMonthlyTargets";
import { Loader2, Edit2, Check, X } from "lucide-react";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import type { MonthlyTarget } from "../types";
import type { User } from "@/features/master-data/user-management/types";

interface TargetMatrixProps {
  initialYear?: number;
  showHeader?: boolean;
  data?: MonthlyTarget[];
  isLoading?: boolean;
  onLoadMore?: () => void;
  hasMore?: boolean;
  isLoadingMore?: boolean;
  searchQuery?: string;
}


export function TargetMatrix({ 
  initialYear = new Date().getFullYear(), 
  showHeader = true,
  data: externalData,
  isLoading: externalIsLoading,
  onLoadMore,
  hasMore,
  isLoadingMore,
  searchQuery
}: Readonly<TargetMatrixProps>) {
  const t = useTranslations("monthlyTargetManagement.planner");
  // `monthsShort` is stored as a comma-separated string in translations (arrays not supported by next-intl)
  const monthsShort = (t("monthsShort") as unknown as string).split(",");
  const [year] = useState(initialYear);
  const [editingCell, setEditingCell] = useState<{ userId: string; month: number } | null>(null);
  const [editValue, setEditValue] = useState("");
  
  // Fetch targets for the year if not provided externally
  const { data: fetchedData, isLoading: isFetching } = useMonthlyTargets({ 
    year, 
    per_page: 500, // Fetch enough to show in grid
    scope: "user",
  }, { 
    enabled: !externalData // Only fetch if external data is not provided
  });
  
  const targets = externalData || fetchedData?.data || [];
  const isLoading = externalIsLoading ?? isFetching;

  const bulkSetTarget = useBulkSetMonthlyTarget();

  // Transform data into matrix structure
  // { userId: { user: UserObj, targets: { 1: targetInJan, 2: targetInFeb ... } } }
  const matrixData = useMemo(() => {
    const map = new Map<string, { user: User; targets: Record<number, MonthlyTarget> }>();

    targets.forEach((t) => {
      if (t.user_id && t.user) {
        if (!map.has(t.user_id)) {
          map.set(t.user_id, { user: t.user as User, targets: {} });
        }
        const entry = map.get(t.user_id)!;
        entry.targets[t.month] = t;
      }
    });

    return Array.from(map.values()).sort((a, b) => 
      (a.user.name || "").localeCompare(b.user.name || "")
    );
  }, [targets]);

  // Apply search filtering if provided
  const filteredMatrixData = useMemo(() => {
    if (!searchQuery || !searchQuery.trim()) return matrixData;
    const q = searchQuery.trim().toLowerCase();
    return matrixData.filter(({ user }) => {
      const name = (user?.name || "").toLowerCase();
      const email = (user?.email || "").toLowerCase();
      return name.includes(q) || email.includes(q) || (user?.id || "").toLowerCase().includes(q);
    });
  }, [matrixData, searchQuery]);

  const handleEditClick = (userId: string, month: number, currentValue: number) => {
    setEditingCell({ userId, month });
    // Initialize with formatted value for display, user can edit it
    // But Input value state should be the raw string or handle formatting on change
    // Using formatted value in state makes editing harder (cursor jumps). 
    // Best UX: Show unformatted or properly masked input. 
    // User asked for "inputnya ada format '.'". 
    // We will simple format it: "1.000.000".
    const rupiahValue = Math.round(currentValue / 100);
    setEditValue(rupiahValue > 0 ? rupiahValue.toLocaleString("id-ID") : "");
  };

  const handleCancelEdit = () => {
    setEditingCell(null);
    setEditValue("");
  };

  const handleSaveEdit = async () => {
    if (!editingCell) return;

    const rupiahAmount = Number.parseInt(editValue.replaceAll(/\D/g, ""), 10);
    if (Number.isNaN(rupiahAmount)) {
      toast.error(t("invalidAmount"));
      return;
    }
    const amount = rupiahAmount * 100;

    try {
      await bulkSetTarget.mutateAsync({
        user_id: editingCell.userId,
        year,
        start_month: editingCell.month,
        end_month: editingCell.month, // Single month update
        target_amount: amount,
      });
      setEditingCell(null);
    } catch (error) {
      toast.error(t("failedUpdate"));
      console.error(error);
    }
  };

  const handleBulkApply = async (userId: string, month: number, amount: number) => {
    if (!confirm(t("confirmApply", { amount: formatCurrency(amount) }))) return;

    try {
      await bulkSetTarget.mutateAsync({
        user_id: userId,
        year,
        start_month: month, // Start from this month
        end_month: 12,      // Apply until December
        target_amount: amount,
      });
    } catch (error) {
        toast.error(t("failedBulkApply"));
        console.error(error);
    }
  };

  if (isLoading && !targets.length) {
    return (
      <div className="flex justify-center p-8">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <Table>
        {showHeader && (
          <TableHeader className="sticky top-0 z-20 bg-background shadow-sm">
            <TableRow>
              <TableHead className="w-[200px] sticky left-0 bg-background z-20 border-r">{t("user")}</TableHead>
              {monthsShort.map((m: string) => (
                <TableHead key={m} className="min-w-[120px] text-right bg-background">{m}</TableHead>
              ))}
              <TableHead className="text-right font-bold bg-background">{t("total")}</TableHead>
            </TableRow>
            {/* Grand Total Row at Top */}
             <TableRow className="bg-muted/50 font-bold border-b-2">
               <TableCell className="sticky left-0 bg-muted/90 z-20 border-r">{t("grandTotal")}</TableCell>
               {Array.from({ length: 12 }).map((_, i) => {
                 const month = i + 1;
                 const monthlyTotal = matrixData.reduce((acc, { targets }) => {
                   return acc + (targets[month]?.target_amount || 0);
                 }, 0);
                 return (
                   <TableCell key={month} className="text-right bg-muted/50">
                     {formatCurrency(monthlyTotal)}
                   </TableCell>
                 );
               })}
               <TableCell className="text-right font-bold bg-muted/50">
                 {formatCurrency(
                   matrixData.reduce((acc, { targets }) => {
                     return acc + Object.values(targets).reduce((sum, t) => sum + (t.target_amount || 0), 0);
                   }, 0)
                 )}
               </TableCell>
             </TableRow>
          </TableHeader>
        )}
        <TableBody>
            {filteredMatrixData.length === 0 ? (
             <TableRow>
                 <TableCell colSpan={14} className="text-center h-24 text-muted-foreground">
                   {t("noTargetsFound")}
                 </TableCell>
             </TableRow>
            ) : (
              <>
                {filteredMatrixData.map(({ user, targets }) => {
                const totalYear = Object.values(targets).reduce((acc, t) => acc + (t.target_amount || 0), 0);
                
                return (
                  <TableRow key={user.id}>
                    <TableCell className="sticky left-0 bg-background font-medium z-10 border-r">
                      {user.name}
                    </TableCell>
                    {Array.from({ length: 12 }).map((_, i) => {
                      const month = i + 1;
                      const target = targets[month];
                      const amount = target?.target_amount || 0;
                      const isEditing = editingCell?.userId === user.id && editingCell?.month === month;

                      return (
                        <TableCell key={month} className="text-right p-2 relative group">
                            {isEditing ? (
                                <div 
                                      className="absolute right-0 top-0 z-20 flex items-center gap-1 p-1 rounded-md border bg-background shadow-md"
                                      style={{ minWidth: "140px", width: "max-content" }}
                                    >
                                    <Input 
                                        type="text" 
                                        className="h-8 min-w-[100px] text-right"
                                        style={{ width: `${Math.max(editValue.length + 2, 8)}ch` }}
                                        value={editValue}
                                        onChange={(e) => {
                                          // Handle dot formatting
                                          const val = e.target.value;
                                          // Remove non-digits to get raw number
                                          const rawInfo = val.replaceAll(/\D/g, "");
                                          if (rawInfo === "") {
                                              setEditValue("");
                                              return;
                                          }
                                          // Format with dots
                                          const formatted = Number(rawInfo).toLocaleString("id-ID");
                                          setEditValue(formatted);
                                        }}
                                        onKeyDown={(e) => {
                                          if (e.key === "Enter") handleSaveEdit();
                                          if (e.key === "Escape") handleCancelEdit();
                                        }}
                                        autoFocus
                                        onFocus={(e) => e.target.select()}
                                    />
                                    <div className="flex shrink-0">
                                      <Button 
                                          size="icon" 
                                          variant="ghost" 
                                          className="h-8 w-8 text-green-600 hover:text-green-700 hover:bg-green-50" 
                                          onClick={handleSaveEdit}
                                      >
                                          <Check className="h-4 w-4" />
                                      </Button>
                                      <Button 
                                          size="icon" 
                                          variant="ghost" 
                                          className="h-8 w-8 text-red-600 hover:text-red-700 hover:bg-red-50" 
                                          onClick={handleCancelEdit}
                                      >
                                          <X className="h-4 w-4" />
                                      </Button>
                                    </div>
                                </div>
                            ) : (
                                <button 
                                    type="button"
                                    className="w-full h-full cursor-pointer hover:bg-muted/50 p-2 pr-8 rounded flex justify-end items-center gap-2 transition-colors bg-transparent border-none"
                                    onClick={() => handleEditClick(user.id, month, amount)}
                                >
                                    {amount > 0 ? formatCurrency(amount) : t("emptyCell")}
                                </button>
                            )}
                            
                            {!isEditing && amount > 0 && (
                                <Button
                                    size="icon"
                                    variant="ghost"
                                    className="h-6 w-6 absolute right-1 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none group-hover:pointer-events-auto"
                                    title={t("applyToRest")}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        handleBulkApply(user.id, month, amount);
                                    }}
                                >
                                    <Edit2 className="h-3 w-3 text-muted-foreground" />
                                </Button>
                            )}
                        </TableCell>
                      );
                    })}
                    <TableCell className="text-right font-bold bg-muted/20">
                      {formatCurrency(totalYear)}
                    </TableCell>
                  </TableRow>
                );
              })}
              
              {/* Load More Trigger */}
              {hasMore && (
                <TableRow>
                   <TableCell colSpan={14} className="p-4 text-center">
                      <Button 
                        variant="ghost" 
                        onClick={onLoadMore} 
                        disabled={isLoadingMore}
                        className="w-full"
                      >
                        {isLoadingMore ? (
                           <>
                             <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                             Loading more...
                           </>
                        ) : (
                           "Load More"
                        )}
                      </Button>
                   </TableCell>
                </TableRow>
              )}
            </>
          )}
        </TableBody>
        {/* Removed TableFooter since totals are now at top */}
      </Table>
    </div>
  );
}
