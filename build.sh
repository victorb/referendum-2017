#! /bin/sh
set -e 
COMMIT=$(git rev-parse HEAD)
VERSION=${1:-$COMMIT}
echo "## Building version '$VERSION'"
# export CGO_ENABLED=0
echo "## Building OSX version..."
GOOS=darwin time go build -o "dist/referendum-2017-darwin-$VERSION" main.go
echo "## Building Windows version"
GOOS=darwin time go build -o "dist/referendum-2017-windows-$VERSION" main.go
