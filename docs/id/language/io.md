# Masukan & Keluaran (I/O)

DAP menyediakan instruksi dasar untuk membaca data dari pengguna dan mencetak hasil ke layar konsol.

---

## 1. Menampilkan Keluaran (`write` / `print`)

DAP mendukung kata kunci `write` dan `print` (serta huruf besar `WRITE` dan `PRINT`).

### Sintaksis
```javascript
write ekspresi1, ekspresi2, ...
print(ekspresi1, ekspresi2, ...)
```

### Contoh
```javascript
program ContohKeluaran
dictionary
    nilai : real
    namaPemain : string
algorithm
    namaPemain <- "Rafi"
    nilai <- 95.5

    {{ Mencetak teks biasa }}
    print "Selamat datang di DAP!"

    {{ Mencetak beberapa variabel dan teks sekaligus }}
    write "Pemain: ", namaPemain, " | Nilai: ", nilai

    {{ Mencetak hasil evaluasi ekspresi }}
    print "Target nilai berikutnya: ", nilai + 10.0
endprogram
```

---

## 2. Membaca Masukan (`read` / `input`)

Membaca nilai masukan dari pengguna melalui standar masukan (stdin). Kata kunci `read` dan `input` (tidak peka huruf besar/kecil) didukung secara penuh.

### Sintaksis
```javascript
read namaVariabel
input(namaVariabel)
```

### Contoh
```javascript
program ContohMasukan
dictionary
    umur : integer
    gaji : real
    nama : string
algorithm
    write "Masukkan nama Anda: "
    read nama

    write "Masukkan umur Anda: "
    read umur

    write "Masukkan estimasi gaji: "
    read gaji

    print "Ringkasan Profil: ", nama, ", ", umur, " tahun, gaji: Rp", gaji
endprogram
```

> [!NOTE]
> DAP secara otomatis mengonversi string masukan menjadi tipe data target variabel yang sesuai (seperti mengonversi teks angka menjadi bilangan bulat atau pecahan saat membaca variabel `integer` atau `real`).
