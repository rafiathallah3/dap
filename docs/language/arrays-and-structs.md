# Arrays & Structures

DAP provides composite data structures: custom-range arrays, 1-indexed string access, and user-defined records (structures).

---

## 1. Custom-Range Arrays

Unlike languages constrained strictly to 0-indexing, DAP lets you define arbitrary index bounds (`start..end`) matching mathematical and algorithmic specifications.

### Array Declaration
```javascript
arrayName : array[start..end] of ElementType
```

### Examples
```javascript
program ArrayExample
dictionary
    {{ 1-indexed array of 5 integers }}
    scores : array[1..5] of integer

    {{ Custom offset array (e.g. tracking years 2020..2025) }}
    sales : array[2020..2025] of real
    i : integer
algorithm
    {{ Assigning elements }}
    scores[1] <- 90
    scores[2] <- 85
    scores[3] <- 95
    scores[4] <- 70
    scores[5] <- 100

    sales[2020] <- 150000.50
    sales[2021] <- 180000.75

    {{ Iterating over the array }}
    for i <- 1 to 5 do
        print "Student ", i, " score = ", scores[i]
    endfor
endprogram
```

---

## 2. String Indexing

Strings in DAP can be indexed directly to read individual characters:

```javascript
program StringIndexDemo
dictionary
    str : string
    firstChar : string
algorithm
    str <- "Algorithm"
    firstChar <- str[1] {{ Returns "A" (1-indexed) }}
    print "First character: ", firstChar
endprogram
```

---

## 3. Custom Structures (`type < ... >`)

Structures allow you to group related fields of different data types into a single composite type.

### Defining a Struct Type
```javascript
type StructName <
    field1 : type1
    field2 : type2
    ...
>
```

### Declaring & Accessing Struct Variables
```javascript
program StructDemo
dictionary
    {{ 1. Define struct layout }}
    type Student <
        id    : integer
        name  : string
        gpa   : real
    >

    {{ 2. Declare variables of this type }}
    student1 : Student
algorithm
    {{ 3. Access and mutate fields using dot notation (.) }}
    student1.id <- 101
    student1.name <- "Alice"
    student1.gpa <- 3.92

    print "ID: ", student1.id
    print "Name: ", student1.name
    print "GPA: ", student1.gpa
endprogram
```

---

## 4. Arrays of Structures

Combine arrays and structures to model collections of complex records:

```javascript
program ArrayOfStructsDemo
dictionary
    type Product <
        name  : string
        price : real
        stock : integer
    >

    inventory : array[1..3] of Product
    i : integer
algorithm
    inventory[1].name <- "Laptop"
    inventory[1].price <- 999.99
    inventory[1].stock <- 10

    inventory[2].name <- "Wireless Mouse"
    inventory[2].price <- 25.50
    inventory[2].stock <- 50

    inventory[3].name <- "Mechanical Keyboard"
    inventory[3].price <- 75.00
    inventory[3].stock <- 30

    for i <- 1 to 3 do
        print "Product: ", inventory[i].name, " | Price: $", inventory[i].price, " | Stock: ", inventory[i].stock
    endfor
endprogram
```
