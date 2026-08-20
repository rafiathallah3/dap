# Types, Variables & Constants

All data types and variables in DAP must be declared inside the `dictionary` block before being utilized in the `algorithm` block.

---

## 1. Primitive Data Types

DAP provides four foundational primitive data types with predictable default zero-values:

| Type | Description | Default Value | Example Literals |
| :--- | :--- | :--- | :--- |
| `integer` | 64-bit signed integer | `0` | `0`, `42`, `-15` |
| `real` | Floating-point real number | `0.0` | `3.14`, `-0.5`, `100.0` |
| `string` | UTF-8 string text | `""` | `"Hello"`, `"DAP Pseudocode"` |
| `char` | Single character | `" "` (space) | `"a"`, `"Z"`, `"9"` |

---

## 2. Variable Declarations

Variables are declared with the format: `name : Type` or comma-separated lists `var1, var2 : Type`.

```javascript
program VariableDemo
dictionary
    {{ Single variable declarations }}
    count : integer
    score : real
    message : string

    {{ Multi-variable declarations }}
    x, y, z : integer
    firstName, lastName : string
algorithm
    x <- 10
    y <- 20
    z <- x + y
    print z
endprogram
```

---

## 3. Constants (`const`)

Constants are fixed values that cannot be reassigned or modified anywhere in the program.

Format: `const name = value` or `const name <- value`

```javascript
program ConstantDemo
dictionary
    const PI = 3.14159
    const MAX_USERS = 100
    const APP_NAME = "DAP Engine"
    
    radius : real
    area   : real
algorithm
    radius <- 5.0
    area <- PI * (radius ^ 2)
    print "Area: ", area
endprogram
```

> [!WARNING]
> Attempting to reassign a `const` variable inside the `algorithm` section (e.g. `PI <- 3.0`) will trigger a compile-time semantic error.

---

## 4. Type Aliases

You can define custom names for existing types using the `type AliasName : TargetType` syntax:

```javascript
program TypeAliasDemo
dictionary
    type Float : real
    type Age : integer

    myAge : Age
    price : Float
algorithm
    myAge <- 21
    price <- 99.99
    print "Age: ", myAge, ", Price: $", price
endprogram
```
