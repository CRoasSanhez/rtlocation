package models

// Account ...
type Account struct{
	UserName string `json:"user_name" bson:"-"`
	Token string `json:"token" bson:"-"`
}