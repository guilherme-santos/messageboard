package mongodb

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/guilherme-santos/messageboard"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

// MessageBoardStorage is a mongodb implementation of messageboard.Storage
type MessageBoardStorage struct {
	client *mongo.Client
	db     *mongo.Database
	coll   *mongo.Collection
}

func NewMessageBoardStorage(client *mongo.Client) *MessageBoardStorage {
	s := &MessageBoardStorage{client: client}
	s.db = s.client.Database("messageboard")
	s.coll = s.db.Collection("messages")
	return s
}

func (s *MessageBoardStorage) Create(ctx context.Context, msg *messageboard.Message) error {
	msg.ID = uuid.New().String()
	msg.CreationTime = time.Now().UTC()
	_, err := s.coll.InsertOne(ctx, msg)
	return err
}

func (s *MessageBoardStorage) List(ctx context.Context, opts *messageboard.ListOptions) (*messageboard.MessageList, error) {
	mgoOpts := options.Find().
		SetLimit(int64(opts.PerPage)).
		SetSkip(int64(opts.PerPage * (opts.Page - 1))).
		SetSort(bson.M{"creation_time": -1})

	list := new(messageboard.MessageList)

	g, ctx := errgroup.WithContext(ctx)
	// Goroutine to get list of results.
	g.Go(func() error {
		cursor, err := s.coll.Find(ctx, bson.D{}, mgoOpts)
		if err != nil {
			return err
		}
		return cursor.All(ctx, &list.Data)
	})
	// Goroutine to get total of results.
	g.Go(func() error {
		total, err := s.coll.CountDocuments(ctx, bson.D{}, options.Count())
		if err != nil {
			return err
		}
		list.Total = uint(total)
		return nil
	})

	// Wait both goroutines and abort in case of error.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	if list.Data == nil {
		list.Data = make([]*messageboard.Message, 0)
	}
	return list, nil
}

func (s *MessageBoardStorage) Get(ctx context.Context, id string) (*messageboard.Message, error) {
	var msg *messageboard.Message
	err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&msg)
	if err == mongo.ErrNoDocuments {
		return nil, messageboard.NewError("not_found", "message was not found")
	}
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *MessageBoardStorage) Update(ctx context.Context, msg *messageboard.Message) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"_id": msg.ID}, bson.M{
		"$set": bson.M{
			"name":  msg.Name,
			"email": msg.Email,
			"text":  msg.Text,
		},
	})
	return err
}

func (s *MessageBoardStorage) LoadCSV(initialCSV string) error {
	f, err := os.Open(initialCSV)
	if err != nil {
		return err
	}

	ctx := context.Background()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 5

	// Remove current collection to load csv from scratch
	err = s.coll.Drop(ctx)
	if err != nil {
		return err
	}

	var i int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		i++
		if i == 1 {
			// It's the header, we can ignore it
			continue
		}

		creationTime, err := time.Parse(time.RFC3339, record[4])
		if err != nil {
			return fmt.Errorf("invalid time on line %d: %v", i, err)
		}

		msg := &messageboard.Message{
			ID:           record[0],
			Name:         record[1],
			Email:        record[2],
			Text:         record[3],
			CreationTime: creationTime,
		}

		_, err = s.coll.InsertOne(ctx, msg)
		if err != nil {
			return err
		}
	}
	return nil
}
