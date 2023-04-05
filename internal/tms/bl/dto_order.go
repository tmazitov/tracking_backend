package bl

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type CreateOrder struct {
	StartAt        time.Time
	OwnerID        int
	Points         []Point
	Helpers        uint8
	Comment        string
	IsFragileCargo bool
}

type OrderStatus int

const (
	OrderCanceled  OrderStatus = 0
	OrderCreated   OrderStatus = 1
	OrderAccepted  OrderStatus = 2
	OrderInProcess OrderStatus = 3
	OrderDone      OrderStatus = 4
)

type R_OrderListItem struct {
	ID             int64     `json:"orderId"`
	StartAt        time.Time `json:"startAt"`
	EndAt          time.Time `json:"endAt,omitempty"`
	StatusID       int       `json:"statusId"`
	Points         []Point   `json:"points"`
	OwnerID        int       `json:"owner_id,omitempty"`
	WorkerID       int       `json:"worker_id,omitempty"`
	ManagerID      int       `json:"manager_id,omitempty"`
	Helpers        uint8     `json:"helpers,omitempty"`
	Comment        string    `json:"comment,omitempty"`
	IsFragileCargo bool      `json:"isFragileCargo,omitempty"`
}

type DB_OrderListItem struct {
	ID             int64
	CreatedAt      time.Time
	StartAt        time.Time
	EndAt          sql.NullTime
	StatusID       int
	PointsID       pq.Int64Array
	OwnerID        int
	WorkerID       sql.NullInt64
	ManagerID      sql.NullInt64
	Helpers        sql.NullInt16
	Comment        sql.NullString
	IsFragileCargo bool
}

type R_CreatableOrder struct {
	StartAt        time.Time `json:"startAt" validate:"max=32"`
	Points         []Point   `json:"points"`
	Helpers        uint8     `json:"helpers,omitempty"`
	Comment        string    `json:"comment,omitempty" validate:"max=256"`
	IsFragileCargo bool      `json:"isFragileCargo,omitempty"`
}

type R_EditableOrder struct {
	StartAt        time.Time `json:"startAt" validate:"max=32"`
	Points         []Point   `json:"points"`
	Helpers        uint8     `json:"helpers" validate:"max=32""`
	Comment        string    `json:"comment,omitempty" validate:"max=256"`
	IsFragileCargo bool      `json:"isFragileCargo,omitempty"`
}

type DB_EditableOrder struct {
	StartAt        time.Time `json:"startAt" validate:"max=32"`
	PointsID       []int64   `json:"points"`
	Helpers        uint8     `json:"helpers" validate:"max=32""`
	Comment        string    `json:"comment,omitempty" validate:"max=256"`
	IsFragileCargo bool      `json:"isFragileCargo,omitempty"`
}

type Point struct {
	ID        int64   `json:"id"`
	StepID    int16   `json:"step_id" validate:"max=256"`
	Title     string  `json:"title" validate:"max=128"`
	Floor     int8    `json:"floor" validate:"max=128"`
	Latitude  float32 `json:"lat" validate:"max=32"`
	Longitude float32 `json:"lon" validate:"max=32"`
}

func (p *Point) ToCreateData() []interface{} {
	var data []interface{}
	data = append(data, p.Title, p.StepID, p.Floor, p.Latitude, p.Longitude)
	return data
}

func (p *Point) ToEditData() []interface{} {
	var data []interface{}
	data = append(data, p.ID, p.Title, p.StepID, p.Floor, p.Latitude, p.Longitude)
	return data
}
