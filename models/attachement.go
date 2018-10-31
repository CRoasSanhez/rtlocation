package models

import (
	"fmt"
	"bytes"
	"io"
	"mime/multipart"
	"os/exec"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"

)

// Attachment model
type Attachment struct {
	ACL           	string                	`bson:"acl"`
	Action         	string                	`bson:"action"`
	CurrentName    	string                	`bson:"current_name"`
	Format         	string                	`bson:"format"`
	OriginalName   	string               	`bson:"original_name"`
	PATH           	string               	`bson:"path"`
	Signing        	string                	`bson:"signing"`
	SigningTime    	time.Time             	`bson:"signing_time"`
	SigningMinutes 	int                   	`bson:"signing_minutes"`
	Size           	int64                 	`bson:"size"`
	URL            	string                	`bson:"file_url"`
	Dimensions     	map[string]Attachment 	`bson:"dimensions"`
	Binary		   	[]byte					`bson:"binary"`

	// Not saved fields
	TMPFile *multipart.FileHeader `bson:"-"`
}

// Init the model
func (m *Attachment) Init(owner DocumentInterface, part *multipart.FileHeader) error {
	// S3 PATH -> /:model/:id/:action/:uuid.ext
	model := owner.GetDocumentName()
	id := owner.GetID().Hex()
	action := m.Action
	uuid, err := m.GenerateUUID()

	if err != nil {
		return err
	}

	nameSplit := strings.Split(strings.ToLower(part.Filename), ".")
	suffix := nameSplit[len(nameSplit)-1]

	// Open multipart file
	file, err := part.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	m.TMPFile = part
	m.ACL = core.S3PublicRead
	m.Format = suffix
	m.CurrentName = uuid + "." + suffix
	m.OriginalName = part.Filename
	m.PATH = model + "/" + id + "/" + action + "/" + m.CurrentName
	m.Size, _ = file.Seek(io.SeekStart, io.SeekEnd)

	// Set binary data of image to save on DB
	if pos := core.FindOnArray([]string{"jpg","png","jpeg"},suffix); pos>0{
		buf := bytes.NewBuffer(nil)
		_,err = io.Copy(buf,file)
		if err != nil{
			return err
		}
		m.Binary = buf.Bytes()
	}

	return nil
}

// Upload to S3
func (m *Attachment) Upload() error {
	// Open multipart file
	file, err := m.TMPFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	if err = core.UploadFile(file, m.PATH, m.ACL, m.Format, m.Size); err != nil {
		return err
	}

	url, err := core.GetS3Object(core.DefaultSigningTime, m.PATH)
	if err != nil {
		return err
	}

	m.URL = url
	m.SigningMinutes = core.DefaultSigningTime
	m.SigningTime = time.Now()

	return nil
}

// UploadBytes initializes attachement
func (m *Attachment) UploadBytes(owner DocumentInterface, byteArray []byte, fileName string) error {
	// S3 PATH -> /:model/:id/:action/:uuid.ext
	model := owner.GetDocumentName()
	id := owner.GetID().Hex()
	action := m.Action
	uuid, err := m.GenerateUUID()

	if err != nil {
		fmt.Printf("ERROR GENERATE UUID: %s \n",err.Error())
		return err
	}

	nameSplit := strings.Split(strings.ToLower(fileName), ".")
	suffix := nameSplit[len(nameSplit)-1]

	m.ACL = core.S3PublicRead
	m.Format = suffix
	m.CurrentName = uuid + "." + suffix
	m.OriginalName = fileName
	m.PATH = model + "/" + id + "/" + action + "/" + m.CurrentName
	m.Size = int64(len(byteArray))

	// Convert byteArray to reader
	reader := bytes.NewReader(byteArray)
	m.URL, err = core.UploadOSFileReader(reader, m.PATH, m.ACL, m.Format, m.Size)
	if err != nil {
		fmt.Printf("ERROR Upload file: %s \n",err.Error())
		return err
	}
	return nil
}

// Remove delete a resource from S3
func (m *Attachment) Remove() error {
	return core.DeleteFile(m.PATH)
}

// HasExpired returns true if Attachamnet url has expired
func (m *Attachment) HasExpired() bool {
	var validTime = m.SigningTime.Add(time.Duration(m.SigningMinutes) * time.Minute)
	return validTime.Unix() < time.Now().Unix()
}

// UpdateURL update signing time and file url
func (m *Attachment) UpdateURL() error {
	url, err := core.GetS3Object(core.DefaultSigningTime, m.PATH)
	if err != nil {
		return err
	}

	m.URL = url
	m.SigningMinutes = core.DefaultSigningTime
	m.SigningTime = time.Now()

	return nil
}

// GetURL returns presig url
func (m *Attachment) GetURL() string {
	url, _ := core.GetS3Object(core.DefaultSigningTime, m.PATH)

	return url
}

// UpdateURLParentField update signing time and file url
func (m *Attachment) UpdateURLParentField(parent string, parentID string, field string) error {
	url, err := core.GetS3Object(core.DefaultSigningTime, m.PATH)
	if err != nil {
		return err
	}

	m.URL = url
	m.SigningMinutes = core.DefaultSigningTime
	m.SigningTime = time.Now()

	Model := app.Mapper.InitModel(parent)
	var selector = bson.M{"_id": bson.ObjectIdHex(parentID)}
	var query = bson.M{"$set": bson.M{field: m}}

	return Model.UpdateQuery(selector, query, false)
}

// UpdateURLParentQuery updates the attachement according to the given query and selector
// {"_id" : ObjectId(),'parent_field.name_element.field_name': "new_value"},
// {"$set":{"parent_field.$.name_element.field_name": "new_value"}}
func (m *Attachment) UpdateURLParentQuery(parent, field, subfield string, resource interface{}) error {
	url, err := core.GetS3Object(core.DefaultSigningTime, m.PATH)
	if err != nil {
		return err
	}

	m.URL = url
	m.SigningMinutes = core.DefaultSigningTime
	m.SigningTime = time.Now()
	Model := app.Mapper.InitModel(parent)

	var selector = bson.M{
		"_id": resource, field + "." + subfield + ".current_name": m.CurrentName,
	}
	var query = bson.M{"$set": bson.M{field + ".$." + subfield + ".file_url": url}}

	return Model.UpdateQuery(selector, query, false)
}

//GenerateUUID generates unique filename to save in s3
func (m *Attachment) GenerateUUID() (string, error) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}

	return strings.Replace(strings.Trim(string(out), "\n"), "-", "_", -1), nil
}
