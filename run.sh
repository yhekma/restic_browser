#!/bin/bash

# Restic Browser - Build and Run Script

set -e

echo "Building restic-browser..."
go build -o restic-browser main.go

echo "Build complete!"
echo ""
echo "Usage: ./restic-browser -repo <repository-path> -password <password> [-port <port>]"
echo ""
echo "Example:"
echo "  ./restic-browser -repo /path/to/repo -password mypassword"
echo "  ./restic-browser -repo /path/to/repo -password mypassword -port 9000"
echo ""
echo "The service will be available at http://localhost:8081 (or your specified port)"
echo ""

# Check if arguments were provided
if [ $# -eq 0 ]; then
    echo "No arguments provided. Please run with -repo and -password flags."
    echo "Run './restic-browser -h' for help."
else
    echo "Starting restic-browser with provided arguments..."
    ./restic-browser "$@"
fi
