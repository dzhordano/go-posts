package v1

import (
	"bytes"
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/dzhordano/go-posts/internal/service"
	mock_service "github.com/dzhordano/go-posts/internal/service/mocks"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsers, user userSignUpInput)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           userSignUpInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","email":"testingapi@mail.ru","password":"13092003qwerty"}`,
			inputUser: userSignUpInput{
				Name:     "Test",
				Email:    "testingapi@mail.ru",
				Password: "13092003qwerty",
			},
			mockBehavior: func(s *mock_service.MockUsers, input userSignUpInput) {
				s.EXPECT().SignUP(context.Background(), input).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "Invalid Name",
			inputBody:          `{"email":"testingapi@mail.ru","password":"13092003qwerty"}`,
			mockBehavior:       func(s *mock_service.MockUsers, input userSignUpInput) {},
			expectedStatusCode: 400,
		},
		{
			name:               "Invalid Email",
			inputBody:          `{"name":"Test","password":"13092003qwerty"}`,
			mockBehavior:       func(s *mock_service.MockUsers, input userSignUpInput) {},
			expectedStatusCode: 400,
		},
		{
			name:               "Invalid Password",
			inputBody:          `{"name":"Test","email":"testingapi@mail.ru","password":"qwerty"}`,
			mockBehavior:       func(s *mock_service.MockUsers, input userSignUpInput) {},
			expectedStatusCode: 400,
		},
		{
			name:      "Service Failure",
			inputBody: `{"name":"Test","email":"testingapi@mail.ru","password":"13092003qwerty"}`,
			inputUser: userSignUpInput{
				Name:     "Test",
				Email:    "testingapi@mail.ru",
				Password: "13092003qwerty",
			},
			mockBehavior: func(s *mock_service.MockUsers, input userSignUpInput) {
				s.EXPECT().SignUP(context.Background(), input).Return(errors.New("service failure")) // TODO: need to swap this thing to more secure?
			},
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			signup := mock_service.NewMockUsers(c)
			testCase.mockBehavior(signup, testCase.inputUser)

			services := &service.Services{Users: signup}
			handler := NewHandler(services, &auth.Manager{})

			// Test server
			r := gin.New()
			r.POST("/sign-up", handler.userSignUp)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

//func TestHandler_getPosts(t *testing.T) {}
// func TestHandler_createPost(t *testing.T) {}
// func TestHandler_updatePost(t *testing.T) {}
// func TestHandler_deletePost(t *testing.T) {}

// func TestHandler_createPostComment(t *testing.T) {}
// func TestHandler_getUserPostComments(t *testing.T) {}
// func TestHandler_likePost(t *testing.T) {}
// func TestHandler_unlikePost(t *testing.T) {}

// func TestHandler_getComments(t *testing.T) {}
// func TestHandler_updateComment(t *testing.T) {}
// func TestHandler_deleteComment(t *testing.T) {}
