package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/guilherme-santos/messageboard"
	mbhttp "github.com/guilherme-santos/messageboard/http"
	"github.com/guilherme-santos/messageboard/mock"
	"github.com/stretchr/testify/assert"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
)

func TestMessageBoardHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mock.NewService(ctrl)
	svc.EXPECT().
		List(gomock.Any(), &messageboard.ListOptions{
			PerPage: 10,
			Page:    2,
		}).
		Return(&messageboard.MessageList{
			Total: 0,
			Data:  make([]*messageboard.Message, 0),
		}, nil)

	router := chi.NewRouter()
	mbhttp.NewMessageBoardHandler(router, svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://localhost/v1/messages?per_page=10&page=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"total": 0,
		"data": []
	}`, w.Body.String())
}

func TestMessageBoardHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqMsg := &messageboard.Message{
		Name:  "Guilherme",
		Email: "xguiga@gmail.com",
		Text:  "My text goes here",
	}

	svc := mock.NewService(ctrl)
	svc.EXPECT().
		Create(gomock.Any(), reqMsg).
		Return(&messageboard.Message{
			ID:           "my-id",
			Name:         "Guilherme",
			Email:        "xguiga@gmail.com",
			Text:         "My text goes here",
			CreationTime: time.Date(2020, time.August, 12, 15, 30, 0, 0, time.UTC),
		}, nil)

	router := chi.NewRouter()
	mbhttp.NewMessageBoardHandler(router, svc)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(reqMsg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "http://localhost/v1/messages", &buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{
		"id": "my-id",
		"name": "Guilherme",
		"email": "xguiga@gmail.com",
		"text": "My text goes here",
		"creation_time": "2020-08-12T15:30:00Z"
	}`, w.Body.String())
}

func TestMessageBoardHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mock.NewService(ctrl)
	svc.EXPECT().
		Get(gomock.Any(), "my-id").
		Return(&messageboard.Message{
			ID:           "my-id",
			Name:         "Guilherme",
			Email:        "xguiga@gmail.com",
			Text:         "My text goes here",
			CreationTime: time.Date(2020, time.August, 12, 15, 30, 0, 0, time.UTC),
		}, nil)

	router := chi.NewRouter()
	mbhttp.NewMessageBoardHandler(router, svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://localhost/v1/messages/my-id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"id": "my-id",
		"name": "Guilherme",
		"email": "xguiga@gmail.com",
		"text": "My text goes here",
		"creation_time": "2020-08-12T15:30:00Z"
	}`, w.Body.String())
}

func TestMessageBoardHandler_GetNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mock.NewService(ctrl)
	svc.EXPECT().
		Get(gomock.Any(), "my-id").
		Return(nil, messageboard.NewError("not_found", "message not found"))

	router := chi.NewRouter()
	mbhttp.NewMessageBoardHandler(router, svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://localhost/v1/messages/my-id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{
		"code": "not_found",
		"message": "message not found"
	}`, w.Body.String())
}

func TestMessageBoardHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqMsg := &messageboard.Message{
		ID:    "my-id",
		Name:  "Guilherme",
		Email: "xguiga@gmail.com",
		Text:  "My text was updated",
	}

	svc := mock.NewService(ctrl)
	svc.EXPECT().
		Get(gomock.Any(), "my-id").
		Return(&messageboard.Message{
			ID:           "my-id",
			Name:         "Guilherme",
			Email:        "xguiga@gmail.com",
			Text:         "My text goes here",
			CreationTime: time.Date(2020, time.August, 12, 15, 30, 0, 0, time.UTC),
		}, nil)
	svc.EXPECT().
		Update(gomock.Any(), reqMsg).
		Return(&messageboard.Message{
			ID:           "my-id",
			Name:         "Guilherme",
			Email:        "xguiga@gmail.com",
			Text:         "My text was updated",
			CreationTime: time.Date(2020, time.August, 12, 15, 30, 0, 0, time.UTC),
		}, nil)

	router := chi.NewRouter()
	mbhttp.NewMessageBoardHandler(router, svc)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(reqMsg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "http://localhost/v1/messages/my-id", &buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{
		"id": "my-id",
		"name": "Guilherme",
		"email": "xguiga@gmail.com",
		"text": "My text was updated",
		"creation_time": "2020-08-12T15:30:00Z"
	}`, w.Body.String())
}

func TestMessageBoardHandler_UpdateNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mock.NewService(ctrl)
	svc.EXPECT().
		Get(gomock.Any(), "my-id").
		Return(nil, messageboard.NewError("not_found", "message not found"))

	router := chi.NewRouter()
	mbhttp.NewMessageBoardHandler(router, svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "http://localhost/v1/messages/my-id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{
		"code": "not_found",
		"message": "message not found"
	}`, w.Body.String())
}
