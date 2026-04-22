import { notFound } from "next/navigation";

export const metadata = {
  title: "Content | Salesview",
};

/**
 * Catch-all route untuk handle unmatched routes di dashboard
 * Ini akan trigger not-found.tsx di level yang sama (dalam route group)
 * yang akan otomatis di-wrap oleh DashboardLayout dari layout.tsx
 */
export default function DashboardCatchAll() {
  notFound();
}

