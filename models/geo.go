package models

// Geo ...
type Geo struct {
	Type string `json:"type" bson:"type" default:"Point"`

	// regex ^[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?),[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?)$
	// long,lat
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}