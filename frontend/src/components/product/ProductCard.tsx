// File: frontend/src/components/product/ProductCard.tsx
"use client";

import Link from "next/link";
import { ShoppingCart, Loader2 } from "lucide-react";
import { useState } from "react";
import api from "@/lib/axios";
import { useRouter } from "next/navigation";

export function ProductCard({ product }: { product: any }) {
  // Format harga ke Rupiah
  const formattedPrice = new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(product.base_price);

  // Ambil gambar utama jika ada, jika tidak pakai placeholder
  const imageUrl = product.image_url || (product.images?.length > 0 
    ? product.images.find((img: any) => img.is_primary)?.image_url || product.images[0].image_url 
    : "https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?auto=format&fit=crop&w=500&q=80");

  const [isAdding, setIsAdding] = useState(false);
  const router = useRouter();

  const handleAddToCart = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    setIsAdding(true);
    try {
      const res = await api.post("/carts/items", {
        product_id: product.id,
        quantity: 1,
      });

      if (res.data.success) {
        alert("Berhasil ditambahkan ke keranjang!");
      }
    } catch (error: any) {
      if (error.response?.status === 401) {
        alert("Silakan login terlebih dahulu untuk menambah ke keranjang.");
        router.push("/login");
      } else {
        alert(error.response?.data?.message || "Gagal menambahkan ke keranjang");
      }
    } finally {
      setIsAdding(false);
    }
  };

  return (
    <div className="group relative flex flex-col overflow-hidden rounded-xl border bg-card text-card-foreground shadow-sm transition-all hover:shadow-md">
      {/* Gambar Produk */}
      <Link href={`/products/${product.slug}`} className="aspect-[4/5] overflow-hidden bg-muted">
        <img
          src={imageUrl}
          alt={product.name}
          className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
        />
      </Link>

      {/* Detail Produk */}
      <div className="flex flex-1 flex-col p-4">
        <div className="flex items-start justify-between gap-2">
          <div>
            <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">
              {product.category?.name || "Kategori"}
            </p>
            <Link href={`/products/${product.slug}`}>
              <h3 className="font-semibold line-clamp-1 group-hover:text-primary transition-colors">
                {product.name}
              </h3>
            </Link>
          </div>
        </div>

        <div className="mt-auto pt-4 flex items-center justify-between">
          <p className="font-bold text-lg">{formattedPrice}</p>
          <button 
            onClick={handleAddToCart}
            disabled={isAdding}
            className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10 text-primary transition-all hover:bg-primary hover:text-primary-foreground hover:scale-110 active:scale-95 disabled:opacity-50 disabled:pointer-events-none"
            title="Tambah ke Keranjang"
          >
            {isAdding ? <Loader2 className="h-5 w-5 animate-spin" /> : <ShoppingCart className="h-5 w-5" />}
          </button>
        </div>
      </div>
    </div>
  );
}
