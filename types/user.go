package types

import (
	"context"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
	// minLastNameLen  = 2
	// minFirstNameLen = 2
	// minPasswordLen  = 6
)

type CreateUserParams struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=100"`
	LastName  string `json:"lastName" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=7"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=100"`
	LastName  string `json:"lastName" validate:"required,min=2,max=100"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.FirstName) != 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) != 0 {
		m["lastName"] = p.LastName
	}
	return m
}

func (p *CreateUserParams) Validate(ctx context.Context) error {
	validate := validator.New()
	if err := validate.StructCtx(ctx, p); err != nil {
		return err
	}
	return nil
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil

}

type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"isAdmin" json:"isAdmin"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)

	if err != nil {
		return nil, err
	}

	return &User{
		ID:                primitive.NewObjectID(),
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
