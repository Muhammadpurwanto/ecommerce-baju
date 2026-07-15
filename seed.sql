USE product_db;

-- Tambah 10 Kategori
INSERT IGNORE INTO categories (id, name, slug, created_at, updated_at) VALUES 
(1, 'Pakaian Pria', 'pakaian-pria', NOW(), NOW()),
(2, 'Pakaian Wanita', 'pakaian-wanita', NOW(), NOW()),
(3, 'Pakaian Anak', 'pakaian-anak', NOW(), NOW()),
(4, 'Sepatu', 'sepatu', NOW(), NOW()),
(5, 'Tas', 'tas', NOW(), NOW()),
(6, 'Aksesoris', 'aksesoris', NOW(), NOW()),
(7, 'Olahraga', 'olahraga', NOW(), NOW()),
(8, 'Pakaian Dalam', 'pakaian-dalam', NOW(), NOW()),
(9, 'Pakaian Tidur', 'pakaian-tidur', NOW(), NOW()),
(10, 'Busana Muslim', 'busana-muslim', NOW(), NOW());

-- Tambah 50 Produk (5 per kategori)
INSERT IGNORE INTO products (id, category_id, name, slug, description, brand, gender, base_price, weight, is_active, created_at, updated_at) VALUES 
-- Pakaian Pria (1)
(1, 1, 'Kemeja Flanel Pria', 'kemeja-flanel-pria', 'Kemeja flanel kotak-kotak bahan katun.', 'Erigo', 'men', 150000.00, 250, 1, NOW(), NOW()),
(2, 1, 'Kaos Polos Hitam', 'kaos-polos-hitam', 'Kaos katun combed 30s sangat nyaman.', 'Zalora', 'men', 50000.00, 150, 1, NOW(), NOW()),
(3, 1, 'Celana Chino Panjang', 'celana-chino-panjang', 'Celana panjang bahan twill stretch.', 'Levis', 'men', 200000.00, 400, 1, NOW(), NOW()),
(4, 1, 'Jaket Hoodie Pria', 'jaket-hoodie-pria', 'Jaket hoodie fleece tebal dan hangat.', 'H&M', 'men', 250000.00, 500, 1, NOW(), NOW()),
(5, 1, 'Celana Pendek Pantai', 'celana-pendek-pantai', 'Celana pendek bahan parasut.', 'Quiksilver', 'men', 100000.00, 200, 1, NOW(), NOW()),

-- Pakaian Wanita (2)
(6, 2, 'Dress Bunga Musim Panas', 'dress-bunga', 'Dress motif bunga bahan rayon.', 'Zalora', 'women', 180000.00, 300, 1, NOW(), NOW()),
(7, 2, 'Blouse Kerja Wanita', 'blouse-kerja', 'Blouse formal untuk ke kantor.', 'Executive', 'women', 150000.00, 250, 1, NOW(), NOW()),
(8, 2, 'Rok Plisket Panjang', 'rok-plisket', 'Rok panjang bahan hyget super.', 'Uniqlo', 'women', 120000.00, 300, 1, NOW(), NOW()),
(9, 2, 'Cardigan Rajut', 'cardigan-rajut', 'Cardigan rajut halus import.', 'H&M', 'women', 160000.00, 350, 1, NOW(), NOW()),
(10, 2, 'Celana Kulot Wanita', 'celana-kulot', 'Celana kulot bahan scuba.', 'Zalora', 'women', 140000.00, 300, 1, NOW(), NOW()),

-- Pakaian Anak (3)
(11, 3, 'Setelan Anak Laki-Laki', 'setelan-anak-cowo', 'Setelan kaos dan celana pendek.', 'Pigeon', 'unisex', 90000.00, 200, 1, NOW(), NOW()),
(12, 3, 'Gaun Anak Perempuan', 'gaun-anak-cewe', 'Gaun pesta anak lucu.', 'Disney', 'women', 150000.00, 250, 1, NOW(), NOW()),
(13, 3, 'Piyama Anak Motif Hewan', 'piyama-anak-hewan', 'Baju tidur anak katun jepang.', 'Velvet Junior', 'unisex', 80000.00, 150, 1, NOW(), NOW()),
(14, 3, 'Jaket Denim Anak', 'jaket-denim-anak', 'Jaket jeans untuk anak.', 'OshKosh', 'unisex', 170000.00, 300, 1, NOW(), NOW()),
(15, 3, 'Kaos Kaki Anak Karakter', 'kaos-kaki-anak', 'Isi 3 pasang kaos kaki.', 'Socksly', 'unisex', 30000.00, 100, 1, NOW(), NOW()),

-- Sepatu (4)
(16, 4, 'Sepatu Sneakers Pria', 'sneakers-pria', 'Sepatu kasual untuk jalan-jalan.', 'Nike', 'men', 800000.00, 800, 1, NOW(), NOW()),
(17, 4, 'Sepatu Lari Wanita', 'sepatu-lari-wanita', 'Sepatu running ringan.', 'Adidas', 'women', 750000.00, 700, 1, NOW(), NOW()),
(18, 4, 'Sandal Selop Kulit', 'sandal-selop-kulit', 'Sandal kulit asli pria.', 'Bata', 'men', 200000.00, 500, 1, NOW(), NOW()),
(19, 4, 'Sepatu Formal Pria', 'sepatu-formal-pria', 'Sepatu pantofel kerja.', 'Buccheri', 'men', 500000.00, 900, 1, NOW(), NOW()),
(20, 4, 'Flat Shoes Wanita', 'flat-shoes-wanita', 'Sepatu flat kasual wanita.', 'Charles & Keith', 'women', 350000.00, 400, 1, NOW(), NOW()),

-- Tas (5)
(21, 5, 'Tas Ransel Laptop', 'tas-ransel-laptop', 'Ransel muat laptop 15 inch.', 'Eiger', 'unisex', 350000.00, 800, 1, NOW(), NOW()),
(22, 5, 'Sling Bag Wanita', 'sling-bag-wanita', 'Tas selempang kecil elegan.', 'Guess', 'women', 450000.00, 400, 1, NOW(), NOW()),
(23, 5, 'Tote Bag Kanvas', 'tote-bag-kanvas', 'Tote bag aesthetic.', 'Nature', 'unisex', 50000.00, 200, 1, NOW(), NOW()),
(24, 5, 'Tas Gunung 60L', 'tas-gunung', 'Tas carrier untuk hiking.', 'Consina', 'unisex', 600000.00, 1500, 1, NOW(), NOW()),
(25, 5, 'Dompet Kulit Pria', 'dompet-kulit', 'Dompet lipat kulit sapi.', 'Bonia', 'men', 250000.00, 150, 1, NOW(), NOW()),

-- Aksesoris (6)
(26, 6, 'Jam Tangan Analog Pria', 'jam-tangan-pria', 'Jam tangan kulit klasik.', 'Fossil', 'men', 1200000.00, 200, 1, NOW(), NOW()),
(27, 6, 'Kacamata Hitam Aviator', 'kacamata-aviator', 'Kacamata anti UV.', 'Ray-Ban', 'unisex', 800000.00, 100, 1, NOW(), NOW()),
(28, 6, 'Topi Baseball Polos', 'topi-baseball', 'Topi kasual bahan drill.', 'Polo', 'unisex', 75000.00, 100, 1, NOW(), NOW()),
(29, 6, 'Ikat Pinggang Kulit', 'ikat-pinggang', 'Sabuk kulit pria.', 'Levis', 'men', 150000.00, 150, 1, NOW(), NOW()),
(30, 6, 'Kalung Titanium Wanita', 'kalung-titanium', 'Kalung anti karat.', 'Swarovski', 'women', 300000.00, 50, 1, NOW(), NOW()),

-- Olahraga (7)
(31, 7, 'Baju Lari Cepat Kering', 'baju-lari-dryfit', 'Baju olahraga bahan dryfit.', 'Nike', 'unisex', 200000.00, 150, 1, NOW(), NOW()),
(32, 7, 'Celana Training Pria', 'celana-training', 'Celana jogger olahraga.', 'Adidas', 'men', 250000.00, 300, 1, NOW(), NOW()),
(33, 7, 'Matras Yoga NBR', 'matras-yoga', 'Matras yoga anti slip 10mm.', 'Kettler', 'unisex', 150000.00, 1000, 1, NOW(), NOW()),
(34, 7, 'Sport Bra Wanita', 'sport-bra', 'Bra olahraga high impact.', 'Puma', 'women', 180000.00, 150, 1, NOW(), NOW()),
(35, 7, 'Jaket Parasut Olahraga', 'jaket-parasut', 'Jaket anti angin untuk lari.', 'Under Armour', 'unisex', 400000.00, 300, 1, NOW(), NOW()),

-- Pakaian Dalam (8)
(36, 8, 'Boxer Katun Pria', 'boxer-katun', 'Celana dalam boxer isi 3.', 'Calvin Klein', 'men', 150000.00, 200, 1, NOW(), NOW()),
(37, 8, 'Brief Wanita Seamless', 'brief-seamless', 'Celana dalam wanita tanpa jahitan.', 'Sorex', 'women', 40000.00, 50, 1, NOW(), NOW()),
(38, 8, 'Singlet Katun Pria', 'singlet-katun', 'Kaos dalam putih pria.', 'Rider', 'men', 35000.00, 100, 1, NOW(), NOW()),
(39, 8, 'Bra Busa Kawat', 'bra-kawat', 'Bra push up dengan kawat.', 'Wacoal', 'women', 120000.00, 100, 1, NOW(), NOW()),
(40, 8, 'Korset Pelangsing', 'korset-pelangsing', 'Korset pembentuk tubuh.', 'Marena', 'women', 200000.00, 150, 1, NOW(), NOW()),

-- Pakaian Tidur (9)
(41, 9, 'Piyama Sutra Wanita', 'piyama-sutra', 'Baju tidur bahan silk mewah.', 'Victoria Secret', 'women', 450000.00, 250, 1, NOW(), NOW()),
(42, 9, 'Setelan Kaos Tidur Pria', 'kaos-tidur-pria', 'Kaos oblong dan celana pendek.', 'H&M', 'men', 150000.00, 300, 1, NOW(), NOW()),
(43, 9, 'Daster Batik Jumbo', 'daster-batik', 'Daster rumahan adem.', 'Kencana Ungu', 'women', 60000.00, 200, 1, NOW(), NOW()),
(44, 9, 'Kimono Mandi Handuk', 'kimono-mandi', 'Jubah mandi bahan handuk tebal.', 'Terry Palmer', 'unisex', 200000.00, 600, 1, NOW(), NOW()),
(45, 9, 'Piyama Katun Lengan Pendek', 'piyama-katun', 'Baju tidur katun motif.', 'Zalora', 'unisex', 120000.00, 250, 1, NOW(), NOW()),

-- Busana Muslim (10)
(46, 10, 'Gamis Syari Set Hijab', 'gamis-syari', 'Setelan gamis dan khimar.', 'Zoya', 'women', 350000.00, 600, 1, NOW(), NOW()),
(47, 10, 'Baju Koko Kurta Pria', 'koko-kurta', 'Koko model kurta bahan toyobo.', 'Rabbani', 'men', 180000.00, 300, 1, NOW(), NOW()),
(48, 10, 'Hijab Pashmina Ceruty', 'pashmina-ceruty', 'Kerudung pashmina babydoll.', 'Dian Pelangi', 'women', 50000.00, 150, 1, NOW(), NOW()),
(49, 10, 'Sarung Tenun Wadimor', 'sarung-tenun', 'Sarung tenun motif songket.', 'Wadimor', 'men', 75000.00, 350, 1, NOW(), NOW()),
(50, 10, 'Mukena Parasut Premium', 'mukena-parasut', 'Mukena travel ringan dan kecil.', 'Tatuis', 'women', 120000.00, 250, 1, NOW(), NOW());

-- Inisialisasi stok 100 dan gambar default untuk seluruh produk demo
UPDATE products SET stock = 100, image_url = 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?auto=format&fit=crop&w=1000&q=80' WHERE stock = 0;

