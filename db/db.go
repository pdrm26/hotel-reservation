package db

const DBURI = "mongodb://localhost:27017"
const DBNAME = "hotel-reservation"
const DBNAMETEST = "hotel-reservation-test"

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
