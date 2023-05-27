package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
)

type Hub struct {
	messagesChan chan OrderUpdateMessage
	clients      map[*Client]bool
	clientsLock  sync.Mutex
	broadcast    chan []byte
	redis        *redis.Client
	jwt          *jwt.JwtStorage
}

func NewHub(redis *redis.Client, jwt *jwt.JwtStorage) *Hub {
	return &Hub{
		messagesChan: make(chan OrderUpdateMessage),
		clients:      make(map[*Client]bool),
		broadcast:    make(chan []byte),
		clientsLock:  sync.Mutex{},
		redis:        redis,
		jwt:          jwt,
	}
}

func (h *Hub) UpdateStartAtFact(orderId int64, data interface{}) {
	h.messagesChan <- OrderUpdateMessage{
		OrderId: orderId,
		Type:    OrderUpdateStartAtFact,
		Data:    data,
	}
}

func (h *Hub) UpdateEndAtFact(orderId int64, data interface{}) {
	h.messagesChan <- OrderUpdateMessage{
		OrderId: orderId,
		Type:    OrderUpdateEndAtFact,
		Data:    data,
	}
}

func (h *Hub) Broadcast(message OrderUpdateMessage) {
	h.messagesChan <- message
}

func (h *Hub) Register(client *Client) {
	h.clientsLock.Lock()
	h.clients[client] = true
	h.clientsLock.Unlock()
}

func (h *Hub) Unregister(client *Client) {
	h.clientsLock.Lock()
	delete(h.clients, client)
	h.clientsLock.Unlock()
}

func (h *Hub) Run() {

	var ctx context.Context = context.Background()

	for {
		select {
		case messageData := <-h.broadcast:
			fmt.Println("broadcast!")
			for client := range h.clients {

				if _, err := client.checkAccess(h.jwt); err != nil {

					client.isAuthorized = false

					fmt.Println("add message to wait list")
					if err = client.waitList.Add(ctx, client, messageData); err != nil {
						log.Println("wrong save message to the wait list", err)
						continue
					}

					if client.isRefreshNeeded {
						fmt.Println("make report to client")
						jsonReport, err := NewResponse(401, "Unauthorized").Marshal()
						if err != nil {
							log.Println("wrong report marshal")
							continue
						}

						client.send <- jsonReport
						fmt.Println("send report to client")
					}
					continue
				}

				select {
				case client.send <- messageData:
				default:
					h.Unregister(client)
				}
			}
		}
	}
}

func (h *Hub) RunOrderDispatcher() {
	for {
		select {
		case message := <-h.messagesChan:
			fmt.Println("message get!")
			j, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				continue
			}
			h.broadcast <- j
		}
	}
}
