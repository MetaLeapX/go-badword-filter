# Go Badword Filter

A high-performance bad word filtering library for Go applications with multiple data sources support.

## Features

- Multiple data sources (Google Sheets, Google Drive, GitHub)
- Concurrent processing
- Configurable filtering options
- Thread-safe operations
- Comprehensive logging

## Installation

```bash
go get github.com/MetaLeapX/go-badword-filter

Quick Start
package main

import (
    "github.com/MetaLeapX/go-badword-filter/filtermanager"
    "github.com/MetaLeapX/go-badword-filter/filtermanager/resourceloader"
    "github.com/rs/zerolog"
)

func main() {
    logger := zerolog.New(nil).With().Timestamp().Logger()
    
    filter, err := filtermanager.NewFilterManager(
        resourceloader.NewSheetsLoader(
            "YOUR_API_KEY",
            "YOUR_SPREADSHEET_ID", 
            logger,
        ),
    )
    if err != nil {
        panic(err)
    }

    result := filter.ReplaceAll("text to check")
    fmt.Println(result)
}
```
## Data Sources
- Google Sheets
```code
  loader := resourceloader.NewSheetsLoader(
  "YOUR_API_KEY",
  "YOUR_SPREADSHEET_ID",
  logger,
  )
```
##Configuration

Google Sheets Setup
1. Enable Google Sheets API in Google Cloud Console
2. Create API key with Sheets API access
3. Make spreadsheet public or configure access
4. Get spreadsheet ID from URL: https://docs.google.com/spreadsheets/d/{SPREADSHEET_ID}/edit

## Custom Configuration
```code
badwordfilter.ReplaceCharacters = "#" // Default is "*"
```

## Testing
```base
func TestFilter(t *testing.T) {
    zerolog.SetGlobalLevel(zerolog.DebugLevel)
    filter, err := filtermanager.NewFilterManager(
        resourceloader.NewSheetsLoader(
            "API_KEY",
            "SPREADSHEET_ID",
            zerolog.New(nil).With().Timestamp().Logger(),
        ),
    )
    
    if err != nil {
        t.Fatal(err)
    }
    
    result := filter.ReplaceAll("text to check")
    t.Log(result)
}
```
## Dependencies

- github.com/rs/zerolog
- google.golang.org/api/sheets/v4

## License
MIT License - see LICENSE file
