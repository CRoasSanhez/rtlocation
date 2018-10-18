package models

import(
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/globalsign/mgo"

	"rtlocation/core"
)

// LocationHistory ...
type LocationHistory struct{
	DocumentInterface 					`json:"-" bson:"-"`
	ID					bson.ObjectId	`json:"id" form:"id" bson:"id"`
	LocationType		string 			`json:"type" form:"type" bson:"location_type"`
	Parent 				bson.ObjectId	`json:"parent" form:"parent" bson:"parent"`
	Geolocation 		Geo 			`json:"geolocation" form:"geolocation" bson:"geolocation"`
	
}

// GetDocumentName returns the DB collection name
func (m *LocationHistory) GetDocumentName()string{
	return "location_history"
}

// Insert ...
func(m *LocationHistory)Insert()error{

	m.ID = bson.NewObjectId()
	
	session ,err:= mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("Insert: Error getting session %s",core.DBUrl)
		return err
	}

	_, err = session.Clone().DB(core.DBName).C(m.GetDocumentName()).UpsertId(m.ID, m)
	if err != nil{
		fmt.Printf("Insert: Error upsert %s",err.Error())
		return err
	}

	defer session.Close()

	return err
}