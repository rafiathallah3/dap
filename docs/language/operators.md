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
| `^` | Exponentiation | Raises to power | `2 ^ 3` | `8` |

---

## 3. Comparison Operators

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

## 4. Logical Operators

Logical operators combine boolean expressions:

| Operator | Name | Description | Example |
| :--- | :--- | :--- | :--- |
| `&&` | Logical AND | Returns `true` only if both operands are true | `(x > 0) && (x < 100)` |
| `\|\|` | Logical OR | Returns `true` if at least one operand is true | `(a == 0) \|\| (b == 0)` |
| `!` | Logical NOT | Inverts the boolean truth value | `!(x == 10)` |

---

## 5. Operator Precedence

From highest precedence to lowest:

1. Parentheses: `( ... )`
2. Member / Array access: `.`, `[ ]`
3. Exponentiation: `^`
4. Unary: `-`, `!`
5. Multiplicative: `*`, `/`, `div`, `%`
6. Additive: `+`, `-`
7. Relational: `<`, `<=`, `>`, `>=`
8. Equality: `==`, `!=`
9. Logical AND: `&&`
10. Logical OR: `||`
11. Assignment: `<-`, `=`
