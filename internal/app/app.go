package app

import (
	"log"
	"pauuser/mood-bot/internal/app/flags"
	"pauuser/mood-bot/internal/repository"
	"pauuser/mood-bot/internal/repository/repository_impl"
	"pauuser/mood-bot/internal/usecases"
	"pauuser/mood-bot/internal/usecases/usecases_impl"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Sqlite    flags.SqliteFlags   `mapstructure:"sqlite"`
	Logger    flags.LoggerFlags   `mapstructure:"logger"`
	Bot       flags.BotFlags      `mapstructure:"bot"`
	Questions map[string][]string `mapstructure:"questions"`
	Schedule  string              `mapstructure:"schedule"`
}

type App struct {
	Config   Config
	Repos    *appReposFields
	UseCases *appUseCasesFields
	bot      *tgbotapi.BotAPI
	Logger   *zap.Logger
}

type appReposFields struct {
	QuestionRepository repository.QuestionRepository
	UserRepository     repository.UserRepository
}

type appUseCasesFields struct {
	BotService             usecases.BotService
	QuestionsBackgroundJob usecases.QuestionsCronJob
	MessageProcessor       usecases.MessageProcessor
	StatisticsCronJob      usecases.StatisticsCronJob
}

func (a *App) initRepos() (*appReposFields, error) {
	db, err := a.Config.Sqlite.InitDB()
	if err != nil {
		return nil, err
	}

	f := &appReposFields{}

	f.QuestionRepository = repository_impl.NewQuestionRepoSqliteImpl(db, a.Logger)
	f.UserRepository = repository_impl.NewUserRepoSqliteImpl(db, a.Logger)

	return f, nil
}

func (a *App) initUseCases(repos *appReposFields) *appUseCasesFields {
	u := &appUseCasesFields{}

	u.BotService = usecases_impl.NewBotServiceUseCaseImpl(a.bot, a.Logger)
	u.QuestionsBackgroundJob = usecases_impl.NewQuestionCronJob(a.Config.Questions, a.Config.Schedule, repos.UserRepository, u.BotService, a.Logger)
	u.MessageProcessor = usecases_impl.NewMessageProcessorImpl(a.Logger, u.BotService, repos.UserRepository, repos.QuestionRepository, a.Config.Questions)
	u.StatisticsCronJob = usecases_impl.NewStatisticsCronJob(repos.UserRepository, repos.QuestionRepository, u.BotService, a.Logger)

	return u
}

func (a *App) ParseConfig(pathToConfig string, configFileName string) error {
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType("json")
	v.AddConfigPath(pathToConfig)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&a.Config)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Init(logger *zap.Logger) error {
	bot, err := a.Config.Bot.NewBot()
	if err != nil {
		logger.Fatal("Cannot init telegram bot")
		return err
	}
	a.bot = bot

	repos, err := a.initRepos()

	if err != nil {
		logger.Fatal("error init repo", zap.Error(err))
		return err
	}

	a.Repos = repos
	a.UseCases = a.initUseCases(repos)

	return nil
}

func (a *App) Run() {
	logger, err := a.Config.Logger.NewZapLogger()
	if err != nil {
		log.Fatal(err)
		return
	}
	a.Logger = logger
	defer func() {
		if err := a.Logger.Sync(); err != nil {
			log.Fatal("error logger sync")
		}
	}()

	err = a.Init(logger)
	if err != nil {
		log.Fatal(err)
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	go a.UseCases.QuestionsBackgroundJob.RunQuestionsCronJob()
	go a.UseCases.StatisticsCronJob.RunSendStatisticsJob()

	updates := a.bot.GetUpdatesChan(u)
	for update := range updates {
		a.UseCases.MessageProcessor.Process(update)
	}
}
