# Installation Guide

You can install DAP via web installers, by downloading pre-compiled standalone binaries, or by compiling from source with Go.

---

## Quick Web Installer (Recommended)

### Windows (PowerShell)
Open PowerShell and run the following command:

=== "PowerShell"
    ```powershell
    iwr -useb https://raw.githubusercontent.com/rafiathallah3/dap/main/install.ps1 | iex
    ```

### Linux & macOS (Bash / Zsh)
Open your terminal and run:

=== "Bash / Zsh"
    ```bash
    curl -fsSL https://raw.githubusercontent.com/rafiathallah3/dap/main/install.sh | bash
    ```

---

## Pre-compiled Binaries (GitHub Releases)

Download pre-compiled binaries from the [GitHub Releases Page](https://github.com/rafiathallah3/dap/releases):

| Operating System | Architecture | Package Archive |
| :--- | :--- | :--- |
| **Windows** | 64-bit (`x86_64` / `amd64`) | `dap-windows-amd64.zip` |
| **Linux** | 64-bit (`x86_64` / `amd64`) | `dap-linux-amd64.tar.gz` |
| **macOS** | Intel (`x86_64`) | `dap-darwin-amd64.tar.gz` |
| **macOS** | Apple Silicon (`arm64` / M1/M2/M3/M4) | `dap-darwin-arm64.tar.gz` |

After downloading, extract the archive and run the included `install.ps1` (Windows) or `./install.sh` (Linux/macOS) script to add `dap` to your system `PATH`.

---

## Build From Source

If you have Go 1.20+ installed on your machine:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/rafiathallah3/dap.git
   cd dap
   ```

2. **Run the local installer:**
   === "Windows"
       ```powershell
       .\install.ps1
       ```
   === "Linux / macOS"
       ```bash
       chmod +x install.sh
       ./install.sh
       ```

3. **Verify installation:**
   Open a new terminal window and run:
   ```bash
   dap
   ```
   If successful, you will enter the DAP interactive console.
