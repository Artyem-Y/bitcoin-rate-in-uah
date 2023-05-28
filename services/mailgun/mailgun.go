package mailgun

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"test/bitcoin-rate-in-uah/conf"
	"test/bitcoin-rate-in-uah/utils"

	mailgun "github.com/mailgun/mailgun-go/v4"
)

var MgApiKey = conf.GetEnvConst("MAILGUN_API_KEY")
var MgDomain = conf.GetEnvConst("MAILGUN_DOMAIN")

func SendMail(sender, recipient, subject, text string) (bool, error) {
	fmt.Printf("Recepient: %s\n", recipient)

	if !utils.ValidateEmail(recipient) {
		return false, errors.New("email address recipient is invalid")
	}
	mg := mailgun.NewMailgun(MgDomain, MgApiKey)
	message := mg.NewMessage(
		sender,
		subject,
		text,
		recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Print(err)
		return false, err
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return true, nil
}
