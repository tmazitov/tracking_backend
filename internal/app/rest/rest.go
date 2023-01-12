package rest

import (
	"internal/core/router"
)

func NewRouter(port string) {
	return router.NewRouter(port)
}
