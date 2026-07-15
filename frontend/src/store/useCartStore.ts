// File: frontend/src/store/useCartStore.ts
import { create } from 'zustand';
import api from '@/lib/axios';

interface CartItem {
  id: number;
  product_id: number;
  quantity: number;
}

interface CartState {
  items: CartItem[];
  productsCache: Record<number, any>;
  isLoading: boolean;
  fetchCart: () => Promise<void>;
  addToCart: (productId: number, quantity: number) => Promise<void>;
  updateItem: (itemId: number, quantity: number) => Promise<void>;
  removeItem: (itemId: number) => Promise<void>;
  selectedItems: number[];
  toggleSelectItem: (itemId: number) => void;
  selectAll: () => void;
  clearSelection: () => void;
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],
  productsCache: {},
  isLoading: false,
  selectedItems: [],

  toggleSelectItem: (itemId) => set((state) => ({
    selectedItems: state.selectedItems.includes(itemId) // apakah id barang sudah ada di selectedItems? jika iya maka hapus, jika tidak maka tambahkan
      ? state.selectedItems.filter((id) => id !== itemId) // hapus id barang dari selectedItems
      : [...state.selectedItems, itemId], // tambahkan id barang ke selectedItems
  })),
  selectAll: () => set((state) => ({ selectedItems: state.items.map((i) => i.id) })), // pilih semua id barang
  clearSelection: () => set({ selectedItems: [] }), // hapus semua id barang

  fetchCart: async () => {
    try {
      set({ isLoading: true });
      const response = await api.get('/carts/'); // ambil data keranjang
      if (response.data.success) {
        set({ items: response.data.data.items || [] }); // simpan data keranjang
        
        // Ambil semua produk untuk cache di frontend (karena microservice memisahkannya)
        try {
          const prodRes = await api.get('/products/'); // ambil data produk
          if (prodRes.data.success) {
             const pMap: Record<number, any> = {};
             (prodRes.data.data || []).forEach((p: any) => pMap[p.id] = p); // simpan data produk 101: { "id": 101, "name": "Kemeja Flanel", "price": 150000 }
             set({ productsCache: pMap }); // simpan data produk
          }
        } catch (e) {
          console.error("Gagal mengambil cache produk", e);
        }
      }
    } catch (error) {
      console.error("Gagal mengambil keranjang:", error);
      // Jika error 401, interceptor axios otomatis akan mengurusnya
    } finally {
      set({ isLoading: false });
    }
  },

  addToCart: async (productId, quantity) => {
    try {
      await api.post('/carts/items', {
        product_id: productId,
        quantity: quantity
      });
      // Setelah berhasil ditambah di backend, tarik ulang data terbarunya
      get().fetchCart();
    } catch (error) {
      console.error("Gagal menambah ke keranjang:", error);
      alert("Gagal menambah barang, pastikan Anda sudah login.");
    }
  },

  updateItem: async (itemId, quantity) => {
    try {
      await api.put(`/carts/items/${itemId}`, { quantity });
      get().fetchCart();
    } catch (error) {
      console.error("Gagal mengubah kuantitas:", error);
    }
  },

  removeItem: async (itemId) => {
    try {
      await api.delete(`/carts/items/${itemId}`);
      get().fetchCart();
    } catch (error) {
      console.error("Gagal menghapus item:", error);
    }
  }
}));
