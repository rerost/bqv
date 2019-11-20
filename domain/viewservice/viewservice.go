package viewservice

import "context"

type ViewService interface {
	Diff(ctx context.Context, src ViewReader, dst ViewReader) (string, error)
	Copy(ctx context.Context, src ViewReader, dst ViewReadWriter) error
}

type View interface {
	DataSet() string
	Name() string
	Query() string
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

type viewServiceImpl struct {
	source      ViewReader
	destination ViewWriter
}

func NewService() ViewService {
	return viewServiceImpl{}
}

func (s viewServiceImpl) Diff(ctx context.Context, src ViewReader, dst ViewReader) (string, error) {
	// TODO
	return "", nil
}

func (s viewServiceImpl) Copy(ctx context.Context, src ViewReader, dst ViewReadWriter) error {
	// TODO
	return nil
}

func diff(source View, destination View) string {
	// TODO
	return ""
}
