# Analisis Skalabilitas & Konfigurasi RNR Blockchain

Berikut adalah analisis mendalam berdasarkan pemeriksaan kode sumber (Source Code Audit) terhadap pertanyaan Anda.

## 1. Feabilitas Blok 2GB dengan Block Time 1 Menit

**Skenario:** Blok berukuran 2GB diproduksi setiap 60 detik (1 menit).
**Perangkat:** Laptop/PC seharga Rp 10-15 Juta.

### A. Penyimpanan (Storage) - ✅ AMAN (Karena Pruning)
Jika node harus menyimpan semua sejarah, 2GB/menit = ~2.8 TB/hari. SSD laptop 1TB akan penuh dalam 8 jam.
**Namun**, kode RNR saat ini memiliki fitur **Aggressive Pruning** hardcoded (lihat `internal/blockchain/blockchain.go` dan `internal/storage/manager.go`).
*   Node hanya menyimpan **25 blok terakhir** (sekitar 25 menit data).
*   Storage yang dibutuhkan hanya sekitar **50 GB** konstan.
*   **Kesimpulan**: SSD bawaan laptop (512GB/1TB NVMe) **SANGAT CUKUP**.

### B. Komputasi (CPU/RAM) - ✅ MUNGKIN (Terbantukan Sharding)
Algoritma konsensus adalah **PoSSR (Sorting Race)**.
*   Total 2GB dibagi ke 10 Shard = 200MB per Shard.
*   **PENTING (Klarifikasi Implementasi):** Dalam kode saat ini (`consensus/engine.go` dan `p2p/gossipsub.go`), pembagian shard masih bersifat **Paralel Lokal**, belum **Terdistribusi Penuh**.
    *   Setiap node memproses ke-10 shard tersebut secara bersamaan menggunakan *multi-threading* (10 core CPU).
    *   Setiap node tetap mendownload *seluruh* 1GB data (Full Block) karena topik gossip P2P hanya satu (`rnr/blocks`).
    *   Konsep "Hanya menghitung shard yang dipilih" adalah desain ideal (Whitepaper), namun di versi kode v1.0 ini, node Anda melakukan validasi untuk SEMUA shard.
*   Meski begitu, RAM 16GB dan CPU modern masih mampu menanganinya karena beban dibagi ke core prosesor yang berbeda.

### C. Jaringan (Bandwidth) - ❌ KRITIS (Bottleneck Utama)
Ini adalah rintangan terbesar.
*   2 GB / 60 detik = **~34 MB/detik (Megabytes per second)**.
*   Dalam satuan kecepatan internet (Mbps): 34 * 8 = **~272 Mbps (Megabits per second)**.
*   Ini adalah kecepatan **Upload & Download stabil** yang dibutuhkan terus menerus.
*   Kebanyakan internet rumahan di Indonesia memiliki kecepatan Download 50-100 Mbps, tapi **Upload seringkali hanya 10-20 Mbps**.
*   **Kesimpulan**: Laptopnya kuat, tapi **Koneksi Internet Rumahan TIDAK AKAN KUAT**. Anda butuh koneksi bisnis/dedicated dengan upload simetris >300 Mbps agar tidak tertinggal sinkronisasi.

---

## 2. Bagaimana RNR Menjadi Blockchain Masa Depan?

Berdasarkan arsitektur kode saat ini, RNR mencapai skalabilitas "supermasif" melalui tiga pilar:

1.  **Verifikasi O(N) vs Kreasi O(N log N)**:
    *   Membuat blok (Sorting) itu berat, tapi memverifikasinya (Linear Scan) sangat ringan. Ini memungkinkan blok besar diproses dengan cepat oleh jaringan.
2.  **Native Sharding (Pemisahan Beban)**:
    *   Konfigurasi `NumShards = 10` (bisa ditingkatkan ke 256). Setiap node tidak perlu memproses *seluruh* 2GB, melainkan hanya bagian shard-nya saja saat proses validasi intensif.
3.  **Stateless/Lightweight Design**:
    *   Dengan hanya menyimpan 25 blok terakhir (`PruningWindow`), RNR mencegah masalah "State Bloat" yang dialami Ethereum/Bitcoin (dimana node butuh TB-an storage). Node RNR tetap ringan selamanya, membuatnya mudah dipasang di jutaan perangkat IoT atau laptop tanpa dedicated server farm.

---

## 3. Apa Insentif untuk Archives Node?

Pertanyaan: *"Apa yang diberikan jaringan jika ada node yang ingin menjadi Archives node?"*

**Temuan Fakta dari Kode:**
1.  **Tidak Didukung Secara Default**:
    *   Logika Pruning di `internal/blockchain/blockchain.go` (baris 114) bersifat **Hardcoded**:
        ```go
        if block.Header.Height > 25 {
            bc.store.PruneOldBlocks(...)
        }
        ```
    *   Artinya, jika Anda menjalankan software standar, node Anda **AKAN** menghapus data lama secara otomatis. Anda tidak bisa menjadi Archives Node tanpa memodifikasi kode programnya (menghapus baris tersebut).

2.  **Tidak Ada Insentif On-Chain (Saat Ini)**:
    *   Di module `internal/economics/supply.go`, reward hanya diberikan untuk **Mining/Sorting** (Block Reward).
    *   Tidak ditemukan kode yang memberikan reward khusus untuk "Storage" atau "Serving Historical Data".

**Jawaban Strategis:**
Saat ini, peran Archives Node di ekosistem RNR bersifat **sukarela atau komersial off-chain** (seperti bisnis Explorer, Data Analytics, atau API Provider ala Infura/Alchemy). Protokol inti belum memberikan reward otomatis untuk peran ini.

**Rekomendasi Pengembangan:**
Jika ingin RNR menopang aplikasi web3 masif, protokol perlu menambahkan "Storage Rent" atau "Data Retrieval Fee" di masa depan agar ada alasan ekonomis bagi seseorang untuk menyimpan data berukuran Petabytes (jika 2GB/blok terus berjalan).
