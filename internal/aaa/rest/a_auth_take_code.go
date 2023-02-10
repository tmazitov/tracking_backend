package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/aaa/bl"
	"github.com/tmazitov/tracking_backend.git/internal/core/conductor"
	"github.com/tmazitov/tracking_backend.git/internal/core/jwt"
	core "github.com/tmazitov/tracking_backend.git/internal/core/request"
)

type AuthUserTakeCode struct {
	Storage   bl.Storage
	Conductor conductor.Conductor
	input     struct {
		Token string
		Code  string
	}
	result struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}
}

func (h *AuthUserTakeCode) Handle(ctx *gin.Context) {
	if err := ctx.BindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	// If token is invalid
	ticket, err := h.Conductor.GetTicket(ctx, h.input.Token)
	if err == conductor.ErrInvalidToken {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if err != nil {
		core.ErrorLog(500, "Internal Server error", err, ctx)
		return
	}

	// If code is invalid
	if err = ticket.ValidateCode(h.input.Code); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	tokens, err := h.Conductor.CreateTokenPair(ctx, jwt.AccessClaims{
		IP: ctx.ClientIP(),
	})
	if err != nil {
		core.ErrorLog(500, "Internal Server error", err, ctx)
		return
	}

	h.result.Access = tokens.Access
	h.result.Refresh = tokens.Refresh

	if err := h.Conductor.DeleteTicket(ctx, h.input.Token); err != nil {
		core.ErrorLog(500, "Internal Server error", err, ctx)
		return
	}

	core.SendResponse(201, h.result, ctx)
}
