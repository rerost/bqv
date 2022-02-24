// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/googleapis/google-cloud-go-testing/bigquery/bqiface (interfaces: Client)

// Package mock_bqiface is a generated GoMock package.
package mock_bqiface

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	bqiface "github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// Dataset mocks base method.
func (m *MockClient) Dataset(arg0 string) bqiface.Dataset {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dataset", arg0)
	ret0, _ := ret[0].(bqiface.Dataset)
	return ret0
}

// Dataset indicates an expected call of Dataset.
func (mr *MockClientMockRecorder) Dataset(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dataset", reflect.TypeOf((*MockClient)(nil).Dataset), arg0)
}

// DatasetInProject mocks base method.
func (m *MockClient) DatasetInProject(arg0, arg1 string) bqiface.Dataset {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatasetInProject", arg0, arg1)
	ret0, _ := ret[0].(bqiface.Dataset)
	return ret0
}

// DatasetInProject indicates an expected call of DatasetInProject.
func (mr *MockClientMockRecorder) DatasetInProject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatasetInProject", reflect.TypeOf((*MockClient)(nil).DatasetInProject), arg0, arg1)
}

// Datasets mocks base method.
func (m *MockClient) Datasets(arg0 context.Context) bqiface.DatasetIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Datasets", arg0)
	ret0, _ := ret[0].(bqiface.DatasetIterator)
	return ret0
}

// Datasets indicates an expected call of Datasets.
func (mr *MockClientMockRecorder) Datasets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Datasets", reflect.TypeOf((*MockClient)(nil).Datasets), arg0)
}

// DatasetsInProject mocks base method.
func (m *MockClient) DatasetsInProject(arg0 context.Context, arg1 string) bqiface.DatasetIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatasetsInProject", arg0, arg1)
	ret0, _ := ret[0].(bqiface.DatasetIterator)
	return ret0
}

// DatasetsInProject indicates an expected call of DatasetsInProject.
func (mr *MockClientMockRecorder) DatasetsInProject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatasetsInProject", reflect.TypeOf((*MockClient)(nil).DatasetsInProject), arg0, arg1)
}

// JobFromID mocks base method.
func (m *MockClient) JobFromID(arg0 context.Context, arg1 string) (bqiface.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JobFromID", arg0, arg1)
	ret0, _ := ret[0].(bqiface.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// JobFromID indicates an expected call of JobFromID.
func (mr *MockClientMockRecorder) JobFromID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JobFromID", reflect.TypeOf((*MockClient)(nil).JobFromID), arg0, arg1)
}

// JobFromIDLocation mocks base method.
func (m *MockClient) JobFromIDLocation(arg0 context.Context, arg1, arg2 string) (bqiface.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JobFromIDLocation", arg0, arg1, arg2)
	ret0, _ := ret[0].(bqiface.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// JobFromIDLocation indicates an expected call of JobFromIDLocation.
func (mr *MockClientMockRecorder) JobFromIDLocation(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JobFromIDLocation", reflect.TypeOf((*MockClient)(nil).JobFromIDLocation), arg0, arg1, arg2)
}

// Jobs mocks base method.
func (m *MockClient) Jobs(arg0 context.Context) bqiface.JobIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Jobs", arg0)
	ret0, _ := ret[0].(bqiface.JobIterator)
	return ret0
}

// Jobs indicates an expected call of Jobs.
func (mr *MockClientMockRecorder) Jobs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Jobs", reflect.TypeOf((*MockClient)(nil).Jobs), arg0)
}

// Location mocks base method.
func (m *MockClient) Location() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Location")
	ret0, _ := ret[0].(string)
	return ret0
}

// Location indicates an expected call of Location.
func (mr *MockClientMockRecorder) Location() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Location", reflect.TypeOf((*MockClient)(nil).Location))
}

// Query mocks base method.
func (m *MockClient) Query(arg0 string) bqiface.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0)
	ret0, _ := ret[0].(bqiface.Query)
	return ret0
}

// Query indicates an expected call of Query.
func (mr *MockClientMockRecorder) Query(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockClient)(nil).Query), arg0)
}

// SetLocation mocks base method.
func (m *MockClient) SetLocation(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLocation", arg0)
}

// SetLocation indicates an expected call of SetLocation.
func (mr *MockClientMockRecorder) SetLocation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLocation", reflect.TypeOf((*MockClient)(nil).SetLocation), arg0)
}

// embedToIncludeNewMethods mocks base method.
func (m *MockClient) embedToIncludeNewMethods() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "embedToIncludeNewMethods")
}

// embedToIncludeNewMethods indicates an expected call of embedToIncludeNewMethods.
func (mr *MockClientMockRecorder) embedToIncludeNewMethods() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "embedToIncludeNewMethods", reflect.TypeOf((*MockClient)(nil).embedToIncludeNewMethods))
}
