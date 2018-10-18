package middleware

import (
	"errors"
	"fmt"

	"rtlocation/models"
	"rtlocation/core"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/kataras/iris"
	"github.com/globalsign/mgo/bson"
)

// Const for auth
const (
	SigningKey = "$3cr3t_k3y!"
)

// Claims Estructura para el manejo de la seccion Payload en el token.
type Claims struct {
	jwt.StandardClaims
	ID     string `json:"id"`
	Action string `json:"action"`
}

// GenerateToken Funcion que genera un token de authenticacion para un usuario.
func GenerateToken(id string, action string) (token string, err error) {
	// Tiempo de expiracion del Token, 1 semana
	// expireToken := time.Now().Add(time.Hour * 24 * 7).Unix()

	// Se crea la estructura del Payload.
	claims := Claims{
		ID:     id,
		Action: action,
	}

	// Se genera el token y se firma.
	tokenKey := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenKey.SignedString([]byte(SigningKey))
	if err!=nil{
		fmt.Printf("Error generating token: %s \n",err.Error())
		return
	}
	fmt.Printf("Token generated: %s \n",token)
	return
}

// VerifyToken ...
func VerifyToken(token, action string) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}

		return []byte(SigningKey), nil
	})
	if err != nil {
		return &Claims{}, err
	}

	// Verificamos la firma del Token, si el valida obtenemos el id del usurio.
	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		if claims.Action != action {
			return claims, errors.New("Action not match")
		}

		return claims, err
	}

	return &Claims{}, errors.New("Invalid token")
}

// AuthenticateByToken funcion para obtener un usuaio con el token.
func AuthenticateByToken(ctx iris.Context) (user models.User, err error) {
	authHeader := ctx.GetHeader("Authorization")

	if !strings.Contains(authHeader, core.Bearer) {
		err = errors.New("Missing token")
		return
	}

	token := strings.Split(authHeader, core.Bearer)[1]

	// Se intenta verificar el Token (algoritmo, y payload)
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}

		return []byte(SigningKey), nil
	})
	if err != nil {
		return user, err
	}

	// Verificamos la firma del Token, si el valida obtenemos el id del usurio.
	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		// Buscamos a el usuario en la base de datos por ID
		session ,err := mgo.Dial(core.DBUrl)
		if err!=nil{
			fmt.Printf("TokenAuth error %s \n",err)
			return user,err
		}

		if !bson.IsObjectIdHex(claims.ID){
			return user,errors.New("Not valid ID: "+claims.ID)
		}

		err = session.DB(core.DBName).C(user.GetDocumentName()).FindId(bson.ObjectIdHex(claims.ID)).One(&user)
		if err != nil{
			fmt.Printf("TokenAuth Error FindOne %s \n",err.Error())
			return user,err
		}
		defer session.Close()
	}

	fmt.Printf("User obtained %s - %s - %s \n",user.ID.Hex(),user.Firstname, user.Username)

	return user, nil
}

// AuthenticateBySession funcion para obtener un usuaio con el token.
func AuthenticateBySession(token string) (user models.User, err error) {
	// Se intenta verificar el Token (algoritmo, y payload)
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}

		return []byte(SigningKey), nil
	})
	if err != nil {
		return user, err
	}

	// Verificamos la firma del Token, si el valida obtenemos el id del usurio.
	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		// Buscamos a el usuario en la base de datos por ID
		session ,err := mgo.Dial(core.DBUrl)
		if err!=nil{
			fmt.Printf("GetByID error %s \n",err)
			return user,err
		}

		err = session.DB(core.DBName).C(user.GetDocumentName()).FindId( bson.ObjectId(claims.ID)).One(&user)
		if err != nil{
			fmt.Printf("FindByID: Error FindOne %s \n",err.Error())
			return user,err
		}
		defer session.Close()
	}

	return user, nil
}

// GetHeaderLanguage gets the language from the header "Accept-Language""
func GetHeaderLanguage(ctx iris.Context) (lang string) {
	lang = ctx.GetHeader("Accept-Language")
	return lang
}
