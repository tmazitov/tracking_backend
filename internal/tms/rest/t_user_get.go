package rest

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
// 	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
// 	core "github.com/tmazitov/tracking_backend.git/pkg/request"
// )

// type GetOrderHandler struct {
// 	Jwt     jwt.JwtStorage
// 	Storage *bl.Storage
// 	result  bl.GetUser
// }

// func (h *GetOrderHandler) Handle(ctx *gin.Context) {
// 	userPayload, err := h.Jwt.ValidateAccess(ctx)
// 	if err != nil {
// 		core.ErrorLog(403, "Bad request", err, ctx)
// 		return
// 	}

// 	userPayload = h
// }
