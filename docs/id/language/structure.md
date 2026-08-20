# Struktur Program & Tata Letak

Setiap program DAP mengikuti struktur yang rapi dan disiplin, terinspirasi dari tata cara penulisan pseudocode algoritma perkuliahan.

---

## Format Standar Program

Berkas sumber DAP terdiri atas tiga bagian utama:

```javascript
program <NamaProgram>
dictionary
    {{ 1. Deklarasi (Tipe bentukan, Variabel, Konstanta) }}
algorithm
    {{ 2. Pernyataan eksekusi (Logika, Perulangan, I/O) }}
endprogram
```

### Penjelasan Bagian Program

1. **`program <NamaProgram>`**
   - Wajib menjadi pernyataan pertama di awal berkas.
   - `<NamaProgram>` adalah sebuah pengenal/identifier (misalnya `DeretFibonacci`, `Kalkulator`).

2. **Bagian `dictionary` (Kamus)**
   - Semua variabel, tipe bentukan (struct), konstanta, dan alias tipe wajib dideklarasikan di bagian ini.
   - Variabel yang digunakan pada bagian `algorithm` namun belum terdaftar di `dictionary` akan memicu galat semantik (compile error).

3. **Bagian `algorithm` (Algoritma)**
   - Berisi seluruh logika prosedural dan struktur kontrol program.
   - Penugasan variabel, percabangan, perulangan, definisi fungsi, dan perintah masukan/keluaran diletakkan di sini.

4. **`endprogram`**
   - Menandai akhir dari program.

---

## Komentar

Pada DAP, komentar ditulis di dalam kurung kurawal: `{{ komentar }}` atau `{ komentar }`.

```javascript
{{ Ini adalah contoh komentar di DAP }}
x <- 10 {{ Komentar di ujung baris }}
```

> [!NOTE]
> Semua teks di dalam `{ ... }` dan `{{ ... }}` akan diabaikan oleh lexer saat proses tokenisasi.

---

## Kata Kunci Penutup Blok

DAP mendukung penutup blok spesifik maupun kata kunci umum `end`:

| Jenis Blok | Penutup Spesifik | Penutup Umum |
| :--- | :--- | :--- |
| Program | `endprogram` | - |
| If / Elif | `endif` | `end` |
| Perulangan While | `endwhile` | `end` |
| Perulangan For | `endfor` | `end` |
| Perulangan Repeat | `until <kondisi>` | - |
| Depend on (Switch) | `enddependon` | `end` |
| Fungsi | `end` / `endif` | `end` |
