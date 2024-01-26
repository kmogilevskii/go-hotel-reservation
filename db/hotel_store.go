package db

import (
	"context"
	"os"

	"github.com/kmogilevskii/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStore interface {
	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, Map, Map) error
	GetHotels(context.Context, Map, *Pagination) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	DBNAME := os.Getenv(MONGO_DBNAME_ENV_VARIABLE_NAME)
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection("hotels"),
	}
}

func (m *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	_, err := m.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	return hotel, nil
}

func (m *MongoHotelStore) Update(ctx context.Context, filter, update Map) error {
	_, err := m.coll.UpdateOne(ctx, filter, update)
	return err
}

// do we even need filter if function should return all?
func (m *MongoHotelStore) GetHotels(ctx context.Context, filter Map, pag *Pagination) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip(int64(pag.Page-1) * pag.Limit)
	opts.SetLimit(int64(pag.Limit))
	var hotels []*types.Hotel
	cursor, err := m.coll.Find(ctx, filter, &opts)

	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &hotels); err != nil { // cursor.Decode(&users) but if nothing to decode will result in error
		return []*types.Hotel{}, nil
	}

	return hotels, nil
}

func (m *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	if err := m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
}
