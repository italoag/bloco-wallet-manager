#!/bin/bash

# Build and run BlocoWallet

echo "Building BlocoWallet..."
go build -o build/blocowallet ./cmd/blocowallet

if [ $? -eq 0 ]; then
    echo "Build successful! Running BlocoWallet..."
    ./build/blocowallet
else
    echo "Build failed!"
    exit 1
fi
