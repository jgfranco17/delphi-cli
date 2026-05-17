# Development scripts

INSTALL_PATH := "${HOME}/.claude/tools"

# Default command
_default:
    @just --list --unsorted

# Sync Go modules
tidy:
    go mod tidy
    @echo "All modules synced, Go workspace ready!"

# CLI local run wrapper
delphi *args:
    @go run . {{ args }}

# Run all BDD tests
test:
    @echo "Running unit tests!"
    @go clean -testcache
    go test -cover -race ./...

# Build the binary
build:
    #!/usr/bin/env bash
    set -euo pipefail
    # Detect OS and architecture
    case "$(uname -s)" in
        Linux*) OS="linux" ;;
        Darwin*) OS="darwin" ;;
        *) echo "Error: Unsupported OS (${OS})"; exit 1 ;;
    esac
    case "$(uname -m)" in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        arm64) ARCH="arm64" ;;
        *) echo "Error: Unsupported architecture (${ENV_ARCH})"; exit 1 ;;
    esac

    echo "Building CLI for ${OS}/${ARCH}..."
    go mod download all
    CGO_ENABLED=0 GOOS="${OS}" GOARCH="${ARCH}" go build -o ./delphi .
    echo "Built CLI for ${OS}/${ARCH} successfully!"

# Update the project dependencies
update-deps:
    @echo "Updating project dependencies..."
    go get -u ./...
    go mod tidy

# Install the binary locally
install-local: build
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Installing CLI locally..."
    BIN_PATH="{{ INSTALL_PATH }}/bin/delphi"
    cp ./delphi "${BIN_PATH}"
    chmod +x "${BIN_PATH}"
    echo "Installed CLI locally: ${BIN_PATH}"

# Remove the local binary
uninstall-local:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Uninstalling CLI..."
    BIN_PATH="{{ INSTALL_PATH }}/bin/delphi"
    rm "${BIN_PATH}"
    echo "Uninstalled CLI from ${BIN_PATH}"
