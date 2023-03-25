package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/aaa/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/conductor"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type AuthUserTakeCode struct {
	Storage   bl.Storage
	Jwt       *jwt.JwtStorage
	Conductor *conductor.Conductor
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

	payload, err := h.Conductor.ValidateCode(ctx, h.input.Token)
	if err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	// If token is invalid
	recordedCode, err := h.Conductor.GetCode(ctx, h.input.Token)
	if err == conductor.ErrInvalidToken {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}
	if err != nil {
		core.ErrorLog(500, "Internal Server error", err, ctx)
		return
	}

	// If code is invalid
	if h.input.Code != recordedCode {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	var (
		userId int
		roleId int
	)
	userId, roleId, err = h.createUserIfNotExists(payload.Email)
	if err != nil {
		core.ErrorLog(500, "Internal Server error", err, ctx)
		return
	}
	tokens, err := h.Jwt.CreateTokenPair(ctx, jwt.AccessClaims{
		RoleId: roleId,
		UserId: userId,
	})
	if err != nil {
		core.ErrorLog(500, "Internal Server error", err, ctx)
		return
	}

	h.result.Access = tokens.Access
	h.result.Refresh = tokens.Refresh

	if err := h.Conductor.DeleteCode(ctx, h.input.Token); err != nil {
		core.ErrorLog(500, "Internal Server error", err, ctx)
		return
	}

	core.SendResponse(201, h.result, ctx)
}

func (h *AuthUserTakeCode) createUserIfNotExists(email string) (int, int, error) {
	var (
		err error
	)

	// Check exist user
	userId, roleId, err := h.Storage.UserStorage().CheckUserByEmail(email)
	if err != nil || userId != 0 {
		return userId, roleId, err
	}

	// Create new user
	userId, err = h.Storage.UserStorage().CreateUser(email)
	if err != nil {
		return userId, 0, err
	}

	return userId, 0, nil
}
