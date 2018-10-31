package models

import(
	"errors"

	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

// User ..
type User struct {
	DocumentBase `json:",inline" bson:",inline"`
	
	ID				bson.ObjectId `json:"id" form:"id" bson:"id"`
	Firstname      	string    	`json:"firstname" form:"firstname" bson:"first_name"`
	Username       	string    	`json:"username" form:"username" bson:"user_name"`
	Gender			string		`json:"gender" form:"gender" bson:"gender"`
	Country			string		`json:"country" form:"country" bson:"country"`
	Device			Device		`json:"-" form:"-" bson:"device"`
	Password 		string		`json:"password" form:"password" bson:"-"`
	HashedPassword 	string    	`json:"-" form:"-" bson:"hashed_password"`
	Geolocation 	Geo			`json:"geolocation" form:"geolocation" bson:"geolocation"`
	ProfilePicture	Attachment	`json:"profile_picture" form:"profile_picture" bson:"profile_picture"`
	Payment			Payment		`json:"-" form:"-" bson:"payment"`
}

// GetDocumentName returns the DB collection name
func (u *User) GetDocumentName()string{
	return "users"
}

// IsValid ...
func (u User) IsValid() bool {
	return true
}

// GeneratePassword ...
func (u *User) GeneratePassword() error {
	if u.Password == "" {
		return errors.New("Invalid password")
	}

	hash, err := MD5Crypt(u.Password)
	if err != nil {
		return err
	}

	u.HashedPassword = string(hash)

	return nil
}

// GeneratePassword generates a hashed password
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword validates if passwords match
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}

// MD5Crypt ...
func MD5Crypt(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}