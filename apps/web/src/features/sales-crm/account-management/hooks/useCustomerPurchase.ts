"use client";

import { useQuery } from "@tanstack/react-query";
import { customerPurchaseService } from "../services/customer-purchase.service";

export function useCustomerPurchaseHistory(
  accountId: string,
  page = 1,
  perPage = 20,
) {
  return useQuery({
    queryKey: ["customer-purchase-history", accountId, page, perPage],
    queryFn: () =>
      customerPurchaseService.getAccountPurchaseHistory(
        accountId,
        page,
        perPage,
      ),
    enabled: !!accountId,
  });
}

export function useCustomerProductAnalytics(accountId: string) {
  return useQuery({
    queryKey: ["customer-product-analytics", accountId],
    queryFn: () =>
      customerPurchaseService.getAccountProductAnalytics(accountId),
    enabled: !!accountId,
  });
}

export function useCustomerPurchaseSummary(accountId: string) {
  return useQuery({
    queryKey: ["customer-purchase-summary", accountId],
    queryFn: () => customerPurchaseService.getAccountPurchaseSummary(accountId),
    enabled: !!accountId,
  });
}
