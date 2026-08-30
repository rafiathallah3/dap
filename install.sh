#!/bin/sh
set -e

# DAP Installer for Linux and macOS
# Downloads pre-built DAP binary or builds from source, then adds it to PATH.

REPO="rafiathallah3/dap"
INSTALL_DIR="$HOME/.dap/bin"
BINARY_NAME="dap"

mkdir -p "$INSTALL_DIR"

# Check if Go is installed and we are in the source directory
if command -v go >/dev/null 2>&1 && [ -f "cmd/dap/main.go" ]; then
    echo "Building DAP from source..."
    go build -o "$INSTALL_DIR/$BINARY_NAME" ./cmd/dap
else
    echo "Downloading DAP pre-built binary..."
    
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"

    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            echo "Error: Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    case "$OS" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            echo "Error: Unsupported operating system: $OS"
            exit 1
            ;;
    esac

    ASSET_NAME="dap-${OS}-${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET_NAME}"

    TMP_DIR="$(mktemp -d)"
    trap 'rm -rf "$TMP_DIR"' EXIT

    echo "Downloading from ${DOWNLOAD_URL}..."
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$ASSET_NAME"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$TMP_DIR/$ASSET_NAME" "$DOWNLOAD_URL"
    else
        echo "Error: Neither curl nor wget is available."
        exit 1
    fi

    tar -xzf "$TMP_DIR/$ASSET_NAME" -C "$TMP_DIR"
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

echo "DAP installed successfully at $INSTALL_DIR/$BINARY_NAME"

# Update Shell Profile
SHELL_NAME="$(basename "$SHELL")"
PROFILE_FILE=""

case "$SHELL_NAME" in
    bash)
        if [ -f "$HOME/.bashrc" ]; then
            PROFILE_FILE="$HOME/.bashrc"
        elif [ -f "$HOME/.bash_profile" ]; then
            PROFILE_FILE="$HOME/.bash_profile"
        fi
        ;;
    zsh)
        PROFILE_FILE="$HOME/.zshrc"
        ;;
    *)
        PROFILE_FILE="$HOME/.profile"
        ;;
esac

PATH_CMD="export PATH=\"\$HOME/.dap/bin:\$PATH\""

if [ -n "$PROFILE_FILE" ]; then
    if ! grep -q '\.dap/bin' "$PROFILE_FILE" 2>/dev/null; then
        echo "" >> "$PROFILE_FILE"
        echo "# DAP Pseudocode Runtime" >> "$PROFILE_FILE"
        echo "$PATH_CMD" >> "$PROFILE_FILE"
        echo "Added $INSTALL_DIR to PATH in $PROFILE_FILE"
    fi
fi

echo ""
echo "Installation complete! Restart your shell or run:"
echo "  export PATH=\"\$HOME/.dap/bin:\$PATH\""
echo "Then test with: dap"
