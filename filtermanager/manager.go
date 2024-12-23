package filtermanager

import (
	"errors"
	"fmt"

	"github.com/MetaLeapX/go-badword-filter/filtermanager/filtermodel"
)

var (
	ErrInvalidSource = errors.New("invalid source")
)

// Resource loader interface
type ResourceLoader interface {
	Load() ([]string, error)
	ValidateSource() bool
}

type FilterController interface {
	GetAll(string) ([]string, error)
	GetAllAsync(string) []string
	ReplaceAll(string) string
	ReplaceAllAsync(string) string
	ContainsBadWords(string) bool
}

type controller interface {
	FilterController
}

// FilterEvent handles the badword filtering operations
type FilterEvent struct {
}

// Create new filter manager with specific loader
func NewFilterManager(loader ResourceLoader) (FilterController, error) {
	if !loader.ValidateSource() {
		return nil, ErrInvalidSource
	}
	resource, err := loader.Load()
	if err != nil {
		return nil, err
	}

	if len(resource) == 0 {
		return nil, errors.New("resource is empty")
	}
	fmt.Println("resource count: ", len(resource))

	filtermodel.SetResource(resource)
	return &FilterEvent{}, nil
}

// GetAll returns all found badwords in the text
func (f FilterEvent) GetAll(text string) ([]string, error) {
	params := filtermodel.FilterParams{
		Text:     text,
		Prefix:   "",
		Postfix:  "",
		MarkOnly: false,
	}
	return params.GetAll()
}

// GetAllAsync asynchronously returns all found badwords
func (f FilterEvent) GetAllAsync(text string) []string {
	resultChan := make(chan []string)
	go func() {
		params := filtermodel.FilterParams{
			Text:     text,
			Prefix:   "",
			Postfix:  "",
			MarkOnly: false,
		}
		all, _ := params.GetAll()
		resultChan <- all
	}()

	return <-resultChan
}

// ReplaceAll replaces all badwords in the text
func (f FilterEvent) ReplaceAll(text string) string {
	params := filtermodel.FilterParams{
		Text:     text,
		Prefix:   "***",
		Postfix:  "***",
		MarkOnly: false,
	}
	return params.ReplaceAll()
}

// ReplaceAllAsync asynchronously replaces all badwords
func (f FilterEvent) ReplaceAllAsync(text string) string {
	resultChan := make(chan string)

	go func() {
		params := filtermodel.FilterParams{
			Text:     text,
			Prefix:   "***",
			Postfix:  "***",
			MarkOnly: false,
		}
		resultChan <- params.ReplaceAll()
	}()

	return <-resultChan
}

func (f FilterEvent) ContainsBadWords(text string) bool {
	params := filtermodel.FilterParams{
		Text:     text,
		Prefix:   "",
		Postfix:  "",
		MarkOnly: false,
	}
	return params.ContainsBadWords()
}
