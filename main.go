package main

import (
	"time"

	"github.com/kataras/iris"

	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/mvc"

	"rtlocation/web/controllers"
	"rtlocation/repositories"
	"rtlocation/services"
	"rtlocation/models"
)

func main() {
	app := iris.New()

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

	//mvc.Configure(app.Party("/"), RootController)

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
