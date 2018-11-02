package controllers

import(
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/mvc"

	"rtlocation/services"
	"rtlocation/models"
	auth "rtlocation/middleware"
)

// PaymentController ...
type PaymentController struct{
	BaseController
	Ctx 		iris.Context
	Service 	services.PaymentService
	Session 	*sessions.Session
	CurrentUser models.User
}

// isLoggedIn ... 
func (c *PaymentController) isLoggedIn() bool {
	if c.GetCurrentUserID()!=""{
		return true
	}
	if user, err := auth.AuthenticateByToken(c.Ctx); err == nil {
		c.CurrentUser = user
		return true
	}

	return false
}

func (c *PaymentController) logout() {
	c.Session.Destroy()
}

var paymentStaticView = mvc.View{
	Name: "payment/index.html",
}

// GetPayments handles POST: http://localhost:8080/payment/payments.
// returns the payment view
func(c *PaymentController)GetPayments()mvc.Result{
	if c.isLoggedIn() {
		c.logout()
	}
	return paymentStaticView
}

// PostRegister ...
func(c *PaymentController)PostRegister()models.BaseResponse{

	if !c.isLoggedIn(){
		c.logout()
	}

	var payment = models.Payment{}
	c.Ctx.ReadJSON(&payment)

	c.Service.AddPayment(c.CurrentUser.ID.Hex(), payment)
	return c.Successresponse(nil,"Payment added successfully",nil)
}

// DeletePayment ...
func (c *PaymentController)DeletePayment()models.BaseResponse{
	if !c.isLoggedIn(){
		c.logout()
	}

	var payment = models.Payment{}
	c.Ctx.ReadJSON(&payment)

	c.Service.AddPayment(c.CurrentUser.ID.Hex(), payment)
	return c.Successresponse(nil,"Payment added successfully",nil)
}


// GetCurrentUserID returns userID from session
func (c *PaymentController) GetCurrentUserID() string {
	userID := c.Session.GetString(userIDKey)
	return userID
}
