package ws

import (
	"context"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

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

func (h *Hub) UpdateWorker(ctx context.Context, order *bl.R_Order, result interface{}) error {
	var message OrderUpdateMessage = OrderUpdateMessage{
		OrderId: order.ID,
		Type:    OrderUpdateWorker,
		Data:    result,
	}

	var messageForWorker OrderUpdateMessage = OrderUpdateMessage{
		OrderId: order.ID,
		Type:    OrderUpdateWorker,
		Data:    order,
	}

	var err error

	// Send updates for all admins
	if err = h.sendByUserRole(ctx, &message, bl.Admin); err != nil {
		return err
	}

	// Send updates to the order owner
	if err = h.sendByUserId(ctx, &message, order.Owner.ID); err != nil {
		return err
	}

	// Send order to the selected worker
	if err = h.sendByUserId(ctx, &messageForWorker, order.Worker.ID); err != nil {
		return err
	}

	return nil
}
