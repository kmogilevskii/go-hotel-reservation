package api

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kmogilevskii/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	store  *db.Store
	client *mongo.Client
}

func (tdb *testdb) teardown(t *testing.T) {
	DBNAME := os.Getenv(db.MONGO_DBNAME_ENV_VARIABLE_NAME)
	if err := tdb.client.Database(DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	MONGO_URI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))

	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	userStore := db.NewMongoUserStore(client)
	bookingStore := db.NewMongoBookingStore(client)

	return &testdb{
		store: &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		},
		client: client,
	}
}
