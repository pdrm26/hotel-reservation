package db

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
