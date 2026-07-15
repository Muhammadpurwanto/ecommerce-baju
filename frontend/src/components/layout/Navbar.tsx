// File: frontend/src/components/layout/Navbar.tsx
"use client"; // Tambahkan ini karena kita akan menggunakan Hook (Zustand & React)

import { useEffect, useState } from "react";
import Link from "next/link";
import { ShoppingCart, User, Search, LogOut } from "lucide-react";
import { useCartStore } from "@/store/useCartStore";
import Cookies from "js-cookie";

import { useRouter } from "next/navigation";
import api from "@/lib/axios";

export function Navbar() {
  const { items, fetchCart } = useCartStore();
  const [isClient, setIsClient] = useState(false);
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const router = useRouter();

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/products?search=${encodeURIComponent(searchQuery)}`);
      setIsSearchOpen(false);
      setSearchQuery("");
    }
  };
  
  // Deteksi apakah user sedang login dengan melihat ketersediaan cookie
  const isLoggedIn = isClient && !!Cookies.get("access_token");

  const [userName, setUserName] = useState("");

  useEffect(() => {
    setIsClient(true);
    // Jika user punya token, tarik data keranjangnya saat navbar dirender
    if (Cookies.get("access_token")) {
      fetchCart();
      // Tarik nama profil secara diam-diam (background)
      api.get("/users/profile")
        .then((res) => {
          if (res.data.success) setUserName(res.data.data.name);
          console.log(userName)
        })
        .catch(() => {});
    }
  }, [fetchCart]);

  // Hitung total kuantitas barang di keranjang
  const totalItems = items.reduce((acc, item) => acc + item.quantity, 0);

  const handleLogout = async () => {
    try {
      const refreshToken = Cookies.get("refresh_token");
      if (refreshToken) {
        await api.post("/auth/logout", { refresh_token: refreshToken });
      }
    } catch (error) {
      console.error("Logout failed on server:", error);
    } finally {
      Cookies.remove("access_token");
      Cookies.remove("refresh_token");
      window.location.href = "/";
    }
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        
        {/* Kiri: Logo & Navigasi Desktop */}
        <div className="flex gap-6 md:gap-10">
          <Link href="/" className="flex items-center space-x-2">
            <span className="font-bold text-xl inline-block tracking-tight">
              E-Commerce<span className="text-primary">Baju</span>
            </span>
          </Link>
          <nav className="hidden md:flex gap-6">
            <Link href="/products" className="flex items-center text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">Produk</Link>
            <Link href="/categories" className="flex items-center text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">Kategori</Link>
          </nav>
        </div>

        {/* Kanan: Icons & Auth */}
        <div className="flex items-center gap-4">
          {isSearchOpen ? (
            <form onSubmit={handleSearch} className="flex items-center animate-in fade-in slide-in-from-right-4">
              <input
                type="text"
                placeholder="Cari produk..."
                autoFocus
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="text-sm px-3 py-1 bg-muted rounded-l-md border-0 focus:ring-1 focus:ring-primary h-8 w-[150px] sm:w-[200px]"
                onBlur={() => !searchQuery && setIsSearchOpen(false)}
              />
              <button type="submit" className="h-8 px-2 bg-primary text-primary-foreground rounded-r-md">
                <Search className="h-4 w-4" />
              </button>
            </form>
          ) : (
            <button onClick={() => setIsSearchOpen(true)} className="text-muted-foreground hover:text-foreground transition-colors">
              <Search className="h-5 w-5" />
            </button>
          )}
          
          <Link href="/cart" className="relative text-muted-foreground hover:text-foreground transition-colors">
            <ShoppingCart className="h-5 w-5" />
            {isClient && totalItems > 0 && (
              <span className="absolute -top-2 -right-2 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] font-bold text-primary-foreground">
                {totalItems}
              </span>
            )}
          </Link>
          
          <div className="h-5 w-px bg-border mx-1 hidden sm:block"></div>
          
          {/* Tampilan kondisional Login / Logout */}
          {isClient && !isLoggedIn ? (
            <Link href="/login" className="hidden sm:flex items-center gap-2 text-sm font-medium hover:text-primary transition-colors">
              <User className="h-5 w-5" />
              <span>Login</span>
            </Link>
          ) : isClient && isLoggedIn ? (
            <div className="hidden sm:flex items-center gap-4">
              {(() => {
                let isAdmin = false;
                try {
                  const token = Cookies.get("access_token");
                  if (token) {
                    const payload = JSON.parse(atob(token.split('.')[1]));
                    if (payload?.role === "admin") isAdmin = true;
                  }
                } catch (e) {}

                return isAdmin ? (
                  <Link href="/dashboard" className="text-sm font-bold text-primary hover:underline transition-colors mr-2">
                    Admin Dashboard
                  </Link>
                ) : null;
              })()}
              <Link href="/profile" className="flex items-center gap-2 text-sm font-medium hover:text-primary transition-colors">
                <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                  <User className="h-4 w-4" />
                </div>
                {/* Menampilkan nama user atau teks "Profil" jika namanya belum termuat */}
                <span className="max-w-[120px] truncate">{userName || "Profil"}</span>
              </Link>
              <button onClick={handleLogout} title="Keluar" className="flex items-center gap-2 text-sm font-medium text-destructive hover:opacity-80 transition-opacity">
                <LogOut className="h-5 w-5" />
              </button>
            </div>
          ) : null}
        </div>

      </div>
    </header>
  );
}
