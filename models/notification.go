package models

import (
	"rtlocation/core"
	"strconv"
	"time"
	"fmt"

	"github.com/globalsign/mgo/bson"

)

// Notification model
type Notification struct {
	DocumentBase `json:"-" bson:",inline"`

	To          bson.ObjectId `bson:"to"`
	From        interface{}   `bson:"from"`
	Resource    interface{}   `bson:"resource"`
	ImageURL    string        `bson:"image_url"`
	Type        string        `bson:"type"`
	Action      string        `bson:"action"`
	Message     string        `bson:"message"`
	FullMessage string        `bson:"full_message"`
	Screen      string        `bson:"screen"`
	Title       string        `bson:"title"`
	Device      string        `bson:"device"`
	IDs         []string      `bson:"ids"`
	Attachment  Attachment    `bson:"attachment"`

	// Not saved fields
	ExtraData map[string]interface{} `json:"extra_data" bson:"-"`
}

// GetDocumentName returns the name of the collection in DB
func (m *Notification) GetDocumentName() string {
	return "notifications"
}

// SendPush sends the push notification to the given ids
func (m Notification) SendPush(to string, ids []string, badge int) (*core.FcmResponseStatus, error) {
	var from = m.ExtraData["from"]
	var resource = m.ExtraData["resource"]
	var notify core.FCMNotification

	if m.Device == "IOS" {
		notify.Title = m.Title
		notify.Body = m.FullMessage
		notify.Badge = strconv.Itoa(badge)
		notify.Sound = "default"
	}

	data := &struct {
		Type      string      `json:"type"`
		Id        string      `json:"id"`
		Image     string      `json:"image"`
		Message   string      `json:"message"`
		Action    string      `json:"action"`
		From      interface{} `json:"actor"`
		Resource  interface{} `json:"resource"`
		CreatedAt int64       `json:"created_at"`
	}{m.Type, m.GetID().Hex(), "", m.Message, m.Action, from, resource, m.CreatedAt.Unix()}

	if m.Attachment.PATH != "" {
		data.Image = m.Attachment.GetURL()
	}

	resp, err := core.NewFCMClient(
		to,
		core.PriorityHigh,
		notify,
		ids,
		data,
	).Send(); 

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SendToTopic send a notification to the topic and users subscribed to topic
func (m *Notification) SendToTopic(topic, id string, regIDs []string) {

	data := struct {
		Type     string      `json:"type"`
		Id       string      `json:"id"`
		Image    string      `json:"image"`
		Message  string      `json:"message"`
		Action   string      `json:"action"`
		From     interface{} `json:"actor"`
		Resource interface{} `json:"resource"`
	}{"notification", "id_notifica", "image_url", "Come and join us", "play", "from_id", id}

	c := core.FCMClient{}
	c.AppengRegIDs(regIDs)
	c.NewFCMMessageTo(topic, data)

	status, err := c.Send()

	if err == nil {
		fmt.Printf("Push Notification  %v \n", status)
	} else {
		fmt.Printf("Push Notification error  %s \n", err.Error())
	}

}

// SendRegsIds atach the devices to the notification
func (m *Notification) SendRegsIds(regIDs []string, device string) {
	fcm := core.NewSimpleFcmClient()
	var from = m.ExtraData["from"]
	var resource = m.Resource
	var notify core.FCMNotification

	if device == "IOS" {
		notify.Title = m.Title
		notify.Body = m.FullMessage
		notify.Badge = "1"
		notify.Sound = "default"
		fcm.SetNotficationMessage(notify)
	}

	var data = &struct {
		Type      string      `json:"type"`
		Id        string      `json:"id"`
		Image     string      `json:"image"`
		Message   string      `json:"message"`
		Action    string      `json:"action"`
		From      interface{} `json:"actor"`
		Resource  interface{} `json:"resource"`
		CreatedAt interface{} `json:"created_at"`
	}{m.Type, m.Resource.(string), "", m.Message, m.Action, from, resource, time.Now().Unix()}

	if m.Attachment.PATH != "" {
		data.Image = m.Attachment.GetURL()
	}

	fcm.NewFcmRegIdsMsg(regIDs[0:1], data)
	if len(regIDs) > 1 {
		fcm.AppendDevices(regIDs[1:(len(regIDs) - 1)])
	}
	fcm.Send()
}
