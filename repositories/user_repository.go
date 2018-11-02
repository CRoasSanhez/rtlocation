package repositories

import(
	"fmt"
	"mime/multipart"

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
	SaveProfilePicture(file *multipart.FileHeader, owner models.User)(error)
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
	user.Status = core.StatusInit
	
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

// SaveProfilePicture ...
func (r *UserBDRepository)SaveProfilePicture(filePart *multipart.FileHeader, user models.User)error{
	var att = models.Attachment{}

	var owner = models.AsDocumentBase(&user)
	if errInit := att.Init(owner,filePart); errInit!=nil{
		fmt.Printf("Error initializing attachment: %s \n",errInit.Error())
		return errInit
	}

	found,errF := r.GetByID(owner.GetID())
	if errF!=nil{
		fmt.Printf("Error User to update: %s \n",errF.Error())
		return errF
	}

	found.ProfilePicture = att

	session ,err:= mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("Update: Error getting session %s\n",core.DBUrl)
		return err
	}
	
	_, err = session.DB(core.DBName).C(found.GetDocumentName()).UpsertId(found.ID, found)
	if err != nil{
		fmt.Printf("Update: Error update %s\n",err)
		return err
	}

	defer session.Close()

	return nil
}