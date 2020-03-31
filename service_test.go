package messageboard_test

import (
	"context"
	"testing"
	"time"

	"github.com/guilherme-santos/messageboard"
	"github.com/guilherme-santos/messageboard/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expMsg := &messageboard.Message{
		ID:           "my-id",
		Name:         "Guilherme",
		Email:        "xguiga@gmail.com",
		Text:         "My long message",
		CreationTime: time.Now().UTC(),
	}
	reqMsg := &messageboard.Message{
		Name:  expMsg.Name,
		Email: expMsg.Email,
		Text:  expMsg.Text,
	}

	storage := mock.NewStorage(ctrl)
	storage.EXPECT().
		Create(gomock.Any(), reqMsg).
		DoAndReturn(func(ctx context.Context, msg *messageboard.Message) error {
			// Set id in the request message
			msg.ID = expMsg.ID
			return nil
		})
	storage.EXPECT().
		Get(gomock.Any(), expMsg.ID).
		Return(expMsg, nil)

	ctx := context.Background()

	svc := messageboard.NewService(storage)
	msg, err := svc.Create(ctx, reqMsg)
	assert.NoError(t, err)
	assert.Equal(t, expMsg, msg)
}

// TODO:: implement test for List, today is just a bypass for storage, but it's not
// been implemented for sake of time.
// func TestService_List(t *testing.T) {
// }

// TODO:: implement test for Get, today is just a bypass for storage, but it's not
// been implemented for sake of time.
// func TestService_Get(t *testing.T) {
// }

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expMsg := &messageboard.Message{
		ID:           "my-id",
		Name:         "Guilherme",
		Email:        "xguiga@gmail.com",
		Text:         "My long message",
		CreationTime: time.Now().UTC(),
	}
	reqMsg := &messageboard.Message{
		ID:    expMsg.ID,
		Name:  expMsg.Name,
		Email: expMsg.Email,
		Text:  expMsg.Text,
	}

	storage := mock.NewStorage(ctrl)
	storage.EXPECT().
		Update(gomock.Any(), reqMsg).
		Return(nil)
	storage.EXPECT().
		Get(gomock.Any(), expMsg.ID).
		Return(expMsg, nil)

	ctx := context.Background()

	svc := messageboard.NewService(storage)
	msg, err := svc.Update(ctx, reqMsg)
	assert.NoError(t, err)
	assert.Equal(t, expMsg, msg)
}
