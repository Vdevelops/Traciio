import { redirect } from "next/navigation";

export default async function LegacyPipelineDetailRedirect({
  params,
}: {
  params: Promise<{ locale: string; id: string }>;
}) {
  const { locale, id } = await params;
  redirect(`/${locale}/pipeline?dealId=${encodeURIComponent(id)}`);
}
