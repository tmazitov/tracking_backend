package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/aaa/bl"
	"github.com/tmazitov/tracking_backend.git/internal/core/conductor"
	core "github.com/tmazitov/tracking_backend.git/internal/core/request"
)

type AuthUserSendCode struct {
	Storage   bl.Storage
	Conductor conductor.Conductor
	input     struct {
		Email string `json:"email" validate:"max=64"`
	}
	result struct {
		Token string
	}
}

func (h *AuthUserSendCode) Handle(c *gin.Context) {
	if err := c.BindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, c)
		return
	}
	var err error
	email := h.input.Email

	// exist, err := h.Storage.UserStorage().CheckUserByEmail(email)
	// if err != nil {
	// 	core.ErrorLog(500, "Internal server error", err, c)
	// 	return
	// }

	// fmt.Println(exist)

	h.result.Token, err = h.Conductor.CreateTicket(c, email)

	if err == conductor.ErrTooManyAttempts {
		core.ErrorLog(429, "Too many attempts", err, c)
		return
	}

	if err != nil {
		core.ErrorLog(500, "Internal server error", err, c)
		return
	}

	core.SendResponse(200, h.result, c)
}
