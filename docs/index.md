# DAP Pseudocode Language

*A friendly, structured pseudocode language designed for algorithmic thinking and learning.*

---

## What is DAP?

**DAP** is an intuitive, structured pseudocode programming language and execution environment implemented in Go. Inspired by classical computer science curricula, DAP bridges the gap between academic algorithm design and executable code.

With DAP, you can write human-readable pseudocode featuring explicit dictionary declarations, intuitive control flow, flexible array ranges, custom data structures, and functions, running them instantly with zero external dependencies.

```javascript
program FactorialCalculation
dictionary
    n, result : integer
algorithm
    write "Enter a number: "
    read n

    result <- 1
    for i = 1 to n do
        {{ Calculate factorial accumulation }}
        result <- result * i
    endfor

    print "Factorial of ", n, " is ", result
endprogram
```

---

## Key Features

- **Structured Program Layout**: Explicit `dictionary` (declarations) and `algorithm` (execution) sections encourage clean variable tracking.
- **Instant Execution**: Statically compiled binary with fast execution and zero external runtime dependencies.
- **Rich Type System**: Primitive types (`integer`, `real`, `string`, `char`), custom structs (`type < ... >`), and type aliases.
- **Flexible Custom Range Arrays**: Declare arrays with arbitrary index bounds like `array[1..10]` or `array[100..200]`.
- **Versatile Control Flow**: `if-elif-else`, `while-do`, `for-to/downto-step-do`, `repeat-until`, and `depend on` (switch/case).
- **Developer Tooling**: Built-in interactive REPL console, token visualizer (`--show-token`), and Abstract Syntax Tree inspector (`--show-ast`).
- **Cross-Platform**: Works on Windows, Linux, and macOS (Intel & Apple Silicon).

---

## Quick Navigation

<div class="grid cards" markdown>

-   **[Installation](getting-started/installation.md)**

    ---

    Install DAP on Windows, macOS, Linux, or build directly from source.

-   **[Quickstart](getting-started/quickstart.md)**

    ---

    Write and execute your first DAP program in under 2 minutes.

-   **[Language Reference](language/structure.md)**

    ---

    Explore DAP syntax, variables, operators, loops, structs, and functions.

-   **[Examples](examples/algorithms.md)**

    ---

    Browse complete implementations of sorting, searching, and math algorithms.

</div>
