#!/bin/bash
# Release helper script
# Usage: ./scripts/release.sh v1.0.0

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

VERSION="$1"

# Validate version format
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must match pattern v*.*.* (e.g., v1.0.0)"
    exit 1
fi

# Check if on main branch
BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$BRANCH" != "main" ]; then
    echo "Warning: Not on main branch (currently on $BRANCH)"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check for uncommitted changes
if [[ -n $(git status -s) ]]; then
    echo "Error: You have uncommitted changes"
    git status -s
    exit 1
fi

# Pull latest changes
echo "Pulling latest changes..."
git pull origin "$BRANCH"

# Run tests
echo "Running tests..."
go test -v ./...

# Build for all platforms to verify
echo "Building for all platforms..."
make build-all

# Clean up build artifacts
echo "Cleaning up..."
make clean

# Create and push tag
echo "Creating tag $VERSION..."
git tag "$VERSION"

echo "Pushing tag to GitHub..."
git push origin "$VERSION"

echo ""
echo "✓ Release tag $VERSION created and pushed!"
echo ""
echo "GitHub Actions will now:"
echo "  1. Run tests"
echo "  2. Build binaries for all platforms"
echo "  3. Create a GitHub release with artifacts"
echo ""
echo "Monitor progress at:"
echo "  https://github.com/pappasjfed/vibesortmedia/actions"
echo ""
echo "Once complete, check the release at:"
echo "  https://github.com/pappasjfed/vibesortmedia/releases/tag/$VERSION"
