package services

import(
	"github.com/globalsign/mgo/bson"

	"rtlocation/models"
	"rtlocation/repositories"
)

// UsersService ...
type UsersService interface{
	GetFriends(userID bson.ObjectId, friendship string)([]models.User,error)
	GetAll()([]models.User,error)
}

// NewUsersService returns the default users service.
func NewUsersService(repo repositories.UsersRepository) UsersService {
	return &mUsersService{
		repo: repo,
	}
}

// mUsersService ...
type mUsersService struct {
	repo repositories.UsersRepository
}

// GetFriends returns a list of user friends based on the friendship
func (s *mUsersService)GetFriends(userID bson.ObjectId, friendship string)([]models.User,error){


	return []models.User{},nil
}

func (s *mUsersService)GetAll()([]models.User,error){
	return s.repo.GetAllUsers()
}