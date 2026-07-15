export const dynamic = "force-dynamic";
import Link from "next/link";
import { ChevronRight } from "lucide-react";
import { ProductClientView } from "@/components/product/ProductClientView";

async function getProductBySlug(slug: string) {
  const apiUrl = process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
  try {
    const res = await fetch(`${apiUrl}/products/${slug}`, {
      cache: "no-store",
    });
    
    if (!res.ok) return null;
    const json = await res.json();
    return json.data;
  } catch (error) {
    console.error("Gagal mengambil data produk:", error);
    return null;
  }
}

export default async function ProductDetailPage({ params }: { params: Promise<{ slug: string }> }) {
  const resolvedParams = await params;
  const product = await getProductBySlug(resolvedParams.slug);

  if (!product) {
    return (
      <div className="container mx-auto px-4 py-20 text-center min-h-screen flex flex-col items-center justify-center">
        <h1 className="text-3xl font-bold mb-4">Produk Tidak Ditemukan</h1>
        <p className="text-muted-foreground mb-8">Maaf, produk yang Anda cari mungkin sudah dihapus atau URL tidak valid.</p>
        <Link href="/products" className="bg-primary text-primary-foreground px-6 py-3 rounded-full font-medium hover:bg-primary/90 transition-colors">
          Kembali ke Katalog
        </Link>
      </div>
    );
  }

  return (
    <div className="bg-slate-50 min-h-screen pb-20">
      {/* Breadcrumb */}
      <div className="bg-white border-b">
        <div className="container mx-auto px-4 py-4 flex items-center text-sm text-muted-foreground">
          <Link href="/" className="hover:text-primary transition-colors">Beranda</Link>
          <ChevronRight className="h-4 w-4 mx-2 opacity-50" />
          <Link href="/products" className="hover:text-primary transition-colors">Produk</Link>
          <ChevronRight className="h-4 w-4 mx-2 opacity-50" />
          <span className="text-foreground font-medium truncate">{product.name}</span>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8 md:py-12">
        <div className="bg-white rounded-2xl shadow-sm border p-4 md:p-8">
          <ProductClientView product={product} />
        </div>
      </div>
    </div>
  );
}
