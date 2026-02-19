# Release Guide

This document explains how to create releases for VibeSort Media.

## Automatic Releases

Releases are automatically created via GitHub Actions when you push a version tag.

### Creating a Release

#### Using the Helper Script (Recommended)

```bash
./scripts/release.sh v1.0.0
```

The script will:
- Validate the version format
- Check for uncommitted changes
- Run tests
- Build for all platforms to verify
- Create and push the tag

#### Manual Release

1. **Ensure all changes are committed and pushed to main**

2. **Create and push a version tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **GitHub Actions will automatically:**
   - Run tests
   - Build binaries for all platforms (Linux, macOS x64/ARM64, Windows)
   - Generate SHA256 checksums
   - Create a GitHub release with all artifacts

### Version Numbering

We use [Semantic Versioning](https://semver.org/):

- **MAJOR** version (v1.0.0 → v2.0.0): Incompatible API changes
- **MINOR** version (v1.0.0 → v1.1.0): New functionality in a backward compatible manner
- **PATCH** version (v1.0.0 → v1.0.1): Backward compatible bug fixes

### Tag Format

Tags must follow the pattern `v*.*.*` (e.g., v1.0.0, v2.1.3)

Examples:
```bash
# First release
git tag v1.0.0
git push origin v1.0.0

# Bug fix release
git tag v1.0.1
git push origin v1.0.1

# New feature release
git tag v1.1.0
git push origin v1.1.0

# Breaking change release
git tag v2.0.0
git push origin v2.0.0
```

## Manual Build

If you need to build locally:

```bash
# Build for all platforms
make build-all

# Build with version information
go build -ldflags="-X main.Version=v1.0.0" -o vibesortmedia
```

## Continuous Integration

The CI workflow runs on every push and pull request to:
- Run all tests with race detection
- Build for all platforms
- Upload artifacts for testing

## Release Artifacts

Each release includes:
- `vibesortmedia-linux-amd64` - Linux x64 binary
- `vibesortmedia-darwin-amd64` - macOS Intel binary
- `vibesortmedia-darwin-arm64` - macOS Apple Silicon binary
- `vibesortmedia-windows-amd64.exe` - Windows x64 binary
- `checksums.txt` - SHA256 checksums for verification

## Deleting a Tag

If you need to delete a tag:

```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin :refs/tags/v1.0.0
```

Note: This will not delete the GitHub release. You'll need to delete that manually from the GitHub UI.
