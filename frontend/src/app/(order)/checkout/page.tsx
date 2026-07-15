// File: frontend/src/app/(customer)/checkout/page.tsx
"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useCartStore } from "@/store/useCartStore";
import api from "@/lib/axios";

export default function CheckoutPage() {
  const router = useRouter();
  const { items, fetchCart, selectedItems, productsCache } = useCartStore();
  const checkoutItems = items.filter(item => selectedItems.includes(item.id));
  const [isLoading, setIsLoading] = useState(false);

  // Form State
  const [addresses, setAddresses] = useState<any[]>([]);
  const [address, setAddress] = useState("");
  const [courier, setCourier] = useState("JNE");
  const [notes, setNotes] = useState("");

  useEffect(() => {
    fetchCart();
    
    // Fetch user addresses
    const fetchAddresses = async () => {
      try {
        const res = await api.get("/users/addresses");
        if (res.data.success && res.data.data) {
          setAddresses(res.data.data);
          // Set default address
          const defaultAddr = res.data.data.find((a: any) => a.is_default);
          if (defaultAddr) {
            setAddress(`${defaultAddr.detail}, ${defaultAddr.district}, ${defaultAddr.city}, ${defaultAddr.province} ${defaultAddr.postal_code}`);
          } else if (res.data.data.length > 0) {
            const firstAddr = res.data.data[0];
            setAddress(`${firstAddr.detail}, ${firstAddr.district}, ${firstAddr.city}, ${firstAddr.province} ${firstAddr.postal_code}`);
          }
        }
      } catch (error) {
        console.error("Gagal mengambil alamat:", error);
      }
    };
    fetchAddresses();
  }, [fetchCart]);

  // Jika tidak ada barang yang dicheckout, kembalikan ke home
  useEffect(() => {
    if (checkoutItems.length === 0 && !isLoading) {
      router.push("/");
    }
  }, [checkoutItems.length, isLoading, router]);

  const handleCheckout = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    // Siapkan array items dengan harga varian produk dinamis dari cache
    const orderItems = checkoutItems.map(item => {
      const product = productsCache[item.product_id];
      const price = product ? product.base_price : 150000;
      return {
        product_id: item.product_id,
        quantity: Number(item.quantity),
        price: price
      };
    });

    const payload = {
      shipping_cost: courier === "JNE" ? 20000 : 15000,
      shipping_address: address,
      courier: courier,
      notes: notes,
      items: orderItems,
    };

    try {
      const res = await api.post("/orders/", payload);
      if (res.data.success) {
        alert(`Pesanan Berhasil Dibuat! Nomor Invoice: ${res.data.data.order_number}`);
        
        // Refresh state keranjang agar jadi 0
        await fetchCart();
        
        // Gunakan payment_url langsung dari response order yang dibuat oleh Saga Orchestration
        const paymentUrl = res.data.data.payment_url;
        if (paymentUrl) {
          window.location.href = paymentUrl;
          return;
        }
        
        // Fallback jika tidak ada link payment
        router.push("/");
      }
    } catch (error: any) {
      console.error(error);
      alert(error.response?.data?.message || "Gagal membuat pesanan.");
    } finally {
      setIsLoading(false);
    }
  };

  if (checkoutItems.length === 0) return null;

  return (
    <div className="container mx-auto px-4 py-12">
      <h1 className="text-3xl font-bold tracking-tight mb-8">Checkout</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
        {/* Form Pengiriman */}
        <div>
          <h2 className="text-xl font-semibold mb-6 border-b pb-2">Informasi Pengiriman</h2>
          <form id="checkout-form" onSubmit={handleCheckout} className="space-y-5">
            <div className="space-y-2">
              <label className="text-sm font-medium">Pilih Alamat Pengiriman</label>
              {addresses.length > 0 ? (
                <select
                  required
                  value={address}
                  onChange={(e) => setAddress(e.target.value)}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                  <option value="" disabled>-- Pilih Alamat Tersimpan --</option>
                  {addresses.map((addr) => {
                    const fullAddress = `${addr.detail}, ${addr.district}, ${addr.city}, ${addr.province} ${addr.postal_code}`;
                    return (
                      <option key={addr.id} value={fullAddress}>
                        {addr.label} - {addr.recipient} ({addr.city})
                      </option>
                    );
                  })}
                </select>
              ) : (
                <div className="text-sm text-amber-600 bg-amber-50 p-4 rounded-md border border-amber-200">
                  <p className="mb-2">Anda belum memiliki alamat tersimpan.</p>
                  <a href="/addresses" className="inline-block bg-amber-600 text-white px-4 py-2 rounded-md font-medium hover:bg-amber-700 transition-colors">
                    Tambah Alamat Sekarang
                  </a>
                </div>
              )}
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Pilih Kurir</label>
              <select
                value={courier}
                onChange={(e) => setCourier(e.target.value)}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              >
                <option value="JNE">JNE (Rp 20.000)</option>
                <option value="SICEPAT">SiCepat (Rp 15.000)</option>
              </select>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Catatan (Opsional)</label>
              <input
                type="text"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Cth: Titip di pos satpam"
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              />
            </div>
          </form>
        </div>

        {/* Ringkasan Belanja */}
        <div className="bg-slate-50/50 p-6 border rounded-xl h-fit">
          <h2 className="text-xl font-semibold mb-6 border-b pb-2">Ringkasan Pesanan</h2>
          
          <div className="space-y-3 mb-6">
            {checkoutItems.map((item) => {
              const product = productsCache[item.product_id];
              const price = product ? product.base_price : 150000;
              return (
                <div key={item.id} className="flex justify-between text-sm">
                  <span>{product?.name || `Produk ID ${item.product_id}`} (x{item.quantity})</span>
                  <span>Rp {(price * Number(item.quantity)).toLocaleString("id-ID")}</span>
                </div>
              );
            })}
          </div>

          <div className="border-t pt-4 space-y-2">
            <div className="flex justify-between text-sm text-muted-foreground">
              <span>Subtotal Produk</span>
              <span>
                Rp {checkoutItems.reduce((acc, item) => {
                  const product = productsCache[item.product_id];
                  const price = product ? product.base_price : 150000;
                  return acc + (price * Number(item.quantity));
                }, 0).toLocaleString("id-ID")}
              </span>
            </div>
            <div className="flex justify-between text-sm text-muted-foreground">
              <span>Ongkos Kirim ({courier})</span>
              <span>Rp {courier === "JNE" ? "20.000" : "15.000"}</span>
            </div>
            <div className="flex justify-between text-lg font-bold mt-4 pt-2 border-t">
              <span>Total Tagihan</span>
              <span>
                Rp {(checkoutItems.reduce((acc, item) => {
                  const product = productsCache[item.product_id];
                  const price = product ? product.base_price : 150000;
                  return acc + (price * Number(item.quantity));
                }, 0) + (courier === "JNE" ? 20000 : 15000)).toLocaleString("id-ID")}
              </span>
            </div>
          </div>

          <button
            type="submit"
            form="checkout-form"
            disabled={isLoading || addresses.length === 0}
            className="w-full mt-8 flex items-center justify-center rounded-md bg-primary px-4 py-3 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors shadow-sm disabled:opacity-50"
          >
            {isLoading ? "Memproses Pesanan..." : "Buat Pesanan Sekarang"}
          </button>
        </div>
      </div>
    </div>
  );
}
