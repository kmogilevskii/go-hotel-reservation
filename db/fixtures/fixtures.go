package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	//"github.com/kmogilevskii/hotel-reservation/api"
	"github.com/kmogilevskii/hotel-reservation/db"
	"github.com/kmogilevskii/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(s db.UserStore, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := s.CreateUser(context.TODO(), user)

	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddHotel(s db.HotelStore, name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		ID:       primitive.NewObjectID(),
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := s.Insert(context.TODO(), &hotel)

	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(s db.RoomStore, hotelID primitive.ObjectID, size string, price float64) *types.Room {
	room := types.Room{
		ID:      primitive.NewObjectID(),
		Size:    size,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := s.InsertRoom(context.TODO(), &room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(s db.BookingStore, userID, roomID primitive.ObjectID, fromDate, tillDate time.Time, numPersons int) *types.Booking {
	booking := types.Booking{
		ID:         primitive.NewObjectID(),
		RoomID:     roomID,
		UserID:     userID,
		FromDate:   fromDate,
		TillDate:   tillDate,
		NumPersons: numPersons,
		Canceled:   false,
	}

	insertedBooking, err := s.Insert(context.TODO(), &booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}
