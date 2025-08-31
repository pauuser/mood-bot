package usecases_impl

import (
	"pauuser/mood-bot/internal/repository"
	"pauuser/mood-bot/internal/usecases"

	"go.uber.org/zap"
	"gopkg.in/robfig/cron.v2"
)

type questionCronJob struct {
	questions      map[string][]string
	userRepository repository.UserRepository
	botService     usecases.BotService
	logger         *zap.Logger
	schedule       string
}

func NewQuestionCronJob(questions map[string][]string,
	schedule string,
	userRepository repository.UserRepository,
	botService usecases.BotService,
	logger *zap.Logger) usecases.QuestionsCronJob {
	return &questionCronJob{
		questions:      questions,
		userRepository: userRepository,
		botService:     botService,
		logger:         logger,
		schedule:       schedule,
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

	s.AddFunc(q.schedule, func() {
		q.sendQuestions()
	})

	s.Start()
}
