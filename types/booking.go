package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RoomID     primitive.ObjectID `json:"roomID,omitempty" bson:"roomID,omitempty"`
	UserID     primitive.ObjectID `json:"userID,omitempty" bson:"userID,omitempty"`
	FromDate   time.Time          `json:"fromDate,omitempty" bson:"fromDate,omitempty"`
	TillDate   time.Time          `json:"tillDate,omitempty" bson:"tillDate,omitempty"`
	NumPersons int                `json:"numPersons,omitempty" bson:"numPersons,omitempty"`
	Canceled   bool               `json:"canceled,omitempty" bson:"canceled,omitempty"`
}

type BookParams struct {
	FromDate   time.Time `json:"fromDate,omitempty" bson:"fromDate,omitempty" validate:"required"`
	TillDate   time.Time `json:"tillDate,omitempty" bson:"tillDate,omitempty" validate:"required"`
	NumPersons int       `json:"numPersons,omitempty" bson:"numPersons,omitempty" validate:"required,gt=0"`
}

func (p *BookParams) Validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book in the past")
	}
	if p.FromDate.After(p.TillDate) {
		return fmt.Errorf("from date should be before till date")
	}
	return nil
}
