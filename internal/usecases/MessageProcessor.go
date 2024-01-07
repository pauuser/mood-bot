package usecases

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MessageProcessor interface {
	Process(update tgbotapi.Update)
}
