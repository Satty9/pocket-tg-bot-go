package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	pocketClient *pocket.Client
	redirectURL  string
}

// внедрение зависимости. Эту поля можно было бы инициализировать в методе "StartBot",
// а не принимать параметры в этой функции. Так мы не зависим от значений внутри бота
func NewBot(newBot *tgbotapi.BotAPI, newPocketClient *pocket.Client, newRedirectUrl string) *Bot {
	return &Bot{
		bot:          newBot,
		pocketClient: newPocketClient,
		redirectURL:  newRedirectUrl,
	}
}

func (b *Bot) StartBot() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates := b.InitChannelUpdates()
	b.HandleUpdates(updates)
	// code below will never execute
	// gracefull shutdown поможет избежать этого. что это?
	return nil
}

func (b *Bot) InitChannelUpdates() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}

// work cycle of bot
func (b *Bot) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // If we not got a message go to next update
			continue
		}

		if update.Message.IsCommand() {
			err := b.HandleCommand(update.Message)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		err := b.HandleMessage(update.Message)
		if err != nil {
			log.Fatal(err)
		}
	}
}
