package ws

import (
	"encoding/json"
	"fmt"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

type ClientMessageType uint8

const (
	UpdateFilters ClientMessageType = 1
)

type ClientMessage struct {
	Access string            `json:"access"`
	Type   ClientMessageType `json:"type"`
	Data   json.RawMessage   `json:"data"`
}

func (c *Client) route(message *ClientMessage) error {

	var err error

	switch message.Type {
	case UpdateFilters:
		err = c.updateFilters(message)
	}

	return err
}

func (c *Client) updateFilters(message *ClientMessage) error {

	var (
		newFilters bl.R_OrderListFilters
		err        error
	)

	if err = json.Unmarshal(message.Data, &newFilters); err != nil {
		fmt.Println(string(message.Data))
		return errBadRequest
	}

	c.filters = newFilters

	return nil
}

func (c *Client) CheckFilters(order *bl.R_Order) bool {

	// Check order date
	filterYear, filterMonth, filterDate := c.filters.Date.Date()
	orderYear, orderMonth, orderDate := order.StartAt.Date()
	if filterDate != orderDate || filterMonth != orderMonth || filterYear != orderYear {
		return false
	}

	// Check order worker id
	if c.filters.WorkerId != 0 && c.filters.WorkerId != order.Worker.ID {
		return false
	}

	// Check order status id
	if len(c.filters.Statuses) != 0 {
		var statusIncluded bool = false
		for _, status := range c.filters.Statuses {
			if status == bl.OrderStatus(order.StatusID) {
				statusIncluded = true
				break
			}
		}
		if !statusIncluded {
			return statusIncluded
		}
	}

	// Check order regular customer
	if c.filters.IsRegularCustomer && c.filters.IsRegularCustomer == order.IsRegularCustomer {
		return false
	}

	return true
}
