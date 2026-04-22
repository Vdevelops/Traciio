import { Suspense } from "react";
import { PageMotion } from "@/components/motion";
import { SalesRepDetailPageClient } from "./sales-rep-detail-page-client";

export const metadata = {
  title: "Sales Rep | Salesview",
};

interface SalesRepDetailPageProps {
  params: Promise<{ userId: string }>;
}

export default async function SalesRepDetailPage({ params }: SalesRepDetailPageProps) {
  const { userId } = await params;

  return (
    <PageMotion className=" ">
      <Suspense fallback={null}>
        <SalesRepDetailPageClient userId={userId} />
      </Suspense>
    </PageMotion>
  );
}

