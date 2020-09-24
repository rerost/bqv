package viewservice

import (
	"context"
	"fmt"
	"os"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/rerost/bqv/domain/viewmanager"
	datatransfer "cloud.google.com/go/bigquery/datatransfer/apiv1"
	datatransferpb "google.golang.org/genproto/googleapis/cloud/bigquery/datatransfer/v1"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type View = viewmanager.View
type ViewReadWriter = viewmanager.ViewReadWriter
type ViewWriter = viewmanager.ViewWriter
type ViewReader = viewmanager.ViewReader

type ViewService interface {
	List(ctx context.Context, src ViewReader) ([]View, error)
	Diff(ctx context.Context, src ViewReader, dst ViewReader) ([]View, error)
	Copy(ctx context.Context, src ViewReader, dst ViewWriter) error
}

type viewServiceImpl struct {
	source      ViewReader
	destination ViewWriter
}

type CachedViewTable interface {
	DataSet() string
	Name() string
	Query() string
	Setting() viewmanager.Setting
}

type cachedviewtable struct {
	view View
}

func (v cachedviewtable) Name() string {
	return "Cached_" + v.view.Name()
}

func (v cachedviewtable) DataSet() string {
	return v.view.DataSet()
}

func (v cachedviewtable) Query() string {
	return v.view.Query()
}

func (v cachedviewtable) Setting() viewmanager.Setting {
	return v.view.Setting()
}

func NewService() ViewService {
	return viewServiceImpl{}
}

func (s viewServiceImpl) List(ctx context.Context, src ViewReader) ([]View, error) {
	return src.List(ctx)
}

func (s viewServiceImpl) Diff(ctx context.Context, src ViewReader, dst ViewReader) ([]View, error) {
	zap.L().Debug("Start Diff")
	srcList, err := src.List(ctx)
	if err != nil {
		zap.L().Debug("Failed to List", zap.String("src type", fmt.Sprintf("%T", src)))
		return nil, errors.WithStack(err)
	}

	diffViews := []View{}
	for _, srcView := range srcList {
		dstView, err := dst.Get(ctx, srcView.DataSet(), srcView.Name())
		if err == viewmanager.NotFoundError {
			err = nil
			diffView, err := diff(srcView, nil)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			diffViews = append(diffViews, diffView)
		}
		diffView, err := diff(srcView, dstView)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		diffViews = append(diffViews, diffView)
	}

	return diffViews, nil
}

func (s viewServiceImpl) copy(ctx context.Context, item viewmanager.View, dst ViewWriter) error {
	zap.L().Debug("Src", zap.String("dataset", item.DataSet()), zap.String("table", item.Name()))
	_, err := dst.Update(ctx, item)
	if err != nil {
		zap.L().Debug("Err", zap.String("err", err.Error()))
	}
	if err == viewmanager.NotFoundError {
		zap.L().Debug("Creating view", zap.String("Dataset", item.DataSet()), zap.String("Table", item.Name()))
		_, err := dst.Create(ctx, item)
		if err != nil {
			zap.L().Debug("Failed to create view", zap.String("Dataset", item.DataSet()), zap.String("Table", item.Name()))
			return errors.WithStack(err)
		}
		// TODO: ビューを作るかどうか
		cached_view_table := cachedviewtable{item}
		_, c_err := dst.Create(ctx, cached_view_table) 
		if c_err != nil {
			zap.L().Debug("Failed to create cached view table", zap.String("Dataset", cached_view_table.DataSet()), zap.String("Table", cached_view_table.Name()))
			return errors.WithStack(c_err)
		}
		// BQ定期ジョブ実行
		ctx := context.Background()
		c, err := datatransfer.NewClient(ctx)
		if err != nil {
			zap.L().Debug("Err", zap.String("err", err.Error()))
		}
		req := &datatransferpb.CreateTransferConfigRequest{
			Parent: datatransfer.ProjectPath(os.Getenv("GOOGLE_APPLICATION_PROJECT_ID")),
			TransferConfig: &datatransferpb.TransferConfig{
				// TODO: 中身作る
			}}
		resp, err := c.CreateTransferConfig(ctx, req)
		if c_err != nil {
			zap.L().Debug("Err", zap.String("err", err.Error()))
		}
		fmt.Println(resp)

	} else if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s viewServiceImpl) Copy(ctx context.Context, src ViewReader, dst ViewWriter) error {
	srcList, err := src.List(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	var errs []error
	for _, srcView := range srcList {
		err := s.copy(ctx, srcView, dst)
		if err != nil {
			zap.L().Debug("Failed to copy view", zap.String("Dataset", srcView.DataSet()), zap.String("Table", srcView.Name()))
			errs = append(errs, errors.WithStack(err))
		}
	}

	return errors.WithStack(multierr.Combine(errs...))
}

func (s viewServiceImpl) DeleteOld(ctx context.Context, src ViewReader, dst ViewReadWriter) error {
	srcList, err := src.List(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	dstList, err := dst.List(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, dstView := range dstList {
		if !matchInclude(dstView, srcList) {
			err := dst.Delete(ctx, dstView)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

type diffView struct {
	dataSet string
	name    string
	query   string
}

func (d diffView) DataSet() string {
	return d.dataSet
}

func (d diffView) Name() string {
	return d.name
}

func (d diffView) Query() string {
	return d.query
}

func (d diffView) Setting() viewmanager.Setting {
	panic("NOT IMPLEMENTED") // TODO(@rerost)
}

func diff(source View, destination View) (diffView, error) {
	if source == nil && destination == nil {
		return diffView{}, nil
	}
	if source == nil {
		return diffView{
			dataSet: destination.DataSet(),
			name:    destination.Name(),
			query:   destination.Query(),
		}, nil
	}
	if destination == nil {
		return diffView{
			dataSet: source.DataSet(),
			name:    source.Name(),
			query:   source.Query(),
		}, nil
	}
	if !match(source, destination) {
		zap.L().Debug(
			"Failed to diff",
			zap.String("source name", source.Name()),
			zap.String("destination name", destination.Name()),
			zap.String("source dataset", source.DataSet()),
			zap.String("destination dataset", destination.DataSet()),
		)
		return diffView{}, errors.New("Failed to diff")
	}
	if equal(source, destination) {
		return diffView{}, nil
	}

	return diffView{
		dataSet: source.DataSet(),
		name:    source.Name(),
		query:   cmp.Diff(source.Query(), destination.Query()),
	}, nil
}

func matchInclude(v View, vs []View) bool {
	for _, vv := range vs {
		if match(v, vv) {
			return true
		}
	}
	return false
}

func match(v1, v2 View) bool {
	if v1 == nil && v1 == v2 {
		return true
	}
	if v1 == nil || v2 == nil {
		return false
	}
	return v1.Name() == v2.Name() && v1.DataSet() == v2.DataSet()
}

func equal(v1, v2 View) bool {
	return match(v1, v2) && v1.Query() == v2.Query()
}
