package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

const (
	REPORT_DIR = "../report"
)

var (
	HIGH   *PRIORITY = cretePriority("HIGH")
	LOW              = cretePriority("LOW")
	MEDIUM           = cretePriority("MEDIUM")
)

type PRIORITY = string

type Report struct {
	Timestamp int       `json:"timestamp" yml:"timestamp"`
	Time      string    `json:"time" yml:"time"`
	Err       string    `json:"error" yml:"error"`
	Reviewed  bool      `json:"reviewed,omitempty" yml:"reviewed,omitempty"`
	Priority  *PRIORITY `json:"priority" yml:"priority"`
	Metadata  string    `json:"meta_data,omitempty" yml:"meta_data,omitempty"`
}

func cretePriority(payload string) *PRIORITY {
	return &payload
}

func (r *Report) setReportTime() {
	r.Time = time.Now().Format(time.DateTime)
}

func CreateNewReport() *Report {
	timestamp := time.Now().Nanosecond()
	report := Report{}
	report.setReportTime()
	report.Timestamp = timestamp
	report.Priority = LOW
	return &report
}

func getReportDir() (string, error) {
	CWD, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(CWD, REPORT_DIR), nil
}

func (report *Report) PersistReport() error {
	var encodedReport []byte
	reportDir, err := getReportDir()
	if err != nil {
		return err
	}
	timestampDir := fmt.Sprintf("%s/%d", reportDir, report.Timestamp)
	_, err = os.Stat(timestampDir)
	if err != nil {
		err = os.MkdirAll(timestampDir, 0755)
		if err != nil {
			return err
		}
	}
	encodedReport, err = json.Marshal(report)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(timestampDir, "report.json"), encodedReport, 0644)
	if err != nil {
		return err
	}
	return nil
}
