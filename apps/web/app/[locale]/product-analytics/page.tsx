import { PageMotion } from "@/components/motion";
import { ProductAnalyticsPageClient } from "./product-analytics-page-client";

export const metadata = {
  title: "Product Analytics | Salesview",
};

export default function ProductAnalyticsPage() {
  return (
    <PageMotion className=" ">
      <ProductAnalyticsPageClient />
    </PageMotion>
  );
}
