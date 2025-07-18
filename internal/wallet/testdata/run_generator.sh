#!/bin/bash

# Navigate to the testdata directory
cd "$(dirname "$0")"

# Run the generator script
echo "Generating keystore test files..."
go run generate_simple_keystores.go

echo "Done!"