// File: frontend/src/app/(user)/profile/page.tsx
"use client";

import { useEffect, useState } from "react";
import { User, MapPin } from "lucide-react";
import Link from "next/link";
import api from "@/lib/axios";

export default function ProfilePage() {
  const [profile, setProfile] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isUpdating, setIsUpdating] = useState(false);
  
  // State Form
  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [isUploading, setIsUploading] = useState(false);

  const fetchProfile = async () => {
    setIsLoading(true);
    try {
      const res = await api.get("/users/profile");
      if (res.data.success) {
        const data = res.data.data;
        setProfile(data);
        setName(data.name || "");
        setPhone(data.phone || "");
      }
    } catch (error) {
      console.error("Gagal mengambil data profil", error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchProfile();
  }, []);

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsUpdating(true);
    try {
      // Menghindari error validasi backend (min=8) jika nomor telepon dikosongkan
      const payload: any = { name };
      if (phone && phone.trim() !== "") {
        payload.phone = phone;
      }

      const res = await api.put("/users/profile", payload);
      
      if (res.data.success) {
        alert("Profil berhasil diperbarui!");
        fetchProfile(); // Refresh tampilan dan foto
      }
    } catch (error: any) {
      alert(error.response?.data?.message || "Gagal memperbarui profil");
    } finally {
      setIsUpdating(false);
    }
  };

  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Optional: Validate file type and size
    if (!file.type.startsWith("image/")) {
      alert("Hanya file gambar yang diperbolehkan");
      return;
    }
    if (file.size > 2 * 1024 * 1024) {
      alert("Ukuran gambar maksimal 2MB");
      return;
    }

    setIsUploading(true);
    const formData = new FormData();
    formData.append("avatar", file);

    try {
      const res = await api.post("/users/profile/avatar", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      if (res.data.success) {
        alert("Avatar berhasil diunggah!");
        fetchProfile(); // Memuat ulang gambar baru
      }
    } catch (error: any) {
      alert(error.response?.data?.message || "Gagal mengunggah avatar");
    } finally {
      setIsUploading(false);
    }
  };

  if (isLoading) {
    return <div className="p-12 text-center text-muted-foreground">Memuat data profil...</div>;
  }

  return (
    <div className="container mx-auto px-4 py-12 max-w-4xl">
      <h1 className="text-3xl font-bold tracking-tight mb-8">Pengaturan Akun</h1>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
        {/* Sidebar Navigasi Profil */}
        <div className="md:col-span-1 space-y-2">
          <Link href="/profile" className="flex items-center gap-2 p-3 bg-primary text-primary-foreground rounded-lg font-medium">
            <User className="h-5 w-5" /> Profil Saya
          </Link>
          <Link href="/orders" className="flex items-center gap-2 p-3 hover:bg-muted text-muted-foreground hover:text-foreground rounded-lg font-medium transition-colors">
            <User className="h-5 w-5" /> Pesanan Saya
          </Link>
          <Link href="/addresses" className="flex items-center gap-2 p-3 hover:bg-muted text-muted-foreground hover:text-foreground rounded-lg font-medium transition-colors">
            <MapPin className="h-5 w-5" /> Buku Alamat
          </Link>
        </div>

        {/* Konten Utama: Update Profil */}
        <div className="md:col-span-3 space-y-6">
          <div className="p-6 border rounded-xl bg-card shadow-sm">
            
            {/* Header Profil (Avatar + Role) */}
            <div className="flex items-center gap-6 mb-8 border-b pb-6">
              <div className="h-20 w-20 rounded-full overflow-hidden bg-muted flex items-center justify-center border-2 border-primary/20">
                {profile?.avatar_url ? (
                  <img src={profile.avatar_url} alt="Profile" className="h-full w-full object-cover" />
                ) : (
                  <User className="h-10 w-10 text-muted-foreground" />
                )}
              </div>
              <div>
                <h2 className="text-2xl font-bold">{profile?.name}</h2>
                <div className="inline-flex items-center rounded-full bg-secondary mt-1 mb-3 px-3 py-0.5 text-xs font-semibold uppercase text-secondary-foreground">
                  {profile?.role || "CUSTOMER"}
                </div>
                <div className="mt-2">
                  <label className="cursor-pointer inline-flex items-center gap-2 rounded-md bg-secondary/50 px-3 py-1.5 text-xs font-medium text-secondary-foreground hover:bg-secondary transition-colors">
                    <span>{isUploading ? "Mengunggah..." : "Ubah Foto"}</span>
                    <input type="file" className="hidden" accept="image/*" onChange={handleAvatarChange} disabled={isUploading} />
                  </label>
                </div>
              </div>
            </div>
            
            <form onSubmit={handleUpdateProfile} className="space-y-5 max-w-lg">
              <div className="space-y-2">
                <label className="text-sm font-medium">Email <span className="text-muted-foreground font-normal">(Tidak bisa diubah)</span></label>
                <input 
                  type="email" 
                  disabled 
                  value={profile?.email || ""} 
                  className="w-full rounded-md border px-3 py-2 text-sm bg-muted cursor-not-allowed" 
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Nama Lengkap</label>
                <input 
                  type="text" 
                  required 
                  value={name} 
                  onChange={(e) => setName(e.target.value)}
                  className="w-full rounded-md border px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" 
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Nomor Telepon</label>
                <input 
                  type="tel" 
                  value={phone} 
                  onChange={(e) => setPhone(e.target.value)}
                  placeholder="Contoh: 081234567890"
                  className="w-full rounded-md border px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" 
                />
              </div>

              {/* Field Avatar URL telah dihapus dan diganti dengan tombol Upload di atas */}

              <div className="pt-4">
                <button 
                  type="submit" 
                  disabled={isUpdating}
                  className="rounded-md bg-primary px-6 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
                >
                  {isUpdating ? "Menyimpan..." : "Simpan Perubahan"}
                </button>
              </div>
            </form>

          </div>
        </div>
      </div>
    </div>
  );
}
