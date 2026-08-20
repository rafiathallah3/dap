# Algorithm Examples

Explore real-world algorithms and data manipulation patterns implemented in DAP.

---

## 1. Linear Search & Min-Max Finding

Reads integers until a sentinel value `-999` is entered, finding the smallest number:

```javascript
program FindMinimumNumber
dictionary
    n, minVal : integer
algorithm
    write "Enter numbers (enter -999 to stop): "
    read n
    minVal <- n

    while n != -999 do
        if n < minVal then
            minVal <- n
        endif
        read n
    endwhile

    if minVal == -999 then
        print "No numbers were entered."
    else
        print "Minimum number is: ", minVal
    endif
endprogram
```

---

## 2. Bubble Sort

Sorts an array of 5 integers in ascending order:

```javascript
program BubbleSortExample
dictionary
    arr : array[1..5] of integer
    i, j, temp : integer
algorithm
    {{ Initialize unsorted array }}
    arr[1] <- 64
    arr[2] <- 25
    arr[3] <- 12
    arr[4] <- 22
    arr[5] <- 11

    print "Original array:"
    for i <- 1 to 5 do
        write arr[i], " "
    endfor
    print ""

    {{ Bubble sort algorithm }}
    for i <- 1 to 4 do
        for j <- 1 to (5 - i) do
            if arr[j] > arr[j + 1] then
                {{ Swap arr[j] and arr[j+1] }}
                temp <- arr[j]
                arr[j] <- arr[j + 1]
                arr[j + 1] <- temp
            endif
        endfor
    endfor

    print "Sorted array (Ascending):"
    for i <- 1 to 5 do
        write arr[i], " "
    endfor
    print ""
endprogram
```

---

## 3. Fibonacci Sequence Generator

Computes the first $N$ Fibonacci numbers:

```javascript
program FibonacciSeries
dictionary
    n, a, b, nextVal, i : integer
algorithm
    write "How many Fibonacci numbers to generate? "
    read n

    a <- 0
    b <- 1

    print "Fibonacci Sequence:"
    for i <- 1 to n do
        print a
        nextVal <- a + b
        a <- b
        b <- nextVal
    endfor
endprogram
```

---

## 4. Temperature Conversion System

Uses `depend on` and arithmetic operations to convert temperatures between Celsius, Fahrenheit, and Kelvin:

```javascript
program TemperatureConverter
dictionary
    tempC, tempF, tempK : real
    choice : integer
algorithm
    print "=== Temperature Converter ==="
    print "1. Celsius to Fahrenheit"
    print "2. Celsius to Kelvin"
    write "Select option (1-2): "
    read choice

    write "Enter temperature in Celsius: "
    read tempC

    depend on (choice)
        1:
            tempF <- (tempC * 9.0 / 5.0) + 32.0
            print tempC, " C = ", tempF, " F"
        2:
            tempK <- tempC + 273.15
            print tempC, " C = ", tempK, " K"
        default:
            print "Invalid option selected."
    enddependon
endprogram
```
