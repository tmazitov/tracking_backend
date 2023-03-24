package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/aaa/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
	"github.com/tmazitov/tracking_backend.git/pkg/token"
)

type RefreshHandler struct {
	Storage bl.Storage
	Jwt     *jwt.JwtStorage
	input   struct {
		Refresh string `json:"refresh" validate:"max=64"`
	}

	result struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}
}

func (h *RefreshHandler) Handle(ctx *gin.Context) {
	if err := ctx.BindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	// Take Access token
	access, err := token.GetAuth(ctx)
	if err != nil {
		core.ErrorLog(401, "unathorized", err, ctx)
		return
	}

	// Check refresh token
	claims, err := h.Jwt.ValidateRefresh(ctx, h.input.Refresh)
	if err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	// Generate new tokens
	newTokens, err := h.Jwt.CreateTokenPair(ctx, *claims)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	h.result.Access = newTokens.Access
	h.result.Refresh = newTokens.Refresh

	// Delete old tokens
	if err = h.Jwt.DeleteTokenPair(ctx, access, h.input.Refresh); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(200, h.result, ctx)
}
