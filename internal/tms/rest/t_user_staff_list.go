package rest

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type StaffListHandler struct {
	Storage bl.Storage
	Jwt     jwt.JwtStorage
	result  []bl.R_GetUser
}

func (h *StaffListHandler) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if userPayload.RoleId != int(bl.Admin) && userPayload.RoleId != int(bl.Manager) {
		core.ErrorLog(403, "Forbidden", errors.New("the manager or the admin only can get the list of workers"), ctx)
		return
	}

	workersInfo, err := h.Storage.UserStorage().UserGetStaffList()
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	h.result = []bl.R_GetUser{}
	for _, workerInfo := range workersInfo {
		h.result = append(h.result, bl.R_GetUser{
			ID:        workerInfo.ID.Int64,
			ShortName: workerInfo.ShortName.String,
			RoleID:    bl.UserRole(workerInfo.RoleID.Int32),
		})
	}

	core.SendResponse(200, h.result, ctx)
}
