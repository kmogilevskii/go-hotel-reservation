package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kmogilevskii/hotel-reservation/api"
	"github.com/kmogilevskii/hotel-reservation/db"
	"github.com/kmogilevskii/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	MONGO_URI := os.Getenv("MONGO_URI")
	DBNAME := os.Getenv("MONGO_DBNAME")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}

func main() {
	admin := fixtures.AddUser(userStore, "admin", "admin", true)
	fmt.Printf("%s -> %s\n", admin.Email, api.CreateTokenFromUser(admin))
	user := fixtures.AddUser(userStore, "Jhon", "Doe", false)
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	hotel := fixtures.AddHotel(hotelStore, "Hilton", "New York", 5)
	room := fixtures.AddRoom(roomStore, hotel.ID, "Single", 99.99)
	fixtures.AddBooking(bookingStore, user.ID, room.ID, time.Now().AddDate(0, 0, 2), time.Now().AddDate(0, 0, 3), 2)
	fixtures.AddBooking(bookingStore, user.ID, room.ID, time.Now().AddDate(0, 0, 4), time.Now().AddDate(0, 0, 5), 2)
	fixtures.AddBooking(bookingStore, user.ID, room.ID, time.Now().AddDate(0, 0, 6), time.Now().AddDate(0, 0, 7), 2)
	fixtures.AddBooking(bookingStore, user.ID, room.ID, time.Now().AddDate(0, 0, 8), time.Now().AddDate(0, 0, 9), 2)
	fixtures.AddBooking(bookingStore, user.ID, room.ID, time.Now().AddDate(0, 0, 9), time.Now().AddDate(0, 0, 10), 2)
	roomSizes := []string{"Single", "Double", "Triple", "Quad", "Queen", "King", "Twin", "Double-double", "Studio", "Master suite", "Mini suite", "President suite"}
	firstNames := []string{"Jhon", "Jane", "Jack", "Jill", "Kate", "Kevin", "Katie", "Karl", "Kris", "Kurt", "Kim", "Kane"}
	lastNames := []string{"Doe", "Black", "White", "Smith", "Johnson", "Williams", "Jones", "Brown", "Davis", "Miller", "Wilson", "Moore"}
	for i := 0; i < 100; i++ {
		fixtures.AddHotel(hotelStore, fmt.Sprintf("Hotel %d", i), fmt.Sprintf("City %d", i), rand.Intn(5)+1)
		fixtures.AddRoom(roomStore, hotel.ID, roomSizes[rand.Intn(len(roomSizes)-1)], 50.+rand.Float64()*100.)
		fixtures.AddUser(userStore, firstNames[rand.Intn(len(firstNames)-1)], lastNames[rand.Intn(len(lastNames)-1)], false)
	}
}
