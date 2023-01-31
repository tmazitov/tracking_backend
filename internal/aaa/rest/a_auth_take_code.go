package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/aaa/bl"
	"github.com/tmazitov/tracking_backend.git/internal/core/conductor"
)

type AuthUserTakeCode struct {
	Storage   bl.Storage
	Conductor conductor.Conductor
}

func (h *AuthUserTakeCode) Handle(ctx *gin.Context) {

}
