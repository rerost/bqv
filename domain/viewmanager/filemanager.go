package viewmanager

import "context"

type FileManager struct {
}

func NewFileManager(dir string) FileManager {
	// TODO
	return FileManager{}
}

func (FileManager) List(ctx context.Context) ([]View, error) {
	return nil, nil
}
func (FileManager) Get(ctx context.Context, dataset string, name string) (View, error) {
	return nil, nil
}
func (FileManager) Create(ctx context.Context, view View) (View, error) {
	return nil, nil
}
func (FileManager) Update(ctx context.Context, view View) (View, error) {
	return nil, nil
}
func (FileManager) Delete(ctx context.Context, view View) error {
	return nil
}
