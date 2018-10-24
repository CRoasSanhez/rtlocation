package controllers

import(
	"fmt"
	"strconv"

	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/pusher/pusher-http-go"

	"rtlocation/services"
	"rtlocation/models"
	"rtlocation/core"
	auth "rtlocation/middleware"
)

// VehicleController ...
type VehicleController struct{
	BaseController
	Ctx iris.Context
	Service services.VehicleService
	Session *sessions.Session
	CurrentUser models.User
}

// BeforeActivation ...
func(c *VehicleController)BeforeActivation(b *mvc.BeforeActivation){

}

// isLoggedIn ... 
func (c *VehicleController) isLoggedIn() bool {
	if user, err := auth.AuthenticateByToken(c.Ctx); err == nil {
		c.CurrentUser = user
		return true
	}

	return false
}

func (c *VehicleController) logout() {
	c.Session.Destroy()
}

// TODO: move statics views to separeta folder
var vehiclesStaticView = mvc.View{
	Name: "vehicle/index.html",
}

var registerVehicleStaticView = mvc.View{
	Name: "vehicle/create.html",
	Data: iris.Map{"Title": "Create vehicle"},
}

// GetIndex handles GET: http://localhost:1337/vehicle/index.
func (c *VehicleController) GetIndex() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}
	vehiclesStaticView.Data= iris.Map{
		"Title": "Vehicles",
		"UserName": c.CurrentUser.Username,
		"PusherKePusherKey": core.PusherKey,
	}

	return vehiclesStaticView
}

// GetVehicles returns a list of vehicles according to the given coordinates
func (c *VehicleController)GetVehicles()models.BaseResponse{
	if !c.isLoggedIn(){
		c.logout()
	}

	lon,errLon := strconv.ParseFloat(c.Ctx.URLParam("lon"),64)
	lat,errLat := strconv.ParseFloat(c.Ctx.URLParam("lat"),64)
	if errLon != nil{ 
		fmt.Printf("Vehicle: Error parsing longitude: %s \n",errLon.Error())
		lon =0
	}
	if errLat!=nil{
		fmt.Printf("Vechicle : Error parsing latitude - %s \n",errLat.Error())
		lat=0
	}

	vehicles, err := c.Service.GetNearVehicles(lon,lat)
	if err != nil{
		return c.InternalErrorResponse(err.Error())
	}
	return c.Successresponse(vehicles,"Vehicles Obtained")
}

// GetRegister ...
func (c *VehicleController) GetRegister()mvc.Result{
	// verify if users is already logged in
	if c.isLoggedIn() {
		c.logout()
	}

	return registerVehicleStaticView
}

// PostRegister handles POST: http://localhost:1337/vehicle/register.
// Creates Vehicle
func (c *VehicleController) PostRegister() mvc.Result {
	var vehicle = models.Vehicle{}
	if c.Ctx.GetHeader("Content-Type") == "application/json"{
		c.Ctx.ReadJSON(&vehicle)
	}else{
		var (
			mark = c.Ctx.FormValue("mark")
			carType  = c.Ctx.FormValue("cartype")
			driverName  = c.Ctx.FormValue("driver")
		)

		lon,errLon := strconv.ParseFloat(c.Ctx.FormValue("lon"),64)
		lat,errLat := strconv.ParseFloat(c.Ctx.FormValue("lat"),64)
		if errLon != nil{ 
			fmt.Printf("Vehicle: Error parsing longitude: %s \n",errLon.Error())
			lon =0
		}
		if errLat!=nil{
			fmt.Printf("Vechicle : Error parsing latitude - %s \n",errLat.Error())
			lat=0
		}
		vehicle.CarType= carType
		vehicle.DriverName= driverName
		vehicle.Mark=mark
		vehicle.Geolocation= models.Geo{Type: "Point", Coordinates: []float64{lon,lat}}
	}

	// create the new vehicle
	newVehicle, err := c.Service.Create(vehicle)
	if err!=nil{
		fmt.Printf("Vehicle Create : Error creating vehicle %v",err)
	}

	// Pusher 
	client := pusher.Client{
		AppId: core.PusherAppID,
		Key: core.PusherKey,
		Secret: core.PusherSecret,
		Cluster: core.PusherCluster,
		Secure: true,
	}

	client.Trigger("vehicles-chan", "vehicles-evnt", newVehicle)

	return mvc.Response{
		Err: err,
		Path: "/vehicle/index",
	}

}

// PatchCoords updates vehicle coordinates
func (c *VehicleController) PatchCoords()models.BaseResponse{
	if !c.isLoggedIn(){
		c.logout()
	}

	var vehicle = models.Vehicle{}
	c.Ctx.ReadJSON(&vehicle)

	uVehicle,err := c.Service.UpdateCoords(vehicle.ID.Hex(),vehicle)
	if err !=nil{
		fmt.Printf("Vehicle Create : Error updating vehicle %s",err.Error())
		return c.InternalErrorResponse(err.Error())
	}

	// Pusher 
	client := pusher.Client{
		AppId: core.PusherAppID,
		Key: core.PusherKey,
		Secret: core.PusherSecret,
		Cluster: core.PusherCluster,
		Secure: true,
	}

	client.Trigger("vehicles-cUpdate", "vehicles-eUpdate", uVehicle)

	return c.Successresponse(uVehicle,"Vehicle updated")
}

// GetCurrentUserID returns userID from session
func (c *VehicleController) GetCurrentUserID() string {
	userID := c.Session.GetString(userIDKey)
	return userID
}
