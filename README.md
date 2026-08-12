# Git-Smart-Commit 🚀

**Git-Smart-Commit** adalah CLI (Command Line Interface) tool lokal yang ditulis menggunakan **Go** untuk membantu developer membuat pesan commit yang mengikuti standar **Conventional Commits**.

Git-Smart-Commit dapat menganalisis perubahan yang telah di-stage menggunakan **Google Gemini AI**, memahami tujuan utama dari perubahan tersebut, lalu menghasilkan commit message secara otomatis.

Tool ini juga tetap menyediakan mode manual sebagai alternatif ketika developer ingin menentukan commit message sendiri.

---

## ✨ Fitur Utama

### 🤖 AI Commit Message

Menggunakan **Google Gemini** untuk menganalisis seluruh perubahan yang telah di-stage melalui:

```bash
git diff --cached
```

Gemini kemudian menentukan tujuan utama perubahan dan menghasilkan Conventional Commit message secara otomatis.

Contoh:

```text
feat(ai): integrate Gemini for commit message generation
```

AI tidak hanya melihat satu file, tetapi menganalisis keseluruhan staged diff sehingga dapat menangani perubahan yang melibatkan banyak file.

---

### 🧠 Analisis Multi-File

Git-Smart-Commit dapat menganalisis perubahan yang melibatkan banyak file sekaligus.

Contoh:

```text
9 files changed
445 insertions(+)
21 deletions(-)
```

Dari perubahan tersebut, AI dapat menyimpulkan tujuan utamanya menjadi:

```text
feat(ai): integrate Gemini for commit message generation
```

Hal ini membantu ketika developer mengalami kesulitan menentukan pesan commit setelah melakukan perubahan besar pada project.

---

### 📝 Conventional Commits

Commit message mengikuti standar Conventional Commits:

```text
<type>(<scope>): <subject>
```

Contoh:

```text
feat(auth): add WebAuthn authentication
```

atau:

```text
fix(api): handle invalid authentication token
```

---

### ✍️ Manual Commit Mode

Selain AI, developer tetap dapat membuat commit message secara manual menggunakan formulir interaktif.

Mode manual menyediakan:

- **Type** — tipe perubahan.
- **Scope** — bagian atau modul yang diubah.
- **Subject** — ringkasan singkat perubahan.
- **Body** — penjelasan detail perubahan secara opsional.
- **Commit confirmation** — konfirmasi sebelum commit.
- **Push confirmation** — pilihan untuk langsung melakukan push.

---

### 🔍 Validasi Git Repository

Sebelum menjalankan proses commit, aplikasi memastikan direktori saat ini merupakan repository Git yang valid.

Jika bukan repository Git, aplikasi akan memberikan peringatan.

---

### 📦 Validasi Staged Changes

Git-Smart-Commit hanya memproses perubahan yang sudah di-stage.

Contoh:

```bash
git add .
```

Jika tidak terdapat staged changes, aplikasi akan menghentikan proses dan meminta developer melakukan `git add` terlebih dahulu.

---

### 🌿 Automatic Branch Detection

Aplikasi dapat mendeteksi branch Git yang sedang digunakan.

Contoh:

```text
feature/ai-commit-message
```

Informasi branch tersebut digunakan ketika melakukan push ke remote repository.

---

### 🚀 Automatic Push

Setelah commit berhasil dibuat, developer dapat memilih untuk langsung melakukan push ke remote repository.

Contoh:

```text
🚀 Sedang melakukan push ke remote repository...

✅ Berhasil push ke origin/feature/ai-commit-message!
```

---

### 🔐 API Key melalui Environment Variable

Google Gemini API key tidak disimpan langsung di dalam source code.

Git-Smart-Commit menggunakan environment variable:

```text
GEMINI_API_KEY
```

Contoh pada PowerShell:

```powershell
$env:GEMINI_API_KEY="YOUR_API_KEY"
```

API key **jangan pernah ditulis langsung di source code atau di-commit ke repository**.

---

## 🔄 Cara Kerja

Ketika developer ingin membuat commit:

```text
git add .
      │
      ▼
┌──────────────────────┐
│ Staged Git Changes   │
└──────────┬───────────┘
           │
           ▼
    git diff --cached
           │
           ▼
┌──────────────────────┐
│    Git-Smart-Commit  │
└──────────┬───────────┘
           │
           ▼
┌──────────────────────┐
│     Google Gemini    │
│   AI Diff Analysis   │
└──────────┬───────────┘
           │
           ▼
 feat(ai): add Gemini
 commit generation
           │
           ▼
      git commit
           │
           ▼
      git push
```

AI berfokus pada **tujuan utama perubahan**, bukan sekadar menyebutkan file yang berubah.

---

## 🛠️ Struktur Folder

Project menggunakan struktur modular Go:

```text
git-smart-commit/
├── cmd/
│   └── root.go
│       # Entry point command dan workflow utama CLI
│
├── internal/
│   ├── ai/
│   │   └── gemini.go
│   │       # Integrasi Google Gemini API
│   │
│   ├── git/
│   │   └── git.go
│   │       # Operasi Git dan staged diff
│   │
│   └── prompt/
│       └── prompt.go
│           # Form dan prompt interaktif
│
├── testing/
│   ├── test-gemini.go
│   └── test-git-diff.go
│       # Testing dan eksperimen integrasi
│
├── main.go
│   # Entry point aplikasi
│
├── go.mod
│   # Deklarasi module dan dependency
│
├── go.sum
│   # Dependency checksums
│
├── .gitignore
│
└── README.md
```

---

## ⚙️ Cara Menjalankan

### Prasyarat

Pastikan sistem sudah memiliki:

- **Go**
- **Git**
- **Google Gemini API Key**

---

### 1. Clone Repository

Clone repository Git-Smart-Commit kemudian masuk ke direktorinya.

```bash
git clone <repository-url>

cd git-smart-commit
```

---

### 2. Install Dependency

Jalankan:

```bash
go mod tidy
```

---

### 3. Konfigurasi Gemini API Key

Set environment variable `GEMINI_API_KEY`.

#### Windows PowerShell

```powershell
$env:GEMINI_API_KEY="YOUR_API_KEY"
```

Kemudian pastikan environment variable tersedia:

```powershell
echo $env:GEMINI_API_KEY
```

> Jangan membagikan API key atau memasukkannya ke dalam repository.

---

### 4. Buat Perubahan pada Project

Lakukan perubahan seperti biasa pada project.

Kemudian stage perubahan:

```bash
git add .
```

Git-Smart-Commit akan menganalisis **hanya perubahan yang sudah di-stage**.

---

### 5. Jalankan Git-Smart-Commit

```bash
go run .
```

Aplikasi kemudian akan:

1. Memeriksa repository Git.
2. Memeriksa staged changes.
3. Mengambil `git diff --cached`.
4. Mengirim staged diff ke Gemini.
5. Menganalisis tujuan utama perubahan.
6. Menghasilkan Conventional Commit message.
7. Membuat commit.
8. Menawarkan proses push ke remote repository.

Contoh:

```text
🤖 Menganalisis perubahan dengan Gemini...

🤖 Suggested commit:
feat(ai): integrate Gemini for commit message generation

✅ Berhasil membuat commit!

🚀 Sedang melakukan push ke remote repository...
✅ Berhasil push ke origin/feature/ai-commit-message!
```

---

## 📦 Build Executable

Jika ingin membuat executable:

```bash
go build -o git-smart.exe
```

Kemudian jalankan:

```powershell
.\git-smart.exe
```

---

## 🧠 Contoh AI Commit Generation

Misalnya developer melakukan perubahan pada:

```text
cmd/root.go
internal/ai/gemini.go
internal/git/git.go
internal/prompt/prompt.go
go.mod
go.sum
README.md
test.js
testing/test-gemini.go
```

Dengan total:

```text
9 files changed
445 insertions(+)
21 deletions(-)
```

Git-Smart-Commit mengambil staged diff dan mengirimkannya kepada Gemini.

Daripada developer harus menentukan sendiri:

```text
feat?
fix?
chore?
scope?
subject?
```

AI menganalisis perubahan tersebut dan menghasilkan:

```text
feat(ai): integrate Gemini for commit message generation
```

Tujuannya adalah membuat commit message yang menggambarkan **perubahan utama secara keseluruhan**, bukan sekadar daftar file yang dimodifikasi.

---

## 📚 Conventional Commits

Format utama:

```text
<type>(<scope>): <subject>
```

Contoh:

```text
feat(auth): add WebAuthn authentication
```

### Tipe Commit yang Didukung

| Type       | Penggunaan                                         |
| ---------- | -------------------------------------------------- |
| `feat`     | Menambahkan fitur baru                             |
| `fix`      | Memperbaiki bug                                    |
| `docs`     | Perubahan dokumentasi                              |
| `style`    | Perubahan formatting/style tanpa perubahan logic   |
| `refactor` | Restrukturisasi kode tanpa mengubah behavior utama |
| `perf`     | Optimasi performa                                  |
| `test`     | Menambahkan atau mengubah testing                  |
| `build`    | Perubahan build system atau dependency             |
| `ci`       | Perubahan konfigurasi CI/CD                        |
| `chore`    | Maintenance atau pekerjaan pendukung               |
| `revert`   | Membatalkan perubahan sebelumnya                   |

---

## 🧪 Testing

Testing integrasi dapat dilakukan melalui folder:

```text
testing/
```

Contoh pengujian:

```bash
go run testing/test-gemini.go
```

Testing tersebut digunakan untuk memastikan komunikasi antara aplikasi Go dan Google Gemini API berjalan dengan baik.

---

## 🎯 Tujuan Project

Git-Smart-Commit dibuat untuk menyelesaikan masalah sederhana yang sering dialami developer:

> **"Saya sudah melakukan banyak perubahan pada project, tetapi bingung harus menulis commit message apa."**

Dengan memanfaatkan AI, developer dapat menyerahkan proses analisis perubahan kepada Gemini dan mendapatkan commit message yang lebih konsisten berdasarkan perubahan aktual pada Git.

Project ini juga menjadi implementasi pembelajaran mengenai:

- Go CLI development
- Git automation
- Conventional Commits
- Google Gemini API
- AI-assisted developer tools
- Modular project architecture
- Environment variable management
- Command execution menggunakan Go

---

## 🚧 Status Project

Project saat ini masih dalam tahap pengembangan.

Fitur utama yang sudah tersedia:

- ✅ Git repository validation
- ✅ Staged changes validation
- ✅ Staged diff extraction
- ✅ Gemini AI integration
- ✅ AI-generated Conventional Commit
- ✅ Multi-file change analysis
- ✅ Manual commit message mode
- ✅ Automatic branch detection
- ✅ Git commit execution
- ✅ Optional automatic push
- ✅ Gemini API key melalui environment variable

---

## 📄 License

Project ini dibuat sebagai project pribadi dan pembelajaran.
