"use client";

import { useEffect, useState, Suspense } from "react";
import { User, MapPin, ShoppingBag, Package, ChevronRight } from "lucide-react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import api from "@/lib/axios";

function OrdersContent() {
  const [orders, setOrders] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const searchParams = useSearchParams();

  // Jika dialihkan dari midtrans
  const transactionStatus = searchParams.get("transaction_status");

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const res = await api.get("/orders/");
        if (res.data.success) {
          setOrders(res.data.data);
        }
      } catch (error) {
        console.error("Gagal mengambil daftar pesanan:", error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchOrders();
  }, []);

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "pending":
        return <span className="px-3 py-1 bg-amber-100 text-amber-700 text-xs font-bold rounded-full">Menunggu Pembayaran</span>;
      case "paid":
        return <span className="px-3 py-1 bg-emerald-100 text-emerald-700 text-xs font-bold rounded-full">Sudah Dibayar</span>;
      case "processing":
        return <span className="px-3 py-1 bg-blue-100 text-blue-700 text-xs font-bold rounded-full">Diproses</span>;
      case "shipped":
        return <span className="px-3 py-1 bg-purple-100 text-purple-700 text-xs font-bold rounded-full">Dikirim</span>;
      case "completed":
        return <span className="px-3 py-1 bg-slate-100 text-slate-700 text-xs font-bold rounded-full">Selesai</span>;
      case "cancelled":
        return <span className="px-3 py-1 bg-rose-100 text-rose-700 text-xs font-bold rounded-full">Dibatalkan</span>;
      default:
        return <span className="px-3 py-1 bg-slate-100 text-slate-700 text-xs font-bold rounded-full">{status}</span>;
    }
  };

  const formatPrice = (price: number) => new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", minimumFractionDigits: 0 }).format(price);

  return (
    <div className="container mx-auto px-4 py-12 max-w-5xl">
      <h1 className="text-3xl font-bold tracking-tight mb-8">Pengaturan Akun</h1>

      {transactionStatus === "settlement" && (
        <div className="mb-8 p-4 bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-xl flex items-center gap-3">
          <div className="h-8 w-8 bg-emerald-100 rounded-full flex items-center justify-center flex-shrink-0">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7"></path></svg>
          </div>
          <div>
            <h3 className="font-bold">Pembayaran Berhasil!</h3>
            <p className="text-sm">Terima kasih, pesanan Anda sedang kami proses.</p>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
        {/* Sidebar Navigasi Profil */}
        <div className="md:col-span-1 space-y-2">
          <Link href="/profile" className="flex items-center gap-2 p-3 hover:bg-muted text-muted-foreground hover:text-foreground rounded-lg font-medium transition-colors">
            <User className="h-5 w-5" /> Profil Saya
          </Link>
          <Link href="/orders" className="flex items-center gap-2 p-3 bg-primary text-primary-foreground rounded-lg font-medium">
            <ShoppingBag className="h-5 w-5" /> Pesanan Saya
          </Link>
          <Link href="/addresses" className="flex items-center gap-2 p-3 hover:bg-muted text-muted-foreground hover:text-foreground rounded-lg font-medium transition-colors">
            <MapPin className="h-5 w-5" /> Buku Alamat
          </Link>
        </div>

        {/* Konten Utama: Daftar Pesanan */}
        <div className="md:col-span-3 space-y-6">
          <div className="p-6 border rounded-xl bg-card shadow-sm">
            <h2 className="text-xl font-bold mb-6 border-b pb-4">Daftar Pesanan</h2>

            {isLoading ? (
              <div className="text-center py-12 text-muted-foreground">Memuat data pesanan...</div>
            ) : orders.length === 0 ? (
              <div className="text-center py-12">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-slate-100 mb-4">
                  <Package className="h-8 w-8 text-slate-400" />
                </div>
                <h3 className="text-lg font-bold text-slate-900 mb-2">Belum ada pesanan</h3>
                <p className="text-slate-500 mb-6">Anda belum pernah melakukan pemesanan.</p>
                <Link href="/products" className="inline-block bg-primary text-primary-foreground px-6 py-2 rounded-lg font-medium hover:bg-primary/90 transition-colors">
                  Mulai Belanja
                </Link>
              </div>
            ) : (
              <div className="space-y-4">
                {orders.map((order) => (
                  <div key={order.id} className="border rounded-xl p-4 sm:p-5 hover:border-slate-300 transition-colors">
                    <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4 pb-4 border-b">
                      <div>
                        <div className="flex items-center gap-3 mb-1">
                          <span className="font-bold text-slate-900">{order.order_number}</span>
                          {getStatusBadge(order.status)}
                        </div>
                        <div className="text-sm text-slate-500">
                          {new Date(order.created_at).toLocaleDateString("id-ID", { day: 'numeric', month: 'long', year: 'numeric' })}
                        </div>
                      </div>
                      <div className="text-left sm:text-right">
                        <div className="text-sm text-slate-500 mb-1">Total Belanja</div>
                        <div className="font-bold text-primary">{formatPrice(order.total_amount)}</div>
                      </div>
                    </div>

                    <div className="flex justify-between items-center">
                      <div className="text-sm text-slate-600">
                        {order.items?.length || 0} Produk • Kurir: <span className="font-medium uppercase">{order.courier}</span>
                      </div>
                      <button className="text-primary font-semibold text-sm flex items-center hover:underline">
                        Lihat Detail <ChevronRight className="h-4 w-4 ml-1" />
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default function OrdersPage() {
  return (
    <Suspense fallback={<div className="text-center py-12 text-muted-foreground">Memuat halaman pesanan...</div>}>
      <OrdersContent />
    </Suspense>
  );
}
