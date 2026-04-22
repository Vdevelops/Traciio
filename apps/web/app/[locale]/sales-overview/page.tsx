import { Suspense } from "react";
import { PageMotion } from "@/components/motion";
import { SalesOverviewPageClient } from "./sales-overview-page-client";

export const metadata = {
  title: "Sales Overview | Salesview",
};

export default function SalesOverviewPage() {
  return (
    <PageMotion className=" ">
      <Suspense fallback={null}>
        <SalesOverviewPageClient />
      </Suspense>
    </PageMotion>
  );
}

