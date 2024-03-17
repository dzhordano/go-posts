// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mocks/mock.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	domain "github.com/dzhordano/go-posts/internal/domain"
	service "github.com/dzhordano/go-posts/internal/service"
	gomock "go.uber.org/mock/gomock"
)

// MockUsers is a mock of Users interface.
type MockUsers struct {
	ctrl     *gomock.Controller
	recorder *MockUsersMockRecorder
}

// MockUsersMockRecorder is the mock recorder for MockUsers.
type MockUsersMockRecorder struct {
	mock *MockUsers
}

// NewMockUsers creates a new mock instance.
func NewMockUsers(ctrl *gomock.Controller) *MockUsers {
	mock := &MockUsers{ctrl: ctrl}
	mock.recorder = &MockUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsers) EXPECT() *MockUsersMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockUsers) GetAll(ctx context.Context) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockUsersMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockUsers)(nil).GetAll), ctx)
}

// GetById mocks base method.
func (m *MockUsers) GetById(ctx context.Context, userId uint) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, userId)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockUsersMockRecorder) GetById(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUsers)(nil).GetById), ctx, userId)
}

// RefreshTokens mocks base method.
func (m *MockUsers) RefreshTokens(ctx context.Context, refreshToken string) (service.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshTokens", ctx, refreshToken)
	ret0, _ := ret[0].(service.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshTokens indicates an expected call of RefreshTokens.
func (mr *MockUsersMockRecorder) RefreshTokens(ctx, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshTokens", reflect.TypeOf((*MockUsers)(nil).RefreshTokens), ctx, refreshToken)
}

// SignIN mocks base method.
func (m *MockUsers) SignIN(ctx context.Context, input domain.UserSignInInput) (service.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIN", ctx, input)
	ret0, _ := ret[0].(service.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIN indicates an expected call of SignIN.
func (mr *MockUsersMockRecorder) SignIN(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIN", reflect.TypeOf((*MockUsers)(nil).SignIN), ctx, input)
}

// SignUP mocks base method.
func (m *MockUsers) SignUP(ctx context.Context, input domain.UserSignUpInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUP", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUP indicates an expected call of SignUP.
func (mr *MockUsersMockRecorder) SignUP(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUP", reflect.TypeOf((*MockUsers)(nil).SignUP), ctx, input)
}

// MockAdmins is a mock of Admins interface.
type MockAdmins struct {
	ctrl     *gomock.Controller
	recorder *MockAdminsMockRecorder
}

// MockAdminsMockRecorder is the mock recorder for MockAdmins.
type MockAdminsMockRecorder struct {
	mock *MockAdmins
}

// NewMockAdmins creates a new mock instance.
func NewMockAdmins(ctrl *gomock.Controller) *MockAdmins {
	mock := &MockAdmins{ctrl: ctrl}
	mock.recorder = &MockAdminsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdmins) EXPECT() *MockAdminsMockRecorder {
	return m.recorder
}

// CensorComment mocks base method.
func (m *MockAdmins) CensorComment(ctx context.Context, commId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CensorComment", ctx, commId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CensorComment indicates an expected call of CensorComment.
func (mr *MockAdminsMockRecorder) CensorComment(ctx, commId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CensorComment", reflect.TypeOf((*MockAdmins)(nil).CensorComment), ctx, commId)
}

// DeleteComment mocks base method.
func (m *MockAdmins) DeleteComment(ctx context.Context, commId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteComment", ctx, commId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteComment indicates an expected call of DeleteComment.
func (mr *MockAdminsMockRecorder) DeleteComment(ctx, commId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComment", reflect.TypeOf((*MockAdmins)(nil).DeleteComment), ctx, commId)
}

// DeleteUser mocks base method.
func (m *MockAdmins) DeleteUser(ctx context.Context, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockAdminsMockRecorder) DeleteUser(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockAdmins)(nil).DeleteUser), ctx, userId)
}

// RefreshTokens mocks base method.
func (m *MockAdmins) RefreshTokens(ctx context.Context, refreshToken string) (service.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshTokens", ctx, refreshToken)
	ret0, _ := ret[0].(service.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshTokens indicates an expected call of RefreshTokens.
func (mr *MockAdminsMockRecorder) RefreshTokens(ctx, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshTokens", reflect.TypeOf((*MockAdmins)(nil).RefreshTokens), ctx, refreshToken)
}

// SignIN mocks base method.
func (m *MockAdmins) SignIN(ctx context.Context, input domain.UserSignInInput) (service.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIN", ctx, input)
	ret0, _ := ret[0].(service.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIN indicates an expected call of SignIN.
func (mr *MockAdminsMockRecorder) SignIN(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIN", reflect.TypeOf((*MockAdmins)(nil).SignIN), ctx, input)
}

// SuspendPost mocks base method.
func (m *MockAdmins) SuspendPost(ctx context.Context, postId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuspendPost", ctx, postId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SuspendPost indicates an expected call of SuspendPost.
func (mr *MockAdminsMockRecorder) SuspendPost(ctx, postId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuspendPost", reflect.TypeOf((*MockAdmins)(nil).SuspendPost), ctx, postId)
}

// SuspendUser mocks base method.
func (m *MockAdmins) SuspendUser(ctx context.Context, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuspendUser", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SuspendUser indicates an expected call of SuspendUser.
func (mr *MockAdminsMockRecorder) SuspendUser(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuspendUser", reflect.TypeOf((*MockAdmins)(nil).SuspendUser), ctx, userId)
}

// UpdateUser mocks base method.
func (m *MockAdmins) UpdateUser(ctx context.Context, input domain.UpdateUserInput, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, input, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockAdminsMockRecorder) UpdateUser(ctx, input, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockAdmins)(nil).UpdateUser), ctx, input, userId)
}

// MockPosts is a mock of Posts interface.
type MockPosts struct {
	ctrl     *gomock.Controller
	recorder *MockPostsMockRecorder
}

// MockPostsMockRecorder is the mock recorder for MockPosts.
type MockPostsMockRecorder struct {
	mock *MockPosts
}

// NewMockPosts creates a new mock instance.
func NewMockPosts(ctrl *gomock.Controller) *MockPosts {
	mock := &MockPosts{ctrl: ctrl}
	mock.recorder = &MockPostsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPosts) EXPECT() *MockPostsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPosts) Create(ctx context.Context, input domain.Post, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, input, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockPostsMockRecorder) Create(ctx, input, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPosts)(nil).Create), ctx, input, userId)
}

// Delete mocks base method.
func (m *MockPosts) Delete(ctx context.Context, postId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, postId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPostsMockRecorder) Delete(ctx, postId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPosts)(nil).Delete), ctx, postId)
}

// DeleteUser mocks base method.
func (m *MockPosts) DeleteUser(ctx context.Context, postId, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, postId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockPostsMockRecorder) DeleteUser(ctx, postId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockPosts)(nil).DeleteUser), ctx, postId, userId)
}

// GetAll mocks base method.
func (m *MockPosts) GetAll(ctx context.Context) ([]domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPostsMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPosts)(nil).GetAll), ctx)
}

// GetAllUser mocks base method.
func (m *MockPosts) GetAllUser(ctx context.Context, userId uint) ([]domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUser", ctx, userId)
	ret0, _ := ret[0].([]domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUser indicates an expected call of GetAllUser.
func (mr *MockPostsMockRecorder) GetAllUser(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUser", reflect.TypeOf((*MockPosts)(nil).GetAllUser), ctx, userId)
}

// GetById mocks base method.
func (m *MockPosts) GetById(ctx context.Context, postId uint) (domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, postId)
	ret0, _ := ret[0].(domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockPostsMockRecorder) GetById(ctx, postId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockPosts)(nil).GetById), ctx, postId)
}

// Like mocks base method.
func (m *MockPosts) Like(ctx context.Context, postId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Like", ctx, postId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Like indicates an expected call of Like.
func (mr *MockPostsMockRecorder) Like(ctx, postId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Like", reflect.TypeOf((*MockPosts)(nil).Like), ctx, postId)
}

// RemoveLike mocks base method.
func (m *MockPosts) RemoveLike(ctx context.Context, postId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLike", ctx, postId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveLike indicates an expected call of RemoveLike.
func (mr *MockPostsMockRecorder) RemoveLike(ctx, postId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLike", reflect.TypeOf((*MockPosts)(nil).RemoveLike), ctx, postId)
}

// Update mocks base method.
func (m *MockPosts) Update(ctx context.Context, input domain.UpdatePostInput, postId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, input, postId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPostsMockRecorder) Update(ctx, input, postId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPosts)(nil).Update), ctx, input, postId)
}

// UpdateUser mocks base method.
func (m *MockPosts) UpdateUser(ctx context.Context, input domain.UpdatePostInput, postId, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, input, postId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockPostsMockRecorder) UpdateUser(ctx, input, postId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockPosts)(nil).UpdateUser), ctx, input, postId, userId)
}

// MockComments is a mock of Comments interface.
type MockComments struct {
	ctrl     *gomock.Controller
	recorder *MockCommentsMockRecorder
}

// MockCommentsMockRecorder is the mock recorder for MockComments.
type MockCommentsMockRecorder struct {
	mock *MockComments
}

// NewMockComments creates a new mock instance.
func NewMockComments(ctrl *gomock.Controller) *MockComments {
	mock := &MockComments{ctrl: ctrl}
	mock.recorder = &MockCommentsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComments) EXPECT() *MockCommentsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockComments) Create(ctx context.Context, input domain.Comment, postId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, input, postId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCommentsMockRecorder) Create(ctx, input, postId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockComments)(nil).Create), ctx, input, postId)
}

// DeleteUser mocks base method.
func (m *MockComments) DeleteUser(ctx context.Context, commId, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, commId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockCommentsMockRecorder) DeleteUser(ctx, commId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockComments)(nil).DeleteUser), ctx, commId, userId)
}

// GetComments mocks base method.
func (m *MockComments) GetComments(ctx context.Context, postId uint) ([]domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComments", ctx, postId)
	ret0, _ := ret[0].([]domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetComments indicates an expected call of GetComments.
func (mr *MockCommentsMockRecorder) GetComments(ctx, postId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComments", reflect.TypeOf((*MockComments)(nil).GetComments), ctx, postId)
}

// GetUserComments mocks base method.
func (m *MockComments) GetUserComments(ctx context.Context, userId uint) ([]domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserComments", ctx, userId)
	ret0, _ := ret[0].([]domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserComments indicates an expected call of GetUserComments.
func (mr *MockCommentsMockRecorder) GetUserComments(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserComments", reflect.TypeOf((*MockComments)(nil).GetUserComments), ctx, userId)
}

// GetUserPostComments mocks base method.
func (m *MockComments) GetUserPostComments(ctx context.Context, userId, postId uint) ([]domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPostComments", ctx, userId, postId)
	ret0, _ := ret[0].([]domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPostComments indicates an expected call of GetUserPostComments.
func (mr *MockCommentsMockRecorder) GetUserPostComments(ctx, userId, postId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPostComments", reflect.TypeOf((*MockComments)(nil).GetUserPostComments), ctx, userId, postId)
}

// UpdateUser mocks base method.
func (m *MockComments) UpdateUser(ctx context.Context, input domain.UpdateCommentInput, commId, userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, input, commId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockCommentsMockRecorder) UpdateUser(ctx, input, commId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockComments)(nil).UpdateUser), ctx, input, commId, userId)
}
