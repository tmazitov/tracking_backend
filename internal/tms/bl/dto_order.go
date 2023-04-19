package bl

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type OrderStatus int

const (
	OrderCanceled  OrderStatus = 0
	OrderCreated   OrderStatus = 1
	OrderAccepted  OrderStatus = 2
	OrderInProcess OrderStatus = 3
	OrderDone      OrderStatus = 4
)

type OrderType int

const (
	OrderInCity        OrderType = 1
	OrderInNearOfCity  OrderType = 2
	OrderInBetweenCity OrderType = 3
)

type R_OrderListItem struct {
	ID                int64     `json:"orderId"`
	Title             string    `json:"title"`
	StartAt           time.Time `json:"startAt"`
	EndAt             time.Time `json:"endAt,omitempty"`
	StatusID          int       `json:"statusId"`
	Points            []Point   `json:"points"`
	OwnerID           int       `json:"owner_id,omitempty"`
	WorkerID          int       `json:"worker_id,omitempty"`
	ManagerID         int       `json:"manager_id,omitempty"`
	Helpers           uint8     `json:"helpers,omitempty"`
	Comment           string    `json:"comment,omitempty"`
	IsFragileCargo    bool      `json:"isFragileCargo,omitempty"`
	IsRegularCustomer bool      `json:"isRegularCustomer,omitempty"`
}

type DB_OrderListItem struct {
	ID                int64
	CreatedAt         time.Time
	Title             string
	StartAt           time.Time
	EndAt             sql.NullTime
	StatusID          int
	PointsID          pq.Int64Array
	OwnerID           int
	WorkerID          sql.NullInt64
	ManagerID         sql.NullInt64
	Helpers           sql.NullInt16
	Comment           sql.NullString
	IsFragileCargo    bool
	IsRegularCustomer bool
}

type R_CreatableOrder struct {
	StartAt           time.Time `json:"startAt" binding:"required" validate:"max=32"`
	Points            []Point   `json:"points"  binding:"required" validate:"dive"`
	Title             string    `json:"title,omitempty"   validate:"max=64"`
	WorkerID          int64     `json:"workerId,omitempty"`
	Helpers           uint8     `json:"helpers,omitempty"`
	OrderType         uint8     `json:"orderType,omitempty"`
	Comment           string    `json:"comment,omitempty" validate:"max=256"`
	IsFragileCargo    bool      `json:"isFragileCargo,omitempty"`
	IsRegularCustomer bool      `json:"isRegularCustomer,omitempty"`
}

type CreateOrder struct {
	OwnerID           int
	WorkerID          int64
	StartAt           time.Time
	Title             string
	Points            []Point
	Helpers           uint8
	OrderType         uint8
	Comment           string
	IsFragileCargo    bool
	IsRegularCustomer bool
}

type R_EditableOrder struct {
	StartAt           time.Time `json:"startAt" binding:"required" validate:"max=32"`
	Points            []Point   `json:"points"  binding:"required" validate:"dive"`
	Title             string    `json:"title,omitempty"   validate:"max=64"`
	WorkerID          int64     `json:"workerId,omitempty"`
	Helpers           uint8     `json:"helpers,omitempty"`
	OrderType         uint8     `json:"orderType,omitempty"`
	Comment           string    `json:"comment,omitempty" validate:"max=256"`
	IsFragileCargo    bool      `json:"isFragileCargo,omitempty"`
	IsRegularCustomer bool      `json:"isRegularCustomer,omitempty"`
}

type DB_EditableOrder struct {
	StartAt        time.Time `json:"startAt" validate:"max=32"`
	PointsID       []int64   `json:"points"`
	Helpers        uint8     `json:"helpers"`
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
