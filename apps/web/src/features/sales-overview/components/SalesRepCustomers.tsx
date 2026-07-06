"use client";

import { useState, useMemo, useEffect, Fragment } from "react";
import { useTranslations } from "next-intl";
import { ChevronDown, ChevronRight } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { AccountDetailModal } from "@/features/sales-crm/account-management/components/account-detail-modal";
import { useDeals } from "@/features/sales-crm/pipeline-management/hooks/useDeals";
import { toBadgeVariant } from "@/lib/badge-variant";
import type { Account } from "@/features/sales-crm/account-management/types";
import type { Deal } from "@/features/sales-crm/pipeline-management/types";
import { accountService } from "@/features/sales-crm/account-management/services/accountService";
import { formatCurrency } from "@/lib/utils";

interface SalesRepCustomersProps {
  readonly userId: string;
  readonly startDate?: string;
  readonly endDate?: string;
}

interface AccountWithData extends Account {
  totalRevenue: number;
  latestStage?: {
    name: string;
    color?: string;
  };
}

export function SalesRepCustomers({ userId, startDate, endDate }: SalesRepCustomersProps) {
  const t = useTranslations("salesOverview.customers");
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());
  const [selectedAccountId, setSelectedAccountId] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);

  // Fetch deals with date range filter
  const { data: dealsResponse, isLoading: dealsLoading } = useDeals(
    {
      assigned_to: userId,
      ...(startDate ? { date_from: startDate } : {}),
      ...(endDate ? { date_to: endDate } : {}),
    },
    1,
    1000
  );

  const filteredDeals = useMemo(
    () => dealsResponse?.data ?? [],
    [dealsResponse?.data],
  );

  // Extract unique account IDs from filtered deals
  const accountIds = useMemo(() => {
    const uniqueIds = new Set<string>();
    filteredDeals.forEach((deal) => {
      if (deal?.account_id) {
        uniqueIds.add(deal.account_id);
      }
    });
    return Array.from(uniqueIds);
  }, [filteredDeals]);

  // Calculate total revenue per account and get latest stage
  const accountDataMap = useMemo(() => {
    const revenueMap = new Map<string, number>();
    const stageMap = new Map<string, { name: string; color?: string }>();
    
    if (!filteredDeals || filteredDeals.length === 0) return { revenueMap, stageMap };

    const dealsByAccount = new Map<string, Deal[]>();
    filteredDeals.forEach((deal) => {
      if (!deal?.account_id) return;
      if (!dealsByAccount.has(deal.account_id)) {
        dealsByAccount.set(deal.account_id, []);
      }
      dealsByAccount.get(deal.account_id)!.push(deal);
    });

    dealsByAccount.forEach((deals, accountId) => {
      let accountRevenue = 0;
      let latestDeal: Deal | null = null;

      deals.forEach((deal: Deal) => {
        if (Array.isArray(deal.product_items) && deal.product_items.length > 0) {
          deal.product_items.forEach((item) => {
            const subtotal =
              (item.unit_price ?? 0) * (item.quantity ?? 0) -
              (item.discount_amount ?? 0);
            accountRevenue += subtotal;
          });
        } else {
          accountRevenue += deal.value ?? 0;
        }

        if (!latestDeal) {
          latestDeal = deal;
        } else {
          const dealDate = deal.actual_close_date
            ? new Date(deal.actual_close_date)
            : new Date(deal.created_at);
          const latestDate = latestDeal.actual_close_date
            ? new Date(latestDeal.actual_close_date)
            : new Date(latestDeal.created_at);
          if (dealDate > latestDate) {
            latestDeal = deal;
          }
        }
      });

      revenueMap.set(accountId, accountRevenue);

      if (latestDeal) {
        const dealStage = (latestDeal as Deal).stage;
        if (dealStage) {
          stageMap.set(accountId, {
            name: dealStage.name ?? "-",
            color: dealStage.color,
          });
        }
      }
    });

    return { revenueMap, stageMap };
  }, [filteredDeals]);

  // Fetch accounts
  const [allAccounts, setAllAccounts] = useState<Account[]>([]);
  const [isFetchingAccounts, setIsFetchingAccounts] = useState(false);

  useEffect(() => {
    if (accountIds.length === 0) {
      setAllAccounts([]);
      setIsFetchingAccounts(false);
      return;
    }

    let isCancelled = false;

    const fetchAccounts = async () => {
      setIsFetchingAccounts(true);
      const fetchedAccounts: Account[] = [];
      const accountIdsSet = new Set(accountIds);
      let pageNum = 1;
      let hasMore = true;
      let foundCount = 0;

      while (hasMore && foundCount < accountIds.length) {
        try {
          const response = await accountService.list({
            per_page: 100,
            page: pageNum,
          });

          if (response.data && Array.isArray(response.data)) {
            const matchingAccounts = response.data.filter((account) =>
              accountIdsSet.has(account.id)
            );
            fetchedAccounts.push(...matchingAccounts);
            foundCount += matchingAccounts.length;

            if (foundCount >= accountIds.length) {
              break;
            }
          }

          const pagination = response.meta?.pagination;
          hasMore = pagination?.has_next ?? false;
          pageNum++;
        } catch (error) {
          break;
        }
      }

      if (!isCancelled) {
        setAllAccounts(fetchedAccounts);
        setIsFetchingAccounts(false);
      }
    };

    fetchAccounts();

    return () => {
      isCancelled = true;
    };
  }, [accountIds]);

  const handleViewAccount = (accountId: string) => {
    setSelectedAccountId(accountId);
    setIsModalOpen(true);
  };

  const toggleRow = (accountId: string) => {
    const newExpanded = new Set(expandedRows);
    if (newExpanded.has(accountId)) {
      newExpanded.delete(accountId);
    } else {
      newExpanded.add(accountId);
    }
    setExpandedRows(newExpanded);
  };

  // Prepare accounts with revenue and stage data
  const accountsWithData: AccountWithData[] = useMemo(() => {
    const accounts = accountIds.length === 0 ? [] : allAccounts;
    return accounts.map((account) => ({
      ...account,
      totalRevenue: accountDataMap.revenueMap.get(account.id) ?? 0,
      latestStage: accountDataMap.stageMap.get(account.id),
    }));
  }, [accountIds.length, allAccounts, accountDataMap]);

  // Client-side pagination
  const paginatedAccounts = useMemo(() => {
    const start = (page - 1) * perPage;
    const end = start + perPage;
    return accountsWithData.slice(start, end);
  }, [accountsWithData, page, perPage]);

  const totalAccounts = accountsWithData.length;
  const totalPages = Math.ceil(totalAccounts / perPage);

  const pagination = useMemo(() => {
    if (!totalAccounts) return undefined;
    return {
      page,
      per_page: perPage,
      total: totalAccounts,
      total_pages: totalPages,
      has_next: page < totalPages,
      has_prev: page > 1,
    };
  }, [page, perPage, totalAccounts, totalPages]);

  const handlePerPageChange = (newPerPage: number) => {
    const limitedPerPage = Math.min(newPerPage, 20);
    setPerPage(limitedPerPage);
    setPage(1);
  };

  const isLoading = dealsLoading || isFetchingAccounts;

  if (isLoading && accountsWithData.length === 0) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-[400px] w-full" />
      </div>
    );
  }

  if (!Array.isArray(accountsWithData) || accountsWithData.length === 0) {
    return (
      <div className="text-center text-muted-foreground py-12 px-4">
        {t("empty")}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Minimalist Table with Expandable Rows */}
      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[40px]"></TableHead>
              <TableHead className="min-w-[200px]">{t("table.name")}</TableHead>
              <TableHead className="w-[120px]">{t("table.category")}</TableHead>
              <TableHead className="w-[120px]">{t("table.status")}</TableHead>
              <TableHead className="w-[140px] text-right">{t("table.totalRevenue")}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {paginatedAccounts.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center text-muted-foreground py-8">
                  {t("empty")}
                </TableCell>
              </TableRow>
            ) : (
              paginatedAccounts.map((account) => {
                const isExpanded = expandedRows.has(account.id);
                const accountDeals = filteredDeals.filter(
                  (deal) => deal?.account_id === account.id
                );

                return (
                  <Fragment key={account.id}>
                    <TableRow key={account.id} className="hover:bg-muted/50">
                      <TableCell>
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          className="h-6 w-6 cursor-pointer"
                          onClick={() => toggleRow(account.id)}
                        >
                          {isExpanded ? (
                            <ChevronDown className="h-4 w-4" />
                          ) : (
                            <ChevronRight className="h-4 w-4" />
                          )}
                        </Button>
                      </TableCell>
                      <TableCell>
                        <button
                          onClick={() => handleViewAccount(account.id)}
                          className="font-medium text-primary hover:underline cursor-pointer text-left"
                        >
                          {account.name ?? "-"}
                        </button>
                      </TableCell>
                      <TableCell>
                        {account.category ? (
                          <Badge
                            variant={toBadgeVariant(account.category.badge_color, "secondary")}
                            className="font-normal text-xs"
                          >
                            {account.category.name ?? "-"}
                          </Badge>
                        ) : (
                          <span className="text-muted-foreground text-sm">-</span>
                        )}
                      </TableCell>
                      <TableCell>
                        {account.latestStage ? (
                          <Badge
                            variant="outline"
                            className="text-xs"
                            style={{
                              backgroundColor: account.latestStage.color
                                ? `${account.latestStage.color}15`
                                : undefined,
                              borderColor: account.latestStage.color ?? undefined,
                              color: account.latestStage.color ?? undefined,
                            }}
                          >
                            {account.latestStage.name}
                          </Badge>
                        ) : (
                          <span className="text-muted-foreground text-sm">-</span>
                        )}
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="font-medium text-foreground text-sm">
                          {formatCurrency(account.totalRevenue)}
                        </div>
                      </TableCell>
                    </TableRow>
                    {isExpanded && (
                      <TableRow>
                        <TableCell colSpan={5} className="bg-muted/20 p-0">
                          <div className="p-4 space-y-4">
                            {accountDeals.length === 0 ? (
                              <div className="text-center py-4 text-sm text-muted-foreground">
                                {t("noDeals")}
                              </div>
                            ) : (
                              (() => {
                                // Group deals by contact
                                const dealsByContact = new Map<string, Deal[]>();
                                accountDeals.forEach((deal) => {
                                  const contactId = deal.contact_id ?? "no-contact";
                                  if (!dealsByContact.has(contactId)) {
                                    dealsByContact.set(contactId, []);
                                  }
                                  dealsByContact.get(contactId)!.push(deal);
                                });

                                return Array.from(dealsByContact.entries()).map(
                                  ([contactId, contactDeals]) => {
                                    const contact = contactDeals[0]?.contact;
                                    const isNoContact = contactId === "no-contact";

                                    return (
                                      <div key={contactId} className="space-y-2">
                                        {/* Contact Header */}
                                        <div className="font-medium text-sm text-foreground pl-2 border-l-2 border-primary">
                                          {isNoContact ? (
                                            <span className="text-muted-foreground">
                                              {t("noContact")}
                                            </span>
                                          ) : (
                                            contact?.name ?? "-"
                                          )}
                                        </div>

                                        {/* Product Items for this contact */}
                                        <div className="pl-4 space-y-1">
                                          {contactDeals.map((deal) => {
                                            const productItems =
                                              Array.isArray(deal.product_items) &&
                                              deal.product_items.length > 0
                                                ? deal.product_items
                                                : [];

                                            if (productItems.length === 0) {
                                              return (
                                                <div
                                                  key={deal.id}
                                                  className="text-xs text-muted-foreground pl-2"
                                                >
                                                  {deal.title ?? "-"} - {t("noProducts")}
                                                </div>
                                              );
                                            }

                                            return productItems.map((item, idx) => {
                                              const subtotal =
                                                (item.unit_price ?? 0) * (item.quantity ?? 0) -
                                                (item.discount_amount ?? 0);

                                              return (
                                                <div
                                                  key={item.id ?? `item-${deal.id}-${idx}`}
                                                  className="text-xs text-muted-foreground pl-2 flex items-center gap-2"
                                                >
                                                  <span className="text-foreground min-w-[200px]">
                                                    {item.product_name ?? "-"}
                                                  </span>
                                                  <span className="text-muted-foreground">
                                                    {item.quantity ?? 0}x
                                                  </span>
                                                  <span className="text-muted-foreground">
                                                    @ {formatCurrency(item.unit_price ?? 0)}
                                                  </span>
                                                  {item.discount_amount &&
                                                    item.discount_amount > 0 && (
                                                      <span className="text-red-600">
                                                        -{formatCurrency(item.discount_amount)}
                                                      </span>
                                                    )}
                                                  <span className="text-foreground font-medium ml-auto">
                                                    = {formatCurrency(subtotal)}
                                                  </span>
                                                </div>
                                              );
                                            });
                                          })}
                                        </div>
                                      </div>
                                    );
                                  }
                                );
                              })()
                            )}
                          </div>
                        </TableCell>
                      </TableRow>
                    )}
                  </Fragment>
                );
              })
            )}
          </TableBody>
        </Table>
      </div>

      {/* Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <div className="border-t bg-muted/30 px-6 py-4">
          <div className="flex flex-col lg:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-3 order-3 lg:order-1">
              <label htmlFor="rows-per-page" className="text-sm whitespace-nowrap">
                {t("rowsPerPage")}
              </label>
              <Select
                value={String(perPage)}
                onValueChange={(value) => handlePerPageChange(Number(value))}
              >
                <SelectTrigger
                  id="rows-per-page"
                  className="w-fit whitespace-nowrap h-9"
                >
                  <SelectValue placeholder="Select rows" />
                </SelectTrigger>
                <SelectContent>
                  {[10, 20].map((option) => (
                    <SelectItem key={option} value={String(option)}>
                      {option}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="flex grow justify-center lg:justify-end text-sm whitespace-nowrap text-muted-foreground order-2">
              <p>
                <span className="text-foreground font-medium">
                  {(pagination.page - 1) * pagination.per_page + 1}-
                  {Math.min(pagination.page * pagination.per_page, pagination.total)}
                </span>{" "}
                {t("paginationOf", {
                  total: pagination.total,
                })}
              </p>
            </div>
            <div className="order-1 lg:order-3 flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(1)}
                disabled={!pagination.has_prev}
              >
                {t("first")}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(Math.max(1, pagination.page - 1))}
                disabled={!pagination.has_prev}
              >
                {t("prev")}
              </Button>
              <span className="text-sm">
                Page {pagination.page} of {pagination.total_pages}
              </span>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(pagination.page + 1)}
                disabled={!pagination.has_next}
              >
                {t("next")}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(pagination.total_pages)}
                disabled={!pagination.has_next}
              >
                {t("last")}
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Account Detail Modal */}
      <AccountDetailModal
        accountId={selectedAccountId}
        open={isModalOpen}
        onOpenChange={(open) => {
          setIsModalOpen(open);
          if (!open) {
            setSelectedAccountId(null);
          }
        }}
        onAccountUpdated={() => {
          // Refresh handled by react-query
        }}
      />
    </div>
  );
}
