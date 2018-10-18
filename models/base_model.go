package models

import (
	"time"

"github.com/globalsign/mgo/bson")

// DocumentInterface Used for common document actions
type DocumentInterface interface {
	After()  // After call
	Before() // Before call
	Update() // Update is calling before document is updated

	Init()
	GetID() bson.ObjectId
	GetDocumentName() string
	GetTimeFormat() string
	SetDocument(doc DocumentInterface)
}

// DocumentBase ...
type DocumentBase struct {
	document DocumentInterface

	ID        bson.ObjectId `json:"-" bson:"_id,omitempty"`
	CreatedAt time.Time     `json:"-" bson:"created_at"`
	UpdatedAt time.Time     `json:"-" bson:"updated_at"`
	Status    string        `json:"-" bson:"status"`
	Deleted   bool          `json:"-" bson:"deleted"`
}

// Init ...
func (d *DocumentBase) Init() {
	d.ID = bson.NewObjectId()
	d.Status = "init"
	d.Deleted = false
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
}

// After ...
func (d *DocumentBase) After() {

}

// Before ...
func (d *DocumentBase) Before() {

}

// GetID ...
func (d *DocumentBase) GetID() bson.ObjectId {
	return d.ID
}

// GetDocumentName ...
func (d *DocumentBase) GetDocumentName() string {
	return ""
}

// SetDocument ...
func (d *DocumentBase) SetDocument(doc DocumentInterface) {
	d.document = doc
}
