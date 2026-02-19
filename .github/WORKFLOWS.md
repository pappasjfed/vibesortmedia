# GitHub Actions CI/CD

This repository uses GitHub Actions for continuous integration and automated releases.

## Workflows

### CI Workflow (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `main` branch
- Push to any `copilot/**` branch
- Pull requests to `main` branch

**Jobs:**
1. **Test**: Runs all tests with race detection and generates coverage report
2. **Build**: Builds binaries for all platforms and uploads as artifacts

**Artifacts:**
- CI builds are retained for 7 days
- Available in the "Actions" tab under each workflow run

### Release Workflow (`.github/workflows/release.yml`)

**Triggers:**
- Push of a tag matching `v*.*.*` pattern (e.g., v1.0.0, v2.1.3)

**Jobs:**
1. Run all tests
2. Build binaries for all platforms with version embedded
3. Generate SHA256 checksums
4. Create GitHub release with all artifacts

**Artifacts in Release:**
- `vibesortmedia-linux-amd64` - Linux x64 binary
- `vibesortmedia-darwin-amd64` - macOS Intel binary
- `vibesortmedia-darwin-arm64` - macOS Apple Silicon binary
- `vibesortmedia-windows-amd64.exe` - Windows x64 binary
- `checksums.txt` - SHA256 checksums for all binaries

## Creating a Release

### Step 1: Prepare the Release

Ensure all changes are committed and pushed to main:

```bash
git checkout main
git pull origin main
```

### Step 2: Create and Push a Version Tag

```bash
# Create a tag (replace with your version)
git tag v1.0.0

# Push the tag to GitHub
git push origin v1.0.0
```

### Step 3: Monitor the Release

1. Go to the "Actions" tab in GitHub
2. Watch the "Release" workflow execute
3. Once complete, check the "Releases" page for the new release

### Step 4: Verify the Release

Download and test a binary:

```bash
# Example: Download and test Linux binary
wget https://github.com/pappasjfed/vibesortmedia/releases/download/v1.0.0/vibesortmedia-linux-amd64
chmod +x vibesortmedia-linux-amd64
./vibesortmedia-linux-amd64 --version
```

## Version Embedding

The build process embeds the version into the binary using Go's `-ldflags`:

```bash
go build -ldflags="-X main.Version=v1.0.0"
```

Users can check the version:

```bash
./vibesortmedia --version
```

## Troubleshooting

### Workflow Fails

1. Check the workflow logs in the "Actions" tab
2. Common issues:
   - Test failures: Fix the failing tests before creating a release
   - Build failures: Ensure code compiles for all platforms
   - Permission issues: Ensure `GITHUB_TOKEN` has write permissions

### Tag Already Exists

If you need to recreate a tag:

```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin :refs/tags/v1.0.0

# Create new tag
git tag v1.0.0
git push origin v1.0.0
```

**Note:** Deleting a tag won't delete the GitHub release. Delete that manually if needed.

### Release Not Created

Ensure:
1. Tag matches the pattern `v*.*.*`
2. Workflow has `contents: write` permission
3. Check workflow logs for errors

## Local Testing

Test the CI workflow locally before pushing:

```bash
# Run tests like CI does
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Build like CI does
GOOS=linux GOARCH=amd64 go build -o vibesortmedia-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o vibesortmedia-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o vibesortmedia-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o vibesortmedia-windows-amd64.exe

# Or use make
make build-all
```

## Security

- Workflows use pinned versions of actions (e.g., `@v4`, `@v5`)
- No secrets are required for public repositories
- `GITHUB_TOKEN` is automatically provided and scoped to the repository
- Binary checksums are generated for verification

## Maintenance

### Updating Actions

Periodically update action versions:

```yaml
# Example: Update from v4 to v5
- uses: actions/checkout@v4  # Old
- uses: actions/checkout@v5  # New
```

### Changing Go Version

Update the Go version in both workflows:

```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.20'  # Change this
```
