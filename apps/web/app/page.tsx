import { redirect } from "next/navigation";
import { routing } from "@/i18n/routing";

export const metadata = {
  title: "Home | Tracio",
};

export default function RootRedirectPage() {
  // Redirect root "/" ke default locale login, contoh: "/en/login"
  redirect(`/${routing.defaultLocale}/login`);
}


