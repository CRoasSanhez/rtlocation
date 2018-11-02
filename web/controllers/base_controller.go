package controllers

import(
	"fmt"
	"errors"
	"encoding/json"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/kataras/iris"

	"rtlocation/models"
	"rtlocation/core"
	"rtlocation/models/serializers"
	//auth "rtlocation/middleware"
	
)

// BaseController ...
type BaseController struct{
}

// BaseControllerHandler for base controller
func BaseControllerHandler(api iris.Party){
		//users.Use(myAuthMiddlewareHandler)

		// http://localhost:8080/users/42/profile
		//users.Get("/{id:uint64}/profile", userProfileHandler)
		// http://localhost:8080/users/messages/1
		//users.Get("/inbox/{id:uint64}", userMessageHandler)

}

// GetCurrentUser returns user's session
func (c *BaseController) GetCurrentUser() bool {
	/*
	if user, err := auth.AuthenticateByToken(c.Request); err == nil {
		c.CurrentUser = user
		return true
	}

	return false
	*/
	return true
}

// Successresponse ...
func (c BaseController) Successresponse(data interface{},msg string, serializer serializers.SerializerInterface)(response models.BaseResponse){
	response.Success = true
	response.Message = msg
	response.Error = nil
	if serializer != nil{
		response.Data = serializers.NewSerializer(data,serializer).Cast(nil)
	}else{
		response.Data = data
	}

	return 
}

// ErrorResponse ...
func(c BaseController) ErrorResponse(code uint16, errMsg, msg string)(response models.BaseResponse){
	response.Success = false
	response.Data = nil
	response.Message = msg
	response.SetError(code, errMsg)
	return 
}

// InternalErrorResponse ...
func (c BaseController) InternalErrorResponse(errMsg string)(response models.BaseResponse){
	response.Success = false
	response.Data = nil
	response.Message = "Internal server error"
	response.SetError(400, errMsg)
	return 
}

// BadRequest ...
func(c BaseController)BadRequest(errMsg string)(response models.BaseResponse){
	response.Success = false
	response.Data = nil
	response.Message = "Bad request"
	response.SetError(400, errMsg)
	return 
}

// ForbiddenResponse ...
func (c BaseController)ForbiddenResponse()(response models.BaseResponse){
	response.Success = false
	response.Data = nil
	response.Message = "Invalid User"
	response.SetError(403, response.Message)
	return 
}

// NewNotification ...
// 't' is notification type = invitation, notification
func (c *BaseController) NewNotification(action string, attachment models.Attachment, message string, screen string, title string, t string, fromDocName string, to models.User, resource interface{}) error {

	notification := models.Notification{
		From:        fromDocName,
		To:          to.ID,
		Resource:    resource,
		Action:      action,
		Type:        t,
		Message:     message,
		FullMessage: message,
		Screen:      screen,
		Title:       title,
		Device:      to.Device.OS,
		Attachment:  attachment,
	}

	resp, err := notification.SendPush(to.Device.MessagingToken, []string{}, 1); 
	if err != nil {
		return err
	} 
	if !resp.Ok {
		fmt.Printf("Send Push error: %v \n",resp.Results)
	}
	return nil
}

// PermitParams ...
func (c BaseController) PermitParams(params map[string]interface{},out interface{}, validate bool, allowed ...string) error {

	for k := range params {
		if core.FindOnArray(allowed, k) < 0 {
			delete(params, k)
		}
	}

	if j, err := json.Marshal(params); err != nil {
		return err
	} else {
		if err := json.Unmarshal(j, out); err != nil {
			return err
		}
	}

	if validate {
		if v, ok := out.(validation.Validatable); ok {
			return v.Validate()
		} else {
			return errors.New("Validatable object expected")
		}
	}

	return nil
}