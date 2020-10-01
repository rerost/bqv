package viewservice

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/iterator"
	"github.com/pkg/errors"
	"github.com/rerost/bqv/domain/viewmanager"
	datatransfer "cloud.google.com/go/bigquery/datatransfer/apiv1"
	datatransferpb "google.golang.org/genproto/googleapis/cloud/bigquery/datatransfer/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type View = viewmanager.View
type ViewReadWriter = viewmanager.ViewReadWriter
type ViewWriter = viewmanager.ViewWriter
type ViewReader = viewmanager.ViewReader
type DataTransferClient = *datatransfer.Client

type ViewService interface {
	List(ctx context.Context, src ViewReader) ([]View, error)
	Diff(ctx context.Context, src ViewReader, dst ViewReader) ([]View, error)
	Copy(ctx context.Context, src ViewReader, dst ViewWriter) error
}

type ViewServiceImpl struct {
	source      ViewReader
	destination ViewWriter
	datatransferClient DataTransferClient
	projectID string
}

type cachedViewTable struct {
	view View
}

func (v cachedViewTable) Name() string {
	return "cached_" + v.view.Name()
}

func (v cachedViewTable) DataSet() string {
	return v.view.DataSet()
}

func (v cachedViewTable) Query() string {
	return v.view.Query()
}

func (v cachedViewTable) Setting() viewmanager.Setting {
	return v.view.Setting()
}

func NewService(datatransferClient DataTransferClient, projectID string) ViewService {
	return ViewServiceImpl{datatransferClient: datatransferClient, projectID: projectID}
}

func (s ViewServiceImpl) List(ctx context.Context, src ViewReader) ([]View, error) {
	return src.List(ctx)
}

func (s ViewServiceImpl) Diff(ctx context.Context, src ViewReader, dst ViewReader) ([]View, error) {
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

func (s ViewServiceImpl) copy(ctx context.Context, item viewmanager.View, dst ViewWriter) error {
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
	} else if err != nil {
		return errors.WithStack(err)
	} 
	return nil
}
func (s ViewServiceImpl) applyCachedView(ctx context.Context, item viewmanager.View, dst ViewWriter, listSchedulingMap map[string]string) error {
	// スケジューリングクエリの一覧を検索し、item.Setting().Metadata()["view_table"]によっては
	// テーブルを削除したり、スケジューリングクエリを削除する
	cachedViewTable := cachedViewTable{item}
	fmt.Println("CachedViewTable", cachedViewTable)

	_, ok := listSchedulingMap[cachedViewTable.DataSet() + "." + cachedViewTable.Name()]
	if ok == false {
		zap.L().Debug("View already exists, but no scheduling query", zap.String("Table", cachedViewTable.Name()), zap.String("Dataset", cachedViewTable.DataSet()))
		if item.Setting().Metadata()["view_table"] == true {
			err := s.cachedviewtable(ctx, item)
			if err != nil {
				zap.L().Debug("Delete Scheduling query err", zap.String("err", err.Error()))
				return errors.WithStack(err)
			}
		}
	} else {
		if item.Setting().Metadata()["view_table"] == false {
			zap.L().Debug("Deleting scheduling query and cached table")

			req := &datatransferpb.DeleteTransferConfigRequest{
				Name: listSchedulingMap[cachedViewTable.DataSet() + "." + cachedViewTable.Name()]}
			err := s.datatransferClient.DeleteTransferConfig(ctx, req)
			if err != nil {
				zap.L().Debug("Delete Scheduling query err", zap.String("err", err.Error()))
				return errors.WithStack(err)

			err = dst.Delete(ctx, cachedViewTable)
			if err != nil {
				zap.L().Debug("Delete Cached Table err", zap.String("err", err.Error()))
				return errors.WithStack(err)
			}
		}
	}
}
	return nil
}
func (s ViewServiceImpl) cachedviewtable(ctx context.Context, item viewmanager.View) error {
	cachedViewTable := cachedViewTable{item}
	zap.L().Debug("creating view table",  zap.String("View Table", cachedViewTable.Name()),  zap.String("Dataset", cachedViewTable.DataSet()))
	// BQ定期ジョブ実行

	m := make(map[string]interface{})
	m["query"] =  "CREATE OR REPLACE TABLE " + cachedViewTable.DataSet() + "." + cachedViewTable.Name() + " AS SELECT * FROM " + item.DataSet() + "." + item.Name()
	structParams, err := structpb.NewStruct(m)
	if err != nil {
		zap.L().Debug("Err", zap.String("err", err.Error()))
		return errors.WithStack(err)
	}
	schedulingTime, ok := item.Setting().Metadata()["scheduling_time"].(string)
	if !ok {
		zap.L().Debug("Err", zap.String("invalid scheduling_time setting", err.Error()))
		return errors.WithStack(err)
	}
	req := &datatransferpb.CreateTransferConfigRequest{
		Parent: datatransfer.ProjectPath(s.projectID),
		TransferConfig: &datatransferpb.TransferConfig{
			Destination: &datatransferpb.TransferConfig_DestinationDatasetId{
				DestinationDatasetId: cachedViewTable.DataSet()},		
			DisplayName: cachedViewTable.DataSet() + "." + cachedViewTable.Name(),
			DataSourceId: "scheduled_query",
			Params: structParams,
			Schedule: schedulingTime,
			ScheduleOptions: &datatransferpb.ScheduleOptions{ StartTime: timestamppb.Now()}}}
	resp, err := s.datatransferClient.CreateTransferConfig(ctx, req)
	if err != nil {
		zap.L().Debug("Err", zap.String("err", err.Error()))
		return errors.WithStack(err)
	}
	zap.L().Debug("scheduling query", zap.String("display_name", resp.GetDisplayName()), zap.String("config_id", resp.GetName()))
	return nil
}

func (s ViewServiceImpl) Copy(ctx context.Context, src ViewReader, dst ViewWriter) error {
	srcList, err := src.List(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	// スケジューリングクエリ収集
	req := &datatransferpb.ListTransferConfigsRequest{
		Parent: datatransfer.ProjectPath(s.projectID),
	}
	it :=  s.datatransferClient.ListTransferConfigs(ctx, req)
	schedulingQueryMap := map[string] string{}
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			zap.L().Debug("Err", zap.String("err", err.Error()))
			return errors.WithStack(err)
		}
		schedulingQueryMap[resp.GetDisplayName()] = resp.GetName()
	}
	var errs []error
	for _, srcView := range srcList {
		err := s.copy(ctx, srcView, dst)
		if err != nil {
			zap.L().Debug("Failed to copy view", zap.String("Dataset", srcView.DataSet()), zap.String("Table", srcView.Name()))
			errs = append(errs, errors.WithStack(err))
		}
		err = s.applyCachedView(ctx, srcView, dst, schedulingQueryMap)
		if err != nil {
			zap.L().Debug("Failed to apply cached view", zap.String("Dataset", srcView.DataSet()), zap.String("Table", srcView.Name()))
			errs = append(errs, errors.WithStack(err))
		}
	
	}

	return errors.WithStack(multierr.Combine(errs...))
}

func (s ViewServiceImpl) DeleteOld(ctx context.Context, src ViewReader, dst ViewReadWriter) error {
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
