# Tipe Data, Variabel & Konstanta

Semua tipe data dan variabel pada DAP wajib dideklarasikan di dalam blok `dictionary` sebelum digunakan pada bagian `algorithm`.

---

## 1. Tipe Data Primitif

DAP menyediakan empat tipe data dasar dengan nilai awal (zero-value) yang terdefinisi jelas:

| Tipe Data | Deskripsi | Nilai Default | Contoh Nilai |
| :--- | :--- | :--- | :--- |
| `integer` | Bilangan bulat bertanda (64-bit) | `0` | `0`, `42`, `-15` |
| `real` | Bilangan pecahan / desimal (float) | `0.0` | `3.14`, `-0.5`, `100.0` |
| `string` | Teks rangkaian karakter (UTF-8) | `""` | `"Halo"`, `"Pseudocode DAP"` |
| `char` | Karakter tunggal | `" "` (spasi) | `"a"`, `"Z"`, `"9"` |

---

## 2. Deklarasi Variabel

Variabel dideklarasikan dengan pola: `namaVariabel : TipeData` atau daftar dipisahkan koma `var1, var2 : TipeData`.

```javascript
program ContohVariabel
dictionary
    {{ Deklarasi variabel tunggal }}
    jumlah : integer
    nilai  : real
    pesan  : string

    {{ Deklarasi banyak variabel sekaligus }}
    x, y, z : integer
    namaDepan, namaBelakang : string
algorithm
    x <- 10
    y <- 20
    z <- x + y
    print z
endprogram
```

---

## 3. Konstanta (`const`)

Konstanta adalah nilai tetap yang tidak dapat diubah atau ditugaskan kembali di mana pun dalam program.

Format: `const namaKonstanta = nilai` atau `const namaKonstanta <- nilai`

```javascript
program ContohKonstanta
dictionary
    const PI = 3.14159
    const MAKS_PENGGUNA = 100
    const NAMA_APLIKASI = "Mesin DAP"
    
    radius : real
    luas   : real
algorithm
    radius <- 5.0
    luas <- PI * (radius ^ 2)
    print "Luas Lingkaran: ", luas
endprogram
```

> [!WARNING]
> Mencoba mengubah nilai variabel `const` di dalam bagian `algorithm` (misalnya `PI <- 3.0`) akan memicu galat semantik saat kompilasi.

---

## 4. Alias Tipe (`type`)

Anda dapat membuat nama kustom untuk tipe data yang sudah ada menggunakan sintaksis `type NamaAlias : TipeTujuan`:

```javascript
program ContohAliasTipe
dictionary
    type Pecahan : real
    type Umur : integer

    umurSaya : Umur
    harga    : Pecahan
algorithm
    umurSaya <- 21
    harga <- 99.99
    print "Umur: ", umurSaya, ", Harga: Rp", harga
endprogram
```
