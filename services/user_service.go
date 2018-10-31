package services

import (
	"mime/multipart"
	"errors"
	"fmt"

	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"


	"rtlocation/repositories"
	"rtlocation/models"
)

// UserService ...
type UserService interface{
	Create(string,models.User)(models.User,error)
	GetByID(string)(models.User,bool)
	GetByUsernameAndPassword(username,pwd string)(models.User,bool)
	SaveProfilePicture(file *multipart.FileHeader, owner models.User)(models.User,bool)
}

// NewUserService returns the default user service.
func NewUserService(repo repositories.UserRepository) UserService {
	return &mUserService{
		repo: repo,
	}
}

// UserService ...
type mUserService struct {
	repo repositories.UserRepository
}


// Create inserts a new User
func (s *mUserService) Create(userPassword string, user models.User) (models.User, error) {
	if userPassword == "" || user.Firstname == "" || user.Username == "" {
		return models.User{}, errors.New("unable to create this user")
	}

	/*
	hashed, err := models.GeneratePassword(userPassword)
	if err != nil {
		return models.User{}, err
	}
	

	user.HashedPassword = hashed
	*/

	err := user.GeneratePassword();
	if err !=nil{
		return models.User{},err
	}

	err = s.repo.Insert(&user)
	if err !=nil{
		fmt.Printf("Service Insert: Error inserting %v",err)
		return user,err
	}
	return user,nil

	//return s.repo.InsertOrUpdate(user)
}

// GetByID ...
func(s *mUserService) GetByID(idHex string)(models.User,bool){
	if(!bson.IsObjectIdHex(idHex)){
		fmt.Printf("GetByID error invalid %s",idHex)
		return models.User{},false
	}
	user, err:= s.repo.GetByID(bson.ObjectIdHex(idHex))
	if err != nil{
		return models.User{},false
	}
	return user,true
}

// GetByUsernameAndPassword ...
func (s *mUserService)GetByUsernameAndPassword(username,pwd string)(models.User, bool){
	if( username == "" && pwd == ""){
		fmt.Printf("GetByUsrname&pwd error invalid %s - %s",username,pwd)
		return models.User{},false
	}
	
	user, err:= s.repo.GetByUsername(username,pwd)
	if err != nil{
		return models.User{},false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(pwd)); err != nil {
		return models.User{},false
	}
	
	return user,true
}

// SaveProfilePicture ...
func (s *mUserService)SaveProfilePicture(file *multipart.FileHeader,owner models.User)(models.User,bool){

	s.repo.SaveProfilePicture(file, owner)
	return models.User{},true
}