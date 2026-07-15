// File: frontend/src/app/auth/callback/page.tsx
"use client";

import { useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Cookies from "js-cookie";

function AuthCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const accessToken = searchParams.get("access_token");
    const refreshToken = searchParams.get("refresh_token");

    if (accessToken && refreshToken) {
      // Simpan JWT ke Cookies
      Cookies.set("access_token", accessToken, { expires: 1 });
      Cookies.set("refresh_token", refreshToken, { expires: 7 });

      // Cek role dari token
      try {
        const payload = JSON.parse(atob(accessToken.split('.')[1]));
        if (payload?.role === "admin") {
          router.push("/dashboard");
        } else {
          router.push("/");
        }
      } catch (e) {
        router.push("/");
      }
      router.refresh();
    } else {
      // Jika token tidak ada, mungkin terjadi error di Google
      alert("Autentikasi Google gagal. Silakan coba lagi.");
      router.push("/login");
    }
  }, [searchParams, router]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="flex flex-col items-center gap-4">
        <div className="h-12 w-12 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
        <p className="text-lg font-medium text-muted-foreground animate-pulse">
          Mengautentikasi dengan Google...
        </p>
      </div>
    </div>
  );
}

export default function AuthCallback() {
  return (
    <Suspense fallback={<div className="flex min-h-screen items-center justify-center"><div className="h-12 w-12 animate-spin rounded-full border-4 border-primary border-t-transparent"></div></div>}>
      <AuthCallbackContent />
    </Suspense>
  );
}
