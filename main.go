package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Version is set during build via -ldflags
var Version = "dev"

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

type MediaType string

const (
	MediaDVD     MediaType = "DVD"
	MediaBluRay  MediaType = "Blu-ray"
	MediaCD      MediaType = "CD"
	MediaUnknown MediaType = "Unknown"
	MediaEmpty   MediaType = "Empty"
)

const (
	// Size thresholds for media type identification (in bytes)
	BluRayMinSize = 20000000000 // 20GB - minimum for Blu-ray discs
	DVDMinSize    = 1000000000  // 1GB - minimum for DVD discs
	CDMinSize     = 1           // Any size above 0 for CDs
)

func classifyMediaBySize(size int64) MediaType {
	if size > BluRayMinSize {
		return MediaBluRay
	} else if size > DVDMinSize {
		return MediaDVD
	} else if size > CDMinSize {
		return MediaCD
	}
	return MediaUnknown
}

func getColorForMediaType(mediaType MediaType) string {
	switch mediaType {
	case MediaDVD:
		return ColorBlue
	case MediaBluRay:
		return ColorPurple
	case MediaCD:
		return ColorYellow
	case MediaUnknown:
		return ColorRed
	case MediaEmpty:
		return ColorWhite
	default:
		return ColorReset
	}
}

func detectOpticalDrives() []string {
	var drives []string

	switch runtime.GOOS {
	case "windows":
		drives = detectWindowsDrives()
	case "darwin":
		drives = detectMacOSDrives()
	case "linux":
		drives = detectLinuxDrives()
	}

	return drives
}

func detectWindowsDrives() []string {
	var drives []string
	cmd := exec.Command("wmic", "logicaldisk", "where", "drivetype=5", "get", "deviceid")
	output, err := cmd.Output()
	if err != nil {
		return drives
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && line != "DeviceID" {
			drives = append(drives, line)
		}
	}

	return drives
}

func detectMacOSDrives() []string {
	var drives []string
	cmd := exec.Command("diskutil", "list")
	output, err := cmd.Output()
	if err != nil {
		return drives
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "(external") && (strings.Contains(strings.ToLower(line), "optical") || strings.Contains(line, "disk")) {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				diskName := fields[0]
				if strings.HasPrefix(diskName, "/dev/disk") {
					drives = append(drives, diskName)
				}
			}
		}
	}

	if len(drives) == 0 {
		cmd = exec.Command("drutil", "status")
		output, err = cmd.Output()
		if err == nil && len(output) > 0 {
			cmd = exec.Command("diskutil", "list")
			output, _ = cmd.Output()
			lines = strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "/dev/disk") {
					fields := strings.Fields(line)
					for _, field := range fields {
						if strings.HasPrefix(field, "/dev/disk") {
							drives = append(drives, field)
							break
						}
					}
				}
			}
		}
	}

	return drives
}

func detectLinuxDrives() []string {
	var drives []string

	entries, err := os.ReadDir("/dev")
	if err != nil {
		return drives
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "sr") || strings.HasPrefix(name, "cdrom") {
			drives = append(drives, "/dev/"+name)
		}
	}

	return drives
}

func ejectDrive(drive string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		powershellScript := fmt.Sprintf(`(New-Object -com Shell.Application).Namespace(17).ParseName('%s').InvokeVerb('Eject')`, drive)
		cmd = exec.Command("powershell", "-Command", powershellScript)
	case "darwin":
		cmd = exec.Command("drutil", "eject")
	case "linux":
		cmd = exec.Command("eject", drive)
	}

	return cmd.Run()
}

func isDriveEmpty(drive string) bool {
	switch runtime.GOOS {
	case "windows":
		return isWindowsDriveEmpty(drive)
	case "darwin":
		return isMacOSDriveEmpty(drive)
	case "linux":
		return isLinuxDriveEmpty(drive)
	}
	return true
}

func isWindowsDriveEmpty(drive string) bool {
	cmd := exec.Command("cmd", "/C", "dir", drive)
	err := cmd.Run()
	return err != nil
}

func isMacOSDriveEmpty(drive string) bool {
	cmd := exec.Command("diskutil", "info", drive)
	output, err := cmd.Output()
	if err != nil {
		return true
	}

	return !strings.Contains(string(output), "Mounted")
}

func isLinuxDriveEmpty(drive string) bool {
	cmd := exec.Command("blkid", drive)
	output, _ := cmd.Output()
	return len(output) == 0
}

func identifyMediaType(drive string) MediaType {
	switch runtime.GOOS {
	case "windows":
		return identifyWindowsMediaType(drive)
	case "darwin":
		return identifyMacOSMediaType(drive)
	case "linux":
		return identifyLinuxMediaType(drive)
	}
	return MediaUnknown
}

func identifyWindowsMediaType(drive string) MediaType {
	cmd := exec.Command("wmic", "logicaldisk", "where", fmt.Sprintf("deviceid='%s'", drive), "get", "volumename,size")
	output, err := cmd.Output()
	if err != nil {
		return MediaUnknown
	}

	outputStr := strings.ToLower(string(output))
	if strings.Contains(outputStr, "bd") || strings.Contains(outputStr, "bluray") {
		return MediaBluRay
	} else if strings.Contains(outputStr, "dvd") {
		return MediaDVD
	} else if strings.Contains(outputStr, "cd") {
		return MediaCD
	}

	cmd = exec.Command("wmic", "logicaldisk", "where", fmt.Sprintf("deviceid='%s'", drive), "get", "size")
	output, err = cmd.Output()
	if err != nil {
		return MediaUnknown
	}

	sizeStr := strings.TrimSpace(string(output))
	lines := strings.Split(sizeStr, "\n")
	if len(lines) > 1 {
		sizeStr = strings.TrimSpace(lines[1])
	}

	if sizeStr != "" && sizeStr != "0" {
		var size int64
		fmt.Sscanf(sizeStr, "%d", &size)
		return classifyMediaBySize(size)
	}

	return MediaUnknown
}

func identifyMacOSMediaType(drive string) MediaType {
	cmd := exec.Command("diskutil", "info", drive)
	output, err := cmd.Output()
	if err != nil {
		return MediaUnknown
	}

	outputStr := strings.ToLower(string(output))
	if strings.Contains(outputStr, "bluray") || strings.Contains(outputStr, "bd") {
		return MediaBluRay
	} else if strings.Contains(outputStr, "dvd") {
		return MediaDVD
	} else if strings.Contains(outputStr, "cd") {
		return MediaCD
	}

	return MediaUnknown
}

func identifyLinuxMediaType(drive string) MediaType {
	cmd := exec.Command("blkid", "-o", "value", "-s", "TYPE", drive)
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		fsType := strings.TrimSpace(string(output))
		if fsType == "udf" || fsType == "iso9660" {
			cmd = exec.Command("isoinfo", "-d", "-i", drive)
			output, err = cmd.Output()
			if err == nil {
				outputStr := strings.ToLower(string(output))
				if strings.Contains(outputStr, "bluray") || strings.Contains(outputStr, "bd") {
					return MediaBluRay
				} else if strings.Contains(outputStr, "dvd") {
					return MediaDVD
				}
			}

			cmd = exec.Command("blockdev", "--getsize64", drive)
			output, err = cmd.Output()
			if err == nil {
				var size int64
				fmt.Sscanf(string(output), "%d", &size)
				return classifyMediaBySize(size)
			}

			return MediaDVD
		}
	}

	cmd = exec.Command("udevadm", "info", "--query=property", "--name="+drive)
	output, err = cmd.Output()
	if err == nil {
		outputStr := strings.ToLower(string(output))
		if strings.Contains(outputStr, "id_cdrom_media_bd") {
			return MediaBluRay
		} else if strings.Contains(outputStr, "id_cdrom_media_dvd") {
			return MediaDVD
		} else if strings.Contains(outputStr, "id_cdrom_media_cd") {
			return MediaCD
		}
	}

	return MediaUnknown
}

func printColored(text string, color string) {
	fmt.Printf("%s%s%s\n", color, text, ColorReset)
}

func waitForUserInput() {
	fmt.Print("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.BoolVar(versionFlag, "v", false, "Print version information (shorthand)")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("VibeSort Media version %s\n", Version)
		fmt.Printf("Built with %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	fmt.Printf("VibeSort Media v%s - Physical Media Detection\n", Version)
	fmt.Println("==========================================")
	fmt.Println()

	for {
		drives := detectOpticalDrives()

		if len(drives) == 0 {
			printColored("No optical drives detected on this system.", ColorRed)
			fmt.Println()
			time.Sleep(5 * time.Second)
			continue
		}

		for _, drive := range drives {
			fmt.Printf("Checking drive: %s\n", drive)

			if isDriveEmpty(drive) {
				printColored("Drive is empty. Ejecting...", ColorWhite)
				err := ejectDrive(drive)
				if err != nil {
					printColored(fmt.Sprintf("Error ejecting drive: %v", err), ColorRed)
				} else {
					printColored("Please insert media into the drive.", ColorCyan)
				}
			} else {
				printColored("Media detected! Identifying...", ColorGreen)
				mediaType := identifyMediaType(drive)

				color := getColorForMediaType(mediaType)
				printColored(fmt.Sprintf("Media Type: %s", mediaType), color)

				time.Sleep(2 * time.Second)

				printColored("Ejecting media...", ColorWhite)
				err := ejectDrive(drive)
				if err != nil {
					printColored(fmt.Sprintf("Error ejecting drive: %v", err), ColorRed)
				}
			}

			fmt.Println()
		}

		time.Sleep(3 * time.Second)
	}
}
