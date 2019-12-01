package annotateparser

import (
	"regexp"

	"github.com/pkg/errors"
)

var (
	annotationsRegexp = regexp.MustCompile(`/\*.*\//`)
	Unspoorted        = errors.New("Unspoorted")
)

type Extractor interface {
	Extract(string) ([]string, error)
}

type extractorImpl struct {
}

func Extract(manifest string) ([]string, error) {
	matched := annotationsRegexp.FindAll([]byte(manifest), -1)
	if len(matched) == 0 {
		return nil, nil
	}

	if len(matched) == 1 {
		return []string{string(matched[0])}, nil
	}

	return nil, Unspoorted
}
