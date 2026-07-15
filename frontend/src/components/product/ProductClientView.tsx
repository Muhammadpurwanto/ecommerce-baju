"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { ShoppingCart, Star, Truck, Plus, Minus, Loader2 } from "lucide-react";
import { useCartStore } from "@/store/useCartStore";

interface Category {
  id: number;
  name: string;
}

interface Product {
  id: number;
  name: string;
  slug: string;
  base_price: number;
  description: string;
  stock: number;
  image_url: string;
  category?: Category;
}

interface ProductClientViewProps {
  product: Product;
}

export function ProductClientView({ product }: ProductClientViewProps) {
  const router = useRouter();
  const { addToCart } = useCartStore();
  
  const defaultImage = product.image_url || "https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?auto=format&fit=crop&w=1000&q=80";
  const [quantity, setQuantity] = useState<number>(1);
  const [isAdding, setIsAdding] = useState(false);

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(price);
  };

  const handleAddToCart = async () => {
    if (product.stock <= 0) {
      alert("Maaf, produk sedang habis.");
      return;
    }

    setIsAdding(true);
    try {
      await addToCart(product.id, quantity);
      alert(`Berhasil menambahkan ${quantity} item ke keranjang!`);
    } catch (err: any) {
      console.error(err);
    } finally {
      setIsAdding(false);
    }
  };

  const maxQuantity = product.stock > 0 ? product.stock : 1;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-8 lg:gap-12 animate-in fade-in duration-300">
      {/* Gallery Section */}
      <div className="space-y-4">
        <div className="aspect-[4/5] md:aspect-square overflow-hidden rounded-2xl bg-slate-50 border relative shadow-sm">
          <img 
            src={defaultImage} 
            alt={product.name} 
            className="w-full h-full object-cover transition-transform duration-300 hover:scale-102"
          />
        </div>
      </div>

      {/* Details & Action Section */}
      <div className="flex flex-col">
        <div className="mb-2 text-sm font-semibold text-primary uppercase tracking-wider">
          {product.category?.name || "Kategori"}
        </div>
        
        <h1 className="text-3xl md:text-4xl font-extrabold tracking-tight mb-4 text-slate-900">
          {product.name}
        </h1>

        <div className="flex items-center gap-4 mb-6">
          <div className="flex items-center text-amber-500">
            <Star className="h-5 w-5 fill-current" />
            <Star className="h-5 w-5 fill-current" />
            <Star className="h-5 w-5 fill-current" />
            <Star className="h-5 w-5 fill-current" />
            <Star className="h-5 w-5 fill-current opacity-30" />
            <span className="text-slate-600 text-sm ml-2 font-medium">4.8 (12 Ulasan)</span>
          </div>
          <div className="w-px h-5 bg-border"></div>
          {product.stock > 0 ? (
            <span className="text-sm font-semibold text-emerald-600 bg-emerald-50 px-3 py-1 rounded-full border border-emerald-100">
              Stok: {product.stock} Tersedia
            </span>
          ) : (
            <span className="text-sm font-semibold text-rose-600 bg-rose-50 px-3 py-1 rounded-full border border-rose-100">
              Stok Habis
            </span>
          )}
        </div>

        {/* Price Display */}
        <div className="text-4xl font-extrabold text-slate-900 mb-8">
          {formatPrice(product.base_price)}
        </div>

        {/* Quantity and Actions */}
        <div className="space-y-4 mt-auto border-t pt-6">
          {product.stock > 0 && (
            <div className="flex items-center gap-4">
              <span className="text-xs font-bold text-slate-500 uppercase tracking-wider">Jumlah</span>
              <div className="flex items-center border border-slate-200 rounded-xl bg-slate-50/50 p-1">
                <button
                  type="button"
                  disabled={quantity <= 1}
                  onClick={() => setQuantity(prev => Math.max(1, prev - 1))}
                  className="p-2 rounded-lg text-slate-500 hover:text-slate-900 hover:bg-white active:scale-95 disabled:opacity-30 transition-all"
                >
                  <Minus className="h-4 w-4" />
                </button>
                <span className="w-12 text-center font-bold text-slate-800">{quantity}</span>
                <button
                  type="button"
                  disabled={quantity >= maxQuantity}
                  onClick={() => setQuantity(prev => Math.min(maxQuantity, prev + 1))}
                  className="p-2 rounded-lg text-slate-500 hover:text-slate-900 hover:bg-white active:scale-95 disabled:opacity-30 transition-all"
                >
                  <Plus className="h-4 w-4" />
                </button>
              </div>
              <span className="text-xs text-slate-400">Maks. {maxQuantity} item</span>
            </div>
          )}

          <div className="flex gap-4">
            <button
              onClick={handleAddToCart}
              disabled={isAdding || product.stock <= 0}
              className="flex-1 bg-primary text-primary-foreground h-14 rounded-xl font-bold text-lg flex items-center justify-center gap-2 hover:bg-primary/90 transition-all shadow-lg shadow-primary/20 hover:shadow-primary/40 active:scale-[0.98] disabled:opacity-50 disabled:pointer-events-none"
            >
              {isAdding ? (
                <Loader2 className="h-5 w-5 animate-spin" />
              ) : (
                <>
                  <ShoppingCart className="h-5 w-5" />
                  {product.stock <= 0 ? "Stok Habis" : "Tambah ke Keranjang"}
                </>
              )}
            </button>
          </div>
        </div>

        {/* Shipping details */}
        <div className="mt-8 pt-8 border-t space-y-4">
          <div className="flex items-center gap-3 text-slate-600">
            <div className="h-10 w-10 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 flex-shrink-0">
              <Truck className="h-5 w-5" />
            </div>
            <div>
              <h4 className="font-semibold text-slate-900 text-sm">Pengiriman Cepat &amp; Aman</h4>
              <p className="text-sm">Bekerjasama dengan kurir terpercaya JNE dan POS</p>
            </div>
          </div>
        </div>

        {/* Description */}
        <div className="mt-8 pt-8 border-t">
          <h3 className="text-lg font-bold text-slate-900 mb-4">Deskripsi Produk</h3>
          <div className="prose prose-sm text-slate-600 max-w-none leading-relaxed">
            {product.description ? (
              <p>{product.description}</p>
            ) : (
              <p>Pakaian berkualitas tinggi yang dirancang untuk kenyamanan maksimal dan gaya yang tak lekang oleh waktu. Cocok digunakan untuk berbagai acara, baik formal maupun kasual.</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
