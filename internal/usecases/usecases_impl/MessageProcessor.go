package usecases_impl

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"pauuser/mood-bot/internal/models"
	"pauuser/mood-bot/internal/repository"
	"pauuser/mood-bot/internal/usecases"
	"time"
)

type messageProcessorImpl struct {
	logger             *zap.Logger
	botService         usecases.BotService
	userRepository     repository.UserRepository
	questionRepository repository.QuestionRepository
	questions          map[string][]string
}

func NewMessageProcessorImpl(logger *zap.Logger,
	botService usecases.BotService,
	userRepository repository.UserRepository,
	questionRepository repository.QuestionRepository,
	questions map[string][]string) usecases.MessageProcessor {
	return &messageProcessorImpl{logger: logger,
		botService:         botService,
		userRepository:     userRepository,
		questionRepository: questionRepository,
		questions:          questions,
	}
}

func (m *messageProcessorImpl) Process(update tgbotapi.Update) {
	if update.Message != nil {
		chatId := update.Message.Chat.ID
		switch update.Message.Text {
		case "/start":
			user, err := m.userRepository.GetUser(chatId)
			if user != nil {
				m.botService.SendMessage(chatId, "Привет! Давно не виделись!")
			} else {
				var user = models.User{
					ID:       0,
					ChatId:   chatId,
					Name:     update.Message.From.FirstName + " " + update.Message.From.LastName,
					Username: update.Message.From.UserName,
				}
				err = m.userRepository.Create(&user)
				if err != nil {
					m.botService.SendMessage(chatId, "Кажется, что-то пошло не так...")
				} else {
					m.botService.SendMessage(chatId, "Привет! Приятно познакомиться!")
				}
			}
		default:
			m.botService.SendMessage(update.Message.Chat.ID, "Не распознал, что Вы хотели!")
			return
		}
	} else if update.CallbackQuery != nil {
		data := update.CallbackQuery.Data
		text := update.CallbackQuery.Message.Text
		callbackId := update.CallbackQuery.ID
		chatId := update.CallbackQuery.Message.Chat.ID

		var question = models.Question{
			ID:           0,
			QuestionText: text,
			Answer:       data,
			Date:         time.Now(),
			FromChatId:   chatId,
		}
		err := m.questionRepository.Create(&question)
		if err != nil {
			m.botService.SendMessage(chatId, "Что-то пошло не так...")
			return
		}

		err = m.botService.Request(callbackId, data)
		if err != nil {
			m.botService.SendMessage(chatId, "Что-то пошло не так...")
			return
		}

		response := "Записал ответ " + data + " на вопрос \"" + text + "\". Спасибо! 😘"
		m.botService.SendMessage(chatId, response)
	}
}
