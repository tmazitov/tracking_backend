package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
)

type Hub struct {
	storage      bl.Storage
	messagesChan chan OrderUpdateMessage
	clients      map[*Client]bool
	clientsLock  sync.Mutex
	broadcast    chan []byte
	redis        *redis.Client
	jwt          *jwt.JwtStorage
}

func NewHub(storage bl.Storage, redis *redis.Client, jwt *jwt.JwtStorage) *Hub {
	return &Hub{
		messagesChan: make(chan OrderUpdateMessage),
		clients:      make(map[*Client]bool),
		broadcast:    make(chan []byte),
		clientsLock:  sync.Mutex{},
		storage:      storage,
		redis:        redis,
		jwt:          jwt,
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
					if err := client.waitList.Add(ctx, client, messageData); err != nil {
						log.Println("wrong save message to the wait list", err)
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

func (h *Hub) sendByUserRole(ctx context.Context, message *OrderUpdateMessage, role bl.UserRole) error {

	j, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	var (
		orderRaw *bl.DB_Order
		order    *bl.R_Order
	)

	if orderRaw, err = h.storage.OrderStorage().OrderGet(message.OrderId); err != nil {
		fmt.Println("invalid user id")
		return err
	}

	order = orderRaw.ToReal()

	for client := range h.clients {
		if client.role != role {
			continue
		}

		if _, err := client.checkAccess(h.jwt); err != nil {
			if err := client.waitList.Add(ctx, client, j); err != nil {
				log.Println("wrong save message to the wait list", err)
				return err
			}
			continue
		}

		if ok := client.CheckFilters(order); !ok {
			continue
		}

		client.send <- j
	}

	return nil
}

func (h *Hub) sendByUserId(ctx context.Context, message *OrderUpdateMessage, userId int64) error {

	j, err := json.Marshal(message)
	if err != nil {
		return err
	}

	var userClients []*Client

	for client := range h.clients {
		if client.userId == userId {
			userClients = append(userClients, client)
		}
	}

	var (
		orderRaw *bl.DB_Order
		order    *bl.R_Order
	)

	if orderRaw, err = h.storage.OrderStorage().OrderGet(message.OrderId); err != nil {
		fmt.Println("invalid user id")
		return err
	}

	order = orderRaw.ToReal()

	for _, client := range userClients {
		if _, err := client.checkAccess(h.jwt); err != nil {
			if err := client.waitList.Add(ctx, client, j); err != nil {
				log.Println("wrong save message to the wait list", err)
			}
			continue
		}

		if ok := client.CheckFilters(order); !ok {
			continue
		}

		client.send <- j
	}

	return nil
}
