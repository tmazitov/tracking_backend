package ws

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type WaitList struct {
	redis   *redis.Client
	archive map[*Client][]string
}

func NewWaitList(redis *redis.Client) *WaitList {
	return &WaitList{
		redis:   redis,
		archive: make(map[*Client][]string),
	}
}

func (w *WaitList) Add(ctx context.Context, client *Client, update []byte) error {
	var uuid string = uuid.New().String()

	if _, ok := w.archive[client]; !ok {
		w.archive[client] = []string{}
	}

	w.archive[client] = append(w.archive[client], uuid)

	return w.redis.Set(ctx, "wt:"+uuid, string(update), time.Minute*5).Err()
}

func (w *WaitList) DelALL(client *Client) {
	w.archive = make(map[*Client][]string)
}

func (w *WaitList) GetAll(ctx context.Context, client *Client) ([][]byte, error) {

	var result [][]byte = [][]byte{}

	for _, uuid := range w.archive[client] {
		jsonMessage, err := w.redis.Get(ctx, "wt:"+uuid).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, err
		}

		result = append(result, []byte(jsonMessage))
	}

	w.archive = make(map[*Client][]string)

	return result, nil
}
