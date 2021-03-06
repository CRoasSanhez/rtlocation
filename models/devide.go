package models

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"rtlocation/core"
)

// Device model
type Device struct {
	Number         string   `json:"number,omitempty" bson:"number"`
	OS             string   `json:"os,omitempty" bson:"os"`
	Version        string   `json:"version,omitempty" bson:"version"`
	UUID           string   `json:"uuid,omitempty" bson:"uuid"`
	MessagingToken string   `json:"messaging_token,omitempty" bson:"messaging_token,omitempty"`
	Code           string   `json:"-" bson:"code"`
	Language       string   `json:"language,omitempty" bson:"language"`
	Status         string 	`json:"status,omitempty" bson:"status"`
}

// Init ....
func (d *Device) Init() {
	d.MessagingToken = strings.Trim(d.MessagingToken, " ")
	d.Status = core.StatusInit
}

// RandomCode ....
func (d *Device) RandomCode() {
	rand.Seed(time.Now().Unix())
	r1 := rand.Intn(99999999) + 1
	d.Code = strconv.Itoa(r1)
}

// GenerateCode ...
func (d *Device) GenerateCode() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var buf bytes.Buffer

	for i := 0; i < 6; i++ {
		buf.WriteString(strconv.Itoa(r1.Intn(9)))
	}

	d.Code = buf.String()
	d.Status = core.StatuspendingConfirmation
}
