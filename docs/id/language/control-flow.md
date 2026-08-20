# Alur Kontrol

DAP menyediakan struktur kendali alur yang lengkap: percabangan kondisi, berbagai macam perulangan, dan blok seleksi kasus (`depend on`).

---

## 1. Pernyataan Percabangan (`if ... then ... else`)

### Bentuk Percabangan Banyak Baris
```javascript
program PenilaianNilai
dictionary
    nilai : integer
algorithm
    read nilai

    if nilai >= 90 then
        write "Predikat A"
    elif nilai >= 80 then
        write "Predikat B"
    elif nilai >= 70 then
        write "Predikat C"
    else
        write "Predikat E"
    endif
endprogram
```

### Bentuk Percabangan Satu Baris
```javascript
if x < 0 then x <- -x
```

---

## 2. Perulangan While (`while ... do ... endwhile`)

Mengulang serangkaian instruksi selama kondisi evaluasi bernilai `true`.

```javascript
program ContohWhile
dictionary
    pencacah : integer
algorithm
    pencacah <- 1
    while pencacah <= 5 do
        print "Hitungan: ", pencacah
        pencacah <- pencacah + 1
    endwhile
endprogram
```

---

## 3. Perulangan For (`for ... to / downto ... step ... do`)

Perulangan `for` pada DAP sangat intuitif dan mudah dibaca:

### Penaikan Standar (`to`)
```javascript
for i <- 1 to 10 do
    print i
endfor
```

### Kelipatan Langkah Kustom (`step`)
```javascript
{{ Langkah per 2: 0, 2, 4, 6, 8, 10 }}
for i <- 0 to 10 step 2 do
    print i
endfor
```

### Perulangan Menurun (`downto`)
```javascript
{{ Hitung mundur: 10, 9, 8 ... 1 }}
for i <- 10 downto 1 do
    print "Hitung Mundur: ", i
endfor
```

---

## 4. Perulangan Repeat-Until (`repeat ... until`)

Perulangan dengan pemeriksaan kondisi di akhir. Blok instruksi dieksekusi **minimal satu kali** dan terus berjalan sampai kondisi menjadi `true`.

```javascript
program ContohRepeatUntil
dictionary
    angka : integer
algorithm
    angka <- 1
    repeat
        print "Nilai: ", angka
        angka <- angka * 2
    until angka > 50
endprogram
```

---

## 5. Seleksi Kasus (`depend on ... enddependon`)

DAP mendukung dua pola penulisan blok `depend on`:

### Pola 1: Pencocokan Nilai Variabel
```javascript
depend on (nomorHari)
    1 : print "Senin"
    2 : print "Selasa"
    3 : print "Rabu"
    default : print "Hari lainnya"
enddependon
```

### Pola 2: Pencocokan Kondisi / Ekspresi
```javascript
depend on
    nilai >= 90 : print "Sangat Memuaskan"
    nilai >= 75 : print "Baik"
    nilai >= 50 : print "Cukup"
    default : print "Perlu Perbaikan"
enddependon
```
