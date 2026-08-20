# Fungsi & Subrutin

Fungsi pada DAP memungkinkan pengorganisasian kode modular yang dapat digunakan berulang kali dengan penerimaan argumen dan nilai kembalian.

---

## 1. Definisi Fungsi Banyak Baris

Mendefinisikan fungsi banyak baris menggunakan sintaksis `function <nama>(argumen) ... end`:

```javascript
program ContohFungsi
dictionary
    a, b, nilaiMaks : integer

    {{ Definisi fungsi }}
    function cariMaks(x, y)
        if x > y then
            return x
        else
            return y
        endif
    end
algorithm
    a <- 42
    b <- 99
    nilaiMaks <- cariMaks(a, b)
    print "Nilai maksimum adalah: ", nilaiMaks
endprogram
```

---

## 2. Fungsi Sebaris (Inline Arrow Function)

Untuk ekspresi matematika atau transformasi sederhana, gunakan sintaksis panah `->`:

```javascript
program ContohFungsiSebaris
dictionary
    angka, duaKali, kuadrat : integer

    function kaliDua(n) -> n * 2
    function pangkatDua(n) -> n ^ 2
algorithm
    angka <- 5
    duaKali <- kaliDua(angka)
    kuadrat <- pangkatDua(angka)

    print "Angka: ", angka
    print "Hasil Kali Dua: ", duaKali
    print "Hasil Pangkat Dua: ", kuadrat
endprogram
```

---

## 3. Pernyataan Kendali Fungsi & Perulangan

DAP menyediakan kata kunci pengendalian eksekusi:

### `return [ekspresi]`
Menghentikan eksekusi fungsi yang sedang berjalan dan mengembalikan nilai ke pemanggil.

### `break`
Menghentikan paksa perulangan terdalam (`for`, `while`, atau `repeat`).

```javascript
for i <- 1 to 100 do
    if i == 5 then
        break {{ Hentikan perulangan saat i mencapai 5 }}
    endif
    print i
endfor
```

### `continue`
Melewati sisa perintah pada iterasi saat ini dan langsung melompat ke iterasi berikutnya.

```javascript
for i <- 1 to 10 do
    if i % 2 == 0 then
        continue {{ Lewati bilangan genap }}
    endif
    print "Bilangan ganjil: ", i
endfor
```
