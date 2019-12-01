package annotateparser

import (
	"context"
	"fmt"
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
			if i != 0 && readedManifestLine == 0 {
				// Flush
				manifests = append(manifests, newManifest(buf, mType))
			}
			// Prepare for next manifest
			readedManifestLine = 0
			mType = strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			buf = ""
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

	fmt.Println(manifests)
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
