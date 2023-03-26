package bl

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type CreateOrder struct {
	StartAt        time.Time
	OwnerID        int
	PointsID       []int64
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

type Point struct {
	Title     string  `json:"title" validate:"max=128"`
	Floor     int8    `json:"floor" validate:"max=128"`
	Latitude  float32 `json:"lat" validate:"max=32"`
	Longitude float32 `json:"lon" validate:"max=32"`
}
