import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";
import { AreaMappingManagement } from "@/features/sales-crm/area-mapping/components";

// Import messages
import enMessages from "@/features/sales-crm/area-mapping/i18n/messages/en.json";
import idMessages from "@/features/sales-crm/area-mapping/i18n/messages/id.json";

export const metadata = {
  title: "Area Mapping | Salesview",
  description: "Manage territories, track location captures, and analyze coverage",
};

interface AreaMappingPageProps {
  params: Promise<{
    locale: string;
  }>;
}

async function AreaMappingPageContent({ locale }: { locale: string }) {
  const messages = await getMessages();

  // Merge area mapping messages with existing messages
  const areaMappingMessages = locale === "id" ? idMessages : enMessages;

  return (
    <NextIntlClientProvider locale={locale} messages={{ ...messages, ...areaMappingMessages }}>
      <AreaMappingManagement />
    </NextIntlClientProvider>
  );
}

export default async function AreaMappingPage({ params }: AreaMappingPageProps) {
  const { locale } = await params;
  
  return <AreaMappingPageContent locale={locale} />;
}
