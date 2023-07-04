package rest

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type StaffWorkTimePut struct {
	Storage bl.Storage
	input   bl.StaffWorkTime
}

func (h *StaffWorkTimePut) Handle(ctx *gin.Context) {
	if err := ctx.ShouldBindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if h.input.StartAt > 1440 || h.input.EndAt > 1440 || h.input.StartAt > h.input.EndAt {
		core.ErrorLog(400, "Bad request", errors.New("update work time : invalid 'from' or 'to'"), ctx)
		return
	}

	if err := h.Storage.UserStorage().StaffWorkTimeUpdate(&h.input); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(200, nil, ctx)
}
