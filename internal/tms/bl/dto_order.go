package bl

import (
	"database/sql"
	"time"
)

type OrderStatus int

const (
	OrderDone      OrderStatus = 1
	OrderCanceled  OrderStatus = 2
	OrderCreated   OrderStatus = 3
	OrderAccepted  OrderStatus = 4
	OrderInProcess OrderStatus = 5
)

func OrderStatusArrayToIntArray(statuses []OrderStatus) []int {
	var intArray []int
	for _, status := range statuses {
		intArray = append(intArray, int(status))
	}
	return intArray
}

type OrderType int

const (
	OrderInCity        OrderType = 1
	OrderInNearOfCity  OrderType = 2
	OrderInBetweenCity OrderType = 4
)

const DB_OrderListRowCount uint = 24

type R_OrderListFilters struct {
	Date              time.Time     `json:"date,omitempty"`
	Title             string        `json:"title,omitempty"`
	Page              uint          `json:"page,omitempty"`
	WorkerId          int64         `json:"workerId,omitempty"`
	Statuses          []OrderStatus `json:"statuses,omitempty"`
	Types             []OrderType   `json:"types,omitempty"`
	IsRegularCustomer bool          `json:"isRegularCustomer,omitempty"`
}

type R_Order struct {
	ID                int64        `json:"orderId"`
	Title             string       `json:"title"`
	StartAt           *time.Time   `json:"startAt"`
	StartAtFact       *time.Time   `json:"startAtFact,omitempty"`
	EndAt             *time.Time   `json:"endAt"`
	EndAtFact         *time.Time   `json:"endAtFact,omitempty"`
	StatusID          int          `json:"statusId"`
	Points            []Point      `json:"points"`
	OrderType         uint8        `json:"orderType,omitempty"`
	Owner             *R_GetUser   `json:"owner,omitempty"`
	Worker            *R_GetUser   `json:"worker,omitempty"`
	Manager           *R_GetUser   `json:"manager,omitempty"`
	Helpers           uint8        `json:"helpers,omitempty"`
	Comment           string       `json:"comment,omitempty"`
	IsFragileCargo    bool         `json:"isFragileCargo,omitempty"`
	IsRegularCustomer bool         `json:"isRegularCustomer,omitempty"`
	Price             *R_OrderBill `json:"price"`
}

type DB_Order struct {
	ID                int64
	CreatedAt         time.Time
	Title             string
	StartAt           sql.NullTime
	StartAtFact       sql.NullTime
	EndAt             sql.NullTime
	EndAtFact         sql.NullTime
	StatusID          int
	OrderType         uint8
	Points            []Point
	Owner             DB_GetUser
	Worker            DB_GetUser
	Manager           DB_GetUser
	Comment           sql.NullString
	Bill              DB_OrderBill
	IsRegularCustomer bool
}

func (o *DB_Order) ToReal() *R_Order {
	var (
		startAt time.Time = o.StartAt.Time
		endAt   time.Time = o.EndAt.Time
	)
	var order R_Order = R_Order{
		ID:                o.ID,
		Title:             o.Title,
		StartAt:           &startAt,
		EndAt:             &endAt,
		StartAtFact:       nil,
		EndAtFact:         nil,
		StatusID:          o.StatusID,
		Points:            o.Points,
		Comment:           o.Comment.String,
		IsRegularCustomer: o.IsRegularCustomer,
		Price: &R_OrderBill{
			CarTypeID:      o.Bill.CarTypeID,
			HelperCount:    uint(o.Bill.HelperCount.Int16),
			HelperHours:    uint(o.Bill.HelperHours.Int16),
			IsFragileCargo: o.Bill.IsFragileCargo,
		},
	}

	owner := R_GetUser{
		ID:        o.Owner.ID.Int64,
		ShortName: o.Owner.ShortName.String,
		RoleID:    UserRole(o.Owner.RoleID.Int32),
	}
	order.Owner = &owner

	if o.Worker.ID.Valid {
		var worker R_GetUser = R_GetUser{
			ID:        o.Worker.ID.Int64,
			ShortName: o.Worker.ShortName.String,
			RoleID:    UserRole(o.Worker.RoleID.Int32),
		}
		order.Worker = &worker
	}

	if o.Manager.ID.Valid {
		var manager R_GetUser = R_GetUser{
			ID:        o.Manager.ID.Int64,
			ShortName: o.Manager.ShortName.String,
			RoleID:    UserRole(o.Manager.RoleID.Int32),
		}
		order.Manager = &manager
	}

	if o.StartAtFact.Valid && !o.StartAtFact.Time.IsZero() {
		var endAtFact time.Time = o.StartAtFact.Time
		order.StartAtFact = &endAtFact
	}

	if o.EndAtFact.Valid && !o.EndAtFact.Time.IsZero() {
		var startAtFact time.Time = o.EndAtFact.Time
		order.EndAtFact = &startAtFact
	}

	return &order
}

type R_OrderBill struct {
	CarTypeID      uint8 `json:"carTypeId"`
	HelperCount    uint  `json:"helperCount,omitempty"`
	HelperPrice    uint  `json:"helperPrice,omitempty"`
	HelperHours    uint  `json:"helperHours,omitempty"`
	CarPrice       uint  `json:"carPrice"`
	CarHours       uint  `json:"carHours"`
	KmPrice        uint  `json:"kmPrice,omitempty"`
	KmCount        uint  `json:"kmCount,omitempty"`
	IsFragileCargo bool  `json:"isFragileCargo,omitempty"`
	Total          uint  `json:"total"`
	TotalInFact    uint  `json:"totalInFact,omitempty"`
}

type DB_OrderBill struct {
	CarTypeID      uint8
	HelperCount    sql.NullInt16
	HelperPrice    sql.NullInt16
	HelperHours    sql.NullInt16
	CarPrice       uint
	CarHours       uint
	KmCount        sql.NullInt16
	KmPrice        sql.NullInt16
	IsFragileCargo bool
	Total          uint
	TotalInFact    sql.NullInt32
}

type R_CreatableOrder struct {
	StartAt           time.Time   `json:"startAt" binding:"required" validate:"max=32"`
	EndAt             time.Time   `json:"endAt"   binding:"required" validate:"max=32"`
	Points            []Point     `json:"points"  binding:"required" validate:"dive"`
	Title             string      `json:"title,omitempty"   validate:"max=64"`
	WorkerID          int64       `json:"workerId,omitempty"`
	OrderType         uint8       `json:"orderType,omitempty"`
	Comment           string      `json:"comment,omitempty" validate:"max=256"`
	IsRegularCustomer bool        `json:"isRegularCustomer,omitempty"`
	Price             R_OrderBill `json:"price"`
}

type CreateOrder struct {
	OwnerID           int64
	WorkerID          int64
	StartAt           time.Time
	EndAt             time.Time
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

type OrderPriceList struct {
	BigCarPrice  uint `json:"bigCarPrice"`
	BigCarTime   uint `json:"bigCarTime"`
	HelperPrice  uint `json:"helperPrice"`
	HelperTime   uint `json:"helperTime"`
	FragilePrice uint `json:"fragilePrice"`
	KM           uint `json:"kmPrice"`
}
