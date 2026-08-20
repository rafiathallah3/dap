# Contoh Algoritma

Kumpulan implementasi algoritma dan manipulasi data menggunakan bahasa DAP.

---

## 1. Pencarian Nilai Minimum (Linear Search / Sentinel)

Membaca bilangan bulat berulang kali hingga pengguna memasukkan nilai sentinel `-999`, lalu mencari nilai terkecil:

```javascript
program CariNilaiTerkecil
dictionary
    n, palingKecil : integer
algorithm
    write "Masukkan angka (ketik -999 untuk selesai): "
    read n
    palingKecil <- n

    while n != -999 do
        if n < palingKecil then
            palingKecil <- n
        endif
        read n
    endwhile

    if palingKecil == -999 then
        print "Tidak ada angka yang dimasukkan."
    else
        print "Angka terkecil adalah: ", palingKecil
    endif
endprogram
```

---

## 2. Pengurutan Gelembung (Bubble Sort)

Mengurutkan array berisi 5 bilangan bulat secara menaik (ascending):

```javascript
program PengurutanBubbleSort
dictionary
    arr : array[1..5] of integer
    i, j, temp : integer
algorithm
    {{ Inisialisasi data array acak }}
    arr[1] <- 64
    arr[2] <- 25
    arr[3] <- 12
    arr[4] <- 22
    arr[5] <- 11

    print "Array awal sebelum diurutkan:"
    for i <- 1 to 5 do
        write arr[i], " "
    endfor
    print ""

    {{ Algoritma Bubble Sort }}
    for i <- 1 to 4 do
        for j <- 1 to (5 - i) do
            if arr[j] > arr[j + 1] then
                {{ Tukar posisi arr[j] dan arr[j+1] }}
                temp <- arr[j]
                arr[j] <- arr[j + 1]
                arr[j + 1] <- temp
            endif
        endfor
    endfor

    print "Array setelah diurutkan (Menaik):"
    for i <- 1 to 5 do
        write arr[i], " "
    endfor
    print ""
endprogram
```

---

## 3. Pembangkit Deret Fibonacci

Menghasilkan $N$ bilangan pertama dari deret Fibonacci:

```javascript
program PembangkitFibonacci
dictionary
    n, a, b, selanjutnya, i : integer
algorithm
    write "Berapa banyak bilangan Fibonacci yang ingin dibuat? "
    read n

    a <- 0
    b <- 1

    print "Deret Fibonacci:"
    for i <- 1 to n do
        print a
        selanjutnya <- a + b
        a <- b
        b <- selanjutnya
    endfor
endprogram
```

---

## 4. Konversi Suhu

Menggunakan `depend on` dan operasi aritmetika untuk mengubah suhu Celcius ke Fahrenheit atau Kelvin:

```javascript
program KonversiSuhu
dictionary
    suhuC, suhuF, suhuK : real
    pilihan : integer
algorithm
    print "=== Sistem Konversi Suhu ==="
    print "1. Celcius ke Fahrenheit"
    print "2. Celcius ke Kelvin"
    write "Pilih menu (1-2): "
    read pilihan

    write "Masukkan nilai suhu Celcius: "
    read suhuC

    depend on (pilihan)
        1:
            suhuF <- (suhuC * 9.0 / 5.0) + 32.0
            print suhuC, " C = ", suhuF, " F"
        2:
            suhuK <- suhuC + 273.15
            print suhuC, " C = ", suhuK, " K"
        default:
            print "Pilihan menu tidak valid."
    enddependon
endprogram
```
