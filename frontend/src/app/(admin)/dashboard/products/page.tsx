"use client";

import { useEffect, useState } from "react";
import api from "@/lib/axios";
import { PlusCircle, Trash2, Edit, Save, X } from "lucide-react";

export default function AdminProducts() {
  const [products, setProducts] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Form State
  const [editingId, setEditingId] = useState<number | null>(null);
  const [newProduct, setNewProduct] = useState({
    name: "",
    category_id: "",
    description: "",
    brand: "",
    gender: "unisex",
    base_price: "",
    weight: "200",
    stock: "10",
    image_url: "",
    is_active: true
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const fetchProducts = async () => {
    try {
      const res = await api.get("/products/");
      if (res.data.success) {
        setProducts(res.data.data || []);
      }
    } catch (error) {
      console.error("Gagal mengambil produk:", error);
    }
  };

  const fetchCategories = async () => {
    try {
      const res = await api.get("/categories/");
      if (res.data.success) {
        setCategories(res.data.data || []);
      }
    } catch (error) {
      console.error("Gagal mengambil kategori:", error);
    }
  };

  useEffect(() => {
    Promise.all([fetchProducts(), fetchCategories()]).finally(() => {
      setIsLoading(false);
    });
  }, []);

  const resetForm = () => {
    setNewProduct({
      name: "",
      category_id: "",
      description: "",
      brand: "",
      gender: "unisex",
      base_price: "",
      weight: "200",
      stock: "10",
      image_url: "",
      is_active: true
    });
    setSelectedFile(null);
    setEditingId(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newProduct.name.trim() || !newProduct.category_id || !newProduct.base_price || !newProduct.stock) {
      alert("Mohon isi semua field yang wajib");
      return;
    }

    setIsSubmitting(true);
    try {
      let finalImageUrl = newProduct.image_url;
      if (selectedFile) {
        const formData = new FormData();
        formData.append("image", selectedFile);
        const uploadRes = await api.post("/products/upload-image", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        });
        if (uploadRes.data.success && uploadRes.data.image_url) {
          finalImageUrl = uploadRes.data.image_url;
        } else {
          throw new Error("Gagal mengupload gambar");
        }
      }

      const payload = {
        name: newProduct.name,
        category_id: parseInt(newProduct.category_id),
        description: newProduct.description,
        brand: newProduct.brand,
        gender: newProduct.gender,
        base_price: parseFloat(newProduct.base_price),
        weight: parseInt(newProduct.weight),
        stock: parseInt(newProduct.stock),
        image_url: finalImageUrl,
        is_active: newProduct.is_active
      };

      if (editingId) {
        // Edit Mode
        const res = await api.put(`/products/${editingId}`, payload);
        if (!res.data.success) throw new Error(res.data.message);
        alert("Produk berhasil diperbarui!");
      } else {
        // Create Mode
        const res = await api.post("/products/", payload);
        if (!res.data.success) throw new Error(res.data.message);
        alert("Produk berhasil ditambahkan!");
      }

      resetForm();
      fetchProducts(); // Refresh table
    } catch (error: any) {
      console.error(error);
      alert(error.response?.data?.message || error.message || "Terjadi kesalahan. Pastikan Anda adalah Admin.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEditClick = (product: any) => {
    setEditingId(product.id);
    setNewProduct({
      name: product.name,
      category_id: (product.category_id ?? product.categoryId ?? "").toString(),
      description: product.description || "",
      brand: product.brand || "",
      gender: product.gender || "unisex",
      base_price: (product.base_price ?? product.basePrice ?? "").toString(),
      weight: (product.weight ?? 200).toString(),
      stock: (product.stock ?? 10).toString(),
      image_url: product.image_url || "",
      is_active: product.is_active
    });
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleDeleteProduct = async (id: number) => {
    if (!window.confirm("Yakin ingin menghapus produk ini?")) return;
    
    try {
      const res = await api.delete(`/products/${id}`);
      if (res.data.success) {
        alert("Produk berhasil dihapus!");
        fetchProducts();
      }
    } catch (error: any) {
      console.error(error);
      alert(error.response?.data?.message || "Gagal menghapus produk.");
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-slate-900">Manajemen Produk</h1>
        <p className="text-muted-foreground mt-1">Tambah, ubah, dan kelola produk toko secara terpusat.</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Form Tambah/Edit Produk */}
        <div className="lg:col-span-1">
          <div className="p-6 border rounded-xl bg-card shadow-sm space-y-6 bg-white">
            <h2 className="text-lg font-bold border-b pb-2 flex items-center justify-between">
              {editingId ? "Edit Produk" : "Tambah Produk Baru"}
              {editingId && (
                <button onClick={resetForm} className="text-muted-foreground hover:text-foreground">
                  <X className="h-4 w-4" />
                </button>
              )}
            </h2>
            
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-semibold text-slate-700">Kategori</label>
                <select
                  required
                  value={newProduct.category_id}
                  onChange={(e) => setNewProduct({ ...newProduct, category_id: e.target.value })}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                  <option value="" disabled>Pilih Kategori</option>
                  {categories.map((cat) => (
                    <option key={cat.id} value={cat.id}>{cat.name}</option>
                  ))}
                </select>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-semibold text-slate-700">Nama Produk</label>
                <input
                  type="text"
                  required
                  value={newProduct.name}
                  onChange={(e) => setNewProduct({ ...newProduct, name: e.target.value })}
                  placeholder="Contoh: Kaos Polos Premium"
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-slate-700">Harga (Rp)</label>
                  <input
                    type="number"
                    required
                    value={newProduct.base_price}
                    onChange={(e) => setNewProduct({ ...newProduct, base_price: e.target.value })}
                    placeholder="100000"
                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-slate-700">Stok Utama</label>
                  <input
                    type="number"
                    required
                    min="0"
                    value={newProduct.stock}
                    onChange={(e) => setNewProduct({ ...newProduct, stock: e.target.value })}
                    placeholder="10"
                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-slate-700">Berat (gram)</label>
                  <input
                    type="number"
                    required
                    value={newProduct.weight}
                    onChange={(e) => setNewProduct({ ...newProduct, weight: e.target.value })}
                    placeholder="200"
                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-slate-700">Gender</label>
                  <select
                    value={newProduct.gender}
                    onChange={(e) => setNewProduct({ ...newProduct, gender: e.target.value })}
                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  >
                    <option value="unisex">Unisex</option>
                    <option value="men">Men</option>
                    <option value="women">Women</option>
                  </select>
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-semibold text-slate-700">Merk (Opsional)</label>
                <input
                  type="text"
                  value={newProduct.brand}
                  onChange={(e) => setNewProduct({ ...newProduct, brand: e.target.value })}
                  placeholder="Contoh: Brand-X"
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-semibold text-slate-700">Upload Gambar Produk</label>
                <input
                  type="file"
                  accept="image/*"
                  onChange={(e) => {
                    if (e.target.files && e.target.files.length > 0) {
                      setSelectedFile(e.target.files[0]);
                    }
                  }}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                />
              </div>

              {/* <div className="space-y-2">
                <label className="text-sm font-semibold text-slate-700">Atau URL Gambar (Opsional)</label>
                <input
                  type="text"
                  value={newProduct.image_url}
                  onChange={(e) => setNewProduct({ ...newProduct, image_url: e.target.value })}
                  placeholder="https://example.com/image.jpg"
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                />
              </div> */}

              <div className="space-y-2">
                <label className="text-sm font-semibold text-slate-700">Deskripsi</label>
                <textarea
                  value={newProduct.description}
                  onChange={(e) => setNewProduct({ ...newProduct, description: e.target.value })}
                  placeholder="Keterangan produk..."
                  rows={4}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                />
              </div>

              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="product_is_active"
                  checked={newProduct.is_active}
                  onChange={(e) => setNewProduct({ ...newProduct, is_active: e.target.checked })}
                  className="h-4 w-4 rounded border-slate-300 text-primary focus:ring-primary"
                />
                <label htmlFor="product_is_active" className="text-sm font-semibold text-slate-700">Produk Aktif</label>
              </div>

              <div className="pt-4 flex gap-2 border-t">
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="flex-1 flex items-center justify-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
                >
                  {isSubmitting ? "Menyimpan..." : (
                    <>
                      {editingId ? <Save className="h-4 w-4" /> : <PlusCircle className="h-4 w-4" />}
                      {editingId ? "Simpan Perubahan" : "Simpan Produk"}
                    </>
                  )}
                </button>
                {editingId && (
                  <button
                    type="button"
                    onClick={resetForm}
                    className="px-4 py-2 rounded-md border border-slate-200 hover:bg-slate-100 text-sm font-medium transition-colors"
                  >
                    Batal
                  </button>
                )}
              </div>
            </form>
          </div>
        </div>

        {/* Tabel Produk */}
        <div className="lg:col-span-2">
          <div className="border rounded-xl bg-card overflow-hidden shadow-sm bg-white">
            <div className="p-4 border-b bg-muted/20">
              <h2 className="text-lg font-semibold text-slate-800">Daftar Produk</h2>
            </div>
            {isLoading ? (
              <div className="p-8 text-center text-slate-500">Memuat data...</div>
            ) : products.length === 0 ? (
              <div className="p-8 text-center text-slate-400">Belum ada produk.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm text-left">
                  <thead className="text-xs text-slate-500 uppercase bg-slate-50 border-b">
                    <tr>
                      <th className="px-6 py-3 font-medium">ID</th>
                      <th className="px-6 py-3 font-medium">Nama Produk</th>
                      <th className="px-6 py-3 font-medium">Harga</th>
                      <th className="px-6 py-3 font-medium">Stok</th>
                      <th className="px-6 py-3 font-medium text-right">Aksi</th>
                    </tr>
                  </thead>
                  <tbody>
                    {products.map((prod, index) => (
                      <tr key={prod.id} className={index !== products.length - 1 ? "border-b" : ""}>
                        <td className="px-6 py-4 text-slate-600 font-medium">{prod.id}</td>
                        <td className="px-6 py-4 font-bold text-slate-800">{prod.name}</td>
                        <td className="px-6 py-4 text-slate-700">Rp {prod.base_price?.toLocaleString('id-ID')}</td>
                        <td className="px-6 py-4 text-slate-700">{prod.stock}</td>
                        <td className="px-6 py-4 text-right">
                          <button 
                            onClick={() => handleEditClick(prod)}
                            className="p-2 text-primary hover:bg-primary/10 rounded-md transition-colors mr-2"
                            title="Edit Produk"
                          >
                            <Edit className="h-4 w-4" />
                          </button>
                          <button 
                            onClick={() => handleDeleteProduct(prod.id)}
                            className="p-2 text-destructive hover:bg-destructive/10 rounded-md transition-colors"
                            title="Hapus Produk"
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
