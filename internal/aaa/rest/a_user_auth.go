package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/aaa/bl"
	core "github.com/tmazitov/tracking_backend.git/internal/core/request"
)

type AuthUser struct {
	Storage bl.Storage
	input   struct {
		email string `json:"email" validate:"max=64"`
	}
	result struct{}
}

func (h *AuthUser) Handle(c *gin.Context) {
	if err := c.BindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, c)
		return
	}

	exist, err := h.Storage.UserStorage().CheckUserByEmail(h.input.email)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, c)
		return
	}

	// if user is exists

	if exist {

	} else {

	}

}
