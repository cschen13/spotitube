// Code generated by MockGen. DO NOT EDIT.
// Source: controllers/track.go

// Package mock_controllers is a generated GoMock package.
package mock_controllers

import (
	models "github.com/cschen13/spotitube/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MocktracksClient is a mock of tracksClient interface
type MocktracksClient struct {
	ctrl     *gomock.Controller
	recorder *MocktracksClientMockRecorder
}

// MocktracksClientMockRecorder is the mock recorder for MocktracksClient
type MocktracksClientMockRecorder struct {
	mock *MocktracksClient
}

// NewMocktracksClient creates a new mock instance
func NewMocktracksClient(ctrl *gomock.Controller) *MocktracksClient {
	mock := &MocktracksClient{ctrl: ctrl}
	mock.recorder = &MocktracksClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MocktracksClient) EXPECT() *MocktracksClientMockRecorder {
	return m.recorder
}

// GetPlaylistInfo mocks base method
func (m *MocktracksClient) GetPlaylistInfo(arg0, arg1 string) (*models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlaylistInfo", arg0, arg1)
	ret0, _ := ret[0].(*models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlaylistInfo indicates an expected call of GetPlaylistInfo
func (mr *MocktracksClientMockRecorder) GetPlaylistInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlaylistInfo", reflect.TypeOf((*MocktracksClient)(nil).GetPlaylistInfo), arg0, arg1)
}

// GetTracks mocks base method
func (m *MocktracksClient) GetTracks(arg0 *models.Playlist) (models.Tracks, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTracks", arg0)
	ret0, _ := ret[0].(models.Tracks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTracks indicates an expected call of GetTracks
func (mr *MocktracksClientMockRecorder) GetTracks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTracks", reflect.TypeOf((*MocktracksClient)(nil).GetTracks), arg0)
}

// MockconvertSrcClient is a mock of convertSrcClient interface
type MockconvertSrcClient struct {
	ctrl     *gomock.Controller
	recorder *MockconvertSrcClientMockRecorder
}

// MockconvertSrcClientMockRecorder is the mock recorder for MockconvertSrcClient
type MockconvertSrcClientMockRecorder struct {
	mock *MockconvertSrcClient
}

// NewMockconvertSrcClient creates a new mock instance
func NewMockconvertSrcClient(ctrl *gomock.Controller) *MockconvertSrcClient {
	mock := &MockconvertSrcClient{ctrl: ctrl}
	mock.recorder = &MockconvertSrcClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockconvertSrcClient) EXPECT() *MockconvertSrcClientMockRecorder {
	return m.recorder
}

// GetPlaylistInfo mocks base method
func (m *MockconvertSrcClient) GetPlaylistInfo(arg0, arg1 string) (*models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlaylistInfo", arg0, arg1)
	ret0, _ := ret[0].(*models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlaylistInfo indicates an expected call of GetPlaylistInfo
func (mr *MockconvertSrcClientMockRecorder) GetPlaylistInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlaylistInfo", reflect.TypeOf((*MockconvertSrcClient)(nil).GetPlaylistInfo), arg0, arg1)
}

// GetTrackByID mocks base method
func (m *MockconvertSrcClient) GetTrackByID(arg0 string) (*models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrackByID", arg0)
	ret0, _ := ret[0].(*models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrackByID indicates an expected call of GetTrackByID
func (mr *MockconvertSrcClientMockRecorder) GetTrackByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrackByID", reflect.TypeOf((*MockconvertSrcClient)(nil).GetTrackByID), arg0)
}

// MockconvertDstClient is a mock of convertDstClient interface
type MockconvertDstClient struct {
	ctrl     *gomock.Controller
	recorder *MockconvertDstClientMockRecorder
}

// MockconvertDstClientMockRecorder is the mock recorder for MockconvertDstClient
type MockconvertDstClientMockRecorder struct {
	mock *MockconvertDstClient
}

// NewMockconvertDstClient creates a new mock instance
func NewMockconvertDstClient(ctrl *gomock.Controller) *MockconvertDstClient {
	mock := &MockconvertDstClient{ctrl: ctrl}
	mock.recorder = &MockconvertDstClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockconvertDstClient) EXPECT() *MockconvertDstClientMockRecorder {
	return m.recorder
}

// GetOwnPlaylistInfo mocks base method
func (m *MockconvertDstClient) GetOwnPlaylistInfo(arg0 string) (*models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOwnPlaylistInfo", arg0)
	ret0, _ := ret[0].(*models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOwnPlaylistInfo indicates an expected call of GetOwnPlaylistInfo
func (mr *MockconvertDstClientMockRecorder) GetOwnPlaylistInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwnPlaylistInfo", reflect.TypeOf((*MockconvertDstClient)(nil).GetOwnPlaylistInfo), arg0)
}

// CreatePlaylist mocks base method
func (m *MockconvertDstClient) CreatePlaylist(arg0 string) (*models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePlaylist", arg0)
	ret0, _ := ret[0].(*models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePlaylist indicates an expected call of CreatePlaylist
func (mr *MockconvertDstClientMockRecorder) CreatePlaylist(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePlaylist", reflect.TypeOf((*MockconvertDstClient)(nil).CreatePlaylist), arg0)
}

// InsertTrack mocks base method
func (m *MockconvertDstClient) InsertTrack(arg0 *models.Playlist, arg1 *models.Track) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTrack", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertTrack indicates an expected call of InsertTrack
func (mr *MockconvertDstClientMockRecorder) InsertTrack(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTrack", reflect.TypeOf((*MockconvertDstClient)(nil).InsertTrack), arg0, arg1)
}