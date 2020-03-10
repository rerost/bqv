package viewmanager

import (
	"context"

	"github.com/pkg/errors"
)

type Setting interface {
	Metadata() map[string]interface{}
}

type View interface {
	DataSet() string
	Name() string
	Query() string
	Setting() Setting
}

type ViewReader interface {
	List(ctx context.Context) ([]View, error)
	Get(ctx context.Context, dataset string, name string) (View, error)
}

type ViewWriter interface {
	Create(ctx context.Context, view View) (View, error)
	Update(ctx context.Context, view View) (View, error)
	Delete(ctx context.Context, view View) error
}

type ViewReadWriter interface {
	ViewReader
	ViewWriter
}

var NotFoundError = errors.New("NotFound")
