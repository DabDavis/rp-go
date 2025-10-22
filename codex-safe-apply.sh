#!/usr/bin/env bash
# Codex-safe apply utility for rp-go
# Supports both unified diff patches and direct new file replacements
# Usage examples:
#   ./codex-safe-apply.sh fix.patch
#   ./codex-safe-apply.sh engine/systems/camera/system.go
#   ./codex-safe-apply.sh fix.patch engine/systems/debug/system.go

set -euo pipefail

if [[ $# -eq 0 ]]; then
    echo "Usage: $0 <file1> [file2 ...]"
    echo "Supports .patch/.diff or any file replacement."
    exit 1
fi

# Check working tree cleanliness
if [[ -n "$(git status --porcelain)" ]]; then
    echo "⚠️ Working tree not clean. Commit or stash your work first."
    exit 1
fi

# Create isolated temp branch
BRANCH="codex-apply-$(date +%Y%m%d-%H%M%S)"
echo "🌿 Creating branch: $BRANCH"
git checkout -b "$BRANCH" >/dev/null

# Apply all provided files
for FILE in "$@"; do
    if [[ ! -f "$FILE" ]]; then
        echo "❌ File not found: $FILE"
        git checkout main >/dev/null
        git branch -D "$BRANCH" >/dev/null 2>&1 || true
        exit 1
    fi

    if [[ "$FILE" == *.patch || "$FILE" == *.diff ]]; then
        echo "📦 Applying diff patch: $FILE"
        if ! git apply --3way "$FILE"; then
            echo "❌ Patch failed to apply cleanly."
            git merge --abort 2>/dev/null || true
            git restore .
            git checkout main >/dev/null
            git branch -D "$BRANCH" >/dev/null 2>&1 || true
            exit 1
        fi
    else
        echo "🆕 Applying full file replacement: $FILE"
        DEST="$FILE"
        mkdir -p "$(dirname "$DEST")"
        cp -f "$FILE" "$DEST"
        git add "$DEST"
    fi
done

echo "🧹 Running go fmt, goimports, and go vet..."
go fmt ./...
if command -v goimports >/dev/null 2>&1; then
    goimports -w .
fi
go vet ./... || true

echo "🔍 Building project..."
if ! go build ./...; then
    echo "❌ Build failed after applying Codex changes."
    git restore .
    git checkout main >/dev/null
    git branch -D "$BRANCH" >/dev/null 2>&1 || true
    exit 1
fi

echo "✅ Build succeeded — committing patch."
git add .
git commit -m "Apply Codex changes ($(date +%Y-%m-%d\ %H:%M))"

echo "🔁 Merging back into main..."
git checkout main >/dev/null
git merge --no-ff "$BRANCH" -m "Merge Codex changes ($(date +%Y-%m-%d))"

echo "🧹 Cleaning up temp branch..."
git branch -d "$BRANCH" >/dev/null

echo "🎉 Codex changes applied safely and merged!"

