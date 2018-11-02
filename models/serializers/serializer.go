package serializers

import (
//	"bytes"
//	"compress/gzip"
//	"crypto/rand"
//	"crypto/rsa"
//	"encoding/base64"
//	"encoding/json"
//	"net/http"
	"reflect"

)

// ENCRYPT flag to encrypt data on serializer
const ENCRYPT = false

// SerializerInterface is an interface for make Alias JSON response
type SerializerInterface interface {
	Cast(data interface{}) interface{}
}

// BaseSerializer ...
type BaseSerializer struct{
	Data 		interface{}
	IsArray		bool
	Serializer 	SerializerInterface
}

// NewSerializer intiantiates a new BaseSerializer
func NewSerializer(data interface{}, serializer SerializerInterface)SerializerInterface{
	
	return &BaseSerializer{
		Data: data,
		Serializer: serializer,
		IsArray: IsArray(data),
	}
}

// Cast realize the specific Serializer casting and returns it
func (s *BaseSerializer)Cast(data interface{})interface{}{
	if s.IsArray{
		return s.CastArray()
	}
	return s.Serializer.Cast(s.Data)
}

// CastArray ...
func (s *BaseSerializer)CastArray()interface{}{

	arr := make([]interface{}, 0)
	slice := reflect.ValueOf(s.Data)
	for i := 0; i < slice.Len(); i++ {
		if s.Serializer != nil {
			arr = append(arr, s.Serializer.Cast(slice.Index(i).Interface()))
		} else {
			arr = append(arr, slice.Index(i).Interface())
		}
	}
	return arr
}

// IsArray verify data type if array
func IsArray(data interface{}) bool {
	r := reflect.TypeOf(data)
	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	if r.Kind() == reflect.Slice {
		return true
	}

	return false
}
