package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
)

type Client struct {
	userId          int64
	role            bl.UserRole
	isRefreshNeeded bool
	isAuthorized    bool
	filters         bl.R_OrderListFilters
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
		message ClientMessage
		ctx     context.Context = context.Background()
	)

	var (
		success             *Response = NewResponse(200, "Success")
		badRequestError     *Response = NewResponse(400, "Bad request")
		unauthorizedError   *Response = NewResponse(401, "Unauthorized")
		internalServerError *Response = NewResponse(500, "Internal Server error")
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
			if err = c.sendResponse(unauthorizedError); err != nil {
				log.Println("wrong report marshal: ", err)
			}
			continue
		}

		// Routing
		err = c.route(&message)
		if err == errBadRequest {
			if err = c.sendResponse(badRequestError); err != nil {
				log.Println("wrong report marshal: ", err)
			}
			continue
		}
		if err != nil {
			if err = c.sendResponse(internalServerError); err != nil {
				log.Println("wrong report marshal: ", err)
			}
			continue
		}

		// Send response success
		c.isRefreshNeeded = true
		if err = c.sendResponse(success); err != nil {
			log.Println("wrong report marshal: ", err)
			continue
		}

		// Send all messages from wait list
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

func (c *Client) sendResponse(response *Response) error {

	var (
		json []byte
		err  error
	)

	json, err = response.Marshal()
	if err != nil {
		return err
	}

	c.send <- json

	return nil
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

		fmt.Printf("user id: %d -- authorized: %t	message:%s \n\tfilters: %s \n", c.userId, c.isAuthorized, string(message), c.filters.Date.String())

		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
