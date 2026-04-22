'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/ui/table';
import {
  useCustomerPurchaseHistory,
  useCustomerProductAnalytics
} from '../hooks/useCustomerPurchase';
import { formatCurrency, formatDate } from '@/lib/utils';
import { Package, ShoppingCart, TrendingUp } from 'lucide-react';

interface CustomerPurchaseHistoryProps {
  accountId: string;
}

export function CustomerPurchaseHistory({ accountId }: CustomerPurchaseHistoryProps) {
  const { data: history, isLoading: isLoadingHistory } = useCustomerPurchaseHistory(accountId);
  const { data: analytics, isLoading: isLoadingAnalytics } = useCustomerProductAnalytics(accountId);

  if (isLoadingHistory || isLoadingAnalytics) {
    return <PurchaseHistorySkeleton />;
  }

  const purchases = history?.data || [];
  const productAnalytics = analytics || [];

  return (
    <div className="space-y-6">
      {/* Product Analytics Summary */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2">
            <TrendingUp size={20} />
            Product Purchase Analytics
          </CardTitle>
        </CardHeader>
        <CardContent>
          {productAnalytics.length === 0 ? (
            <p className="text-muted-foreground text-center py-8">
              No purchase data available yet
            </p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {productAnalytics.slice(0, 6).map((product) => (
                <Card key={product.product_id} className="bg-muted/50">
                  <CardContent className="p-4">
                    <h4 className="font-medium truncate">{product.product_name}</h4>
                    <p className="text-sm text-muted-foreground">
                      {product.product_category_name}
                    </p>
                    <div className="mt-3 space-y-1">
                      <div className="flex justify-between text-sm">
                        <span>Quantity:</span>
                        <span className="font-medium">{product.total_quantity_purchased}</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span>Revenue:</span>
                        <span className="font-medium">{product.total_amount_formatted}</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span>Purchases:</span>
                        <span className="font-medium">{product.purchase_count}x</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Purchase History Table */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2">
            <ShoppingCart size={20} />
            Purchase History
          </CardTitle>
        </CardHeader>
        <CardContent>
          {purchases.length === 0 ? (
            <p className="text-muted-foreground text-center py-8">
              No purchase history available
            </p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Purchase #</TableHead>
                  <TableHead>Date</TableHead>
                  <TableHead>Products</TableHead>
                  <TableHead>Total Amount</TableHead>
                  <TableHead>Sales Rep</TableHead>
                  <TableHead>Source</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {purchases.map((purchase) => (
                  <TableRow key={purchase.id}>
                    <TableCell>
                      <Badge variant="outline">#{purchase.purchase_number}</Badge>
                    </TableCell>
                    <TableCell>{formatDate(purchase.purchase_date)}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Package size={16} className="text-muted-foreground" />
                        <span>{purchase.total_items} items</span>
                      </div>
                    </TableCell>
                    <TableCell className="font-medium">
                      {purchase.total_amount_formatted}
                    </TableCell>
                    <TableCell>{purchase.sales_rep_name || '-'}</TableCell>
                    <TableCell>
                      <Badge variant="secondary">
                        {purchase.source_type.replace(/_/g, ' ')}
                      </Badge>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

function PurchaseHistorySkeleton() {
  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <div className="h-6 w-48 bg-gray-200 rounded animate-pulse" />
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 gap-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-32 bg-gray-200 rounded animate-pulse" />
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
