# DAP

[![Documentation](https://img.shields.io/badge/docs-GitHub_Pages-blue.svg)](https://rafiathallah3.github.io/dap/)
[![Release](https://img.shields.io/github/v/release/rafiathallah3/dap)](https://github.com/rafiathallah3/dap/releases)

A friendly pseudocode language and interpreter designed for algorithmic thinking and computer science education, inspired by classical academic curriculum.

[Documentation Website](https://rafiathallah3.github.io/dap/) | [Dokumentasi Bahasa Indonesia](https://rafiathallah3.github.io/dap/id/)

## Example Code

```javascript
program ThisIsAProgram
dictionary
    n : integer
algorithm
    {{ Inisialisasi nilai n }}
    n <- 1

    for i = 1 to 10 do
        n <- n * i
    endfor

    print n
endprogram
```

## Installation

### Quick Web Installer (Recommended)

#### Windows (PowerShell)
Open PowerShell and run:
```powershell
iwr -useb https://raw.githubusercontent.com/rafiathallah3/dap/main/install.ps1 | iex
```

#### Linux & macOS (Bash / Zsh)
Open Terminal and run:
```bash
curl -fsSL https://raw.githubusercontent.com/rafiathallah3/dap/main/install.sh | bash
```

---

### Manual Download from GitHub Releases
Download the pre-compiled archive for your OS from the [GitHub Releases Page](https://github.com/rafiathallah3/dap/releases):
- `dap-windows-amd64.zip` (Windows 64-bit)
- `dap-linux-amd64.tar.gz` (Linux 64-bit)
- `dap-darwin-amd64.tar.gz` (macOS Intel)
- `dap-darwin-arm64.tar.gz` (macOS Apple Silicon)

Extract the archive and run `install.ps1` (Windows) or `./install.sh` (Linux/macOS).

---

### Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/rafiathallah3/dap.git
   cd dap
   ```

2. Run the installer script:
   - Windows: `.\install.ps1`
   - Linux / macOS: `./install.sh`

3. Open a new terminal window and run:
   ```bash
   dap
   ```

## Usage
- Enter Console Mode: `dap`
- Check Version: `dap --version` (or `dap -v`)
- Update to Latest Release: `dap --update` (or `dap -u`)
- Run a File: `dap program.dap`
- Show Tokens: `dap program.dap --show-token`
- Show AST: `dap program.dap --show-ast`
- Show AST (JSON): `dap program.dap --show-ast-json`
- Show Help: `dap --help` (or `dap -h`)

### Interactive Console Commands
Inside the interactive REPL (`dap`):
- `version` - Check current version and check for GitHub updates
- `update` - Self-update DAP to the latest release
- `help` - View available commands and usage
- `exit` / `exit()` - Exit the interactive console

## Building for Linux (amd64)

You can cross-compile a standalone Linux binary from any platform with Go installed:

1. Build the binary into the `dist/` folder:
   ```powershell
   $env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -o dist/dap-linux-amd64 .
   ```

   On macOS/Linux:
   ```bash
   GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/dap-linux-amd64 .
   ```

2. Copy `dist/dap-linux-amd64` to your Linux machine, then run:
   ```bash
   chmod +x dap-linux-amd64
   ./dap-linux-amd64 program.dap
   ```
