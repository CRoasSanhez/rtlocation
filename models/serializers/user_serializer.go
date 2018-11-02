package serializers

import (

	"rtlocation/models"
)

// UserSerializer ...
type UserSerializer struct {
	BaseSerializer			`json:"-"`
	ID          string      `json:"id"`
	ImageURL    string      `json:"image_url"`
	Geolocation models.Geo  `json:"geolocation"`
	UserName 	string 		`json:"user_name"`
}

// Cast ...
func (s UserSerializer) Cast(data interface{}) interface{} {
	serializer := new(UserSerializer)

	if model, ok := data.(models.User); ok {
		serializer.ID = model.GetID().Hex()
		serializer.UserName = model.Username
		serializer.ImageURL = model.ProfilePicture.URL
		serializer.Geolocation = model.Geolocation
	}

	return serializer
}

