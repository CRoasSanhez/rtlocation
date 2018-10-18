package controllers

import(
	"github.com/kataras/iris/mvc"

)

// APIController ...
type APIController struct{}

func MvcRoot(app *mvc.Application){
	app.Handle(new (APIController))
}


func(c *APIController) BeforeActivation(b mvc.BeforeActivation){

	b.Handle("GET","/vehicles","VehicleController")
}
