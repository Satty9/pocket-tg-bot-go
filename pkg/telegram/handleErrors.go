package telegram

import (
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errURLNotValid          = errors.New("url is invalid")
	errUserNotAuthorized    = errors.New("user is not authorized")
	errPocketAnswerNegative = errors.New("unable to save URL in pocket")
)

const (
	URLNotValid          = "Ссылка неверна. Проверьте её"
	UserNotAuthorized    = "Вы не авторизированы. Используйте команду /start"
	PocketAnswerNegative = "Сервис Pocket не принял запрос. Попробуйте позже"
	UnknowError          = "Произошла неизвестная ошибка"
)

func (b *Bot) HandleErrors(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, "")

	switch err {
	case errURLNotValid:
		msg.Text = URLNotValid
	case errUserNotAuthorized:
		msg.Text = UserNotAuthorized
	case errPocketAnswerNegative:
		msg.Text = PocketAnswerNegative
	default:
		msg.Text = UnknowError
	}

	if _, err = b.bot.Send(msg); err != nil {
		log.Fatal(err)
	}

}
