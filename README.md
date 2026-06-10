# Git-Smart-Commit 🚀

**Git-Smart-Commit** adalah sebuah CLI (Command Line Interface) tool lokal yang ditulis menggunakan bahasa pemrograman Go. Aplikasi ini dirancang untuk membantu developer menulis pesan commit Git secara interaktif dan memastikan pesan tersebut mengikuti standar **Conventional Commits** yang rapi dan konsisten.

Sangat cocok untuk pemula yang ingin meningkatkan portofolio GitHub mereka dengan proyek Go yang fungsional dan modular!

---

## ✨ Fitur Utama

- 🔍 **Validasi Repositori**: Memastikan direktori saat ini adalah repositori Git yang valid sebelum memproses commit.
- 🗂️ **Pengecekan Staging Area**: Memberikan peringatan jika belum ada perubahan file yang di-stage (`git add`), mencegah pembuatan commit kosong secara tidak sengaja.
- 💬 **Formulir Interaktif**: Menggunakan prompt interaktif untuk mengumpulkan informasi:
  - **Type**: Memilih tipe commit (`feat`, `fix`, `docs`, `refactor`, dll.).
  - **Scope**: Menentukan lingkup kode yang diubah (opsional).
  - **Subject**: Deskripsi singkat perubahan (wajib).
  - **Body**: Deskripsi detail perubahan (opsional).
  - **Footer**: Menambahkan informasi breaking changes atau menutup isu GitHub (opsional).
- 📝 **Pratinjau & Konfirmasi**: Menampilkan hasil format pesan commit sebelum dieksekusi ke Git.

---

## 🛠️ Struktur Folder

Proyek ini terstruktur secara modular dengan mengikuti standar layout proyek Go yang bersih:

```text
git-smart-commit/
├── cmd/
│   └── root.go           # Logika utama CLI & konfigurasi Cobra
├── internal/
│   ├── prompt/
│   │   └── prompt.go     # Logika pertanyaan interaktif menggunakan Survey
│   └── git/
│       └── git.go        # Logika eksekusi perintah Git lokal (os/exec)
├── main.go               # Entry point aplikasi (menjalankan command root)
├── go.mod                # Deklarasi dependensi modul Go
├── go.sum                # Checksum untuk dependensi Go
└── README.md             # Dokumentasi proyek
```

---

## ⚙️ Cara Menjalankan

### Prasyarat
- [Go](https://go.dev/doc/install) (versi 1.16 atau lebih baru)
- [Git](https://git-scm.com/downloads) terinstal di sistem operasi Anda.

### 1. Klon / Buat Folder Proyek
Pastikan Anda berada di direktori proyek:
```bash
cd git-smart-commit
```

### 2. Jalankan Perintah Git Add
Pastikan ada perubahan file yang siap untuk di-commit:
```bash
git add .
```

### 3. Jalankan Aplikasi
Anda dapat langsung menjalankan aplikasi tanpa kompilasi terlebih dahulu:
```bash
go run main.go
```

Atau, Anda bisa melakukan build menjadi file executable agar bisa dijalankan di mana saja:
```bash
# Lakukan kompilasi
go build -o git-smart.exe

# Jalankan executable
.\git-smart.exe
```

---

## 📚 Referensi Conventional Commits

Format commit yang dihasilkan akan berbentuk:
```text
<type>(<scope>): <subject>

<body>

<footer>
```

Tipe commit yang didukung oleh tool ini:
- `feat`: Menambahkan fitur baru ke dalam codebase.
- `fix`: Memperbaiki bug.
- `docs`: Perubahan pada dokumentasi.
- `style`: Perubahan format penulisan kode (semicolon, spasi, format, dll) tanpa merubah logika.
- `refactor`: Restrukturisasi kode tanpa menambah fitur baru atau memperbaiki bug.
- `perf`: Optimasi performa kode.
- `test`: Menambahkan atau mengubah unit test.
- `build`: Perubahan pada build system atau package manager (go.mod, npm, dll).
- `ci`: Perubahan konfigurasi CI/CD (GitHub Actions, dll).
- `chore`: Tugas pembantu ringan lainnya (maintenance).
- `revert`: Membatalkan commit sebelumnya.
