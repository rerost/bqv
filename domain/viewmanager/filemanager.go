package viewmanager

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type FileManager struct {
	dir string
}

type fileView struct {
	dataSet string
	name    string
	query   string
	setting fileSetting
}

type fileSetting struct {
	Metadata_ map[string]interface{} `yaml:"metadata"`
}

func (fs fileSetting) Metadata() map[string]interface{} {
	return fs.Metadata_
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

func (f fileView) Setting() Setting {
	return f.setting
}

func NewFileManager(dir string) FileManager {
	return FileManager{dir: dir}
}

func (f FileManager) List(ctx context.Context) ([]View, error) {
	zap.L().Debug("Open file", zap.String("dir", f.dir))
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
				zap.L().Info("Unexpected dir found", zap.String("dir name", file.Name()))
				continue
			}

			if !strings.HasSuffix(file.Name(), ".sql") {
				zap.L().Info("Not sql file found", zap.String("file name", file.Name()))
				continue
			}

			name := strings.TrimSuffix(file.Name(), ".sql")
			inCompleteFileView := fileView{dataSet: dataSet, name: name}

			bquery, err := ioutil.ReadFile(f.Path(inCompleteFileView))
			if err != nil {
				return nil, errors.WithStack(err)
			}

			sSetting, err := ioutil.ReadFile(f.SettingPath(inCompleteFileView))
			if err != nil {
				return nil, errors.WithStack(err)
			}
			var setting fileSetting
			if err := yaml.Unmarshal(sSetting, &setting); err != nil {
				return nil, errors.WithMessagef(err, "Failed to parse %s", f.SettingPath(inCompleteFileView))
			}

			v := fileView{
				dataSet: dataSet,
				name:    name,
				query:   string(bquery),
				setting: setting,
			}
			views = append(views, v)
		}
	}

	return views, nil
}

func (f FileManager) Get(ctx context.Context, dataset string, name string) (View, error) {
	if _, err := os.Stat(f.Path(fileView{dataSet: dataset, name: name})); err != nil {
		if os.IsNotExist(err) {
			return nil, NotFoundError
		}
		return nil, err
	}
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
	if _, err := os.Stat(f.DatasetPath(view)); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(f.DatasetPath(view), 0755)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			return nil, err
		}
	}

	// Create sql file.
	{
		file, err := os.Create(f.Path(view))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		defer file.Close()

		_, err = file.WriteString((view.Query()))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// Create yml file
	{
		file, err := os.Create(f.SettingPath(view))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		fv := f.convertToFileView(view)
		out, err := yaml.Marshal(fv.setting)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		_, err = file.Write(out)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return f.convertToFileView(view), nil
}

func (f FileManager) Update(ctx context.Context, view View) (View, error) {
	if _, err := os.Stat(f.Path(view)); err != nil {
		if os.IsNotExist(err) {
			return nil, NotFoundError
		}
		return nil, err
	}
	{
		file, err := os.OpenFile(f.Path(view), os.O_WRONLY, 0222)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		defer file.Close()

		_, err = file.WriteString(view.Query())
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	{
		file, err := os.OpenFile(f.SettingPath(view), os.O_WRONLY, 0222)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		defer file.Close()

		fv := f.convertToFileView(view)
		out, err := yaml.Marshal(fv.setting)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if _, err := file.Write(out); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return f.Get(ctx, view.DataSet(), view.Name())
}
func (f FileManager) Delete(ctx context.Context, view View) error {
	err := os.Remove(f.Path(view))
	if err != nil {
		return errors.WithStack(err)
	}

	if err := os.Remove(f.SettingPath(view)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (f FileManager) Path(view View) string {
	return path.Join(f.dir, view.DataSet(), view.Name()+".sql")
}

func (f FileManager) DatasetPath(view View) string {
	return path.Join(f.dir, view.DataSet())
}

func (f FileManager) SettingPath(view View) string {
	return path.Join(f.dir, view.DataSet(), view.Name()+"yml")
}

func (f FileManager) convertToFileView(view View) fileView {
	return fileView{
		dataSet: view.DataSet(),
		name:    view.Name(),
		query:   view.Query(),
		setting: fileSetting{
			Metadata_: view.Setting().Metadata(),
		},
	}
}
