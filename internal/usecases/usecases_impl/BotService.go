package usecases_impl

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"pauuser/mood-bot/internal/usecases"
)

type botServiceUseCaseImpl struct {
	bot    *tgbotapi.BotAPI
	logger *zap.Logger
}

func (b botServiceUseCaseImpl) Request(callbackId string, data string) error {
	callback := tgbotapi.NewCallback(callbackId, data)
	_, err := b.bot.Request(callback)

	return err
}

func (b botServiceUseCaseImpl) SendMessage(toChatId int64, message string) {
	telegramMessage := tgbotapi.NewMessage(toChatId, message)
	_, err := b.bot.Send(telegramMessage)
	if err != nil {
		b.logger.Error("Telegram message send failed")
	}
}

func (b botServiceUseCaseImpl) SendQuestion(toChatId int64, message string, buttons []string) {
	i := 0
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	curRow := make([]tgbotapi.InlineKeyboardButton, 0)
	for i < len(buttons) {
		button := tgbotapi.NewInlineKeyboardButtonData(buttons[i], buttons[i])
		curRow = append(curRow, button)
		i++

		if i%3 == 0 || i >= len(buttons) {
			row := tgbotapi.NewInlineKeyboardRow(curRow...)
			rows = append(rows, row)
			curRow = curRow[:0]
		}
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	telegramMessage := tgbotapi.NewMessage(toChatId, message)
	telegramMessage.ReplyMarkup = keyboard

	_, err := b.bot.Send(telegramMessage)
	if err != nil {
		b.logger.Error("Telegram message send failed")
	}
}

func NewBotServiceUseCaseImpl(bot *tgbotapi.BotAPI, logger *zap.Logger) usecases.BotService {
	return &botServiceUseCaseImpl{bot, logger}
}
