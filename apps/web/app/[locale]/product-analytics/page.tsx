import { PageMotion } from "@/components/motion";
import { ProductAnalyticsPageClient } from "./product-analytics-page-client";

export const metadata = {
  title: "Product Analytics | Tracio",
};

export default function ProductAnalyticsPage() {
  return (
    <PageMotion className=" ">
      <ProductAnalyticsPageClient />
    </PageMotion>
  );
}
