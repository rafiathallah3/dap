# Panduan Instalasi

Anda dapat memasang DAP melalui skrip installer web otomatis, mengunduh berkas binary dari rilis GitHub, atau melakukan kompilasi langsung dari kode sumber menggunakan Go.

---

## Installer Web Otomatis (Direkomendasikan)

### Windows (PowerShell)
Buka PowerShell dan jalankan perintah berikut:

=== "PowerShell"
    ```powershell
    iwr -useb https://raw.githubusercontent.com/rafiathallah3/dap/main/install.ps1 | iex
    ```

### Linux & macOS (Bash / Zsh)
Buka terminal dan jalankan:

=== "Bash / Zsh"
    ```bash
    curl -fsSL https://raw.githubusercontent.com/rafiathallah3/dap/main/install.sh | bash
    ```

---

## Berkas Binary Siap Pakai (GitHub Releases)

Unduh berkas arsip sesuai sistem operasi Anda melalui [Halaman Rilis GitHub](https://github.com/rafiathallah3/dap/releases):

| Sistem Operasi | Arsitektur | Berkas Arsip |
| :--- | :--- | :--- |
| **Windows** | 64-bit (`x86_64` / `amd64`) | `dap-windows-amd64.zip` |
| **Linux** | 64-bit (`x86_64` / `amd64`) | `dap-linux-amd64.tar.gz` |
| **macOS** | Intel (`x86_64`) | `dap-darwin-amd64.tar.gz` |
| **macOS** | Apple Silicon (`arm64` / M1/M2/M3/M4) | `dap-darwin-arm64.tar.gz` |

Setelah mengunduh, ekstrak berkas arsip dan jalankan skrip `install.ps1` (Windows) atau `./install.sh` (Linux/macOS) untuk mendaftarkan `dap` ke `PATH` sistem Anda secara otomatis.

---

## Kompilasi dari Kode Sumber (Build from Source)

Jika Anda telah menginstal Go versi 1.20+:

1. **Kloning repositori:**
   ```bash
   git clone https://github.com/rafiathallah3/dap.git
   cd dap
   ```

2. **Jalankan skrip instalasi lokal:**
   === "Windows"
       ```powershell
       .\install.ps1
       ```
   === "Linux / macOS"
       ```bash
       chmod +x install.sh
       ./install.sh
       ```

3. **Verifikasi hasil instalasi:**
   Buka jendela terminal baru lalu ketik:
   ```bash
   dap
   ```
   Jika berhasil, Anda akan masuk ke dalam konsol interaktif DAP.
