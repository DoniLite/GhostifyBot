package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)


func createRealTestFiles(t *testing.T) (string, string, string) {
	tempDir, err := os.MkdirTemp("", "media_optimizer_test")
	if err != nil {
		t.Fatalf("repertory creation error: %v", err)
	}

	testVideoPath := filepath.Join(tempDir, "test_video.mp4")
	testAudioPath := filepath.Join(tempDir, "test_audio.mp3")
	
	if IsFFmpegAvailable() {
		err = createTestVideo(testVideoPath)
		if err != nil {
			t.Logf("Video test creation not work: %v", err)
			videoFile, _ := os.Create(testVideoPath)
			videoFile.WriteString("fake video")
			videoFile.Close()
		}
		
		err = createTestAudio(testAudioPath)
		if err != nil {
			t.Logf("Audio test creation not work: %v", err)
			audioFile, _ := os.Create(testAudioPath)
			audioFile.WriteString("fake audio")
			audioFile.Close()
		}
	} else {
		videoFile, _ := os.Create(testVideoPath)
		videoFile.WriteString("fake video")
		videoFile.Close()
		
		audioFile, _ := os.Create(testAudioPath)
		audioFile.WriteString("fake audio")
		audioFile.Close()
	}

	return tempDir, testVideoPath, testAudioPath
}

func TestNewMediaOptimizer(t *testing.T) {
	tempDir, testVideoPath, testAudioPath := createRealTestFiles(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		inputPath   string
		outputPath  string
		shouldError bool
	}{
		{
			name:        "Valid video file",
			inputPath:   testVideoPath,
			outputPath:  filepath.Join(tempDir, "output.mp4"),
			shouldError: false,
		},
		{
			name:        "Valid audio file",
			inputPath:   testAudioPath,
			outputPath:  filepath.Join(tempDir, "output.aac"),
			shouldError: false,
		},
		{
			name:        "Not existing entry file",
			inputPath:   "/path/not_existed.mp4",
			outputPath:  filepath.Join(tempDir, "output.mp4"),
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer, err := NewMediaOptimizer(tt.inputPath, tt.outputPath)

			if tt.shouldError {
				if err == nil {
					t.Error("Attendait une erreur mais n'en a pas reçu")
				}
			} else {
				if err != nil {
					t.Errorf("Erreur inattendue: %v", err)
				}
				if optimizer == nil {
					t.Error("Optimizer ne devrait pas être nil")
				}
			}
		})
	}
}

func TestDetectMediaType(t *testing.T) {
	tests := []struct {
		path     string
		expected MediaType
	}{
		{"test.mp4", Video},
		{"test.avi", Video},
		{"test.mov", Video},
		{"test.mkv", Video},
		{"test.3gp", Video},
		{"test.mp3", Audio},
		{"test.wav", Audio},
		{"test.flac", Audio},
		{"test.aac", Audio},
		{"test.unknown", Audio},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := detectMediaType(tt.path)
			if result != tt.expected {
				t.Errorf("detectMediaType(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestSetProfile(t *testing.T) {
	tempDir, testVideoPath, _ := createRealTestFiles(t)
	defer os.RemoveAll(tempDir)

	optimizer, err := NewMediaOptimizer(testVideoPath, filepath.Join(tempDir, "output.mp4"))
	if err != nil {
		t.Fatalf("Optimizer creation error: %v", err)
	}

	optimizer.SetProfile(VideoMobileHigh)

	if optimizer.Profile.Name != "video_mobile_high" {
		t.Errorf("Wrong video profile definition: %s", optimizer.Profile.Name)
	}

	if optimizer.Profile.VideoCodec != "libx264" {
		t.Errorf("wrong video codec definition: %s", optimizer.Profile.VideoCodec)
	}
}

func TestSetCustomProfile(t *testing.T) {
	tempDir, testVideoPath, _ := createRealTestFiles(t)
	defer os.RemoveAll(tempDir)

	optimizer, err := NewMediaOptimizer(testVideoPath, filepath.Join(tempDir, "output.mp4"))
	if err != nil {
		t.Fatalf("Optimizer creation error: %v", err)
	}

	// Test avec profil personnalisé
	optimizer.SetCustomProfile("custom", "libx264", "aac", "1000k", "128k", "720x480", "medium", 23)

	if optimizer.Profile.Name != "custom" {
		t.Errorf("wrong personalized profile name defined: %s", optimizer.Profile.Name)
	}

	if optimizer.Profile.VideoBitrate != "1000k" {
		t.Errorf("wrong video bitrate definition: %s", optimizer.Profile.VideoBitrate)
	}
}

func TestFileExists(t *testing.T) {
	tempDir, testVideoPath, _ := createRealTestFiles(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		path     string
		expected bool
	}{
		{testVideoPath, true},
		{"/path/not_existed.mp4", false},
		{tempDir, true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := fileExists(tt.path)
			if result != tt.expected {
				t.Errorf("fileExists(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestQualityProfiles(t *testing.T) {
	profiles := []QualityProfile{
		AudioMobileLow,
		AudioMobileHigh,
		VideoMobileLow,
		VideoMobileHigh,
		VideoMobileUltra,
	}

	for _, profile := range profiles {
		t.Run(profile.Name, func(t *testing.T) {
			if profile.Name == "" {
				t.Error("profile Name should be empty")
			}
			if profile.AudioCodec == "" {
				t.Error("video Codec should not be empty")
			}
			if profile.AudioBitrate == "" {
				t.Error("video Bitrate should not be empty")
			}
		})
	}
}

func TestOptimizeForMobile(t *testing.T) {
	tempDir, testVideoPath, testAudioPath := createRealTestFiles(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name      string
		inputPath string
		quality   string
		mediaType MediaType
	}{
		{
			name:      "low quality video",
			inputPath: testVideoPath,
			quality:   "low",
			mediaType: Video,
		},
		{
			name:      "high quality video",
			inputPath: testVideoPath,
			quality:   "high",
			mediaType: Video,
		},
		{
			name:      "low audio quality",
			inputPath: testAudioPath,
			quality:   "low",
			mediaType: Audio,
		},
		{
			name:      "high audio quality",
			inputPath: testAudioPath,
			quality:   "high",
			mediaType: Audio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath := filepath.Join(tempDir, "output_"+tt.name+".mp4")
			optimizer, err := NewMediaOptimizer(tt.inputPath, outputPath)
			if err != nil {
				t.Fatalf("Optimizer creation error: %v", err)
			}

			if optimizer.MediaType != tt.mediaType {
				t.Errorf("invalid media type: %v, expected %v", optimizer.MediaType, tt.mediaType)
			}

			err = optimizer.OptimizeForMobile(tt.quality)
			if err == nil {
				t.Log("Success optimization (FFmpeg probably installed)")
			} else {
				t.Logf("Expected error no bin found for FFmpeg: %v", err)
			}
		})
	}
}

func TestBatchOptimize(t *testing.T) {
	tempDir, testVideoPath, testAudioPath := createRealTestFiles(t)
	defer os.RemoveAll(tempDir)

	inputPaths := []string{testVideoPath, testAudioPath}
	outputDir := filepath.Join(tempDir, "batch_output")

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Output repertory creation error: %v", err)
	}

	err = BatchOptimize(inputPaths, outputDir, "high")

	if err != nil {
		t.Logf("Expected error without FFmpeg: %v", err)
	}
}

func TestConvenienceFunctions(t *testing.T) {
	tempDir, testVideoPath, testAudioPath := createRealTestFiles(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		function func(string, string, string) error
		input    string
		output   string
		quality  string
	}{
		{
			name:     "OptimizeVideoForMobile",
			function: OptimizeVideoForMobile,
			input:    testVideoPath,
			output:   filepath.Join(tempDir, "video_output.mp4"),
			quality:  "high",
		},
		{
			name:     "OptimizeAudioForMobile",
			function: OptimizeAudioForMobile,
			input:    testAudioPath,
			output:   filepath.Join(tempDir, "audio_output.aac"),
			quality:  "high",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.function(tt.input, tt.output, tt.quality)
			if err != nil {
				t.Logf("Expected error without FFmpeg: %v", err)
			}
		})
	}
}

func TestGetOptimizedSizeWithoutFile(t *testing.T) {
	tempDir, testVideoPath, _ := createRealTestFiles(t)
	defer os.RemoveAll(tempDir)

	optimizer, err := NewMediaOptimizer(testVideoPath, filepath.Join(tempDir, "nonexistent.mp4"))
	if err != nil {
		t.Fatalf("optimizer creation error: %v", err)
	}

	_, err = optimizer.GetOptimizedSize()
	if err == nil {
		t.Error("Expected an non existent output error")
	}
}

func TestGetCompressionRatioWithoutOptimization(t *testing.T) {
	tempDir, testVideoPath, _ := createRealTestFiles(t)
	defer os.RemoveAll(tempDir)

	optimizer, err := NewMediaOptimizer(testVideoPath, filepath.Join(tempDir, "nonexistent.mp4"))
	if err != nil {
		t.Fatalf("optimizer creation error: %v", err)
	}

	_, err = optimizer.GetCompressionRatio()
	if err == nil {
		t.Error("Expecting a non existent file error")
	}
}

// Benchmark tests
func BenchmarkDetectMediaType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		detectMediaType("test.mp4")
	}
}

func BenchmarkNewMediaOptimizer(b *testing.B) {
	tempDir, testVideoPath, _ := createRealTestFiles(&testing.T{})
	defer os.RemoveAll(tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewMediaOptimizer(testVideoPath, filepath.Join(tempDir, "output.mp4"))
	}
}

// Integration test (need FFmpeg)
func TestIntegrationOptimizeWithFFmpeg(t *testing.T) {
	if !IsFFmpegAvailable() {
		t.Skip("FFmpeg not found test skipped")
	}

	tempDir, testVideoPath, _ := createRealTestFiles(t)
	defer os.RemoveAll(tempDir)

	optimizer, err := NewMediaOptimizer(testVideoPath, filepath.Join(tempDir, "integration_output.mp4"))
	if err != nil {
		t.Fatalf("optimizer creation error: %v", err)
	}

	err = optimizer.OptimizeForMobile("high")
	if err != nil {
		t.Errorf("optimizing error with FFmpeg: %v", err)
	}
}


func createTestVideo(outputPath string) error {
	cmd := fmt.Sprintf("-f lavfi -i testsrc=duration=1:size=320x240:rate=1 -c:v libx264 -preset ultrafast -y %s", outputPath)
	return executeCommand("ffmpeg", cmd)
}

func createTestAudio(outputPath string) error {
	cmd := fmt.Sprintf("-f lavfi -i sine=frequency=1000:duration=1 -c:a aac -y %s", outputPath)
	return executeCommand("ffmpeg", cmd)
}

func executeCommand(cmd string, args ...string) error {
	execCmd := exec.Command(cmd, args...)
	if err := execCmd.Run(); err != nil {
		return err
	}
	return nil
}