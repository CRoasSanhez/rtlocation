package repositories

import(

	"github.com/globalsign/mgo"

)

// BaseRepository ...
type BaseRepository struct{
	session mgo.Session
}
