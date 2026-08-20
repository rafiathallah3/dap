# Input & Output

DAP provides built-in primitives for console input and output.

---

## 1. Output Statements (`write` / `print`)

DAP supports both `write` and `print` keywords (as well as uppercase `WRITE` and `PRINT`).

### Syntax
```javascript
write expression1, expression2, ...
print(expression1, expression2, ...)
```

### Examples
```javascript
program OutputDemo
dictionary
    score : real
    player : string
algorithm
    player <- "Rafi"
    score <- 95.5

    {{ Basic print }}
    print "Welcome to DAP!"

    {{ Multiple comma-separated arguments }}
    write "Player: ", player, " | Score: ", score

    {{ Expressions evaluated on the fly }}
    print "Next level target: ", score + 10.0
endprogram
```

---

## 2. Input Statements (`read` / `input`)

Read user input directly from standard input (stdin) into variables. Both `read` and `input` keywords (case-insensitive) are supported.

### Syntax
```javascript
read variableName
input(variableName)
```

### Examples
```javascript
program InputDemo
dictionary
    userAge : integer
    salary  : real
    name    : string
algorithm
    write "Enter your name: "
    read name

    write "Enter your age: "
    read userAge

    write "Enter your desired salary: "
    read salary

    print "Profile Summary: ", name, ", ", userAge, " years old, salary: ", salary
endprogram
```

> [!NOTE]
> DAP automatically converts input strings to the appropriate target variable type (e.g. parsing integer or real numbers when reading into numeric variables).
