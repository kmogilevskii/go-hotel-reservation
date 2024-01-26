package db

import (
	"context"
	"os"

	"github.com/kmogilevskii/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, Map, *Pagination) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	coll       *mongo.Collection
	hotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	DBNAME := os.Getenv(MONGO_DBNAME_ENV_VARIABLE_NAME)
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection("rooms"),
		hotelStore: hotelStore,
	}
}

func (m *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	_, err := m.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	filter := Map{"_id": room.HotelID}
	update := Map{"$push": bson.M{"rooms": room.ID}}
	if err := m.hotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}

func (m *MongoRoomStore) GetRooms(ctx context.Context, filter Map, pag *Pagination) ([]*types.Room, error) {
	opts := options.FindOptions{}
	opts.SetSkip(int64(pag.Page-1) * pag.Limit)
	opts.SetLimit(int64(pag.Limit))
	resp, err := m.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
