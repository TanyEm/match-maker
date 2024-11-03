// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/TanyEm/match-maker/v2/internal/lobby (interfaces: Lobbier)
//
// Generated by this command:
//
//	mockgen -destination=./lobby_mock.go -package=lobby github.com/TanyEm/match-maker/v2/internal/lobby Lobbier
//

// Package lobby is a generated GoMock package.
package lobby

import (
	reflect "reflect"

	player "github.com/TanyEm/match-maker/v2/internal/player"
	gomock "go.uber.org/mock/gomock"
)

// MockLobbier is a mock of Lobbier interface.
type MockLobbier struct {
	ctrl     *gomock.Controller
	recorder *MockLobbierMockRecorder
	isgomock struct{}
}

// MockLobbierMockRecorder is the mock recorder for MockLobbier.
type MockLobbierMockRecorder struct {
	mock *MockLobbier
}

// NewMockLobbier creates a new mock instance.
func NewMockLobbier(ctrl *gomock.Controller) *MockLobbier {
	mock := &MockLobbier{ctrl: ctrl}
	mock.recorder = &MockLobbierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLobbier) EXPECT() *MockLobbierMockRecorder {
	return m.recorder
}

// AddPlayer mocks base method.
func (m *MockLobbier) AddPlayer(p player.Player) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddPlayer", p)
}

// AddPlayer indicates an expected call of AddPlayer.
func (mr *MockLobbierMockRecorder) AddPlayer(p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPlayer", reflect.TypeOf((*MockLobbier)(nil).AddPlayer), p)
}

// GetMatchByJoinID mocks base method.
func (m *MockLobbier) GetMatchByJoinID(joinID string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchByJoinID", joinID)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetMatchByJoinID indicates an expected call of GetMatchByJoinID.
func (mr *MockLobbierMockRecorder) GetMatchByJoinID(joinID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchByJoinID", reflect.TypeOf((*MockLobbier)(nil).GetMatchByJoinID), joinID)
}

// Run mocks base method.
func (m *MockLobbier) Run() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run")
}

// Run indicates an expected call of Run.
func (mr *MockLobbierMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockLobbier)(nil).Run))
}

// Stop mocks base method.
func (m *MockLobbier) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockLobbierMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockLobbier)(nil).Stop))
}
