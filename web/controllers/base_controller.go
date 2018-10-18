package controllers

import(
	"fmt"
	"errors"

	"github.com/go-ozzo/ozzo-validation"

	"rtlocation/models"
	"encoding/json"
	"rtlocation/core"
	//auth "rtlocation/middleware"
	
)

// BaseController ...
type BaseController struct{
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