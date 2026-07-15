// File: frontend/src/app/(admin)/dashboard/page.tsx
"use client";

import { useEffect, useState } from "react";
import api from "@/lib/axios";
import { Package, Layers, Users, ShoppingCart } from "lucide-react";

export default function AdminDashboard() {
  const [categoriesCount, setCategoriesCount] = useState(0);
  const [productsCount, setProductsCount] = useState(0);
  const [usersCount, setUsersCount] = useState(0);
  const [ordersCount, setOrdersCount] = useState(0);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [catRes, prodRes, userRes, orderRes] = await Promise.all([
          api.get("/categories/").catch(() => ({ data: { data: [] } })),
          api.get("/products/").catch(() => ({ data: { data: [] } })),
          api.get("/users/").catch(() => ({ data: { data: [] } })),
          api.get("/orders/all").catch(() => ({ data: { data: [] } }))
        ]);

        setCategoriesCount(catRes.data?.data?.length || 0);
        setProductsCount(prodRes.data?.data?.length || 0);
        setUsersCount(userRes.data?.data?.length || 0);
        setOrdersCount(orderRes.data?.data?.length || 0);
      } catch (error) {
        console.error("Gagal mengambil statistik:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchStats();
  }, []);

  if (isLoading) {
    return <div className="p-12 text-center text-muted-foreground">Memuat statistik...</div>;
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Ikhtisar Toko</h1>
        <p className="text-muted-foreground mt-1">Selamat datang di panel admin. Berikut adalah ringkasan data Anda hari ini.</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Stat Cards */}
        <div className="p-6 border rounded-xl bg-card shadow-sm flex items-center gap-4">
          <div className="p-3 bg-purple-100 text-purple-700 rounded-lg">
            <ShoppingCart className="h-6 w-6" />
          </div>
          <div>
            <p className="text-sm font-medium text-muted-foreground">Total Pesanan</p>
            <p className="text-2xl font-bold">{ordersCount}</p>
          </div>
        </div>

        <div className="p-6 border rounded-xl bg-card shadow-sm flex items-center gap-4">
          <div className="p-3 bg-blue-100 text-blue-700 rounded-lg">
            <Layers className="h-6 w-6" />
          </div>
          <div>
            <p className="text-sm font-medium text-muted-foreground">Total Kategori</p>
            <p className="text-2xl font-bold">{categoriesCount}</p>
          </div>
        </div>
        
        <div className="p-6 border rounded-xl bg-card shadow-sm flex items-center gap-4">
          <div className="p-3 bg-emerald-100 text-emerald-700 rounded-lg">
            <Package className="h-6 w-6" />
          </div>
          <div>
            <p className="text-sm font-medium text-muted-foreground">Total Produk</p>
            <p className="text-2xl font-bold">{productsCount}</p>
          </div>
        </div>

        <div className="p-6 border rounded-xl bg-card shadow-sm flex items-center gap-4">
          <div className="p-3 bg-amber-100 text-amber-700 rounded-lg">
            <Users className="h-6 w-6" />
          </div>
          <div>
            <p className="text-sm font-medium text-muted-foreground">Total Pengguna</p>
            <p className="text-2xl font-bold">{usersCount}</p>
          </div>
        </div>
      </div>

      <div className="p-8 border rounded-xl bg-muted/20 text-center mt-8">
        <h3 className="text-lg font-semibold mb-2">Pilih menu di sebelah kiri</h3>
        <p className="text-muted-foreground">
          Gunakan bilah navigasi untuk mengelola kategori, menambahkan produk baru, atau melihat daftar pelanggan.
        </p>
      </div>
    </div>
  );
}
