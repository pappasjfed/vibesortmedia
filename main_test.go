package main

import (
	"testing"
)

func TestGetColorForMediaType(t *testing.T) {
	tests := []struct {
		mediaType MediaType
		expected  string
	}{
		{MediaDVD, ColorBlue},
		{MediaBluRay, ColorPurple},
		{MediaCD, ColorYellow},
		{MediaUnknown, ColorRed},
		{MediaEmpty, ColorWhite},
	}

	for _, tt := range tests {
		t.Run(string(tt.mediaType), func(t *testing.T) {
			result := getColorForMediaType(tt.mediaType)
			if result != tt.expected {
				t.Errorf("getColorForMediaType(%s) = %s; want %s", tt.mediaType, result, tt.expected)
			}
		})
	}
}

func TestDetectOpticalDrives(t *testing.T) {
	drives := detectOpticalDrives()
	t.Logf("Detected %d optical drive(s)", len(drives))
	for i, drive := range drives {
		t.Logf("Drive %d: %s", i+1, drive)
	}
}

func TestMediaTypeConstants(t *testing.T) {
	if MediaDVD != "DVD" {
		t.Errorf("MediaDVD constant is incorrect")
	}
	if MediaBluRay != "Blu-ray" {
		t.Errorf("MediaBluRay constant is incorrect")
	}
	if MediaCD != "CD" {
		t.Errorf("MediaCD constant is incorrect")
	}
	if MediaUnknown != "Unknown" {
		t.Errorf("MediaUnknown constant is incorrect")
	}
	if MediaEmpty != "Empty" {
		t.Errorf("MediaEmpty constant is incorrect")
	}
}
