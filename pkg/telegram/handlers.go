package telegram

import (
	"context"
	"log"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	commandStart           = "start"
	replyStart             = "Привет! Чтобы сохранять ссылки в своём Pocket аккаунте, вам нужно дать доступ к этому аккаунту. Для этого, пожалуйста, перейдите по ссылке \n%s"
	replyAlreadyAuthorized = "Вы уже авторизированы и готовы к работе. Присылайте ссылку для сохранения в аккаунт Pocket"
	URLSuccessSaved        = "Ссылка сохранена в Pocket аккаунт"
)

func (b *Bot) HandleMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		return errURLNotValid
	}

	accessTocken, err := b.GetAccessToken(message.Chat.ID)
	if err != nil {
		return errUserNotAuthorized
	}
	// Here we save URL
	// pocket.AddInput - the struct, to send values
	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessTocken,
		URL:         message.Text,
	}); err != nil {
		return errPocketAnswerNegative
	}

	// if all Okay
	msg := tgbotapi.NewMessage(message.Chat.ID, URLSuccessSaved)
	//msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	//msg.ReplyToMessageID = message.Message.MessageID
	_, err = b.bot.Send(msg)
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
	_, err := b.GetAccessToken(message.Chat.ID)
	// if user not authorized, start authorization process
	if err != nil {
		return b.InitAuthorizationProcess(message)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID,
		replyAlreadyAuthorized)
	b.bot.Send(msg)
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
