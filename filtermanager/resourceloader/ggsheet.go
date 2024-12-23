package resourceloader

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsLoader struct {
	apiKey        string
	spreadsheetID string
	logger        zerolog.Logger
}

func NewSheetsLoader(apiKey string, spreadsheetID string, logger zerolog.Logger) ResourceLoader {
	return &SheetsLoader{
		apiKey:        apiKey,
		spreadsheetID: spreadsheetID,
		logger:        logger,
	}
}

func (s *SheetsLoader) Load() ([]string, error) {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx,
		option.WithAPIKey(s.apiKey),
		option.WithScopes(sheets.SpreadsheetsReadonlyScope))
	if err != nil {
		s.logger.Error().Err(err).Msg("Unable to create sheets client")
		return nil, err
	}

	// 1. Lấy metadata để biết có những sheets nào
	spreadsheet, err := srv.Spreadsheets.Get(s.spreadsheetID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get spreadsheet metadata: %v", err)
	}

	var result []string
	// 2. Với mỗi sheet, lấy values
	for _, sheet := range spreadsheet.Sheets {
		range_ := fmt.Sprintf("%s!A:Z", sheet.Properties.Title)
		resp, err := srv.Spreadsheets.Values.Get(s.spreadsheetID, range_).Do()
		if err != nil {
			return nil, fmt.Errorf("failed to get values from sheet %s: %v",
				sheet.Properties.Title, err)
		}

		// Xử lý values từ resp.Values
		for _, row := range resp.Values {
			for _, cell := range row {
				if str, ok := cell.(string); ok && str != "" {
					result = append(result, str)
				}
			}
		}
	}
	s.logger.Info().
		Int("words_count", len(result)).
		Msg("Loaded bad words from Google Sheets")

	return result, nil
}

func (s *SheetsLoader) ValidateSource() bool {
	return s.spreadsheetID != ""
}
