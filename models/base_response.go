package models

// BaseResponse ...
type BaseResponse struct{
	Message	string			`json:"message" form:"-" bson:"-"`
	Data	interface{}		`json:"data" form:"-" bson:"-"`
	Error	interface{}		`json:"errors" form:"-" bson:"-"`
	Success bool			`json:"success" form:"-" bson:"-"`
}

// ErrorResponse ...
type ErrorResponse struct{
	Code 	uint16		`json:"code" form:"-" bson:"-"`
	Message	string		`json:"message" form:"-" bson:"-"`	
}

// SetError ...
func(m *BaseResponse) SetError(code uint16,msg string)*BaseResponse{
	m.Error = ErrorResponse{
		Code: code,
		Message: msg,
	}
	return m
}
