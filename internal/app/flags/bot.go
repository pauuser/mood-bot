package flags

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotFlags struct {
	Token string `mapstructure:"time_key"`
}

func (b *BotFlags) NewBot() (*tgbotapi.BotAPI, error) {
	return tgbotapi.NewBotAPI(b.Token)
}
