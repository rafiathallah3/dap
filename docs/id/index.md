# Bahasa Pseudocode DAP

*Bahasa pseudocode yang terstruktur dan mudah dipahami, dirancang untuk pembelajaran logika dan perancangan algoritma.*

---

## Apa itu DAP?

**DAP** adalah bahasa pemrograman dan interpreter pseudocode yang intuitif serta terstruktur, diimplementasikan menggunakan bahasa Go. Terinspirasi dari kurikulum dasar ilmu komputer dan algoritma perkuliahan, DAP menjembatani perancangan algoritma di atas kertas dengan kode yang dapat langsung dieksekusi di komputer.

Dengan DAP, Anda dapat menulis pseudocode dengan deklarasi kamus (dictionary) yang rapi, struktur kontrol yang intuitif, rentang array yang fleksibel, tipe bentukan (struct), serta fungsi—dan langsung menjalankannya tanpa konfigurasi yang rumit.

```javascript
program HitungFaktorial
dictionary
    n, hasil : integer
algorithm
    write "Masukkan angka: "
    read n

    hasil <- 1
    for i = 1 to n do
        {{ Menghitung akumulasi faktorial }}
        hasil <- hasil * i
    endfor

    print "Faktorial dari ", n, " adalah ", hasil
endprogram
```

---

## Fitur Utama

- **Struktur Program Teratur**: Pemisahan jelas antara bagian deklarasi (`dictionary`) dan eksekusi (`algorithm`) untuk melatih kedisiplinan pelacakan variabel.
- **Eksekusi Cepat**: Binary mandiri (statically compiled) dengan eksekusi instan tanpa memerlukan runtime tambahan.
- **Sistem Tipe Data Lengkap**: Tipe data dasar (`integer`, `real`, `string`, `char`), tipe bentukan (`type < ... >`), dan alias tipe.
- **Rentang Array Fleksibel**: Mendukung indeks awal dan akhir bebas, contohnya `array[1..10]` atau `array[100..200]`.
- **Struktur Kontrol Lengkap**: `if-elif-else`, `while-do`, `for-to/downto-step-do`, `repeat-until`, dan `depend on` (switch/case).
- **Alat Bantu Pengembang (CLI)**: Dilengkapi REPL interaktif, visualisasi token (`--show-token`), dan penampil Abstract Syntax Tree (`--show-ast`).
- **Lintas Platform**: Berjalan lancar di Windows, Linux, dan macOS (Intel & Apple Silicon).

---

## Navigasi Cepat

<div class="grid cards" markdown>

-   **[Panduan Instalasi](getting-started/installation.md)**

    ---

    Cara memasang DAP di Windows, macOS, Linux, atau kompilasi dari kode sumber.

-   **[Memulai Cepat](getting-started/quickstart.md)**

    ---

    Tulis dan jalankan program DAP pertama Anda dalam 2 menit.

-   **[Panduan Bahasa](language/structure.md)**

    ---

    Pelajari sintaksis DAP, variabel, operator, perulangan, bentukan, dan fungsi.

-   **[Contoh Algoritma](examples/algorithms.md)**

    ---

    Kumpulan implementasi lengkap seperti Bubble Sort, Fibonacci, dan pencarian.

</div>
