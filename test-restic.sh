#!/bin/bash

# Test script to verify restic repository access
# Usage: ./test-restic.sh <repo-path> <password>

set -e

if [ $# -ne 2 ]; then
    echo "Usage: $0 <repo-path> <password>"
    echo "Example: $0 /path/to/repo mypassword"
    echo "Example: $0 sftp:user@host:/path/to/repo mypassword"
    exit 1
fi

REPO_PATH="$1"
PASSWORD="$2"

echo "Testing restic repository access..."
echo "Repository: $REPO_PATH"
echo "----------------------------------------"

# Set password environment variable
export RESTIC_PASSWORD="$PASSWORD"

echo "1. Testing repository connectivity..."
if restic -r "$REPO_PATH" cat config > /dev/null 2>&1; then
    echo "✅ Repository is accessible"
else
    echo "❌ Repository is not accessible"
    echo "Error details:"
    restic -r "$REPO_PATH" cat config
    exit 1
fi

echo ""
echo "2. Testing password..."
if restic -r "$REPO_PATH" check --read-data-subset=1% > /dev/null 2>&1; then
    echo "✅ Password is correct"
else
    echo "❌ Password might be incorrect or repository is corrupted"
    echo "Error details:"
    restic -r "$REPO_PATH" check --read-data-subset=1%
    exit 1
fi

echo ""
echo "3. Testing snapshots command..."
echo "Available snapshots:"
if restic -r "$REPO_PATH" snapshots; then
    echo "✅ Snapshots command works"
else
    echo "❌ Snapshots command failed"
    exit 1
fi

echo ""
echo "4. Testing JSON output..."
if restic -r "$REPO_PATH" snapshots --json > /dev/null; then
    echo "✅ JSON output works"
else
    echo "❌ JSON output failed"
    exit 1
fi

echo ""
echo "5. Testing repository info..."
echo "Repository info:"
restic -r "$REPO_PATH" stats

echo ""
echo "----------------------------------------"
echo "✅ All tests passed! Repository is ready for use with restic-browser."
echo ""
echo "You can now run:"
echo "./restic-browser -repo \"$REPO_PATH\" -password \"$PASSWORD\""
