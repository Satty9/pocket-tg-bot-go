package telegram

import (
	"context"
	"fmt"
)

func (b *Bot) GenerateAuthorizationLink(chatID int64) (string, error) {
	redirectURL := b.GenerateRedirectURL(chatID)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), b.redirectURL)
	if err != nil {
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
