package messageboard

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

// Message represents a message inside of the system.
type Message struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Text         string    `json:"text"`
	CreationTime time.Time `json:"creation_time"`
}

//go:generate mockgen -package mock -mock_names Service=Service -destination mock/service.go github.com/guilherme-santos/messageboard Service

// Service defines an interface which implements a CRUD for message.
type Service interface {
	Create(context.Context, *Message) (*Message, error)
	List(context.Context, *ListOptions) (*MessageList, error)
	Get(_ context.Context, id string) (*Message, error)
	Update(context.Context, *Message) (*Message, error)
}

//go:generate mockgen -package mock -mock_names Storage=Storage -destination mock/storage.go github.com/guilherme-santos/messageboard Storage

// Storage defines an interface to access messages from a arbitrary storage.
type Storage interface {
	Create(context.Context, *Message) error
	List(context.Context, *ListOptions) (*MessageList, error)
	Get(_ context.Context, id string) (*Message, error)
	Update(context.Context, *Message) error
}

// MessageList is a struct containing the list of messages requested with some
// more informations, perhaps statistics, etc.
type MessageList struct {
	Total uint       `json:"total"`
	Data  []*Message `json:"data"`
}

// ListOptions is passed to Storage.List to filter/sort/paginate the results.
type ListOptions struct {
	PerPage uint
	Page    uint
}

const DefaultPerPage = 30

// Load loads values from query string into ListOptions.
func (opts *ListOptions) Load(values url.Values) {
	// PerPage
	if v := values.Get("per_page"); v != "" {
		perPage, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			perPage = DefaultPerPage
		}
		opts.PerPage = uint(perPage)
	} else {
		opts.PerPage = DefaultPerPage // Set default
	}
	// Page
	if v := values.Get("page"); v != "" {
		page, _ := strconv.ParseUint(v, 10, 32)
		opts.Page = uint(page)
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
}
