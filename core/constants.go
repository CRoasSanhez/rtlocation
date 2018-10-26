package core

// DBName is the Database name
var DBName = "rtlocation"
var (
	ActionAuth      = "auth"
	DBUrl = "mongodb://127.0.0.1:27017/"+DBName

	PusherAppID = "622037"
	PusherKey = "f8a8b497e58a374fd8d6"
	PusherSecret = "cd4b60a942fa204edd8d"
	PusherCluster = "us2"

	GoogleMapsAPIKey = "AIzaSyDI8Osr696gx40sqjWaNg9f6EDvcFLgYGQ"
	MapsHereAPIKey = "U58Kc7tHD86pXPjLfAsq"
	MapsHereAPPCode = "G3HRnaT6zc97mg3TEEFHSg"

	LocationTypeVehicle = "car"
	LocationTypeUser = "user"

	MaxDistanceFromDriver = 1000
	Bearer       = "Bearer "

	StatusInit = "init"
	StatusActive = "active"
	StatuspendingConfirmation = "pending_confirmation"
	StatusDeleted = "deleted"

	ResponseStatusSuccess = 200
	ResponseStatusBadRequest = 400
	ResponseStatusForbidden = 403
	ResponseStatusNotFound = 404
	ResponseStatusInternalServerError = 500
)