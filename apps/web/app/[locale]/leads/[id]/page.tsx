import PageClient from "./page.client";

export const metadata = {
  title: "Lead Details | Tracio",
};

export default async function LeadDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <PageClient id={id} />;
}
