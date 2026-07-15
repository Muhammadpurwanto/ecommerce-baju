export const dynamic = "force-dynamic";
import Link from "next/link";
import { Layers } from "lucide-react";

async function getCategories() {
  const apiUrl = process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
  try {
    const res = await fetch(`${apiUrl}/categories/`, {
      cache: "no-store",
    });
    
    if (!res.ok) return [];
    const json = await res.json();
    return json.data || [];
  } catch (error) {
    console.error("Gagal mengambil data kategori:", error);
    return [];
  }
}

export default async function CategoriesPage() {
  const categories = await getCategories();

  return (
    <div className="container mx-auto px-4 py-12 md:py-20 min-h-screen">
      <div className="flex flex-col items-center mb-12 text-center">
        <h1 className="text-4xl font-bold tracking-tight mb-4">Kategori Pakaian</h1>
        <div className="h-1 w-24 bg-primary rounded-full mb-4"></div>
        <p className="text-muted-foreground max-w-2xl">
          Temukan koleksi berdasarkan gaya dan kebutuhan Anda. Kami telah menyusunnya untuk mempermudah pencarian Anda.
        </p>
      </div>

      {categories.length === 0 ? (
        <div className="text-center py-20 border rounded-2xl bg-muted/20">
          <p className="text-muted-foreground text-lg">Maaf, belum ada kategori yang tersedia saat ini.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {categories.map((cat: any) => (
            <Link 
              key={cat.id} 
              href={`/products?category=${cat.id}`}
              className="group relative flex flex-col p-8 border rounded-2xl bg-card hover:border-primary/50 transition-all hover:shadow-md overflow-hidden"
            >
              <div className="absolute top-0 right-0 p-4 opacity-5 group-hover:opacity-10 transition-opacity">
                <Layers className="h-24 w-24" />
              </div>
              <div className="h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center text-primary mb-6 group-hover:scale-110 transition-transform">
                <Layers className="h-6 w-6" />
              </div>
              <h3 className="text-xl font-bold mb-2 group-hover:text-primary transition-colors">{cat.name}</h3>
              <p className="text-sm text-muted-foreground line-clamp-2">
                {cat.description || "Jelajahi berbagai pilihan menarik di kategori ini."}
              </p>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
