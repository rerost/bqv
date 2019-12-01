package annotateparser

import (
	"context"
	"strings"
)

type Manifest interface {
	Type() string
	Body() string
}

type Parser interface {
	Parse(ctx context.Context, annotation string) ([]Manifest, error)
}

type parserImpl struct {
}

func NewParser() Parser {
	return parserImpl{}
}

func (p parserImpl) Parse(ctx context.Context, annotation string) ([]Manifest, error) {
	manifests := []Manifest{}
	lines := strings.Split(annotation, "\n")

	buf := ""
	mType := ""
	readedManifestLine := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "[") {
			// NOTE:
			if i != 0 {
				// Flush
				manifests = append(manifests, newManifest(buf, mType))
			}
			// Prepare for next manifest
			readedManifestLine = 0
			buf = ""
			mType = strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			continue
		}
		if readedManifestLine != 0 {
			buf += "\n"
		}
		buf += line
		readedManifestLine++

		if i == len(lines)-1 {
			manifests = append(manifests, newManifest(buf, mType))
		}
	}

	return manifests, nil
}

type manifestImpl struct {
	body  string
	mtype string
}

func newManifest(body string, mType string) Manifest {
	return manifestImpl{
		body:  body,
		mtype: mType,
	}
}

func (m manifestImpl) Body() string {
	return m.body
}

func (m manifestImpl) Type() string {
	return m.mtype
}
