# VibeSort Media

A cross-platform physical media detection and identification tool that automatically detects DVD, Blu-ray, and CD media.

## Features

- **Cross-platform support**: Works on Windows, macOS, and Linux
- **Automatic detection**: Continuously monitors optical drives for media
- **Media identification**: Identifies DVD, Blu-ray, and CD media types
- **Color-coded output**: Different media types are displayed in different colors:
  - 🔵 DVD (Blue)
  - 🟣 Blu-ray (Purple)
  - 🟡 CD (Yellow)
  - ⚪ Empty drive (White)
  - 🔴 Unknown (Red)
- **Automatic ejection**: Ejects empty drives and prompts for media insertion
- **Single executable**: Easy to deploy and run

## Installation

### Prerequisites
- Go 1.20 or later (for building from source)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/pappasjfed/vibesortmedia.git
cd vibesortmedia

# Build for your current platform
make build

# Or build for all platforms
make build-all
```

### Pre-built Binaries

Download the appropriate binary for your platform from the [releases page](https://github.com/pappasjfed/vibesortmedia/releases).

Each release includes pre-compiled binaries for:
- Linux (x64)
- macOS (Intel and Apple Silicon)
- Windows (x64)

## Usage

Check the version:

```bash
./vibesortmedia --version
# or
./vibesortmedia -v
```

Simply run the executable:

```bash
# Linux/macOS
./vibesortmedia

# Windows
vibesortmedia.exe
```

The program will:
1. Detect all optical drives on your system
2. Check if drives are empty
3. Eject empty drives and prompt for media insertion
4. Identify and display media type when media is inserted
5. Eject the media after identification
6. Loop continuously

Press `Ctrl+C` to exit the program.

## Platform-Specific Notes

### Linux
- May require `eject`, `blkid`, `udevadm`, and `blockdev` commands to be available
- Some operations may require sudo privileges

### macOS
- Uses `drutil` and `diskutil` for drive operations
- Works best with built-in optical drives

### Windows
- Uses PowerShell for eject operations
- Requires Windows PowerShell to be available

## Development

### Running Tests

```bash
# Run tests
go test -v ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Building

```bash
# Run the program
make run

# Build for current platform
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

### Continuous Integration

The project uses GitHub Actions for CI/CD:
- **CI Workflow**: Runs tests and builds on every push and pull request
- **Release Workflow**: Automatically creates releases when version tags are pushed

See [RELEASE.md](RELEASE.md) for information on creating releases.

## License

MIT License
