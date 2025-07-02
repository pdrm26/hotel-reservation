package db

import "os"

var DBNAME string
var DBURL string

func init() {
	DBNAME = os.Getenv("MONGO_DB_NAME")
	DBURL = os.Getenv("MONGO_DB_URL")
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type PaginateFilter struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}
