// Automatically generated by MockGen. DO NOT EDIT!
// Source: service.go

package mock_auth

import (
	session "2019_2_IBAT/pkg/app/auth/session"
	context "context"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// Mock of ServiceClient interface
type MockServiceClient struct {
	ctrl     *gomock.Controller
	recorder *_MockServiceClientRecorder
}

// Recorder for MockServiceClient (not exported)
type _MockServiceClientRecorder struct {
	mock *MockServiceClient
}

func NewMockServiceClient(ctrl *gomock.Controller) *MockServiceClient {
	mock := &MockServiceClient{ctrl: ctrl}
	mock.recorder = &_MockServiceClientRecorder{mock}
	return mock
}

func (_m *MockServiceClient) EXPECT() *_MockServiceClientRecorder {
	return _m.recorder
}

func (_m *MockServiceClient) CreateSession(ctx context.Context, in *session.Session, opts ...grpc.CallOption) (*session.CreateSessionInfo, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "CreateSession", _s...)
	ret0, _ := ret[0].(*session.CreateSessionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceClientRecorder) CreateSession(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateSession", _s...)
}

func (_m *MockServiceClient) DeleteSession(ctx context.Context, in *session.Cookie, opts ...grpc.CallOption) (*session.Bool, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "DeleteSession", _s...)
	ret0, _ := ret[0].(*session.Bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceClientRecorder) DeleteSession(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteSession", _s...)
}

func (_m *MockServiceClient) GetSession(ctx context.Context, in *session.Cookie, opts ...grpc.CallOption) (*session.GetSessionInfo, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GetSession", _s...)
	ret0, _ := ret[0].(*session.GetSessionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceClientRecorder) GetSession(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSession", _s...)
}
