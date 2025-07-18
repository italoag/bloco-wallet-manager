#!/bin/bash

# Navigate to the testdata directory
cd "$(dirname "$0")"

# Run the generator script
echo "Generating additional keystore test files..."
go run generate_additional_test_files.go

echo "Done!"