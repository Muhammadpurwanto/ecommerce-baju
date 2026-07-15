"use client";

import { useEffect, useState } from "react";
import api from "@/lib/axios";
import { User } from "lucide-react";

export default function AdminUsers() {
  const [users, setUsers] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const res = await api.get("/users/");
        if (res.data.success) {
          setUsers(res.data.data || []);
        }
      } catch (error) {
        console.error("Gagal mengambil pengguna:", error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchUsers();
  }, []);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Daftar Pengguna</h1>
        <p className="text-muted-foreground mt-1">Lihat seluruh pengguna terdaftar di sistem.</p>
      </div>

      <div className="border rounded-xl bg-card overflow-hidden shadow-sm">
        <div className="p-4 border-b bg-muted/20">
          <h2 className="text-lg font-semibold flex items-center gap-2">
            <User className="h-5 w-5" /> Pengguna Terdaftar
          </h2>
        </div>
        {isLoading ? (
          <div className="p-8 text-center text-muted-foreground">Memuat data pengguna...</div>
        ) : users.length === 0 ? (
          <div className="p-8 text-center text-muted-foreground">Belum ada pengguna.</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="text-xs text-muted-foreground uppercase bg-muted/10 border-b">
                <tr>
                  <th className="px-6 py-3 font-medium">ID / Role</th>
                  <th className="px-6 py-3 font-medium">Nama</th>
                  <th className="px-6 py-3 font-medium">Email</th>
                </tr>
              </thead>
              <tbody>
                {users.map((user, index) => (
                  <tr key={user.id} className={index !== users.length - 1 ? "border-b" : ""}>
                    <td className="px-6 py-4">
                      <span className="font-mono text-xs text-muted-foreground">{user.id.substring(0, 8)}...</span>
                      <div className="mt-1">
                        <span className={`px-2 py-0.5 rounded text-[10px] font-bold uppercase ${user.role === 'admin' ? 'bg-primary/20 text-primary' : 'bg-muted text-muted-foreground'}`}>
                          {user.role || 'customer'}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 font-medium">{user.name}</td>
                    <td className="px-6 py-4 text-muted-foreground">{user.email}</td>
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
