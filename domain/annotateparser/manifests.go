package annotateparser

import "context"

import "github.com/pkg/errors"

// TODO(@rerost) Rethink name
type Manifests interface {
	Manifests(ctx context.Context, target string) ([]Manifest, error)
}

type manifestsImpl struct {
	extractor Extractor
	parser    Parser
}

func NewManifests(extractor Extractor, parser Parser) Manifests {
	return manifestsImpl{
		extractor: extractor,
		parser:    parser,
	}
}

func (m manifestsImpl) Manifests(ctx context.Context, target string) ([]Manifest, error) {
	annotation, err := m.extractor.Extract(target)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ms, err := m.parser.Parse(ctx, annotation[0])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ms, nil
}
