# Operator & Ekspresi

DAP menyediakan kumpulan operator aritmetika, perbandingan, penugasan, dan logika yang lengkap.

---

## 1. Operator Penugasan (Assignment)

| Operator | Sintaksis | Deskripsi | Contoh |
| :--- | :--- | :--- | :--- |
| `<-` | `x <- nilai` | Simbol penugasan **utama** (panah kiri) | `total <- 100` |
| `=` | `x = nilai` | Simbol penugasan standar | `total = 100` |

> [!TIP]
> Dalam penulisan pseudocode DAP yang baku, simbol `<-` sangat dianjurkan agar berbeda jelas dengan operator pembanding kesetaraan (`==`).
> Operator penugasan majemuk (seperti `+=` atau `*=`) tidak digunakan dalam DAP standar; gunakan bentuk ekspansi jelas (`x <- x + 1`).

---

## 2. Operator Aritmetika

| Operator | Nama | Penjelasan | Contoh | Hasil |
| :--- | :--- | :--- | :--- | :--- |
| `+` | Penjumlahan / Penggabungan | Jumlah angka atau gabung teks | `5 + 3` atau `"Halo " + "Dunia"` | `8` atau `"Halo Dunia"` |
| `-` | Pengurangan / Negasi | Selisih nilai atau nilai minus | `10 - 4` atau `-x` | `6` |
| `*` | Perkalian | Mengalikan dua bilangan | `4 * 3` | `12` |
| `/` | Pembagian Real | Pembagian desimal/pecahan | `7 / 2` | `3.5` |
| `div` / `DIV` | Pembagian Bulat | Pembagian bilangan bulat (tanpa sisa) | `7 div 2` | `3` |
| `%` | Modulo | Sisa hasil bagi | `7 % 3` | `1` |
| `^` | Pemangkatan | Pangkat bilangan | `2 ^ 3` | `8` |

---

## 3. Operator Perbandingan (Relasional)

Operator perbandingan membandingkan dua ekspresi dan mengembalikan nilai kebenaran (`true` atau `false`):

| Operator | Deskripsi | Contoh | Hasil jika `a=5, b=10` |
| :--- | :--- | :--- | :--- |
| `==` | Sama dengan | `a == b` | `false` |
| `!=` | Tidak sama dengan | `a != b` | `true` |
| `<` | Kurang dari | `a < b` | `true` |
| `<=` | Kurang dari atau sama dengan | `a <= 5` | `true` |
| `>` | Lebih dari | `a > b` | `false` |
| `>=` | Lebih dari atau sama dengan | `b >= 10` | `true` |

---

## 4. Operator Logika

Operator logika menggabungkan dua atau lebih kondisi boolean:

| Operator | Nama | Deskripsi | Contoh |
| :--- | :--- | :--- | :--- |
| `&&` | Logika DAN (AND) | Menghasilkan `true` jika kedua kondisi bernilai benar | `(x > 0) && (x < 100)` |
| `\|\|` | Logika ATAU (OR) | Menghasilkan `true` jika salah satu kondisi bernilai benar | `(a == 0) \|\| (b == 0)` |
| `!` | Logika BUKAN (NOT) | Membalik nilai kebenaran suatu kondisi | `!(x == 10)` |

---

## 5. Prioritas Operator (Precedence)

Urutan evaluasi operator dari prioritas tertinggi ke terendah:

1. Tanda kurung: `( ... )`
2. Akses anggota record / elemen array: `.`, `[ ]`
3. Pemangkatan: `^`
4. Unary: `-`, `!`
5. Perkalian & pembagian: `*`, `/`, `div`, `%`
6. Penjumlahan & pengurangan: `+`, `-`
7. Relasional: `<`, `<=`, `>`, `>=`
8. Kesetaraan: `==`, `!=`
9. Logika DAN: `&&`
10. Logika ATAU: `||`
11. Penugasan: `<-`, `=`
