// File: frontend/src/components/layout/Footer.tsx

import Link from "next/link";

export function Footer() {
  return (
    <footer className="border-t border-border/40 bg-muted/20">
      <div className="container mx-auto flex flex-col items-center justify-between gap-4 py-10 md:h-24 md:flex-row md:py-0 px-4">
        <div className="flex flex-col items-center gap-4 md:flex-row md:gap-2 md:px-0">
          <p className="text-center text-sm leading-loose text-muted-foreground md:text-left">
            Dibuat dengan ❤️ Hak Cipta &copy; {new Date().getFullYear()} E-Commerce Baju.
          </p>
        </div>
        <div className="flex gap-4">
          <Link href="/terms" className="text-sm font-medium text-muted-foreground hover:text-foreground underline-offset-4 hover:underline">
            Syarat & Ketentuan
          </Link>
          <Link href="/privacy" className="text-sm font-medium text-muted-foreground hover:text-foreground underline-offset-4 hover:underline">
            Privasi
          </Link>
        </div>
      </div>
    </footer>
  );
}
