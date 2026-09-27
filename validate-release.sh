#!/bin/bash
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
VERSION_FILE="$PROJECT_DIR/.version"

echo "🧬 Launching Go Release Validation Gatekeeper..."

if [ ! -f "$VERSION_FILE" ]; then
  echo "1.0.0" > "$VERSION_FILE"
fi

CURRENT_VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
echo "📍 Current Version: $CURRENT_VERSION"

# Run tests
echo "🧪 Running Go Unit Tests..."
if [ -d "$PROJECT_DIR/services" ]; then
  # Monorepo setup: loop through all service directories containing a go.mod file
  for d in "$PROJECT_DIR"/services/*; do
    if [ -d "$d" ] && [ -f "$d/go.mod" ]; then
      echo "📦 Testing service: $(basename "$d")..."
      if ! (cd "$d" && go test ./...); then
        echo "❌ Error: Go unit tests failed in $(basename "$d")!"
        exit 1
      fi
    fi
  done
else
  if ! go test ./...; then
    echo "❌ Error: Go unit tests failed!"
    exit 1
  fi
fi
echo "✓ Go Unit Tests: PASS"

IFS='.' read -r major minor patch <<< "$CURRENT_VERSION"
NEXT_PATCH=$((patch + 1))
NEXT_VERSION="$major.$minor.$NEXT_PATCH"

echo "✨ Bumping version to: $NEXT_VERSION"
echo "$NEXT_VERSION" > "$VERSION_FILE"
echo "🚀 Validation Succeeded."
exit 0
