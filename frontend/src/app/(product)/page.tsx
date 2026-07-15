// File: frontend/src/app/page.tsx
export const dynamic = "force-dynamic";
import { ProductCard } from "@/components/product/ProductCard";

// Fungsi untuk mengambil data produk langsung dari server (SSR)
async function getProducts() {
  // SSR: pakai API_URL (internal Docker network) karena kode ini jalan di server, bukan di browser
  // NEXT_PUBLIC_API_URL hanya fallback untuk development lokal (tanpa Docker)
  const apiUrl = process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
  try {
    const res = await fetch(`${apiUrl}/products/`, {
      cache: "no-store", // Mengambil data terbaru (bisa diganti "force-cache" untuk kecepatan maksimal)
    });
    
    if (!res.ok) return [];
    const json = await res.json();
    return json.data || [];
  } catch (error) {
    console.error("Gagal mengambil data produk:", error);
    return [];
  }
}

export default async function Home() {
  const products = await getProducts();

  return (
    <div className="flex flex-col min-h-screen">
      
      {/* Hero Section */}
      <section className="relative w-full bg-slate-900 text-white overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-r from-slate-900 to-slate-800/50 z-10"></div>
        {/* Gambar background hero */}
        <img 
          src="https://images.unsplash.com/photo-1441984904996-e0b6ba687e04?auto=format&fit=crop&w=2000&q=80" 
          alt="Hero" 
          className="absolute inset-0 w-full h-full object-cover opacity-50"
        />
        <div className="container relative z-20 mx-auto px-4 py-24 md:py-32 lg:py-40 flex flex-col items-center text-center">
          <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-6 max-w-3xl">
            Pakaian Berkualitas,<br/> <span className="text-primary">Gaya Tanpa Batas.</span>
          </h1>
          <p className="text-lg md:text-xl text-slate-300 max-w-2xl mb-10 leading-relaxed">
            Temukan koleksi terbaru kami untuk melengkapi penampilan Anda. 
            Dibuat dengan bahan premium untuk kenyamanan maksimal sepanjang hari.
          </p>
          <a href="#koleksi-terbaru" className="inline-flex h-12 items-center justify-center rounded-full bg-primary px-8 text-sm font-medium text-primary-foreground transition-transform hover:scale-105 shadow-lg">
            Belanja Sekarang
          </a>
        </div>
      </section>

      {/* Daftar Produk */}
      <section id="koleksi-terbaru" className="container mx-auto px-4 py-16 md:py-24">
        <div className="flex flex-col items-center mb-12 text-center">
          <h2 className="text-3xl font-bold tracking-tight mb-4">Koleksi Terbaru</h2>
          <div className="h-1 w-20 bg-primary rounded-full mb-4"></div>
          <p className="text-muted-foreground max-w-2xl">
            Produk-produk pilihan terbaik bulan ini. Stok terbatas, dapatkan sekarang sebelum kehabisan!
          </p>
        </div>

        {products.length === 0 ? (
          <div className="text-center py-20 border rounded-2xl bg-muted/20">
            <p className="text-muted-foreground text-lg">Belum ada produk yang tersedia saat ini.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 md:gap-8">
            {products.map((product: any) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        )}
      </section>
      
    </div>
  );
}
