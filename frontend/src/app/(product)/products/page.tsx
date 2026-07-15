export const dynamic = "force-dynamic";
import { ProductCard } from "@/components/product/ProductCard";

// Fungsi untuk mengambil data produk langsung dari server (SSR)
async function getProducts() {
  const apiUrl = process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
  try {
    const res = await fetch(`${apiUrl}/products/`, {
      cache: "no-store",
    });
    
    if (!res.ok) return [];
    const json = await res.json();
    return json.data || [];
  } catch (error) {
    console.error("Gagal mengambil data produk:", error);
    return [];
  }
}

export default async function ProductsPage({
  searchParams,
}: {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}) {
  const products = await getProducts();
  const params = await searchParams;
  const search = typeof params?.search === 'string' ? params.search.toLowerCase() : "";
  const categoryId = typeof params?.category === 'string' ? params.category : "";

  let filteredProducts = products;
  
  if (search) {
    filteredProducts = filteredProducts.filter((p: any) => 
      p.name.toLowerCase().includes(search) || 
      (p.description && p.description.toLowerCase().includes(search))
    );
  }

  if (categoryId) {
    filteredProducts = filteredProducts.filter((p: any) => p.category_id?.toString() === categoryId);
  }

  // Get category name if filtering by category
  let categoryName = "";
  if (categoryId && filteredProducts.length > 0) {
    categoryName = filteredProducts[0].category?.name || "Kategori Terpilih";
  }

  return (
    <div className="container mx-auto px-4 py-12 md:py-20 min-h-screen">
      <div className="flex flex-col items-center mb-12 text-center">
        <h1 className="text-4xl font-bold tracking-tight mb-4">
          {search ? `Hasil Pencarian: "${params?.search}"` : categoryId ? `Kategori: ${categoryName}` : "Semua Produk"}
        </h1>
        <div className="h-1 w-24 bg-primary rounded-full mb-4"></div>
        <p className="text-muted-foreground max-w-2xl">
          {search || categoryId ? `Menemukan ${filteredProducts.length} produk yang sesuai.` : "Jelajahi seluruh koleksi pakaian kami. Dari gaya santai hingga formal, temukan pilihan terbaik untuk mengekspresikan diri Anda."}
        </p>
      </div>

      {filteredProducts.length === 0 ? (
        <div className="text-center py-20 border rounded-2xl bg-muted/20">
          <p className="text-muted-foreground text-lg">Maaf, belum ada produk yang sesuai dengan pencarian Anda.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 md:gap-8">
          {filteredProducts.map((product: any) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}
    </div>
  );
}
