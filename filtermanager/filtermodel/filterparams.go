package filtermodel

import (
	"errors"
	"regexp"

	badwordfilter "github.com/MetaLeapX/go-badword-filter"
	"github.com/MetaLeapX/go-badword-filter/common/helper"
)

var (
	ErrTextIsNullOrEmpty = errors.New("Param Text is null or empty string! ")
)

var _resource []string

type FilterParams struct {
	Text     string
	Prefix   string
	Postfix  string
	MarkOnly bool
}
type capture struct {
	Value string
	Index int
}

func SetResource(resource []string) {
	_resource = resource
}

func (f FilterParams) GetAll() ([]string, error) {
	var result []string
	var hasBadWords bool
	if f.Text == "" {
		return nil, ErrTextIsNullOrEmpty
	} else {
		_text := f.Text
		for _, pattern := range _resource {
			r, _ := regexp.Compile(pattern)
			for _, value := range r.FindAllString(_text, -1) {
				result = append(result, value)
				hasBadWords = true
			}
		}
	}

	if hasBadWords {
		return result, nil
	}
	return nil, nil
}

func (f FilterParams) ReplaceAll() string {
	if f.Text == "" {
		return ""
	}

	result := f.Text
	_text := result
	offset := 0

	for _, pattern := range _resource {
		r, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		matches := r.FindAllStringIndex(_text, -1)
		for _, match := range matches {
			start := match[0] + offset
			end := match[1] + offset
			badWord := result[start:end]

			var replacement string
			if f.MarkOnly {
				replacement = f.Prefix + badWord + f.Postfix
			} else {
				replacement = f.Prefix + helper.CreateString(badwordfilter.ReplaceCharacters, len(badWord)) + f.Postfix
			}

			result = result[:start] + replacement + result[end:]
			offset += len(replacement) - (end - start)
		}
	}

	return result
}

func replaceCapture(text, prefix, postfix string, cap capture, markOnly bool, offset int) string {
	var runes []rune = []rune(text)
	var replacement string

	if markOnly {
		replacement = prefix + cap.Value + postfix
	} else {
		replacement = prefix + helper.CreateString(badwordfilter.ReplaceCharacters, len(cap.Value)) + postfix
	}

	index := cap.Index + offset
	runes = append(runes[:index], runes[index+len(cap.Value):]...)
	pre := make([]rune, len(runes[:index]))
	post := make([]rune, len(runes[index:]))
	copy(pre, runes[:index])
	copy(post, runes[index:])
	result := append(append(pre, []rune(replacement)...), post...)
	return string(result)
}

func (f FilterParams) ContainsBadWords() bool {
	if f.Text == "" {
		return false
	}

	_text := f.Text
	for _, pattern := range _resource {
		r, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		if r.MatchString(_text) {
			return true
		}
	}
	return false
}
