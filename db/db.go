package db

var MONGO_DBNAME_ENV_VARIABLE_NAME = "MONGO_DBNAME"

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type Pagination struct {
	Page  int64
	Limit int64
}

type HotelQueryParams struct {
	Pagination
	Rating int //`query:"rating"`
}

type RoomQueryParams struct {
	Pagination
	Size string
}

type UserQueryParams struct {
	Pagination
	FirstName string
}

type BookingQueryParams struct {
	Pagination
	UserID string
}

type ResourceResponse struct {
	Results int   `json:"results"`
	Data    any   `json:"data"`
	Page    int64 `json:"page"`
}
