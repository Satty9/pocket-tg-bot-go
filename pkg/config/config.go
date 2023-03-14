package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerURL     string
	TelegramBotURL    string `mapstructure:"bot_url"` // key "bot_url"
	DBPath            string `mapstructure:"db_file"` // key "db_file"

	Messages Messages
}

type Messages struct {
	Errors    Errors
	Responses Responses
}

type Errors struct {
	ErrDefault                string `mapstructure:"err_default"`
	ErrURLNotValid            string `mapstructure:"err_url_not_valid"`
	ErrUserNotAuthorized      string `mapstructure:"err_user_not_authorized"`
	ErrPocketResponseNegative string `mapstructure:"err_pocket_response_negative"`
}

type Responses struct {
	ReplyStart             string `mapstructure:"reply_start"`
	ReplyAlreadyAuthorized string `mapstructure:"reply_already_authorized"`
	ReplyURLSuccessSaved   string `mapstructure:"reply_url_success_saved"`
	ReplyUnknowCommand     string `mapstructure:"reply_unknow_command"`
}

func InitConfig() (*Config, error) {
	viper.AddConfigPath("configs") // directory
	//viper.SetConfigType("yaml")
	viper.SetConfigName("main") // file

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.reponses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := ParseEnvVariables(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ParseEnvVariables(cfg *Config) error {
	// искусственно создать переменные окружения
	os.Setenv("BOT_TOKEN", "5709016802:AAEFLOXWNAqyCOVVnghXDC8kQmbCZBZkyCk")
	os.Setenv("POCKET_CONSUMER_KEY", "104497-4978de7197dee0ac0bd2fcf")
	os.Setenv("AUTHORIZATION_SERVER_URL", "http://localhost")

	// bind viper key to env variable
	if err := viper.BindEnv("bot_token"); err != nil { // auto change to up case
		return err
	}

	if err := viper.BindEnv("POCKET_CONSUMER_KEY"); err != nil {
		return err
	}

	if err := viper.BindEnv("AUTHORIZATION_SERVER_URL"); err != nil {
		return err
	}

	// write env variables to struct
	cfg.TelegramToken = viper.GetString("BOT_TOKEN")
	cfg.PocketConsumerKey = viper.GetString("POCKET_CONSUMER_KEY")
	cfg.AuthServerURL = viper.GetString("AUTHORIZATION_SERVER_URL")

	return nil
}
