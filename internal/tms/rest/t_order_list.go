package rest

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/binary"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OrderListHandler struct {
	Storage bl.Storage
	Jwt     jwt.JwtStorage
	params  bl.R_OrderListFilters
	result  []bl.R_Order
}

func (h *OrderListHandler) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	h.params, err = fillFiltersByParams(ctx)
	if err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if h.result, err = getOrderList(userPayload.UserId, userPayload.RoleId, h.params, h.Storage); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(200, h.result, ctx)
}

func getOrderList(userId int64, roleId int, filters bl.R_OrderListFilters, storage bl.Storage) ([]bl.R_Order, error) {
	var (
		orders []bl.DB_Order
		result []bl.R_Order = []bl.R_Order{}
		err    error
	)

	if orders, err = storage.OrderStorage().OrderList(userId, roleId, filters); err != nil {
		return result, err
	}

	var resultOrder bl.R_Order

	for _, order := range orders {

		var (
			startAt time.Time = order.StartAt.Time
			endAt   time.Time = order.EndAt.Time
		)

		resultOrder = bl.R_Order{
			ID:                order.ID,
			Title:             order.Title,
			StartAt:           &startAt,
			EndAt:             &endAt,
			StartAtFact:       nil,
			EndAtFact:         nil,
			StatusID:          order.StatusID,
			Points:            order.Points,
			Comment:           order.Comment.String,
			IsRegularCustomer: order.IsRegularCustomer,
			Price: &bl.R_OrderBill{
				CarTypeID:      order.Bill.CarTypeID,
				CarPrice:       order.Bill.CarPrice,
				CarHours:       order.Bill.CarHours,
				HelperCount:    uint(order.Bill.HelperCount.Int16),
				HelperHours:    uint(order.Bill.HelperHours.Int16),
				HelperPrice:    uint(order.Bill.HelperPrice.Int16),
				KM:             uint(order.Bill.KM.Int16),
				IsFragileCargo: order.Bill.IsFragileCargo,
				Total:          order.Bill.Total,
				TotalInFact:    uint(order.Bill.TotalInFact.Int32),
			},
		}

		owner := bl.R_GetUser{
			ID:        order.Owner.ID.Int64,
			ShortName: order.Owner.ShortName.String,
			RoleID:    bl.UserRole(order.Owner.RoleID.Int32),
		}
		resultOrder.Owner = &owner

		if order.Worker.ID.Valid {
			var worker bl.R_GetUser = bl.R_GetUser{
				ID:        order.Worker.ID.Int64,
				ShortName: order.Worker.ShortName.String,
				RoleID:    bl.UserRole(order.Worker.RoleID.Int32),
			}
			resultOrder.Worker = &worker
			fmt.Println(order.Worker.ShortName.String)
		}

		if order.Manager.ID.Valid {
			var manager bl.R_GetUser = bl.R_GetUser{
				ID:        order.Manager.ID.Int64,
				ShortName: order.Manager.ShortName.String,
				RoleID:    bl.UserRole(order.Manager.RoleID.Int32),
			}
			resultOrder.Manager = &manager
		}

		if order.StartAtFact.Valid && !order.StartAtFact.Time.IsZero() {
			var endAtFact time.Time = order.StartAtFact.Time
			resultOrder.StartAtFact = &endAtFact
		}

		if order.EndAtFact.Valid && !order.EndAtFact.Time.IsZero() {
			var startAtFact time.Time = order.EndAtFact.Time
			resultOrder.EndAtFact = &startAtFact
		}

		result = append(result, resultOrder)
	}

	return result, err
}

func fillFiltersByParams(ctx *gin.Context) (bl.R_OrderListFilters, error) {
	var (
		filters bl.R_OrderListFilters = bl.R_OrderListFilters{}
		params  url.Values            = ctx.Request.URL.Query()
		err     error
	)

	if !params.Has("d") {
		return filters, errors.New("date is not defined in query")
	}
	filters.Date, err = time.Parse("2006-01-02", params.Get("d"))
	if err != nil {
		return filters, errors.New("invalid date parameter")
	}

	if params.Has("p") {
		value, err := strconv.ParseUint(params.Get("p"), 10, 32)
		if err != nil {
			return filters, errors.New("invalid page parameter")
		}
		filters.Page = uint(value)
	}

	if params.Has("w") {
		filters.WorkerId, err = strconv.ParseInt(params.Get("w"), 10, 32)
		if err != nil {
			return filters, errors.New("invalid worker_id parameter")
		}
	}

	if params.Has("s") {
		value, err := strconv.Atoi(params.Get("s"))
		if err != nil {
			return filters, errors.New("invalid status_id parameter")
		}
		var statuses []int = binary.PowerOfTwo(value)
		for _, statusItem := range statuses {
			filters.Statuses = append(filters.Statuses, bl.OrderStatus(math.Log2(float64(statusItem))+1))
		}
	}

	if params.Has("t") {
		value, err := strconv.Atoi(params.Get("t"))
		if err != nil {
			return filters, errors.New("invalid type_id parameter")
		}
		if value == 1 || value == 2 || value == 4 {
			filters.Types = append(filters.Types, bl.OrderType(value))
		} else if value == 3 || value == 5 || value == 6 {
			var types []int = binary.PowerOfTwo(value)
			for _, typeItem := range types {
				filters.Types = append(filters.Types, bl.OrderType(typeItem))
			}
			filters.Types = append(filters.Types, bl.OrderType(value))
		}
	}

	if params.Has("is_reg") {
		filters.IsRegularCustomer = true
	}

	return filters, nil
}
