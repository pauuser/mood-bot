package usecases_impl

import (
	"go.uber.org/zap"
	"gopkg.in/robfig/cron.v2"
	"pauuser/mood-bot/internal/models"
	"pauuser/mood-bot/internal/repository"
	"pauuser/mood-bot/internal/usecases"
	"strconv"
	"strings"
	"time"
)

type statisticsCronJob struct {
	userRepository     repository.UserRepository
	questionRepository repository.QuestionRepository
	botService         usecases.BotService
	logger             *zap.Logger
}

func NewStatisticsCronJob(userRepository repository.UserRepository,
	questionRepository repository.QuestionRepository,
	botService usecases.BotService,
	logger *zap.Logger) usecases.StatisticsCronJob {
	return &statisticsCronJob{
		userRepository:     userRepository,
		questionRepository: questionRepository,
		botService:         botService,
		logger:             logger,
	}
}

func (s *statisticsCronJob) RunSendStatisticsJob() {
	cr := cron.New()

	cr.AddFunc("@every 12s", func() {
		s.sendStatistics()
	})

	cr.Start()
}

func (s *statisticsCronJob) sendStatistics() {
	answers, err := s.questionRepository.GetAll()
	if err != nil {
		s.logger.Error("Could not get all questions from database")
		return
	}

	results := s.GetStatistics(answers)
	for chatId, questions := range results {
		user, _ := s.userRepository.GetUser(chatId)

		var sb strings.Builder
		sb.WriteString("Пришло время посчитать статистику за месяц!")
		for questionText, countOverall := range questions.overallCount {
			sb.WriteString("\n\nВопрос: ")
			sb.WriteString(questionText)

			var totalAverage = getAverage(questions.overallSum[questionText], countOverall)
			sb.WriteString("\n\nСреднее за все время: ")
			sb.WriteString(strconv.FormatFloat(totalAverage, 'f', 2, 64))

			var countLastMonth = questions.lastMonthCount[questionText]
			if countLastMonth > 0 {
				var monthAverage = getAverage(questions.lastMonthSum[questionText], countLastMonth)
				sb.WriteString("\nСреднее за месяц: ")
				sb.WriteString(strconv.FormatFloat(monthAverage, 'f', 2, 64))
			}
		}

		s.botService.SendMessage(user.ChatId, sb.String())
	}
}

func getAverage(sum int64, count int64) float64 {
	return float64(sum) / float64(count)
}

func (s *statisticsCronJob) GetStatistics(answers []*models.Question) map[int64]*userResult {
	currentMonth := getMonthNumber(time.Now())

	var result map[int64]*userResult = make(map[int64]*userResult)
	for _, answer := range answers {
		answerMonth := getMonthNumber(answer.Date)
		answerScore, err := strconv.ParseInt(answer.Answer, 10, 32)
		if err != nil {
			s.logger.Error("Could not convert response")
			continue
		}

		val, ok := result[answer.FromChatId]
		if !ok {
			result[answer.FromChatId] = newUserResult()
			val = result[answer.FromChatId]
		}

		val.overallSum[answer.QuestionText] += answerScore
		val.overallCount[answer.QuestionText] += 1
		if answerMonth == currentMonth {
			val.lastMonthSum[answer.QuestionText] += answerScore
			val.lastMonthCount[answer.QuestionText] += 1
		}
	}

	return result
}

type userResult struct {
	overallSum     map[string]int64
	overallCount   map[string]int64
	lastMonthSum   map[string]int64
	lastMonthCount map[string]int64
}

func newUserResult() *userResult {
	return &userResult{
		overallSum:     make(map[string]int64),
		overallCount:   make(map[string]int64),
		lastMonthSum:   make(map[string]int64),
		lastMonthCount: make(map[string]int64),
	}
}

func getMonthNumber(time time.Time) int {
	return int(time.Month())
}
