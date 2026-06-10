const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
const STORAGE_PUBLIC_URL = process.env.NEXT_PUBLIC_STORAGE_PUBLIC_URL || "";

function joinUrl(baseUrl: string, path: string): string {
  const base = baseUrl.replace(/\/+$/, "");
  const normalizedPath = path.replace(/^\/+/, "");
  const baseLastSegment = base.split("/").filter(Boolean).at(-1);
  const [pathFirstSegment, ...pathRest] = normalizedPath.split("/");

  if (baseLastSegment && pathFirstSegment && baseLastSegment === pathFirstSegment) {
    return `${base}/${pathRest.join("/")}`;
  }

  return `${base}/${normalizedPath}`;
}

export function getVisitReportPhotoUrl(photoUrl: string): string {
  const proxyBaseUrl = joinUrl(API_BASE_URL, "/api/v1/files/image");

  if (photoUrl.startsWith("http://") || photoUrl.startsWith("https://")) {
    try {
      const url = new URL(photoUrl);
      const apiUrl = new URL(API_BASE_URL);
      const storageUrl = STORAGE_PUBLIC_URL ? new URL(STORAGE_PUBLIC_URL) : null;

      if (
        (url.origin === apiUrl.origin || url.origin === storageUrl?.origin) &&
        url.pathname.startsWith("/uploads/")
      ) {
        return joinUrl(proxyBaseUrl, url.pathname);
      }
    } catch {
      return photoUrl;
    }

    return photoUrl;
  }

  const cleanUrl = photoUrl.startsWith("/") ? photoUrl : `/${photoUrl}`;

  if (cleanUrl.startsWith("/uploads/")) {
    return joinUrl(proxyBaseUrl, cleanUrl);
  }

  return joinUrl(API_BASE_URL, cleanUrl);
}
