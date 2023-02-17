package main

import (
	"log"

	// к сожалению эта библиотека не предоставляет интерфейсов, а построена на структурах
	// чтобы её замокать данный клиент, чтобы покрыть бота тестами - это головная боль.
	// нужно ли переписывать эту библиотеку, либо писать свою, используя интерфейсы и
	// которую можно будет потом лего обработать моком
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" //tg bot
	"github.com/satty9/pocket-tg-bot-go/pkg/telegram"             // my git
	"github.com/zhashkevych/go-pocket-sdk"                        // pocket
)

func main() {
	// токены нельзя вставлять в код. Отрубать руки. Токен должен храниться безопасно
	bot, err := tgbotapi.NewBotAPI("5709016802:AAEFLOXWNAqyCOVVnghXDC8kQmbCZBZkyCk")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	// // токены нельзя вставлять в код. Отрубать руки. Токен должен храниться безопасно
	pocketClient, err := pocket.NewClient("104497-4978de7197dee0ac0bd2fcf")
	if err != nil {
		log.Fatal(err)
	}

	tgBot := telegram.NewBot(bot, pocketClient, "http://localhost")
	if err := tgBot.StartBot(); err != nil {
		log.Fatal(err)
	}
}
