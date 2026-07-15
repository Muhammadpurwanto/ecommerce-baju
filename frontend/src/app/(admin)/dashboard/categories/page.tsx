"use client";

import { useEffect, useState } from "react";
import api from "@/lib/axios";
import { PlusCircle, Trash2 } from "lucide-react";

export default function AdminCategories() {
  const [categories, setCategories] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Form State untuk Tambah Kategori
  const [newCategory, setNewCategory] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchCategories = async () => {
    try {
      const res = await api.get("/categories/");
      if (res.data.success) {
        setCategories(res.data.data || []);
      }
    } catch (error) {
      console.error("Gagal mengambil kategori:", error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchCategories();
  }, []);

  const handleAddCategory = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCategory.trim()) return;

    setIsSubmitting(true);
    try {
      const res = await api.post("/categories/", { name: newCategory });
      if (res.data.success) {
        alert("Kategori berhasil ditambahkan!");
        setNewCategory("");
        fetchCategories(); // Refresh tabel
      }
    } catch (error: any) {
      console.error(error);
      alert(error.response?.data?.message || "Gagal menambah kategori. Pastikan Anda adalah Admin.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeleteCategory = async (id: number) => {
    if (!window.confirm("Yakin ingin menghapus kategori ini?")) return;
    
    try {
      const res = await api.delete(`/categories/${id}`);
      if (res.data.success) {
        alert("Kategori berhasil dihapus!");
        fetchCategories();
      }
    } catch (error: any) {
      console.error(error);
      alert(error.response?.data?.message || "Gagal menghapus kategori.");
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Manajemen Kategori</h1>
        <p className="text-muted-foreground mt-1">Tambah, lihat, dan kelola kategori produk.</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Form Tambah Kategori */}
        <div className="lg:col-span-1">
          <div className="p-6 border rounded-xl bg-card shadow-sm sticky top-24">
            <h2 className="text-lg font-semibold mb-4 border-b pb-2">Tambah Kategori Baru</h2>
            <form onSubmit={handleAddCategory} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Nama Kategori</label>
                <input
                  type="text"
                  required
                  value={newCategory}
                  onChange={(e) => setNewCategory(e.target.value)}
                  placeholder="Contoh: Kemeja Pria"
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                />
              </div>
              <button
                type="submit"
                disabled={isSubmitting}
                className="w-full flex items-center justify-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
              >
                <PlusCircle className="h-4 w-4" />
                {isSubmitting ? "Menyimpan..." : "Simpan Kategori"}
              </button>
            </form>
          </div>
        </div>

        {/* Tabel Kategori */}
        <div className="lg:col-span-2">
          <div className="border rounded-xl bg-card overflow-hidden shadow-sm">
            <div className="p-4 border-b bg-muted/20">
              <h2 className="text-lg font-semibold">Daftar Kategori</h2>
            </div>
            {isLoading ? (
              <div className="p-8 text-center text-muted-foreground">Memuat data...</div>
            ) : categories.length === 0 ? (
              <div className="p-8 text-center text-muted-foreground">Belum ada kategori.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm text-left">
                  <thead className="text-xs text-muted-foreground uppercase bg-muted/10 border-b">
                    <tr>
                      <th className="px-6 py-3 font-medium">ID</th>
                      <th className="px-6 py-3 font-medium">Nama Kategori</th>
                      <th className="px-6 py-3 font-medium">Slug</th>
                      <th className="px-6 py-3 font-medium text-right">Aksi</th>
                    </tr>
                  </thead>
                  <tbody>
                    {categories.map((cat, index) => (
                      <tr key={cat.id} className={index !== categories.length - 1 ? "border-b" : ""}>
                        <td className="px-6 py-4">{cat.id}</td>
                        <td className="px-6 py-4 font-medium">{cat.name}</td>
                        <td className="px-6 py-4 text-muted-foreground">{cat.slug}</td>
                        <td className="px-6 py-4 text-right">
                          <button 
                            onClick={() => handleDeleteCategory(cat.id)}
                            className="p-2 text-destructive hover:bg-destructive/10 rounded-md transition-colors"
                            title="Hapus Kategori"
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
