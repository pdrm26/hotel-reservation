package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (p CreateUserParams) Validate() error {
	if len(p.FirstName) < minFirstNameLen {
		return fmt.Errorf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(p.LastName) < minLastNameLen {
		return fmt.Errorf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(p.Password) < minPasswordLen {
		return fmt.Errorf("password length should be at least %d characters", minPasswordLen)
	}
	if !isValidEmail(p.Email) {
		return fmt.Errorf("email is invalid")
	}
	return nil
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil

}
