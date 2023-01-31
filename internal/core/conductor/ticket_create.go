package conductor

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/core/jwt"
	"gopkg.in/gomail.v2"
)

const MAX_ATTEND_COUNT = 6

func (c *Conductor) CreateTicket(ctx *gin.Context, email string) (string, error) {

	// Check attempt count
	attemptCount, err := c.checkAuthAttempts(ctx, email)
	if err != nil {
		return "", err
	}

	// Create ticket message
	ticketCode := createCode()
	ticket := ticketMessageTemplate(ticketCode)

	// Send ticket to the email
	c.emailChan <- emailMessage{Email: email, Ticket: ticket}

	// Create token for ticket
	token, err := c.jwt.NewCheckToken(&jwt.CheckClaims{
		Email: email,
		Ip:    ctx.ClientIP(),
	})

	if err != nil {
		return "", err
	}

	// Create ticket record to check its later
	err = c.redis.Set(ctx, "che:"+token, ticketCode, 5*time.Minute).Err()
	if err != nil {
		return "", err
	}

	c.updateAuthAttempts(ctx, email, attemptCount+1)

	return token, nil
}

func (c *Conductor) sendToEmail(email string, ticket string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", c.config.Email)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "\"Удачные перевозки\" : Код доступа")
	msg.SetBody("text/html", ticket)

	n := gomail.NewDialer("smtp.gmail.com", 587, c.config.Email, c.config.Pass)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func ticketMessageTemplate(code string) string {
	return `<div>
				<h1>Код : ` + code + `</h1>
				<p>Используйте этот код для авторизации на сайте.</p>
			</div>`
}

func createCode() string {

	// example : 123 456

	var code string = ""
	var codeItem string

	rand.Seed(time.Now().UnixNano())

	for index := 0; index < 6; index++ {
		codeItem = strconv.Itoa(rand.Intn(10))
		code += codeItem
	}

	return code
}
