package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type StaffRemoveHandler struct {
	Storage bl.Storage
	input   struct {
		UserId int64 `json:"userId"`
	}
}

func (h *StaffRemoveHandler) Handle(ctx *gin.Context) {

	if err := ctx.ShouldBindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if err := h.Storage.UserStorage().StaffRemove(h.input.UserId); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	core.SendResponse(201, nil, ctx)
}
