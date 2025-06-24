package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomID    primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	UserID    primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	NumGuests int                `bson:"numGuests,omitempty" json:"numGuests,omitempty"`
	StartDate time.Time          `bson:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate   time.Time          `bson:"endDate,omitempty" json:"endDate,omitempty"`
}
