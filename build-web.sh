#!/usr/bin/env bash
set -euo pipefail

echo "==> Checking Go..."
go version || { echo "ERROR: Go >= 1.21 required: https://go.dev/dl/"; exit 1; }

echo "==> Building frpc-web..."
go build -ldflags "-s -w" -o frpc-web ./cmd/frpc-web/

echo ""
echo "✓  Done: ./frpc-web"
echo ""
echo "Usage:"
echo "  ./frpc-web                        # listens on 127.0.0.1:7777"
echo "  ./frpc-web -listen 0.0.0.0:8080   # custom address"
echo ""
echo "Then open http://127.0.0.1:7777 in your browser."
echo "Default login: admin / admin"
echo "Data is stored in ./data/ next to the binary."
