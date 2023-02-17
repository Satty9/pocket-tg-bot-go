package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart       = "start"
	replyStartTemplate = "Привет! Чтобы сохранять ссылки в своём Pocket аккаунте, для начала тебе необходимо дать мне на это доступ. Для этого, пожалуйста, перейдите по ссылке \n%s"
)

func (b *Bot) HandleMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	//msg.ReplyToMessageID = message.Message.MessageID

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) HandleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.HandleStartCommand(message)
	default:
		return b.HandleUnkownCommand(message)
	}
}

func (b *Bot) HandleStartCommand(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, "sent command /start")
	authLink, err := b.GenerateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err

	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(replyStartTemplate, authLink))
	_, err = b.bot.Send(msg)
	return err
}

/*
//msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	//msg := tgbotapi.NewMessage(message.Chat.ID,

str := fmt.Sprintf("%s, хелоу. Ваши координаты: долгота - %f,  широта - %f. Погрешность координат - %f метров",
message.From.UserName, message.Location.Longitude, message.Location.Latitude, message.Location.HorizontalAccuracy)
fmt.Printf("Lol")

*/

func (b *Bot) HandleUnkownCommand(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, "sent unknown command!")
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.Text = "Вы ввели неизвестную команду. Пожалуйста, выберите команду из списка"
	_, err := b.bot.Send(msg)
	return err
}
