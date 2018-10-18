package models

import (
	"github.com/globalsign/mgo/bson"
)

// Vehicle ...
type Vehicle struct{
	DocumentInterface `json:"-" bson:"-"`

	ID 			bson.ObjectId 	`json:"id" form:"id" bson:"id"`
	Mark 		string 			`json:"mark" form:"mark" bson:"mark"`
	DriverName 	string 			`json:"driver_name" form:"drivername" bson:"driver_name"`
	CarType 	string 			`json:"car_type" form:"cartype" bson:"car_type"`	
	Geolocation Geo 			`json:"geolocation" form:"geolocation" bson:"geolocation"`
}

// GetDocumentName returns document name in DB
func(m *Vehicle)GetDocumentName()string{
	return "vehicles"
}