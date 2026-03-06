#!/usr/bin/env bash
set -euo pipefail

BINARY="frpc"
OUTDIR="./dist"

echo "==> Checking Go..."
if ! command -v go &>/dev/null; then
  echo "ERROR: Go not found. Install Go >= 1.21 from https://go.dev/dl/"
  exit 1
fi
echo "    $(go version)"

echo "==> Tidying modules..."
go mod tidy

echo "==> Building $BINARY with embedded panel..."
mkdir -p "$OUTDIR"

GOOS=${GOOS:-$(go env GOOS)}
GOARCH=${GOARCH:-$(go env GOARCH)}
OUT="$OUTDIR/${BINARY}-${GOOS}-${GOARCH}"
[ "$GOOS" = "windows" ] && OUT="${OUT}.exe"

go build -ldflags "-s -w" -o "$OUT" ./cmd/frpc/

echo ""
echo "✓  Built: $OUT"
echo ""
echo "Usage:"
echo "  # frpc.toml must have webServer configured:"
echo "  #   [webServer]"
echo "  #   addr = \"127.0.0.1\""
echo "  #   port = 7400"
echo "  #   user = \"admin\""
echo "  #   password = \"admin\""
echo ""
echo "  ./$OUT -c frpc.toml"
echo ""
echo "  Then open: http://127.0.0.1:7400/panel/"
echo "  Default panel login: admin / admin  (stored in frpc-panel-auth.json)"
