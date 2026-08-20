# Array & Record (Tipe Bentukan)

DAP menyediakan struktur data majemuk: array dengan rentang indeks bebas, pengaksesan karakter string 1-indeks, dan tipe bentukan kustom (record / struct).

---

## 1. Array dengan Rentang Indeks Kustom

Berbeda dengan bahasa yang membatasi indeks mulai dari 0, DAP memungkinkan Anda menentukan batas indeks awal dan akhir secara bebas (`awal..akhir`) sesuai dengan notasi matematika dan algoritma buku teks.

### Deklarasi Array
```javascript
namaArray : array[awal..akhir] of TipeElemen
```

### Contoh
```javascript
program ContohArray
dictionary
    {{ Array 1-indeks berisi 5 integer }}
    daftarNilai : array[1..5] of integer

    {{ Array rentang tahun 2020..2025 }}
    penjualan : array[2020..2025] of real
    i : integer
algorithm
    {{ Mengisi elemen array }}
    daftarNilai[1] <- 90
    daftarNilai[2] <- 85
    daftarNilai[3] <- 95
    daftarNilai[4] <- 70
    daftarNilai[5] <- 100

    penjualan[2020] <- 150000.50
    penjualan[2021] <- 180000.75

    {{ Menampilkan isi array }}
    for i <- 1 to 5 do
        print "Nilai mahasiswa ke-", i, " = ", daftarNilai[i]
    endfor
endprogram
```

---

## 2. Pengindeksan String

String pada DAP dapat diakses karakternya menggunakan tanda kurung siku `[ ]`:

```javascript
program ContohIndeksString
dictionary
    teks : string
    hurufPertama : string
algorithm
    teks <- "Algoritma"
    hurufPertama <- teks[1] {{ Menghasilkan "A" (1-indexed) }}
    print "Huruf pertama: ", hurufPertama
endprogram
```

---

## 3. Tipe Bentukan Kustom / Record (`type < ... >`)

Tipe bentukan memungkinkan pengelompokan beberapa atribut dengan tipe data berbeda ke dalam satu tipe data majemuk.

### Mendefinisikan Tipe Bentukan
```javascript
type NamaBentukan <
    field1 : tipe1
    field2 : tipe2
    ...
>
```

### Deklarasi & Akses Anggota Record
```javascript
program ContohStruct
dictionary
    {{ 1. Mendefinisikan struktur record }}
    type Mahasiswa <
        nim  : integer
        nama : string
        ipk  : real
    >

    {{ 2. Mendeklarasikan variabel bertipe Mahasiswa }}
    mhs1 : Mahasiswa
algorithm
    {{ 3. Mengisi dan mengakses data menggunakan notasi titik (.) }}
    mhs1.nim <- 130121001
    mhs1.nama <- "Budi"
    mhs1.ipk <- 3.85

    print "NIM  : ", mhs1.nim
    print "Nama : ", mhs1.nama
    print "IPK  : ", mhs1.ipk
endprogram
```

---

## 4. Array dari Record (Array of Structures)

Anda dapat menggabungkan array dan tipe bentukan untuk mengelola koleksi data:

```javascript
program ContohTabelData
dictionary
    type Produk <
        nama  : string
        harga : real
        stok  : integer
    >

    inventaris : array[1..3] of Produk
    i : integer
algorithm
    inventaris[1].nama <- "Laptop"
    inventaris[1].harga <- 9999000.0
    inventaris[1].stok <- 10

    inventaris[2].nama <- "Mouse Nirkabel"
    inventaris[2].harga <- 150000.0
    inventaris[2].stok <- 50

    inventaris[3].nama <- "Keyboard Mekanikal"
    inventaris[3].harga <- 750000.0
    inventaris[3].stok <- 30

    for i <- 1 to 3 do
        print "Produk: ", inventaris[i].nama, " | Harga: Rp", inventaris[i].harga, " | Stok: ", inventaris[i].stok
    endfor
endprogram
```
