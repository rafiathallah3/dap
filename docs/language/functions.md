# Functions & Subroutines

Functions in DAP allow you to structure reusable logic with parameter passing and return values.

---

## 1. Multi-line Function Definition

Define multi-line functions with `function <name>(args) ... end`:

```javascript
program FunctionExample
dictionary
    a, b, maxVal : integer

    {{ Function definition }}
    function findMax(x, y)
        if x > y then
            return x
        else
            return y
        endif
    end
algorithm
    a <- 42
    b <- 99
    maxVal <- findMax(a, b)
    print "The maximum value is: ", maxVal
endprogram
```

---

## 2. Inline Arrow Functions

For concise mathematical expressions or transformations, use the arrow `->` syntax:

```javascript
program InlineFunctionExample
dictionary
    num, doubled, squared : integer

    function doubleVal(n) -> n * 2
    function squareVal(n) -> n ^ 2
algorithm
    num <- 5
    doubled <- doubleVal(num)
    squared <- squareVal(num)

    print "Number: ", num
    print "Doubled: ", doubled
    print "Squared: ", squared
endprogram
```

---

## 3. Control Flow Statements in Functions & Loops

DAP supports standard subroutine and loop control keywords:

### `return [expression]`
Exits the current function and optionally returns a value back to the caller.

### `break`
Immediately terminates the innermost active loop (`for`, `while`, or `repeat`).

```javascript
for i <- 1 to 100 do
    if i == 5 then
        break {{ Stop the loop when i reaches 5 }}
    endif
    print i
endfor
```

### `continue`
Skips the remainder of the current iteration and proceeds to the next cycle.

```javascript
for i <- 1 to 10 do
    if i % 2 == 0 then
        continue {{ Skip even numbers }}
    endif
    print "Odd number: ", i
endfor
```
