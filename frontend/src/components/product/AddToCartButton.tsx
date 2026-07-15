"use client";

import { useState } from "react";
import { ShoppingCart, Loader2 } from "lucide-react";
import api from "@/lib/axios";
import { useRouter } from "next/navigation";

interface AddToCartButtonProps {
  product: any;
}

export function AddToCartButton({ product }: AddToCartButtonProps) {
  const [isAdding, setIsAdding] = useState(false);
  const router = useRouter();

  const handleAddToCart = async () => {
    setIsAdding(true);
    try {
      const res = await api.post("/carts/items", {
        product_id: product.id,
        quantity: 1,
      });

      if (res.data.success) {
        alert("Produk berhasil ditambahkan ke keranjang!");
      } else {
        alert(res.data.message || "Gagal menambahkan ke keranjang");
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
    <button 
      onClick={handleAddToCart}
      disabled={isAdding}
      className="flex-1 bg-primary text-primary-foreground h-14 rounded-xl font-bold text-lg flex items-center justify-center gap-2 hover:bg-primary/90 transition-all shadow-lg shadow-primary/20 hover:shadow-primary/40 active:scale-[0.98] disabled:opacity-70 disabled:pointer-events-none"
    >
      {isAdding ? (
        <Loader2 className="h-5 w-5 animate-spin" />
      ) : (
        <>
          <ShoppingCart className="h-5 w-5" />
          Tambah ke Keranjang
        </>
      )}
    </button>
  );
}
