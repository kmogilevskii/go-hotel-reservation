package db

import (
	"context"
	"fmt"
	"os"

	"github.com/kmogilevskii/hotel-reservation/errors"
	"github.com/kmogilevskii/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Map map[string]any

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper
	GetUserByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context, Map, *Pagination) ([]*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, types.UpdateUserParams) error
}

type MongoUserStore struct {
	client *mongo.Client
	dbname string
	coll   *mongo.Collection
}

func (m *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, errors.ErrInvalidID()
	}

	var user types.User
	if err := m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := m.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *MongoUserStore) GetUsers(ctx context.Context, filter Map, pag *Pagination) ([]*types.User, error) {
	opts := options.FindOptions{}
	opts.SetSkip(int64(pag.Page-1) * pag.Limit)
	opts.SetLimit(int64(pag.Limit))
	var users []*types.User
	cursor, err := m.coll.Find(ctx, filter, &opts)

	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &users); err != nil { // cursor.Decode(&users) but if nothing to decode will result in error
		return []*types.User{}, nil
	}

	return users, nil
}

func (m *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	_, err := m.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.ErrInvalidID()
	}
	// TODO: check if we deleted anything
	_, err = m.coll.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return err
	}

	return nil
}

func (m *MongoUserStore) UpdateUser(ctx context.Context, id string, params types.UpdateUserParams) error {
	update := bson.M{"$set": params.ToBSON()}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.ErrInvalidID()
	}
	_, err = m.coll.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping user collection ---")
	return m.coll.Drop(ctx)
	//return m.client.Database(m.dbname).Collection(userColl).Drop(ctx)
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	DBNAME := os.Getenv(MONGO_DBNAME_ENV_VARIABLE_NAME)
	coll := client.Database(DBNAME).Collection("users")
	return &MongoUserStore{
		client: client,
		dbname: DBNAME,
		coll:   coll,
	}
}
