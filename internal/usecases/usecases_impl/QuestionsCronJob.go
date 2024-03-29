package usecases_impl

import (
	"go.uber.org/zap"
	"gopkg.in/robfig/cron.v2"
	"pauuser/mood-bot/internal/repository"
	"pauuser/mood-bot/internal/usecases"
)

type questionCronJob struct {
	questions      map[string][]string
	userRepository repository.UserRepository
	botService     usecases.BotService
	logger         *zap.Logger
}

func NewQuestionCronJob(questions map[string][]string,
	userRepository repository.UserRepository,
	botService usecases.BotService,
	logger *zap.Logger) usecases.QuestionsCronJob {
	return &questionCronJob{
		questions:      questions,
		userRepository: userRepository,
		botService:     botService,
		logger:         logger,
	}
}

func (q *questionCronJob) sendQuestions() {
	questions := q.questions
	users, err := q.userRepository.GetAll()
	if err != nil {
		q.logger.Error("Error getting users")
	}

	for _, user := range users {
		q.botService.SendMessage(user.ChatId, "Привет! Пора ответить на пару вопросов 🥺")
		for question, answers := range questions {
			q.botService.SendQuestion(user.ChatId, question, answers)
		}
	}
}

func (q *questionCronJob) RunQuestionsCronJob() {
	s := cron.New()

	s.AddFunc("0 0 22 * * ?", func() {
		q.sendQuestions()
	})

	s.Start()
}
