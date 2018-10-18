package services

import (
	"github.com/globalsign/mgo/bson"
	"fmt"

	"rtlocation/repositories"
	"rtlocation/models"
)

// VehicleService ...
type VehicleService interface{
	Create(vehicle models.Vehicle)(models.Vehicle,error)
	GetByID(id string)(models.Vehicle,bool)
	UpdateCoords(id string, model models.Vehicle)(models.Vehicle, error)
	GetNearVehicles(lon,lat float64)([]models.Vehicle,error)
}

// NewVehicleService returns the default user service.
func NewVehicleService(repo repositories.VehicleRepository) VehicleService {
	return &mVehicleService{
		repo: repo,
	}
}

// UserService ...
type mVehicleService struct {
	repo repositories.VehicleRepository
}


// Create inserts a new User
func (s *mVehicleService) Create(model models.Vehicle) (models.Vehicle, error) {
	err := s.repo.Insert(&model)
	if err !=nil{
		return model,err
	}
	return model,nil
}

// GetByID ...
func(s *mVehicleService) GetByID(idHex string)(models.Vehicle,bool){
	if(!bson.IsObjectIdHex(idHex)){
		fmt.Printf("GetByID error invalid %s",idHex)
		return models.Vehicle{},false
	}
	model, err:= s.repo.GetByID(bson.ObjectIdHex(idHex))
	if err != nil{
		return models.Vehicle{},false
	}
	return model,true
}

// GetNearVehicles ...
func(s *mVehicleService)GetNearVehicles(lon,lat float64)([]models.Vehicle,error){
	vehicles,err := s.repo.GetNearVehicles(lon,lat)
	if err!=nil{
	return []models.Vehicle{},err
	}
	return vehicles,nil
}

// UpdateCoords ...
func(s *mVehicleService)UpdateCoords(id string, model models.Vehicle)(models.Vehicle,error){
	if !bson.IsObjectIdHex(id){
		fmt.Printf("Update error invalid %s",id) 
	}

	err := s.repo.UpdateCoords(&model)
	if err!=nil{
		return models.Vehicle{},err
	}
	return model,nil
}