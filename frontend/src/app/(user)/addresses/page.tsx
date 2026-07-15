// File: frontend/src/app/(user)/addresses/page.tsx
"use client";

import { useEffect, useState } from "react";
import { User, MapPin, Plus, Trash2, Edit, Star } from "lucide-react";
import Link from "next/link";
import api from "@/lib/axios";

export default function AddressesPage() {
  const [addresses, setAddresses] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Form State
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);
  
  const [formData, setFormData] = useState({
    label: "",
    recipient: "",
    phone: "",
    province: "",
    city: "",
    district: "",
    postal_code: "",
    detail: "",
    is_default: false,
  });

  const fetchAddresses = async () => {
    setIsLoading(true);
    try {
      const res = await api.get("/users/addresses");
      if (res.data.success) {
        setAddresses(res.data.data || []);
      }
    } catch (error) {
      console.error("Gagal mengambil data alamat", error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchAddresses();
  }, []);

  const openAddForm = () => {
    setFormData({ 
      label: "", 
      recipient: "", 
      phone: "", 
      province: "", 
      city: "", 
      district: "", 
      postal_code: "", 
      detail: "", 
      is_default: false 
    });
    setIsEditing(false);
    setIsFormOpen(true);
  };

  const openEditForm = (addr: any) => {
    setFormData({
      label: addr.label || "",
      recipient: addr.recipient || "",
      phone: addr.phone || "",
      province: addr.province || "",
      city: addr.city || "",
      district: addr.district || "",
      postal_code: addr.postal_code || "",
      detail: addr.detail || "",
      is_default: addr.is_default || false,
    });
    setEditId(addr.id);
    setIsEditing(true);
    setIsFormOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (isEditing && editId) {
        const res = await api.put(`/users/addresses/${editId}`, formData);
        if (res.data.success) alert("Alamat berhasil diperbarui!");
      } else {
        const res = await api.post("/users/addresses", formData);
        if (res.data.success) alert("Alamat berhasil ditambahkan!");
      }
      setIsFormOpen(false);
      fetchAddresses();
    } catch (error: any) {
      alert(error.response?.data?.message || "Gagal menyimpan alamat");
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm("Apakah Anda yakin ingin menghapus alamat ini?")) return;
    try {
      const res = await api.delete(`/users/addresses/${id}`);
      if (res.data.success) fetchAddresses();
    } catch (error: any) {
      alert(error.response?.data?.message || "Gagal menghapus alamat");
    }
  };

  if (isLoading) return <div className="p-12 text-center text-muted-foreground">Memuat data alamat...</div>;

  return (
    <div className="container mx-auto px-4 py-12 max-w-4xl">
      <h1 className="text-3xl font-bold tracking-tight mb-8">Pengaturan Akun</h1>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
        {/* Sidebar Navigasi Profil */}
        <div className="md:col-span-1 space-y-2">
          <Link href="/profile" className="flex items-center gap-2 p-3 hover:bg-muted text-muted-foreground hover:text-foreground rounded-lg font-medium transition-colors">
            <User className="h-5 w-5" /> Profil Saya
          </Link>
          <Link href="/orders" className="flex items-center gap-2 p-3 hover:bg-muted text-muted-foreground hover:text-foreground rounded-lg font-medium transition-colors">
            <User className="h-5 w-5" /> Pesanan Saya
          </Link>
          <Link href="/addresses" className="flex items-center gap-2 p-3 bg-primary text-primary-foreground rounded-lg font-medium">
            <MapPin className="h-5 w-5" /> Buku Alamat
          </Link>
        </div>

        {/* Konten Utama: Manajemen Alamat */}
        <div className="md:col-span-3 space-y-6">
          
          <div className="flex items-center justify-between border-b pb-4">
            <h2 className="text-xl font-bold">Daftar Alamat Pengiriman</h2>
            {!isFormOpen && (
              <button onClick={openAddForm} className="flex items-center gap-1 text-sm bg-primary text-primary-foreground px-4 py-2 rounded-md hover:bg-primary/90 transition-colors">
                <Plus className="h-4 w-4" /> Tambah Alamat
              </button>
            )}
          </div>

          {/* Form Tambah/Edit */}
          {isFormOpen && (
            <div className="p-6 border rounded-xl bg-slate-50/50">
              <h3 className="text-lg font-semibold mb-6">{isEditing ? "Edit Alamat" : "Alamat Baru"}</h3>
              <form onSubmit={handleSubmit} className="space-y-5">
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Label Alamat (Contoh: Rumah, Kantor)</label>
                    <input required type="text" className="w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-primary outline-none" value={formData.label} onChange={e => setFormData({...formData, label: e.target.value})} />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Nama Penerima</label>
                    <input required type="text" className="w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-primary outline-none" value={formData.recipient} onChange={e => setFormData({...formData, recipient: e.target.value})} />
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Nomor Telepon Penerima</label>
                  <input required type="tel" className="w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-primary outline-none" value={formData.phone} onChange={e => setFormData({...formData, phone: e.target.value})} />
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Provinsi</label>
                    <input required type="text" className="w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-primary outline-none" value={formData.province} onChange={e => setFormData({...formData, province: e.target.value})} />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Kota/Kabupaten</label>
                    <input required type="text" className="w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-primary outline-none" value={formData.city} onChange={e => setFormData({...formData, city: e.target.value})} />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Kecamatan</label>
                    <input required type="text" className="w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-primary outline-none" value={formData.district} onChange={e => setFormData({...formData, district: e.target.value})} />
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Kode Pos</label>
                    <input required type="text" className="w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-primary outline-none" value={formData.postal_code} onChange={e => setFormData({...formData, postal_code: e.target.value})} />
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Detail Alamat Lengkap (Jalan, RT/RW, Blok, Patokan)</label>
                  <textarea required rows={3} className="w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-primary outline-none" value={formData.detail} onChange={e => setFormData({...formData, detail: e.target.value})} />
                </div>

                <div className="flex items-center space-x-2 pt-2">
                  <input 
                    type="checkbox" 
                    id="isDefault" 
                    className="rounded border-gray-300 text-primary focus:ring-primary h-4 w-4"
                    checked={formData.is_default}
                    onChange={(e) => setFormData({...formData, is_default: e.target.checked})}
                  />
                  <label htmlFor="isDefault" className="text-sm font-medium cursor-pointer">
                    Jadikan sebagai alamat utama
                  </label>
                </div>

                <div className="flex justify-end gap-2 pt-4 border-t mt-4">
                  <button type="button" onClick={() => setIsFormOpen(false)} className="px-5 py-2 text-sm font-medium border rounded-md hover:bg-muted transition-colors">Batal</button>
                  <button type="submit" className="px-5 py-2 text-sm font-medium bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors">
                    {isEditing ? "Update Alamat" : "Simpan Alamat"}
                  </button>
                </div>
              </form>
            </div>
          )}

          {/* List Alamat */}
          {!isFormOpen && addresses.length === 0 && (
            <div className="p-12 text-center border border-dashed rounded-xl text-muted-foreground bg-slate-50/50">
              <MapPin className="h-8 w-8 mx-auto mb-3 opacity-50" />
              <p>Anda belum menyimpan alamat apa pun.</p>
            </div>
          )}

          {!isFormOpen && addresses.length > 0 && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              {addresses.map((addr) => (
                <div key={addr.id} className={`relative flex flex-col p-5 border rounded-xl bg-card shadow-sm transition-all ${addr.is_default ? 'border-primary ring-1 ring-primary/20' : 'hover:border-primary/50'}`}>
                  
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex items-center gap-2">
                      <h3 className="font-semibold text-lg">{addr.label}</h3>
                      {addr.is_default && (
                        <span className="inline-flex items-center rounded-full bg-primary/10 px-2.5 py-0.5 text-[10px] font-bold text-primary">
                          UTAMA
                        </span>
                      )}
                    </div>
                  </div>

                  <div className="text-sm text-muted-foreground leading-relaxed mb-4 flex-1">
                    <p className="font-medium text-foreground mb-1">{addr.recipient} <span className="text-muted-foreground font-normal">({addr.phone})</span></p>
                    <p>{addr.detail}</p>
                    <p>{addr.district}, {addr.city}</p>
                    <p>{addr.province}, {addr.postal_code}</p>
                  </div>

                  <div className="flex items-center gap-4 border-t pt-3 mt-auto">
                    <button onClick={() => openEditForm(addr)} className="text-sm flex items-center gap-1.5 font-medium text-primary hover:opacity-80 transition-opacity">
                      <Edit className="h-4 w-4" /> Edit
                    </button>
                    <div className="w-px h-4 bg-border"></div>
                    <button onClick={() => handleDelete(addr.id)} className="text-sm flex items-center gap-1.5 font-medium text-destructive hover:opacity-80 transition-opacity">
                      <Trash2 className="h-4 w-4" /> Hapus
                    </button>
                  </div>

                </div>
              ))}
            </div>
          )}

        </div>
      </div>
    </div>
  );
}
