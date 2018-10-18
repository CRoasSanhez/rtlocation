package repositories

import(
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	
	"rtlocation/models"
	"rtlocation/core"
)

// LocationHistoryRepository ...
type LocationHistoryRepository interface{
	Insert(*models.LocationHistory)error
	GetByID(bson.ObjectId)(models.LocationHistory,error)
	Update(*models.LocationHistory)error
}

// NewLocationHistoryRepository returns a new user repository,
func NewLocationHistoryRepository()LocationHistoryRepository{
	return &LocationBDHistoryRepository{}
}

// LocationBDHistoryRepository is the DB entity for users
type LocationBDHistoryRepository struct{
	BaseRepository
	History models.LocationHistory
}

// Insert ...
func (r *LocationBDHistoryRepository) Insert(model *models.LocationHistory)error{
	if !bson.ObjectId.Valid(model.ID) {
		model.ID = bson.NewObjectId()
	}
	
	session ,err:= mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("Insert: Error getting session %s",core.DBUrl)
		return err
	}

	_, err = session.DB(core.DBName).C(model.GetDocumentName()).UpsertId(model.ID, model)
	if err != nil{
		fmt.Printf("Insert: Error upsert %s",err.Error())
		return err
	}

	defer session.Close()

	return err
}

// GetByID ...
func (r *LocationBDHistoryRepository) GetByID(id bson.ObjectId)(models.LocationHistory, error){
	var model models.LocationHistory
	session ,err := mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("GetByID error %s",err)
		return model,err
	}

	err = session.DB(core.DBName).C(model.GetDocumentName()).FindId(id).One(&model)
	if err != nil{
		fmt.Printf("FindByID: Error FindOne %v",err)
		return model,err
	}
	defer session.Close()

	return model, nil
}

// Update ...
func (r * LocationBDHistoryRepository)Update(model *models.LocationHistory)error{
	session ,err:= mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("Insert: Error getting session %s",core.DBUrl)
		return err
	}
	
	_, err = session.DB(core.DBName).C(model.GetDocumentName()).UpsertId(model.ID, model)
	if err != nil{
		fmt.Printf("Insert: Error upsert %s",err)
		return err
	}

	defer session.Close()

	return err
}