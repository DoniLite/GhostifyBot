package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Ensure ffmpeg/ffprobe exist before running
func requireFFmpeg(t *testing.T) {
	if !IsFFmpegAvailable() {
		t.Skip("Skipping: FFmpeg not available on system")
	}
	cmd := exec.Command("ffprobe", "-version")
	if err := cmd.Run(); err != nil {
		t.Skip("Skipping: ffprobe not available on system")
	}
}

func TestIntegration_OptimizeProfiles(t *testing.T) {
	requireFFmpeg(t)

	// Use real fixtures under testdata/
	videoInput := filepath.Join("testdata", "sample.mp4")
	audioInput := filepath.Join("testdata", "sample.mp3")

	if _, err := os.Stat(videoInput); os.IsNotExist(err) {
		t.Fatalf("Missing fixture: %s", videoInput)
	}
	if _, err := os.Stat(audioInput); os.IsNotExist(err) {
		t.Fatalf("Missing fixture: %s", audioInput)
	}

	tempDir := t.TempDir()

	videoQualities := []string{"low", "high", "ultra"}
	for _, q := range videoQualities {
		t.Run(fmt.Sprintf("Video-%s", q), func(t *testing.T) {
			output := filepath.Join(tempDir, "output", fmt.Sprintf("video_%s.mp4", q))
			optimizer, err := NewMediaOptimizer(videoInput, output)
			if err != nil {
				t.Fatalf("Failed to init optimizer: %v", err)
			}
			if err := optimizer.OptimizeForMobile(q); err != nil {
				t.Fatalf("OptimizeForMobile failed: %v", err)
			}

			info, _ := getMediaInfo(output)
			if !strings.Contains(info, "libx264") {
				t.Errorf("Expected libx264 codec, got: %s", info)
			}
		})
	}

	audioQualities := []string{"low", "high"}
	for _, q := range audioQualities {
		t.Run(fmt.Sprintf("Audio-%s", q), func(t *testing.T) {
			output := filepath.Join(tempDir, fmt.Sprintf("audio_%s.aac", q))
			optimizer, err := NewMediaOptimizer(audioInput, output)
			if err != nil {
				t.Fatalf("Failed to init optimizer: %v", err)
			}
			if err := optimizer.OptimizeForMobile(q); err != nil {
				t.Fatalf("OptimizeForMobile failed: %v", err)
			}

			info, _ := getMediaInfo(output)
			if !strings.Contains(info, "aac") {
				t.Errorf("Expected aac codec, got: %s", info)
			}
		})
	}
}

func TestIntegration_CompressionRatio(t *testing.T) {
	requireFFmpeg(t)

	videoInput := filepath.Join("testdata", "sample.mp4")
	tempDir := t.TempDir()
	output := filepath.Join(tempDir, "compressed.mp4")

	optimizer, err := NewMediaOptimizer(videoInput, output)
	if err != nil {
		t.Fatalf("init error: %v", err)
	}

	if err := optimizer.OptimizeForMobile("low"); err != nil {
		t.Fatalf("optimize error: %v", err)
	}

	ratio, err := optimizer.GetCompressionRatio()
	if err != nil {
		t.Fatalf("GetCompressionRatio error: %v", err)
	}
	if ratio >= 1.0 {
		t.Errorf("Expected compressed output < input, got ratio %.2f", ratio)
	}
}

// Use ffprobe to inspect media info
func getMediaInfo(path string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_streams", "-of", "default=noprint_wrappers=1", path)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
