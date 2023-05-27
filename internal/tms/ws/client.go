package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
)

type Client struct {
	userId          int64
	isRefreshNeeded bool
	isAuthorized    bool
	conn            *websocket.Conn
	send            chan []byte
	access          string
	waitList        *WaitList
}

func NewClient(conn *websocket.Conn, redis *redis.Client) *Client {
	return &Client{
		conn:         conn,
		send:         make(chan []byte, 256),
		isAuthorized: false,
		waitList:     NewWaitList(redis),
	}
}

func (c *Client) readPump(hub *Hub, jwt *jwt.JwtStorage) {

	var (
		message      AuthMessage
		ctx          context.Context = context.Background()
		jsonResponse []byte
	)

	defer c.conn.Close()
	for {
		_, byteArray, err := c.conn.ReadMessage()
		if err != nil {
			hub.Unregister(c)
			break
		}

		if err := json.Unmarshal(byteArray, &message); err != nil {
			log.Println("unexpected message: ", string(byteArray))
			continue
		}

		// Authorization check
		_, err = c.authMiddleware(jwt, message.Access)
		if err != nil {
			jsonResponse, err := NewResponse(401, "Unauthorized").Marshal()
			if err != nil {
				log.Println("wrong report marshal")
				continue
			}
			c.send <- jsonResponse
			continue
		}

		c.isRefreshNeeded = true
		jsonResponse, err = NewResponse(200, "Success").Marshal()
		if err != nil {
			log.Println("wrong report marshal")
			continue
		}

		c.send <- jsonResponse

		waitListMessages, err := c.waitList.GetAll(ctx, c)
		if err != nil {
			log.Println("wrong get messages from wait list", err)
			continue
		}
		for _, message := range waitListMessages {
			c.send <- message
		}
	}
}

func (c *Client) writePump(jwt *jwt.JwtStorage) {

	defer c.conn.Close()
	for {
		message, ok := <-c.send
		if !ok {
			err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			if err != nil {
				log.Println(err)
			}
			return
		}

		fmt.Printf("user id: %d -- authorized: %t	message:%s \n", c.userId, c.isAuthorized, string(message))

		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
