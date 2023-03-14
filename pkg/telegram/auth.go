package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/satty9/pocket-tg-bot-go/pkg/repository"
)

func (b *Bot) GenerateAuthorizationLink(chatID int64) (string, error) {
	redirectURL := b.GenerateRedirectURL(chatID)

	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), b.redirectURL)
	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

/*
func (b *Bot) GeneratePocketLink() (string, error) {
	requestTocken, err := b.pocketClient.GetRequestToken(context.Background(), b.redirectURL)
	if err != nil {
		return "", err
	}
	return b.pocketClient.GetAuthorizationURL(requestTocken, b.redirectURL)
}
*/

func (b *Bot) GenerateRedirectURL(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}

func (b *Bot) GetAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessTokens)
}

func (b *Bot) InitAuthorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.GenerateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err

	}

	msg := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf(replyStart, authLink))

	_, err = b.bot.Send(msg)
	return err
}
