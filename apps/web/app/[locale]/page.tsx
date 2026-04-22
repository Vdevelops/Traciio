import { redirect } from "next/navigation";

// Ketika user buka "/en" atau "/id", redirect ke "/[locale]/login"
export const metadata = {
  title: "Home | Salesview",
};
export default async function LocaleRootRedirect({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  redirect(`/${locale}/login`);
}


