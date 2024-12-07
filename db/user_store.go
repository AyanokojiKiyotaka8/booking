package db

import (
	"context"
	"fmt"

	"github.com/AyanokojiKiyotaka8/booking.git/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper

	GetUser(context.Context, bson.M) (*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	GetUsers(context.Context, bson.M) ([]*types.User, error)
	DeleteUser(context.Context, bson.M) error
	UpdateUser(context.Context, bson.M, bson.M) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client, dbname string) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll: client.Database(dbname).Collection("users"),
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("---- dropping the user collection ----")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) GetUser(ctx context.Context, filter bson.M) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context, filter bson.M) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, filter bson.M) error {
	_, err := s.coll.DeleteOne(ctx, filter)
	return err
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}