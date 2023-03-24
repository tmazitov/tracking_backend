package conductor

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
)

const MAX_ATTEND_COUNT = 6

func (c *Conductor) CreateCode(ctx *gin.Context, payload CodePayload) (string, error) {

	// Check attempt count
	attemptCount, err := c.checkAuthAttempts(ctx, payload.Email)
	if err != nil {
		return "", err
	}

	// Create ticket message
	code := createCodeString()

	title := fmt.Sprintf("%s - ваш код доступа для сервиса \"Удачные перевозки\"", code)
	message := ticketMessageTemplate(code)

	// Send ticket to the email
	c.emailChan <- emailMessage{Email: payload.Email, Title: title, HTMLMessage: message}

	// Create token for ticket
	token, err := c.jwt.NewCheckToken(&jwt.CheckClaims{
		Email: payload.Email,
		Ip:    payload.Ip,
	})

	if err != nil {
		return "", err
	}

	// Create ticket record to check its later
	err = c.redis.Set(ctx, "che:"+token, code, 5*time.Minute).Err()
	if err != nil {
		return "", err
	}

	c.updateAuthAttempts(ctx, payload.Email, attemptCount+1)

	return token, nil
}

func ticketMessageTemplate(code string) string {
	return `<div>
				<h1>Код : ` + code + `</h1>
				<p>Используйте этот код для авторизации на сайте.</p>
			</div>`
}

func createCodeString() string {

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
