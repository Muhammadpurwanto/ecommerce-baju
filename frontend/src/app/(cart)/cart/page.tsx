"use client";

import { useEffect, useMemo } from "react";
import { useCartStore } from "@/store/useCartStore";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Trash2, Minus, Plus, ShoppingBag } from "lucide-react";

export default function CartPage() {
  const router = useRouter();
  const { items, productsCache, isLoading, fetchCart, removeItem, updateItem, selectedItems, toggleSelectItem, selectAll, clearSelection } = useCartStore();

  useEffect(() => {
    fetchCart();
  }, [fetchCart]);

  const { totalItems, grandTotal } = useMemo(() => {
    let tItems = 0;
    let gTotal = 0;
    items.forEach(item => {
      if (selectedItems.includes(item.id)) {
        tItems += Number(item.quantity);
        const product = productsCache[item.product_id];
        if (product) {
          const price = product.base_price;
          gTotal += price * Number(item.quantity);
        }
      }
    });
    return { totalItems: tItems, grandTotal: gTotal };
  }, [items, productsCache, selectedItems]);

  const formatPrice = (price: number) => new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", minimumFractionDigits: 0 }).format(price);

  return (
    <div className="bg-slate-50 min-h-screen pb-20">
      <div className="container mx-auto px-4 py-8 md:py-12">
        <div className="flex items-center gap-3 mb-8">
          <div className="h-12 w-12 bg-primary/10 text-primary flex items-center justify-center rounded-2xl">
            <ShoppingBag className="h-6 w-6" />
          </div>
          <h1 className="text-3xl font-extrabold tracking-tight text-slate-900">Keranjang Belanja</h1>
        </div>

        {isLoading && items.length === 0 ? (
          <div className="flex justify-center py-20">
            <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
          </div>
        ) : items.length === 0 ? (
          <div className="flex flex-col items-center justify-center p-12 bg-white border border-dashed border-slate-300 rounded-3xl shadow-sm text-center">
            <div className="h-24 w-24 bg-slate-100 text-slate-400 rounded-full flex items-center justify-center mb-6">
              <ShoppingBag className="h-10 w-10" />
            </div>
            <h2 className="text-2xl font-bold text-slate-900 mb-2">Keranjang Anda Masih Kosong</h2>
            <p className="text-slate-500 max-w-md mb-8">Sepertinya Anda belum menambahkan produk apapun ke keranjang. Mari mulai berbelanja dan temukan gaya terbaik Anda!</p>
            <Link href="/products" className="rounded-xl bg-primary px-8 py-4 font-bold text-primary-foreground transition-all hover:bg-primary/90 shadow-lg shadow-primary/20 hover:scale-105 active:scale-95">
              Mulai Belanja Sekarang
            </Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            
            {/* List Barang */}
            <div className="lg:col-span-2 flex flex-col gap-4">
              {items.length > 0 && (
                <div className="flex items-center gap-3 p-4 border rounded-2xl bg-white shadow-sm">
                  <input
                    type="checkbox"
                    className="h-5 w-5 rounded border-slate-300 text-primary focus:ring-primary"
                    checked={selectedItems.length === items.length && items.length > 0}
                    onChange={(e) => {
                      if (e.target.checked) selectAll();
                      else clearSelection();
                    }}
                  />
                  <span className="font-semibold text-slate-700">Pilih Semua ({items.length})</span>
                </div>
              )}
              {items.map((item) => {
                const product = productsCache[item.product_id];
                // Jika produk masih loading dari cache
                if (!product) {
                  return (
                    <div key={item.id} className="flex p-4 border rounded-2xl bg-white shadow-sm animate-pulse">
                      <div className="w-24 h-24 bg-slate-200 rounded-xl mr-4"></div>
                      <div className="flex-1 py-2"><div className="h-4 bg-slate-200 rounded w-3/4 mb-3"></div><div className="h-3 bg-slate-200 rounded w-1/2"></div></div>
                    </div>
                  );
                }

                const price = product.base_price;
                const imageUrl = product.image_url || "https://via.placeholder.com/150";

                return (
                  <div key={item.id} className="flex flex-col sm:flex-row items-start sm:items-center p-4 border rounded-2xl bg-white shadow-sm transition-all hover:shadow-md gap-4">
                    {/* Checkbox */}
                    <div className="pt-2 sm:pt-0">
                      <input
                        type="checkbox"
                        className="h-5 w-5 rounded border-slate-300 text-primary focus:ring-primary"
                        checked={selectedItems.includes(item.id)}
                        onChange={() => toggleSelectItem(item.id)}
                      />
                    </div>

                    {/* Gambar */}
                    <div className="h-24 w-24 sm:h-32 sm:w-32 flex-shrink-0 overflow-hidden rounded-xl bg-slate-100 border">
                      <img src={imageUrl} alt={product.name} className="h-full w-full object-cover" />
                    </div>

                    {/* Info */}
                    <div className="flex-1 flex flex-col min-w-0">
                      <Link href={`/products/${product.slug}`} className="font-bold text-lg text-slate-900 line-clamp-2 hover:text-primary transition-colors">
                        {product.name}
                      </Link>
                      <p className="text-sm text-slate-500 mt-1">
                        Brand: <span className="font-medium text-slate-700">{product.brand || "Unbranded"}</span>
                      </p>
                      <div className="font-bold text-lg text-slate-900 mt-2">{formatPrice(price)}</div>
                    </div>

                    {/* Aksi Kuantitas & Hapus */}
                    <div className="flex sm:flex-col items-center sm:items-end justify-between w-full sm:w-auto mt-4 sm:mt-0 gap-4">
                      <button 
                        onClick={() => removeItem(item.id)}
                        className="p-2 text-rose-500 bg-rose-50 hover:bg-rose-100 hover:text-rose-600 rounded-xl transition-colors order-2 sm:order-1"
                        title="Hapus dari keranjang"
                      >
                        <Trash2 className="h-5 w-5" />
                      </button>

                      <div className="flex items-center border rounded-xl bg-slate-50 order-1 sm:order-2">
                        <button 
                          onClick={() => updateItem(item.id, Math.max(1, Number(item.quantity) - 1))}
                          disabled={item.quantity <= 1}
                          className="p-2 text-slate-500 hover:text-slate-900 disabled:opacity-30 transition-colors"
                        >
                          <Minus className="h-4 w-4" />
                        </button>
                        <span className="w-8 text-center font-semibold text-sm">{item.quantity}</span>
                        <button 
                          onClick={() => updateItem(item.id, Number(item.quantity) + 1)}
                          className="p-2 text-slate-500 hover:text-slate-900 transition-colors"
                        >
                          <Plus className="h-4 w-4" />
                        </button>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>

            {/* Ringkasan Belanja (Checkout Box) */}
            <div className="flex flex-col p-6 lg:p-8 border rounded-3xl bg-white shadow-lg shadow-slate-200/50 h-fit sticky top-24">
              <h2 className="text-xl font-bold mb-6 text-slate-900">Ringkasan Belanja</h2>
              
              <div className="space-y-4 mb-6">
                <div className="flex justify-between text-slate-600">
                  <span>Total Harga ({totalItems} barang)</span>
                  <span className="font-medium">{formatPrice(grandTotal)}</span>
                </div>
                <div className="flex justify-between text-slate-600">
                  <span>Total Diskon Barang</span>
                  <span className="font-medium text-emerald-500">- Rp 0</span>
                </div>
              </div>
              
              <div className="border-t border-dashed my-4"></div>
              
              <div className="flex justify-between items-center mb-8">
                <span className="font-bold text-slate-900">Total Tagihan</span>
                <span className="font-extrabold text-2xl text-primary">{formatPrice(grandTotal)}</span>
              </div>
              
              <button 
                onClick={() => {
                  if (totalItems > 0) router.push("/checkout");
                }}
                disabled={totalItems === 0}
                className="w-full flex items-center justify-center rounded-xl bg-primary px-4 py-4 font-bold text-lg text-primary-foreground hover:bg-primary/90 transition-all shadow-lg shadow-primary/20 hover:shadow-primary/40 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Beli ({totalItems})
              </button>
            </div>

          </div>
        )}
      </div>
    </div>
  );
}
