package messageboard

import "context"

type service struct {
	storage Storage
}

// NewService returns the default (and likely the only) implementation of messageboard.Service.
//
// For the current use-case this implementation will be really simple, basicaly a proxy for
// storage, but could be more complex, such as emit events, call another services, etc.
//
func NewService(storage Storage) Service {
	return &service{
		storage: storage,
	}
}

func (s *service) Create(ctx context.Context, msg *Message) (*Message, error) {
	err := msg.Validate()
	if err != nil {
		return nil, err
	}

	err = s.storage.Create(ctx, msg)
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, msg.ID)
}

func (s *service) List(ctx context.Context, opts *ListOptions) (*MessageList, error) {
	return s.storage.List(ctx, opts)
}

func (s *service) Get(ctx context.Context, id string) (*Message, error) {
	return s.storage.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, msg *Message) (*Message, error) {
	err := msg.Validate()
	if err != nil {
		return nil, err
	}

	err = s.storage.Update(ctx, msg)
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, msg.ID)
}
