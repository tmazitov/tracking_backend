package conductor

import "log"

type emailMessage struct {
	Email  string
	Ticket string
}

func (c *Conductor) senderWorker(messageChan chan emailMessage) {
	defer close(messageChan)

	for {
		select {
		case m := <-messageChan:
			if err := c.sendToEmail(m.Email, m.Ticket); err != nil {
				log.Println(err.Error())
			}
		}
	}
}
