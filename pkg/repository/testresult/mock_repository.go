// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kubeshop/testkube/internal/pkg/api/repository/testresult (interfaces: Repository)

// Package testresult is a generated GoMock package.
package testresult

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	testkube "github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// DeleteAll mocks base method.
func (m *MockRepository) DeleteAll(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAll indicates an expected call of DeleteAll.
func (mr *MockRepositoryMockRecorder) DeleteAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockRepository)(nil).DeleteAll), arg0)
}

// DeleteByTestSuite mocks base method.
func (m *MockRepository) DeleteByTestSuite(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByTestSuite", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByTestSuite indicates an expected call of DeleteByTestSuite.
func (mr *MockRepositoryMockRecorder) DeleteByTestSuite(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByTestSuite", reflect.TypeOf((*MockRepository)(nil).DeleteByTestSuite), arg0, arg1)
}

// DeleteByTestSuites mocks base method.
func (m *MockRepository) DeleteByTestSuites(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByTestSuites", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByTestSuites indicates an expected call of DeleteByTestSuites.
func (mr *MockRepositoryMockRecorder) DeleteByTestSuites(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByTestSuites", reflect.TypeOf((*MockRepository)(nil).DeleteByTestSuites), arg0, arg1)
}

// EndExecution mocks base method.
func (m *MockRepository) EndExecution(arg0 context.Context, arg1 testkube.TestSuiteExecution) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EndExecution", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EndExecution indicates an expected call of EndExecution.
func (mr *MockRepositoryMockRecorder) EndExecution(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndExecution", reflect.TypeOf((*MockRepository)(nil).EndExecution), arg0, arg1)
}

// Get mocks base method.
func (m *MockRepository) Get(arg0 context.Context, arg1 string) (testkube.TestSuiteExecution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(testkube.TestSuiteExecution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), arg0, arg1)
}

// GetByNameAndTestSuite mocks base method.
func (m *MockRepository) GetByNameAndTestSuite(arg0 context.Context, arg1, arg2 string) (testkube.TestSuiteExecution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameAndTestSuite", arg0, arg1, arg2)
	ret0, _ := ret[0].(testkube.TestSuiteExecution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameAndTestSuite indicates an expected call of GetByNameAndTestSuite.
func (mr *MockRepositoryMockRecorder) GetByNameAndTestSuite(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameAndTestSuite", reflect.TypeOf((*MockRepository)(nil).GetByNameAndTestSuite), arg0, arg1, arg2)
}

// GetExecutions mocks base method.
func (m *MockRepository) GetExecutions(arg0 context.Context, arg1 Filter) ([]testkube.TestSuiteExecution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExecutions", arg0, arg1)
	ret0, _ := ret[0].([]testkube.TestSuiteExecution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExecutions indicates an expected call of GetExecutions.
func (mr *MockRepositoryMockRecorder) GetExecutions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExecutions", reflect.TypeOf((*MockRepository)(nil).GetExecutions), arg0, arg1)
}

// GetExecutionsTotals mocks base method.
func (m *MockRepository) GetExecutionsTotals(arg0 context.Context, arg1 ...Filter) (testkube.ExecutionsTotals, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetExecutionsTotals", varargs...)
	ret0, _ := ret[0].(testkube.ExecutionsTotals)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExecutionsTotals indicates an expected call of GetExecutionsTotals.
func (mr *MockRepositoryMockRecorder) GetExecutionsTotals(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExecutionsTotals", reflect.TypeOf((*MockRepository)(nil).GetExecutionsTotals), varargs...)
}

// GetLatestByTestSuite mocks base method.
func (m *MockRepository) GetLatestByTestSuite(arg0 context.Context, arg1, arg2 string) (testkube.TestSuiteExecution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestByTestSuite", arg0, arg1, arg2)
	ret0, _ := ret[0].(testkube.TestSuiteExecution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestByTestSuite indicates an expected call of GetLatestByTestSuite.
func (mr *MockRepositoryMockRecorder) GetLatestByTestSuite(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestByTestSuite", reflect.TypeOf((*MockRepository)(nil).GetLatestByTestSuite), arg0, arg1, arg2)
}

// GetLatestByTestSuites mocks base method.
func (m *MockRepository) GetLatestByTestSuites(arg0 context.Context, arg1 []string, arg2 string) ([]testkube.TestSuiteExecution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestByTestSuites", arg0, arg1, arg2)
	ret0, _ := ret[0].([]testkube.TestSuiteExecution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestByTestSuites indicates an expected call of GetLatestByTestSuites.
func (mr *MockRepositoryMockRecorder) GetLatestByTestSuites(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestByTestSuites", reflect.TypeOf((*MockRepository)(nil).GetLatestByTestSuites), arg0, arg1, arg2)
}

// GetTestSuiteMetrics mocks base method.
func (m *MockRepository) GetTestSuiteMetrics(arg0 context.Context, arg1 string, arg2, arg3 int) (testkube.ExecutionsMetrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTestSuiteMetrics", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(testkube.ExecutionsMetrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTestSuiteMetrics indicates an expected call of GetTestSuiteMetrics.
func (mr *MockRepositoryMockRecorder) GetTestSuiteMetrics(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTestSuiteMetrics", reflect.TypeOf((*MockRepository)(nil).GetTestSuiteMetrics), arg0, arg1, arg2, arg3)
}

// Insert mocks base method.
func (m *MockRepository) Insert(arg0 context.Context, arg1 testkube.TestSuiteExecution) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockRepositoryMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockRepository)(nil).Insert), arg0, arg1)
}

// StartExecution mocks base method.
func (m *MockRepository) StartExecution(arg0 context.Context, arg1 string, arg2 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartExecution", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartExecution indicates an expected call of StartExecution.
func (mr *MockRepositoryMockRecorder) StartExecution(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartExecution", reflect.TypeOf((*MockRepository)(nil).StartExecution), arg0, arg1, arg2)
}

// Update mocks base method.
func (m *MockRepository) Update(arg0 context.Context, arg1 testkube.TestSuiteExecution) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), arg0, arg1)
}
