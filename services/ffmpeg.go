package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xfrr/goffmpeg/transcoder"
)

// MediaType representing the media type here
type MediaType int

const (
	Audio MediaType = iota
	Video
)

// QualityProfile define the media file quality
type QualityProfile struct {
	Name         string
	VideoCodec   string
	AudioCodec   string
	VideoBitrate string
	AudioBitrate string
	Resolution   string
	Preset       string
	CRF          int    // Constant Rate Factor for the quality
	MaxSize      string // Max file size
}

// Main media transcription struct
type MediaOptimizer struct {
	InputPath  string
	OutputPath string
	MediaType  MediaType
	Profile    QualityProfile
	transcoder *transcoder.Transcoder
}

var (
	// Audio profile
	AudioMobileLow = QualityProfile{
		Name:         "audio_mobile_low",
		AudioCodec:   "aac",
		AudioBitrate: "64k",
		Preset:       "fast",
	}

	AudioMobileHigh = QualityProfile{
		Name:         "audio_mobile_high",
		AudioCodec:   "aac",
		AudioBitrate: "128k",
		Preset:       "fast",
	}

	// Video profile
	VideoMobileLow = QualityProfile{
		Name:         "video_mobile_low",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		VideoBitrate: "500k",
		AudioBitrate: "64k",
		Resolution:   "480x360",
		Preset:       "fast",
		CRF:          20,
	}

	VideoMobileHigh = QualityProfile{
		Name:         "video_mobile_high",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		VideoBitrate: "1000k",
		AudioBitrate: "128k",
		Resolution:   "720x480",
		Preset:       "medium",
		CRF:          23,
	}

	VideoMobileUltra = QualityProfile{
		Name:         "video_mobile_ultra",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		VideoBitrate: "2000k",
		AudioBitrate: "128k",
		Resolution:   "1280x720",
		Preset:       "medium",
		CRF:          20,
	}
)

// Create a new instance of the media transcription service
func NewMediaOptimizer(inputPath, outputPath string) (*MediaOptimizer, error) {
	if !fileExists(inputPath) {
		return nil, errors.New("not found input file")
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return nil, fmt.Errorf("can't create output dir: %v", err)
	}

	mediaType := detectMediaType(inputPath)

	return &MediaOptimizer{
		InputPath:  inputPath,
		OutputPath: outputPath,
		MediaType:  mediaType,
		transcoder: new(transcoder.Transcoder),
	}, nil
}

// Setting the quality profile
func (m *MediaOptimizer) SetProfile(profile QualityProfile) {
	m.Profile = profile
}

// Setting a personalized quality profile
func (m *MediaOptimizer) SetCustomProfile(name, videoCodec, audioCodec, videoBitrate, audioBitrate, resolution, preset string, crf int) {
	m.Profile = QualityProfile{
		Name:         name,
		VideoCodec:   videoCodec,
		AudioCodec:   audioCodec,
		VideoBitrate: videoBitrate,
		AudioBitrate: audioBitrate,
		Resolution:   resolution,
		Preset:       preset,
		CRF:          crf,
	}
}

// Automatic optimization for mobile device depending on the provided quality
// The quality type can be `low` | `high` or `ultra` for video content
func (m *MediaOptimizer) OptimizeForMobile(quality string) error {
	switch m.MediaType {
	case Audio:
		switch quality {
		case "low":
			m.SetProfile(AudioMobileLow)
		case "high":
			m.SetProfile(AudioMobileHigh)
		default:
			m.SetProfile(AudioMobileHigh)
		}
	case Video:
		switch quality {
		case "low":
			m.SetProfile(VideoMobileLow)
		case "high":
			m.SetProfile(VideoMobileHigh)
		case "ultra":
			m.SetProfile(VideoMobileUltra)
		default:
			m.SetProfile(VideoMobileHigh)
		}
	}

	return m.Optimize()
}

// Run the optimization
func (m *MediaOptimizer) Optimize() error {
	err := m.transcoder.Initialize(m.InputPath, m.OutputPath)
	if err != nil {
		return fmt.Errorf("transcoder initialization error: %v", err)
	}

	m.configureTranscoder()

	done := m.transcoder.Run(true)

	progress := m.transcoder.Output()
	go func() {
		for p := range progress {
			fmt.Printf("Progression: %b\n", p.Progress)
		}
	}()

	result := <-done
	if result.Error() != "" {
		return fmt.Errorf("transcription error: %v", result.Error)
	}

	fmt.Printf("Optimization finished: %s -> %s\n", m.InputPath, m.OutputPath)
	return nil
}

// Running the optimization process with a callback func to take the progress
func (m *MediaOptimizer) OptimizeWithCallback(progressCallback func(float64)) error {
	err := m.transcoder.Initialize(m.InputPath, m.OutputPath)
	if err != nil {
		return fmt.Errorf("transcoder initialization error: %v", err)
	}

	m.configureTranscoder()

	done := m.transcoder.Run(true)
	progress := m.transcoder.Output()

	go func() {
		for p := range progress {
			if progressCallback != nil {
				progressCallback(p.Progress)
			}
		}
	}()

	result := <-done
	if result.Error() != "" {
		return fmt.Errorf("transcription error: %v", result.Error)
	}

	return nil
}

// Configuring the transcoder based on the profile configuration
func (m *MediaOptimizer) configureTranscoder() {
	mediaFile := m.transcoder.MediaFile()

	// Based preset
	if m.Profile.Preset != "" {
		mediaFile.SetPreset(m.Profile.Preset)
	}

	// Video Config
	if m.MediaType == Video {
		if m.Profile.VideoCodec != "" {
			mediaFile.SetVideoCodec(m.Profile.VideoCodec)
		}
		if m.Profile.VideoBitrate != "" {
			mediaFile.SetVideoBitRate(m.Profile.VideoBitrate)
		}
		if m.Profile.Resolution != "" {
			mediaFile.SetResolution(m.Profile.Resolution)
		}
		if m.Profile.CRF != 0 {
			mediaFile.SetCRF(uint32(m.Profile.CRF))
		}
	}

	// Audio config
	if m.Profile.AudioCodec != "" {
		mediaFile.SetAudioCodec(m.Profile.AudioCodec)
	}
	if m.Profile.AudioBitrate != "" {
		mediaFile.SetAudioBitRate(m.Profile.AudioBitrate)
	}

	// Specific optimization for mobile
	mediaFile.SetMovFlags("+faststart") // For streaming
	mediaFile.SetPixFmt("yuv420p")      // Max compatibility
}

// Get the optimized file size
func (m *MediaOptimizer) GetOptimizedSize() (int64, error) {
	if !fileExists(m.OutputPath) {
		return 0, errors.New("cannot find output file")
	}

	info, err := os.Stat(m.OutputPath)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

// Calculating Compression ration
func (m *MediaOptimizer) GetCompressionRatio() (float64, error) {
	inputInfo, err := os.Stat(m.InputPath)
	if err != nil {
		return 0, err
	}

	outputSize, err := m.GetOptimizedSize()
	if err != nil {
		return 0, err
	}

	return float64(outputSize) / float64(inputInfo.Size()), nil
}

// Optimizing multiple inputs files
func BatchOptimize(inputPaths []string, outputDir string, quality string) error {
	for _, inputPath := range inputPaths {
		filename := filepath.Base(inputPath)
		ext := filepath.Ext(filename)
		nameWithoutExt := strings.TrimSuffix(filename, ext)

		var outputPath string
		if detectMediaType(inputPath) == Video {
			outputPath = filepath.Join(outputDir, nameWithoutExt+"_optimized.mp4")
		} else {
			outputPath = filepath.Join(outputDir, nameWithoutExt+"_optimized.aac")
		}

		optimizer, err := NewMediaOptimizer(inputPath, outputPath)
		if err != nil {
			fmt.Printf("Optimizer creation error for %s: %v\n", inputPath, err)
			continue
		}

		err = optimizer.OptimizeForMobile(quality)
		if err != nil {
			fmt.Printf("Optimization error %s: %v\n", inputPath, err)
			continue
		}

		fmt.Printf("Optimized: %s\n", inputPath)
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func detectMediaType(path string) MediaType {
	ext := strings.ToLower(filepath.Ext(path))
	videoExts := []string{".mp4", ".avi", ".mov", ".mkv", ".webm", ".3gp", ".flv", ".wmv"}

	for _, vExt := range videoExts {
		if ext == vExt {
			return Video
		}
	}

	return Audio
}

func OptimizeVideoForMobile(inputPath, outputPath, quality string) error {
	optimizer, err := NewMediaOptimizer(inputPath, outputPath)
	if err != nil {
		return err
	}

	return optimizer.OptimizeForMobile(quality)
}

func OptimizeAudioForMobile(inputPath, outputPath, quality string) error {
	optimizer, err := NewMediaOptimizer(inputPath, outputPath)
	if err != nil {
		return err
	}

	return optimizer.OptimizeForMobile(quality)
}

// Compatibility func with the existed codebase
func Transcode() {
	optimizer, err := NewMediaOptimizer("../fixtures/input.3gp", "../test_results/ultrafast-output.mp4")
	if err != nil {
		panic(err)
	}

	err = optimizer.OptimizeForMobile("high")
	if err != nil {
		panic(err)
	}
}
