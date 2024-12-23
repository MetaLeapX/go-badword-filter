package _test

import (
	"fmt"
	"testing"
	"time"

	"github.com/MetaLeapX/go-badword-filter/filtermanager"
	"github.com/MetaLeapX/go-badword-filter/filtermanager/resourceloader"
	"github.com/rs/zerolog"
)

func TestFilter(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	start := time.Now()
	filter, err := filtermanager.NewFilterManager(resourceloader.NewSheetsLoader(
		"AI...",
		"1OVr1wyyxkQ0eWxH6NkKagcE612LI9mgsfY0rE8vTDWE",
		zerolog.New(nil).With().Timestamp().Logger()))

	if err != nil {
		t.Fatal(err)
	}
	load := time.Since(start)
	fmt.Println(load)

	result := filter.ContainsBadWords("admin")
	t.Log(result)
}
