package db

const DBURI = "mongodb://localhost:27017"
const DBNAME = "hotel-reservation"
const TestDBName = "hotel-reservation-test"

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
