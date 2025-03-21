#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Remove any old coverage.html file
rm -f coverage.html

# Run tests and generate the HTML coverage report directly
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
rm coverage.out  # Remove the coverage.out file after generating the HTML report

# Mark the file as generated by adding a comment at the beginning
echo "<!-- Generated by test coverage script -->" | cat - coverage.html > temp.html && mv temp.html coverage.html

# Open the HTML report in the default browser based on the OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    xdg-open coverage.html  # Linux
elif [[ "$OSTYPE" == "darwin"* ]]; then
    open coverage.html  # macOS
elif [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    start coverage.html  # Windows
else
    echo "Unsupported OS. Please open the coverage.html file manually."
fi

# Output success message
echo "Code coverage report generated: coverage.html"

# Make the file read-only
chmod 444 coverage.html  # Read-only for all users
