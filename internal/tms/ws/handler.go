package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type WS_OrderListHandler struct {
	Storage bl.Storage
	Jwt     *jwt.JwtStorage
	Hub     *Hub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *WS_OrderListHandler) Handle(ctx *gin.Context) {

	log.Println("attempt to con!")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	log.Println("conn made!")
	defer conn.Close()
	var client *Client = NewClient(conn, h.Hub.redis)
	log.Println("client is created!")
	h.Hub.Register(client)

	log.Println("success conn!")

	go client.writePump(h.Jwt)
	client.readPump(h.Hub, h.Jwt)
}
