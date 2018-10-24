package main

import (
	"time"
	"fmt"

	"github.com/kataras/iris"

	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/mvc"
	"github.com/globalsign/mgo"

	"rtlocation/web/controllers"
	"rtlocation/repositories"
	"rtlocation/services"
	"rtlocation/models"
	"rtlocation/core"
)

func main() {
	app := iris.New()

	err := InitializeIndexes()
	if err!=nil{
	}

	// Load the template files.
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)

	app.StaticWeb("/public", "./web/public")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	// Handler for /api routes
	//app.PartyFunc("/api", controllers.BaseControllerHandler)

	// ---- Serve our controllers. ----
	userRepo := repositories.NewUserRepository(models.User{})
	userService := services.NewUserService(userRepo)

	vehicleRepo := repositories.NewVehicleRepository()
	vehicleService := services.NewVehicleService(vehicleRepo)

	// "/user" based mvc application.
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})
	user := mvc.New(app.Party("/user"))
	user.Register(
		userService,
		sessManager.Start,
	)

	user.Handle(new(controllers.UserController))

	// Hanldes api endpoints
	vehicles := mvc.New(app.Party("/vehicle"))
	vehicles.Register(
		vehicleService,
		sessManager.Start,
	)

	vehicles.Handle(new(controllers.VehicleController))

	// http://localhost:1337/noexist
	// and all controller's methods like
	// http://localhost:1337/users/1
	app.Run(
		// Starts the web server at localhost:1337
		iris.Addr("localhost:1337"),
		// Disables the updater.
		iris.WithoutVersionChecker,
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
}

// InitializeIndexes initializes db indexes
func InitializeIndexes()error{

	session ,err:= mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("Insert: Error getting session %s \n",core.DBUrl)
		return err
	}

	cUsers := session.DB(core.DBName).C("users")

	// Users Index
	index := mgo.Index{
		Key:        []string{"user_name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = cUsers.EnsureIndex(index)
	if err != nil {
		fmt.Printf("Error creating Users index: %s \n",err.Error())
		return err
	}

	err = cUsers.EnsureIndex(Index2d())
	if err != nil {
		fmt.Printf("Error creating Users 2dindex: %s \n",err.Error())
		return err
	}

	//Vehicles Index
	cVehicles := session.DB(core.DBName).C("vehicles")

	err = cVehicles.EnsureIndex(Index2d())
	if err != nil {
		fmt.Printf("Error creating Vehicles 2dindex: %s \n",err.Error())
		return err
	}
	return nil

}

// Index2d ...
func Index2d()mgo.Index{
	return mgo.Index{
		Key: []string{"$2d:geolocation"},
		Bits: 26,
	}
}