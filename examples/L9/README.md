# Question Example

![soal text](Soal.png "Question")

Answer in DAP language:
```javascript
program MinHowManyAreYou
dictionary
    n, palingKecil, total : integer
algorithm
    read n
    palingKecil <- n
    total <- 1

    while n != -241231 do
        if palingKecil > n then
            palingKecil <- n
            total <- 0
        endif

        if palingKecil == n then
            total <- total + 1
        endif

        read n
    endwhile

    if palingKecil == -241231 then
        write "NONE"
    else
        write total
    endif
endprogram
```