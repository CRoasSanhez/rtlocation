package repositories

import(
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	
	"rtlocation/models"
	"rtlocation/core"
)

// VehicleRepository ...
type VehicleRepository interface{
	Insert(*models.Vehicle)error
	GetByID(bson.ObjectId)(models.Vehicle,error)
	UpdateCoords(*models.Vehicle)error
	GetNearVehicles(lon,lat float64)([]models.Vehicle,error)
}

// NewVehicleRepository returns a new user repository,
func NewVehicleRepository()VehicleRepository{
	return &VehicleBDRepository{}
}

// VehicleBDRepository is the DB entity for users
type VehicleBDRepository struct{
	BaseRepository
	Vehicle models.Vehicle
}

// Insert ...
func (r *VehicleBDRepository) Insert(model *models.Vehicle)error{
	if !bson.ObjectId.Valid(model.ID) {
		model.ID = bson.NewObjectId()
	}

	model.Geolocation.Type = "Point"
	if len(model.Geolocation.Coordinates) != 2{
		model.Geolocation.Coordinates = []float64{-99.16378, 19.410272}
	}
	
	session ,err:= mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("Insert: Error getting session %s \n",core.DBUrl)
		return err
	}

	_, err = session.DB(core.DBName).C("vehicles").UpsertId(model.ID, model)
	if err != nil{
		fmt.Printf("Insert: Error upsert %s \n",err.Error())
		return err
	}

	if err = SaveVehicleLocationHistory(model);err != nil{
		fmt.Printf("Save Vehicle History: Error --- %s \n", err.Error())
	}

	defer session.Close()

	return err
}

// GetByID ...
func (r *VehicleBDRepository) GetByID(id bson.ObjectId)(models.Vehicle, error){
	var model models.Vehicle
	session ,err := mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("GetByID error %s \n",err)
		return model,err
	}

	err = session.DB(core.DBName).C(model.GetDocumentName()).FindId(id).One(&model)
	if err != nil{
		fmt.Printf("FindByID: Error FindOne %s \n",err.Error())
		return model,err
	}
	defer session.Close()

	return model, nil
}

// GetNearVehicles ...
func (r *VehicleBDRepository) GetNearVehicles(lon,lat float64)([]models.Vehicle,error){

	var response []models.Vehicle
	var m models.Vehicle
	var query = bson.M{
					"geolocation": bson.M{
						"$nearSphere": bson.M{
								"$geometry": bson.M{ 
											"type": "Point", 
											"coordinates": []float64{lon, lat},
								},
								"$maxDistance" : core.MaxDistanceFromDriver,
							},
						},
					}

	session ,err := mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("GetNearVehicles error session %s \n",err)
		return response,err
	}

	err = session.DB(core.DBName).C(m.GetDocumentName()).Find(query).All(&response)
	if err != nil{
		fmt.Printf("GetNearVehicles: Error Pipe %s \n",err.Error())
		return response,err
	}

	defer session.Close()

	return response,nil
}

// UpdateCoords ...
func (r * VehicleBDRepository)UpdateCoords(model *models.Vehicle)error{

	found,errF := r.GetByID(bson.ObjectId(model.ID))
	if errF!=nil{
		return errF
	}

	found.Geolocation = model.Geolocation
	model.DriverName = found.DriverName
	model.Mark = found.Mark
	model.CarType = found.CarType

	session ,err:= mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("Update: Error getting session %s",core.DBUrl)
		return err
	}
	
	_, err = session.DB(core.DBName).C(found.GetDocumentName()).UpsertId(model.ID, found)
	if err != nil{
		fmt.Printf("Update: Error update %s",err)
		return err
	}

	if err = SaveVehicleLocationHistory(&found);err != nil{
		fmt.Printf("Save updated Vehicle History: Error --- %s \n", err.Error())
	}

	defer session.Close()

	return nil
}

// SaveVehicleLocationHistory ...
func SaveVehicleLocationHistory(model *models.Vehicle)error{
	var history = models.LocationHistory{
		LocationType: core.LocationTypeVehicle,
		Parent: model.ID,
		Geolocation: model.Geolocation,
	}

	return history.Insert()
}