# Alat Bantu & Panduan CLI

Executable `dap` dilengkapi sejumlah alat bantu untuk mengeksekusi kode, memeriksa token lexer, melihat Abstract Syntax Tree (AST), dan kompilasi lintas platform.

---

## Ringkasan Perintah

```text
Penggunaan:
  dap [opsi] [berkas]

Perintah & Opsi:
  dap                       Membuka konsol interaktif (REPL)
  dap <berkas.dap>          Mengeksekusi berkas kode DAP
  dap --version, -v         Memeriksa versi DAP dan update dari GitHub
  dap --update, -u          Memperbarui DAP otomatis ke rilis terbaru
  dap <berkas.dap> --show-token Memeriksa alur token dari lexer
  dap <berkas.dap> --show-ast   Memeriksa struktur Abstract Syntax Tree (AST)
  dap <berkas.dap> --show-ast-json Memeriksa struktur AST dalam format JSON
  dap --help, -h            Menampilkan bantuan dan panduan penggunaan
```

---

## 1. Menjalankan Berkas Kode

Untuk mengeksekusi berkas sumber `.dap`:

```bash
dap program.dap
```

Jika terdapat kesalahan sintaksis atau semantik pada kode Anda, DAP akan menampilkan pesan galat lengkap dengan lokasi baris dan kolom yang bersangkutan.

---

## 2. Memeriksa Alur Token (`--show-token`)

Untuk melihat bagaimana lexer DAP memecah kode sumber menjadi token:

```bash
dap program.dap --show-token
```

### Contoh Keluaran
```text
[Line 1, Col 1]  TOKEN_PROGRAM: "program"
[Line 1, Col 9]  TOKEN_IDENTIFIER: "ProgramKu"
[Line 2, Col 1]  TOKEN_DICTIONARY: "dictionary"
[Line 3, Col 5]  TOKEN_IDENTIFIER: "x"
[Line 3, Col 7]  TOKEN_COLON: ":"
[Line 3, Col 9]  TOKEN_TYPE: "integer"
...
```

---

## 3. Memeriksa Struktur AST (`--show-ast`)

Untuk melihat representasi hierarki pohon sintaksis abstrak (AST) yang dihasilkan oleh parser:

```bash
dap program.dap --show-ast
```

---

## 4. Kompilasi Lintas Platform (Cross-Compilation)

Karena DAP ditulis dalam bahasa Go murni tanpa dependensi CGO, Anda dapat melakukan kompilasi untuk berbagai sistem operasi lain langsung dari komputer Anda.

### Kompilasi untuk Linux (amd64)
=== "Windows (PowerShell)"
    ```powershell
    $env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -ldflags="-s -w" -o dist/dap-linux-amd64 .
    ```

=== "macOS / Linux (Bash)"
    ```bash
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/dap-linux-amd64 .
    ```

### Kompilasi untuk macOS (Apple Silicon arm64)
```bash
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/dap-darwin-arm64 .
```

### Kompilasi untuk Windows (amd64)
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/dap-windows-amd64.exe .
```
