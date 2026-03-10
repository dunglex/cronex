#!/bin/bash
set -e

# Build the binary
echo "Building cronex..."
go build -o cronex .