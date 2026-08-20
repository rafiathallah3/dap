# Program Structure & Layout

Every DAP program follows a clear, disciplined structure inspired by classical computer science pseudocode.

---

## The Standard Template

A DAP source file consists of three mandatory parts:

```javascript
program <ProgramName>
dictionary
    {{ 1. Declarations (Types, Variables, Constants) }}
algorithm
    {{ 2. Executable statements (Logic, Loops, I/O) }}
endprogram
```

### Sections Explained

1. **`program <ProgramName>`**
   - Must be the first statement of the file.
   - `<ProgramName>` is an identifier (e.g. `FibonacciSequence`, `Calculator`).

2. **`dictionary` Section**
   - All variables, custom struct types, constants, and type aliases must be declared here.
   - Undeclared variables used in the `algorithm` section will raise a compile-time semantic error.

3. **`algorithm` Section**
   - Contains all procedural and control flow logic.
   - Variable assignments, loop constructs, conditionals, function definitions, and I/O statements are written here.

4. **`endprogram`**
   - Marks the termination of the program.

---

## Comments

In DAP, comments are enclosed using curly braces: `{{ comment }}` or `{ comment }`.

```javascript
{{ This is a comment in DAP }}
x <- 10 {{ This is an inline comment }}
```

> [!NOTE]
> The lexer skips everything inside curly braces `{ ... }` and `{{ ... }}` during tokenization.

---

## Block Terminators

DAP supports both specific end keywords and generic `end` keywords for readability:

| Block Type | Specific Ending | Generic Ending |
| :--- | :--- | :--- |
| Program | `endprogram` | - |
| If / Elif | `endif` | `end` |
| While loop | `endwhile` | `end` |
| For loop | `endfor` | `end` |
| Repeat loop | `until <condition>` | - |
| Depend on (Switch) | `enddependon` | `end` |
| Function | `end` / `endif` | `end` |
