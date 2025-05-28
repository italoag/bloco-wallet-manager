#!/bin/bash

# Build and run BlocoWallet

echo "Building BlocoWallet..."
go build -o bin/blocowallet ./cmd/blocowallet

if [ $? -eq 0 ]; then
    echo "Build successful! Running BlocoWallet..."
    ./bin/blocowallet
else
    echo "Build failed!"
    exit 1
fi
