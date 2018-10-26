package repositories

import(
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/globalsign/mgo"

	"rtlocation/models"
	"rtlocation/core"
)

// PaymentRepository ...
type PaymentRepository interface{
	AddPayment(string, models.Payment)(models.Payment,error)
	DeletePayment(string, models.Payment)(error)
}

// NewPaymentRepository ...
func NewPaymentRepository(models.Payment)PaymentRepository{
	return &PaymentBDRepository{}
}

// PaymentBDRepository is the DB entity for users
type PaymentBDRepository struct{
	BaseRepository
	Payment models.Payment
}

// AddPayment adds user's payment method
func(r *PaymentBDRepository)AddPayment(userIDHex string, payment models.Payment)(models.Payment,error){

	var response models.Payment
	user := models.User{}

	payment.Status = core.StatusActive
/*
	var queryFind = bson.M{
		"_id": bson.ObjectId(userIDHex),
		"payments": bson.M{ 
			"$elemMatch": bson.M{ 
				"card_number": payment.CardNumber,
				},
		},
	}
*/
	var queryAdd = bson.M{ 
			"$addToSet": bson.M{  
				"payments": payment,
			},
	}
/*
	var queryUpdate = bson.M{
		"$set":bson.M{
			"payment": bson.M{
				"payment_type":"credit_card",
				"card_number": "xxxxxxxxxxxxxxxx",
				"cvv":"xxx",
				"end_date": "01/19",
				"user_name": nil,
				},
			},
		}
		*/

	session ,err := mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("AddPayment error session %s \n",err)
		return response,err
	}
/*
	// Find user with payment in DB
	err = session.DB(core.DBName).C(user.GetDocumentName()).FindId(bson.ObjectId(userIDHex)).One(&user)
	if err != nil{
		fmt.Printf("AddPayment: Error Finding user %s \n",err.Error())
		return response,err
	}
	*/

	// Appends payment in user model
	err = session.DB(core.DBName).C(user.GetDocumentName()).UpdateId(bson.ObjectIdHex(userIDHex),queryAdd)
	if err != nil{
		fmt.Printf("AddPayment: Error updating %s \n",err.Error())
		return response,err
	}

	defer session.Close()

	return payment,nil
}

// DeletePayment ...
func(r *PaymentBDRepository)DeletePayment(userIDHex string, payment models.Payment)(error){
	var user = models.User{}

	var querySelector = bson.M{ 
		"_id" : bson.ObjectId(userIDHex), 
		"payments.card_number" : payment.CardNumber, 
	}
	var queryUpdate = bson.M{ 
			"$set" : bson.M{ "payments.$.status" : core.StatusDeleted } ,
		}
	
	session ,err := mgo.Dial(core.DBUrl)
	if err!=nil{
		fmt.Printf("DeletePayment error session %s \n",err)
		return err
	}

	err = session.DB(core.DBName).C(user.GetDocumentName()).Update(querySelector,queryUpdate)
	if err != nil{
		fmt.Printf("DeletePayment: Error deleting %s \n",err.Error())
		return err
	}
	return nil
}