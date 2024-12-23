package resourceloader

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

// Google Drive constants
const (
	googleDriveViewURL = "https://drive.google.com/uc?export=view&id="
)

// Google Drive Adapter
type DriveLoader struct {
	wg      *sync.WaitGroup
	sources []string
	logger  zerolog.Logger
}

func NewDriveLoader(sources []string, logger zerolog.Logger) ResourceLoader {
	return &DriveLoader{
		sources: sources,
		wg:      &sync.WaitGroup{},
		logger:  logger,
	}
}

func (d *DriveLoader) Load() ([]string, error) {
	resource := []string{}
	for _, source := range d.sources {
		d.wg.Add(1)
		go func(source string) {
			directURL := d.getGoogleDriveDirectURL(source)
			res, err := d.loadSource(directURL)
			if err == nil {
				resource = append(resource, res...)
			} else {
				d.logger.Error().Msgf("Error downloading file: %v\n", err)
			}
			d.wg.Done()
		}(source)

	}

	d.wg.Wait()
	return resource, nil
}

func (d *DriveLoader) ValidateSource() bool {
	for _, source := range d.sources {
		directURL := d.getGoogleDriveDirectURL(source)
		if directURL == "" {
			return false
		}
	}

	return true
}

// getGoogleDriveDirectURL converts share URL to direct download URL
func (d *DriveLoader) getGoogleDriveDirectURL(shareURL string) string {
	fileID := ""
	if strings.Contains(shareURL, "drive.google.com/file/d/") {
		parts := strings.Split(shareURL, "/")
		for i, part := range parts {
			if part == "d" && i+1 < len(parts) {
				fileID = parts[i+1]
				break
			}
		}
	} else if strings.Contains(shareURL, "drive.google.com/open?id=") {
		fileID = strings.Split(shareURL, "=")[1]
	}

	if fileID == "" {
		return ""
	}

	return googleDriveViewURL + fileID
}

// loadSource loads badword patterns from Google Drive
func (d *DriveLoader) loadSource(directURL string) ([]string, error) {
	var resource []string

	client := &http.Client{}
	req, err := http.NewRequest("GET", directURL, nil)
	if err != nil {

		return nil, fmt.Errorf("error creating request: %v\n", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error downloading file: %v\n", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			resource = append(resource, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error scanning file: %v\n", err)
	}

	return resource, nil
}
