package repositories

import(
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	
	"rtlocation/models"
	"rtlocation/core"
)

// UserRepository ...
type UserRepository interface{
	Insert(*models.User)error
	GetByID(bson.ObjectId)(models.User,error)
	GetByUsername(username, pwd string)(models.User,error)
}

// NewUserRepository returns a new user repository,
func NewUserRepository(source models.User)UserRepository{
	return &UserBDRepository{}
}

// UserBDRepository is the DB entity for users
type UserBDRepository struct{
	BaseRepository
	User models.User
}

// Insert ...
func (r *UserBDRepository) Insert(user *models.User)error{
	if len(user.ID) == 0 {
		user.ID = bson.NewObjectId()
	}
	
	session ,err:= mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("Insert: Error getting session %s",core.DBUrl)
		return err
	}
	
	_, err = session.DB(core.DBName).C(user.GetDocumentName()).UpsertId(user.ID, user)
	if err != nil{
		fmt.Printf("Insert: Error upsert %s",err)
		return err
	}

	defer session.Close()

	return err
}

// GetByID ...
func (r *UserBDRepository) GetByID(id bson.ObjectId)(models.User, error){
	var user models.User
	session ,err := mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("GetByID error %s",err)
		return user,err
	}

	err = session.DB(core.DBName).C(user.GetDocumentName()).FindId(id).One(&user)
	if err != nil{
		fmt.Printf("FindByID: Error FindOne %v",err)
		return user,err
	}
	defer session.Close()

	return user, nil
}

// GetByUsername ...
func (r *UserBDRepository) GetByUsername(username, pwd string)(models.User,error){
	var user models.User
	session ,err := mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("GetByUsername error %s",err)
		return user,err
	}

	err = session.DB(core.DBName).C(user.GetDocumentName()).Find(bson.M{ "user_name":username }).One(&user)
	if err != nil{
		fmt.Printf("GetByUsername: Error FindOne %v",err)
		return user,err
	}
	defer session.Close()

	return user, nil
}