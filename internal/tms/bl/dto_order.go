package bl

import "time"

type CreateOrder struct {
	StartAt        time.Time
	OwnerID        int
	PointsID       []int64
	Helpers        uint8
	Comment        string
	IsFragileCargo bool
}

type Point struct {
	Title     string  `json:"title" validate:"max=128"`
	Floor     int8    `json:"floor" validate:"max=128"`
	Latitude  float32 `json:"lat" validate:"max=32"`
	Longitude float32 `json:"lon" validate:"max=32"`
}
