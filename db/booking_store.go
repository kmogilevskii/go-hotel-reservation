package db

import (
	"context"
	"os"

	"github.com/kmogilevskii/hotel-reservation/errors"
	"github.com/kmogilevskii/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookingStore interface {
	Insert(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, Map, *Pagination) ([]*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	DBNAME := os.Getenv(MONGO_DBNAME_ENV_VARIABLE_NAME)
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(DBNAME).Collection("bookings"),
	}
}

func (m *MongoBookingStore) Insert(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	_, err := m.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (m *MongoBookingStore) GetBookings(ctx context.Context, filter Map, pag *Pagination) ([]*types.Booking, error) {
	if filter["userID"] != nil {
		oid, err := primitive.ObjectIDFromHex(filter["userID"].(string))
		if err != nil {
			return nil, errors.ErrInvalidID()
		}
		filter["userID"] = oid
	}
	opts := options.FindOptions{}
	opts.SetSkip(int64(pag.Page-1) * pag.Limit)
	opts.SetLimit(int64(pag.Limit))
	var bookings []*types.Booking
	cursor, err := m.coll.Find(ctx, filter, &opts)

	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &bookings); err != nil {
		return []*types.Booking{}, nil
	}

	return bookings, nil
}

func (m *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	var booking types.Booking
	if err := m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (m *MongoBookingStore) UpdateBooking(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"canceled": true}}
	_, err = m.coll.UpdateOne(ctx, filter, update)
	// _, err = m.coll.UpdateByID(ctx, oid, update)
	return err
}
