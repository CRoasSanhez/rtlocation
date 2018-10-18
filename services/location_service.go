package services

import (
	"github.com/globalsign/mgo/bson"
	"fmt"

	"rtlocation/repositories"
	"rtlocation/models"
)

// LocationHistoryService ...
type LocationHistoryService interface{
	Create(vehicle models.LocationHistory)(models.LocationHistory,error)
	GetByID(id string)(models.LocationHistory,bool)
	Update(id string, model models.LocationHistory)(models.LocationHistory, error)
}

// NewLocationHistoryService returns the default user service.
func NewLocationHistoryService(repo repositories.LocationHistoryRepository) LocationHistoryService {
	return &mLocationHistoryService{
		repo: repo,
	}
}

// mLocationHistoryService ...
type mLocationHistoryService struct {
	repo repositories.LocationHistoryRepository
}


// Create inserts a new LocationHistory
func (s *mLocationHistoryService) Create(model models.LocationHistory) (models.LocationHistory, error) {
	err := s.repo.Insert(&model)
	if err !=nil{
		fmt.Printf("Service Insert: Error inserting %v",err)
		return model,err
	}
	return model,nil
}

// GetByID ...
func(s *mLocationHistoryService) GetByID(idHex string)(models.LocationHistory,bool){
	if(!bson.IsObjectIdHex(idHex)){
		fmt.Printf("GetByID error invalid %s",idHex)
		return models.LocationHistory{},false
	}
	model, err:= s.repo.GetByID(bson.ObjectIdHex(idHex))
	if err != nil{
		return models.LocationHistory{},false
	}
	return model,true
}

// Update ...
func(s *mLocationHistoryService)Update(id string, model models.LocationHistory)(models.LocationHistory,error){
	if !bson.IsObjectIdHex(id){
		fmt.Printf("Update error invalid %s",id) 
	}
	err := s.repo.Update(&model)
	if err!=nil{
		fmt.Printf("Update error %v",err) 
		return models.LocationHistory{},err
	}
	return model,nil
}