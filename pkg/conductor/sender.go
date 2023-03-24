package conductor

import (
	"log"

	"gopkg.in/gomail.v2"
)

type emailMessage struct {
	Email       string
	Title       string
	HTMLMessage string
}

func (c *Conductor) senderWorker(messageChan chan emailMessage) {
	defer close(messageChan)

	for {
		select {
		case m := <-messageChan:
			if err := c.sendToEmail(m.Email, m.Title, m.HTMLMessage); err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func (c *Conductor) sendToEmail(toEmail string, title string, HTMLMessage string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", c.config.Email)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", title)
	msg.SetBody("text/html", HTMLMessage)

	n := gomail.NewDialer("smtp.gmail.com", 587, c.config.Email, c.config.Pass)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
