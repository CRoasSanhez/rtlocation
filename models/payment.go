package models

// Payment ...
type Payment struct{
	PaymentType string `json:"type" form:"type" bson:"payment_type"`
	CardNumber 	string `json:"card_number" form:"card_number" bson:"card_number"`
	CVV 		string `json:"cvv" form:"cvv" bson:"cvv"`
	EndDate 	string `json:"end_date" form:"end_date" bson:"end_date"`
	UserName 	string `json:"user_name" form:"user_name" bson:"user_name"`
	Status		string `json:"-" form:"-" bson:"status"`
}
