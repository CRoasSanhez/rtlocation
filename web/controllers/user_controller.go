package controllers

import(
	"fmt"
	"io/ioutil"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"

	auth "rtlocation/middleware"
	"rtlocation/services"
	"rtlocation/models"
	"rtlocation/core"
)

// UserController ..
type UserController struct{
	BaseController
	Ctx 		iris.Context
	Service 	services.UserService
	Session 	*sessions.Session
	CurrentUser models.User
}

const userIDKey = "UserID"

// isLoggedIn ... 
func (c *UserController) isLoggedIn() bool {
	if c.GetCurrentUserID()!=""{
		return true
	}
	if user, err := auth.AuthenticateByToken(c.Ctx); err == nil {
		c.CurrentUser = user
		return true
	}

	return false
}

func (c *UserController) logout() {
	c.Session.Destroy()
}

// TODO: move statics views to separeta folder
var registerStaticView = mvc.View{
	Name: "user/register.html",
	Data: iris.Map{"Title": "User Registration"},
}

// GetRegister handles GET: http://localhost:8080/user/register.
// returns the user register VIEW
func (c *UserController) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}

	return registerStaticView
}

// PostRegister handles POST: http://localhost:8080/user/register.
// Creates User
func (c *UserController) PostRegister()(models.BaseResponse){
	var user = models.User{}
	var password string
	if c.Ctx.GetHeader("Content-Type") == "application/json"{
		err := c.Ctx.ReadJSON(&user)
		if err!=nil{
			fmt.Printf("Error parsing JSON - %s",err.Error())
		}
		password = user.Password

	}else{

		password  = c.Ctx.FormValue("password")
		user.Username = c.Ctx.FormValue("username")
		user.Firstname = c.Ctx.FormValue("firstname")
	}

	// create the new user, the password will be hashed by the service.
	newUser, err := c.Service.Create(password, user)
	if err!=nil{
		fmt.Printf("UC Create : Error creating user %v",err)
		return c.BadRequest(err.Error())
	}

	// set the user's id to this session even if err != nil
	c.Session.Set(userIDKey, newUser.ID.Hex())

	// Generate user token for JWT
	token, err := auth.GenerateToken(newUser.ID.Hex(), core.ActionAuth)
	if err != nil {
		fmt.Printf("UC Create : Error generating token %s \n",err.Error())
		return c.InternalErrorResponse(err.Error())
	}

	var response = models.Account{
		UserName: newUser.Username,
		Token: token,
	}
	return c.Successresponse(response,"User created")

}

var loginStaticView = mvc.View{
	Name: "user/login.html",
}

// GetLogin handles GET: http://localhost:8080/user/login.
// returns static view LogIn
func (c *UserController) GetLogin() mvc.Result {
	// verify if users is already logged in
	if c.isLoggedIn() {
		c.logout()
	}

	loginStaticView.Data = iris.Map{
		"Title": "User Login",
		"APIKey": core.MapsHereAPIKey,
		"APPCode": core.MapsHereAPPCode,
	}

	return loginStaticView
}

// PostLogin handles POST: http://localhost:8080/user/register.
// Authenticates User
func (c *UserController) PostLogin() models.BaseResponse {
	
	var user = models.User{}
	var password string
	if c.Ctx.GetHeader("Content-Type") == "application/json"{
		err := c.Ctx.ReadJSON(&user)
		if err!=nil{
			fmt.Printf("Error parsing JSON - %s \n",err.Error())
		}
		password = user.Password

	}else{
		password  = c.Ctx.FormValue("password")
		user.Username = c.Ctx.FormValue("username")
	}

	newUser, found := c.Service.GetByUsernameAndPassword(user.Username, password)
	if !found {
		return c.BadRequest("User not found")
	}

	// Generate user token for JWT
	token, err := auth.GenerateToken(newUser.ID.Hex(), core.ActionAuth)
	if err != nil {
		fmt.Printf("UC Create : Error Generating token %s \n",err.Error())
		return c.InternalErrorResponse(err.Error())
	}

	c.Session.Set(userIDKey, newUser.ID.Hex())

	var response = models.Account{
		UserName: newUser.Username,
		Token: token,
	}

	return c.Successresponse(response,"User found")
}

// GetCurrentUserID returns userID from session
func (c *UserController) GetCurrentUserID() string {
	userID := c.Session.GetString(userIDKey)
	//userID := c.Session.GetInt64Default(userIDKey, 0)
	return userID
}

// GetMe handles GET: http://localhost:8080/user/me.
func (c *UserController) GetMe() mvc.Result {
	if !c.isLoggedIn() {
		return mvc.Response{Path: "/user/login"}
	}

	u, found := c.Service.GetByID(c.GetCurrentUserID())
	if !found {
		c.logout()
		return c.GetMe()
	}

	fmt.Printf("user: %v \n",c.CurrentUser)

	return mvc.View{
		Name: "user/me.html",
		Data: iris.Map{
			"PusherKePusherKey": core.PusherKey,
			"Title": "Profile of " + u.Username,
			"userLat":  u.Geolocation.Coordinates[1],
			"userLon": u.Geolocation.Coordinates[0],
		},
	}
}

// PostPicture updates user profile picture
func (c *UserController) PostPicture()models.BaseResponse{

	if !c.isLoggedIn(){
		c.logout()
	}
	
	// Get the max post value size passed via iris.WithPostMaxMemory.
	maxSize := c.Ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()

	err := c.Ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_, fileHeader, errF := c.Ctx.FormFile("picture")
	if errF!=nil{
		return c.BadRequest(err.Error())
	}

	_, ok := c.Service.SaveProfilePicture(fileHeader, c.CurrentUser)
	if !ok {
		return c.InternalErrorResponse("Error updating User profile")
	}

	return c.Successresponse("Success","File uploaded")
}

// GetPicture returns user profile picture
func (c *UserController) GetPicture()models.BaseResponse{
	if !c.isLoggedIn(){
		c.logout()
	}

	var id = c.CurrentUser.ID.Hex()

	var path = "/tmp/" + id + ".png"

	fmt.Println(c.CurrentUser.ProfilePicture.Binary)

	if err := ioutil.WriteFile(path, c.CurrentUser.ProfilePicture.Binary, 0644); err != nil {
		return c.InternalErrorResponse(err.Error())
	}

	c.Ctx.SendFile(path,"profile_picture")

	return c.Successresponse("success","File downloaded")

}

/*
// AnyLogout handles All/Any HTTP Methods for: http://localhost:8080/user/logout.
func (c *UserController) AnyLogout() {
	if c.isLoggedIn() {
		c.logout()
	}
	c.Ctx.Redirect("/user/login")
}
*/