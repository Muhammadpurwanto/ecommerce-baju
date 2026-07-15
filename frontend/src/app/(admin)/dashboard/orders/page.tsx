"use client";

import { useEffect, useState } from "react";
import api from "@/lib/axios";
import { ShoppingCart } from "lucide-react";

export default function AdminOrders() {
  const [orders, setOrders] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchOrders = async () => {
    try {
      const res = await api.get("/orders/all");
      if (res.data.success) {
        setOrders(res.data.data || []);
      }
    } catch (error) {
      console.error("Gagal mengambil pesanan:", error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, []);

  const handleUpdateStatus = async (orderId: number, currentStatus: string, newStatus: string) => {
    if (currentStatus === newStatus) return;
    
    try {
      const res = await api.put(`/orders/${orderId}/status`, { status: newStatus });
      if (res.data.success) {
        alert("Status pesanan berhasil diperbarui!");
        fetchOrders();
      }
    } catch (error: any) {
      console.error(error);
      alert(error.response?.data?.message || "Gagal memperbarui status pesanan.");
      fetchOrders(); // reset select element if failed
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'paid': return 'bg-blue-100 text-blue-800 border-blue-200';
      case 'processing': return 'bg-indigo-100 text-indigo-800 border-indigo-200';
      case 'shipped': return 'bg-purple-100 text-purple-800 border-purple-200';
      case 'completed': return 'bg-emerald-100 text-emerald-800 border-emerald-200';
      case 'cancelled': return 'bg-red-100 text-red-800 border-red-200';
      default: return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Manajemen Pesanan</h1>
        <p className="text-muted-foreground mt-1">Pantau dan kelola status pesanan dari pelanggan.</p>
      </div>

      <div className="border rounded-xl bg-card overflow-hidden shadow-sm">
        <div className="p-4 border-b bg-muted/20">
          <h2 className="text-lg font-semibold flex items-center gap-2">
            <ShoppingCart className="h-5 w-5" /> Daftar Pesanan
          </h2>
        </div>
        {isLoading ? (
          <div className="p-8 text-center text-muted-foreground">Memuat data pesanan...</div>
        ) : orders.length === 0 ? (
          <div className="p-8 text-center text-muted-foreground">Belum ada pesanan.</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="text-xs text-muted-foreground uppercase bg-muted/10 border-b">
                <tr>
                  <th className="px-6 py-3 font-medium">Nomor Invoice</th>
                  <th className="px-6 py-3 font-medium">Pelanggan ID</th>
                  <th className="px-6 py-3 font-medium">Total Harga</th>
                  <th className="px-6 py-3 font-medium">Tanggal</th>
                  <th className="px-6 py-3 font-medium text-right">Status</th>
                </tr>
              </thead>
              <tbody>
                {orders.map((order, index) => (
                  <tr key={order.id} className={index !== orders.length - 1 ? "border-b" : ""}>
                    <td className="px-6 py-4 font-mono text-xs font-medium">{order.order_number}</td>
                    <td className="px-6 py-4 text-xs font-mono text-muted-foreground">{order.user_id.substring(0, 8)}...</td>
                    <td className="px-6 py-4 font-medium">Rp {order.total_amount?.toLocaleString('id-ID')}</td>
                    <td className="px-6 py-4 text-muted-foreground">
                      {new Date(order.created_at).toLocaleDateString('id-ID', {
                        day: 'numeric', month: 'short', year: 'numeric'
                      })}
                    </td>
                    <td className="px-6 py-4 text-right">
                      <select
                        value={order.status}
                        onChange={(e) => handleUpdateStatus(order.id, order.status, e.target.value)}
                        className={`text-xs font-bold px-2 py-1.5 rounded-md border ${getStatusColor(order.status)} focus:outline-none focus:ring-2 focus:ring-primary/50 cursor-pointer`}
                      >
                        <option value="pending">Pending</option>
                        <option value="paid">Paid</option>
                        <option value="processing">Processing</option>
                        <option value="shipped">Shipped</option>
                        <option value="completed">Completed</option>
                        <option value="cancelled">Cancelled</option>
                      </select>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
