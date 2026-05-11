# DAP

A friendly pseudocode that I implemented, inspired by my university professor, Mr. Jimmy. To prove that I can be like him. One day.

## Example code:

```javascript
program ThisIsAProgram
dictionary
    n : integer
algorithm
    n <- 1

    for i = 1 to 10 do
        n <- n * i
    endfor

    print n
endprogram
```

## Installation

1. Clone the repository:
   ```powershell
   git clone https://github.com/rafiathallah3/dap.git
   cd dap
   ```

2. Run the installer script (Windows):
   ```powershell
   .\install.ps1
   ```

3. Open a **new** terminal window and run:
   ```powershell
   dap
   ```

## Usage
- Enter Console Mode: `dap`
- Run a File: `dap program.dap`
- Show Tokens: `dap program.dap --show-token`
- Show AST: `dap program.dap --show-ast`
