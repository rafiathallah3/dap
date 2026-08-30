# Operators & Expressions

DAP provides a complete suite of arithmetic, comparison, assignment, and logical operators.

---

## 1. Assignment Operators

| Operator | Syntax | Description | Example |
| :--- | :--- | :--- | :--- |
| `<-` | `x <- val` | Preferred arrow assignment | `total <- 100` |
| `=` | `x = val` | Standard assignment | `total = 100` |

> [!TIP]
> In idiomatic DAP pseudocode, `<-` is the recommended assignment symbol to clearly distinguish assignment from equality comparisons (`==`).
> Compound assignment operators (like `+=` or `*=`) are not used in standard DAP; use explicit expansion (`x <- x + 1`).

---

## 2. Arithmetic Operators

| Operator | Name | Behavior | Example | Result |
| :--- | :--- | :--- | :--- | :--- |
| `+` | Addition / Concatenation | Numeric sum or string joining | `5 + 3` or `"Hello " + "World"` | `8` or `"Hello World"` |
| `-` | Subtraction / Negation | Difference or unary negative | `10 - 4` or `-x` | `6` |
| `*` | Multiplication | Multiplies two numbers | `4 * 3` | `12` |
| `/` | Real Division | Performs floating-point division | `7 / 2` | `3.5` |
| `div` / `DIV` | Integer Division | Truncates decimal part | `7 div 2` | `3` |
| `%` | Modulo | Remainder of division | `7 % 3` | `1` |
| `^` / `**` | Exponentiation | Raises to power | `2 ** 3` or `2 ^ 3` | `8` |

---

## 3. Bitwise Operators

DAP provides full support for low-level bitwise operations on integer operands:

| Operator | Name | Behavior | Example | Result |
| :--- | :--- | :--- | :--- | :--- |
| `&` | Bitwise AND | Sets each bit to 1 if both bits are 1 | `12 & 10` (1100 & 1010) | `8` (1000) |
| `\|` | Bitwise OR | Sets each bit to 1 if either bit is 1 | `12 \| 10` (1100 \| 1010) | `14` (1110) |
| `^` | Bitwise XOR | Sets each bit to 1 if only one bit is 1 | `12 ^ 10` (1100 ^ 1010) | `6` (0110) |
| `~` | Bitwise NOT | Inverts all the bits (unary) | `~5` | `-6` |
| `<<` | Shift Left | Shifts bits to the left, zero-filled | `3 << 4` | `48` |
| `>>` | Shift Right | Shifts bits to the right | `64 >> 3` | `8` |

---

## 4. Comparison Operators

Comparison operators evaluate expressions and return boolean values (`true` or `false`):

| Operator | Description | Example | Result for `a=5, b=10` |
| :--- | :--- | :--- | :--- |
| `==` | Equal to | `a == b` | `false` |
| `!=` | Not equal to | `a != b` | `true` |
| `<` | Less than | `a < b` | `true` |
| `<=` | Less than or equal to | `a <= 5` | `true` |
| `>` | Greater than | `a > b` | `false` |
| `>=` | Greater than or equal to | `b >= 10` | `true` |

---

## 5. Logical Operators

Logical operators combine boolean expressions:

| Operator | Name | Description | Example |
| :--- | :--- | :--- | :--- |
| `and` / `&&` | Logical AND | Returns `true` only if both operands are true | `(x > 0) and (x < 100)` |
| `or` / `\|\|` | Logical OR | Returns `true` if at least one operand is true | `(a == 0) or (b == 0)` |
| `not` / `!` | Logical NOT | Inverts the boolean truth value | `not (x == 10)` |

---

## 6. Operator Precedence

From highest precedence to lowest:

1. Parentheses: `( ... )`
2. Member / Array access: `.`, `[ ]`
3. Exponentiation: `**`, `^`
4. Unary operators: `-`, `!`, `not`, `~` (Bitwise NOT)
5. Multiplicative: `*`, `/`, `div`, `%`
6. Additive: `+`, `-`
7. Bitwise Shifts: `<<`, `>>`
8. Bitwise AND: `&`
9. Bitwise XOR: `^`
10. Bitwise OR: `|`
11. Relational: `<`, `<=`, `>`, `>=`
12. Equality: `==`, `!=`
13. Logical AND: `and`, `&&`
14. Logical OR: `or`, `||`
15. Assignment: `<-`, `=`
