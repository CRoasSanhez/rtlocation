package repositories

import(
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"rtlocation/core"
	"rtlocation/models"
)

// UsersRepository ...
type UsersRepository interface{
	GetAllUsers()([]models.User,error)
	GetActiveUsers()([]models.User,error)

}

// NewUsersRepository returns a new user repository,
func NewUsersRepository()UsersRepository{
	return &UsersBDRepository{}
}

// UsersBDRepository is the DB entity for users
type UsersBDRepository struct{
	BaseRepository
}

// GetActiveUsers returns a list of Active users
func (r *UsersBDRepository)GetActiveUsers()(response []models.User,err error){
	return []models.User{},nil
}

// GetAllUsers returns a list of all users in DB
func (r *UsersBDRepository)GetAllUsers()(response []models.User, err error){
	var user models.User
	session ,err := mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("GetAllUsers error %s",err)
		return 
	}

	err = session.DB(core.DBName).C(user.GetDocumentName()).Find(bson.M{ "status":core.StatusActive }).All(&response)
	if err != nil{
		fmt.Printf("GetAllUsers: Error FindAll %v",err)
		return 
	}
	defer session.Close()

	return 
}