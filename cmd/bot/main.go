package main

import (
	"log"

	// к сожалению эта библиотека не предоставляет интерфейсов, а построена на структурах
	// чтобы её замокать данный клиент, чтобы покрыть бота тестами - это головная боль.
	// нужно ли переписывать эту библиотеку, либо писать свою, используя интерфейсы и
	// которую можно будет потом лего обработать моком
	"github.com/boltdb/bolt"                                      // DB
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // tg bot
	"github.com/satty9/pocket-tg-bot-go/pkg/config"
	"github.com/satty9/pocket-tg-bot-go/pkg/repository"
	"github.com/satty9/pocket-tg-bot-go/pkg/repository/boltdb" // DB
	"github.com/satty9/pocket-tg-bot-go/pkg/server"
	"github.com/satty9/pocket-tg-bot-go/pkg/telegram" // my git
	"github.com/zhashkevych/go-pocket-sdk"            // pocket
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	// токены нельзя вставлять в код. Отрубать руки. Токен должен храниться безопасно
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	// // токены нельзя вставлять в код. Отрубать руки. Токен должен храниться безопасно
	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := InitDB()
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewBoltDB(db)

	tgBot := telegram.NewBot(bot, pocketClient, *tokenRepository, cfg.AuthServerURL)

	authorizationServer := server.NewAuthorizationServer(pocketClient, *tokenRepository, cfg.TelegramBotURL)

	// anonymous function
	go func() {
		// func "StartBot()" block execution. Use goroutine
		if err := tgBot.StartBot(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func InitDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	})

	return db, err
}
