.PHONY: build build-all clean run

# Build for current platform
build:
	go build -o vibesortmedia

# Build for all platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o vibesortmedia-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o vibesortmedia-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o vibesortmedia-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -o vibesortmedia-windows-amd64.exe

# Clean build artifacts
clean:
	rm -f vibesortmedia vibesortmedia-* *.exe

# Run the program
run: build
	./vibesortmedia
