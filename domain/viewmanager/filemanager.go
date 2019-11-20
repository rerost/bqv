package viewmanager

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

type FileManager struct {
	dir string
}

type fileView struct {
	dataSet string
	name    string
	query   string
}

func (f fileView) DataSet() string {
	return f.dataSet
}

func (f fileView) Name() string {
	return f.name
}

func (f fileView) Query() string {
	return f.query
}

func NewFileManager(dir string) FileManager {
	return FileManager{dir: dir}
}

func (f FileManager) List(ctx context.Context) ([]View, error) {
	dir := f.dir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	views := []View{}
	for _, file := range files {
		if !file.IsDir() {
			return nil, errors.Wrap(errors.New("Unexpected file found"), file.Name())
		}

		dataSet := file.Name()
		files, err := ioutil.ReadDir(path.Join(dir, file.Name()))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		for _, file := range files {
			if file.IsDir() {
				return nil, errors.Wrap(errors.New("Unexpected dir found"), file.Name())
			}

			name := file.Name()
			bquery, err := ioutil.ReadFile(f.Path(fileView{dataSet: dataSet, name: name}))
			if err != nil {
				return nil, errors.WithStack(err)
			}

			v := fileView{
				dataSet: dataSet,
				name:    name,
				query:   string(bquery),
			}
			views = append(views, v)
		}
	}

	return nil, nil
}
func (f FileManager) Get(ctx context.Context, dataset string, name string) (View, error) {
	bquery, err := ioutil.ReadFile(f.Path(fileView{dataSet: dataset, name: name}))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return fileView{
		dataSet: dataset,
		name:    name,
		query:   string(bquery),
	}, nil
}
func (f FileManager) Create(ctx context.Context, view View) (View, error) {
	file, err := os.Create(f.Path(view))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()

	_, err = file.WriteString((view.Query()))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return fileView{
		dataSet: view.DataSet(),
		name:    view.Name(),
		query:   view.Query(),
	}, nil
}
func (f FileManager) Update(ctx context.Context, view View) (View, error) {
	file, err := os.OpenFile(f.Path(view), os.O_WRONLY, 0222)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()

	_, err = file.WriteString(view.Query())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return f.Get(ctx, view.DataSet(), view.Name())
}
func (f FileManager) Delete(ctx context.Context, view View) error {
	err := os.Remove(f.Path(view))
	return errors.WithStack(err)
}

func (f FileManager) Path(view View) string {
	return path.Join(f.dir, view.DataSet(), view.Name()+".sql")
}
