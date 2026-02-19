# Example Usage

## Running the Program

After building the program, simply run it:

```bash
# Linux/macOS
./vibesortmedia

# Windows
vibesortmedia.exe
```

## Example Output

When no optical drives are present:
```
VibeSort Media - Physical Media Detection
==========================================

No optical drives detected on this system.
```

When an optical drive is detected but empty:
```
VibeSort Media - Physical Media Detection
==========================================

Checking drive: D:
Drive is empty. Ejecting...
Please insert media into the drive.
```

When media is detected and identified:
```
VibeSort Media - Physical Media Detection
==========================================

Checking drive: D:
Media detected! Identifying...
Media Type: DVD        [displayed in BLUE]
Ejecting media...

Checking drive: D:
Media detected! Identifying...
Media Type: Blu-ray    [displayed in PURPLE]
Ejecting media...

Checking drive: D:
Media detected! Identifying...
Media Type: CD         [displayed in YELLOW]
Ejecting media...
```

## Color Coding

The program uses ANSI color codes to display different media types:
- **DVD**: Blue (🔵)
- **Blu-ray**: Purple (🟣)
- **CD**: Yellow (🟡)
- **Empty**: White (⚪)
- **Unknown**: Red (🔴)

## Stopping the Program

Press `Ctrl+C` to stop the program at any time.

## Troubleshooting

### Linux
If you get permission errors, you may need to run with sudo:
```bash
sudo ./vibesortmedia
```

### macOS
If the eject command fails, ensure you have necessary permissions and the drive is not in use by other applications.

### Windows
If PowerShell commands fail, ensure PowerShell is available in your system PATH and you have the necessary permissions.
