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

## Building for Linux (amd64)

You can cross-compile a standalone Linux binary from any platform (Windows, macOS, Linux) with Go installed. No CGO or external libraries are required.

1. Build the binary into the `dist/` folder:
   ```powershell
   $env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -o dist/dap-linux-amd64 .
   ```

   On macOS/Linux (bash):
   ```bash
   GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/dap-linux-amd64 .
   ```

2. (Optional) Strip debug info to make a smaller binary by adding `-ldflags="-s -w"`:
   ```bash
   GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/dap-linux-amd64 .
   ```

3. Copy `dist/dap-linux-amd64` to your Linux machine, then run:
   ```bash
   chmod +x dap-linux-amd64
   ./dap-linux-amd64 program.dap
   ```

The result is a statically linked ELF executable that runs on any modern Linux distribution without extra dependencies. For ARM machines (e.g. Raspberry Pi, ARM servers), set `GOARCH=arm64` instead.
