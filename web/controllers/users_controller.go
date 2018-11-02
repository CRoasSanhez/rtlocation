package controllers

import(
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"

	"rtlocation/models"
	"rtlocation/services"
	"rtlocation/models/serializers"
	auth "rtlocation/middleware"
)

// UsersController ..
type UsersController struct{
	BaseController
	Ctx 		iris.Context
	Service 	services.UsersService
	Session 	*sessions.Session
	CurrentUser models.User
}

// GetCurrentUserID returns userID from session
func (c *UsersController) GetCurrentUserID() string {
	userID := c.Session.GetString(userIDKey)
	//userID := c.Session.GetInt64Default(userIDKey, 0)
	return userID
}

// isLoggedIn ... 
func (c *UsersController) isLoggedIn() bool {
	if c.GetCurrentUserID()!=""{
		return true
	}
	if user, err := auth.AuthenticateByToken(c.Ctx); err == nil {
		c.CurrentUser = user
		return true
	}

	return false
}

func (c *UsersController) logout() {
	c.Session.Destroy()
}

// GetAll returns a list of registered users
func (c *UsersController)GetAll()models.BaseResponse{
	if !c.isLoggedIn(){
		c.logout()
		return c.ForbiddenResponse()
	}

	users, err := c.Service.GetAll()
	if err != nil{
		return c.InternalErrorResponse(err.Error())
	}

	return c.Successresponse(users,"sucess",serializers.UserSerializer{})
}