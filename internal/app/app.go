package app

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"pauuser/mood-bot/internal/app/flags"
	"pauuser/mood-bot/internal/repository"
	"pauuser/mood-bot/internal/repository/impl"
)

type Config struct {
	Postgres flags.SqliteFlags `mapstructure:"sqlite"`
	Logger   flags.LoggerFlags `mapstructure:"logger"`
	Bot      flags.BotFlags    `mapstructure:"bot"`
}

type App struct {
	Config   Config
	repos    *appReposFields
	useCases *appUseCasesFields
	bot      *tgbotapi.BotAPI
	logger   *zap.Logger
}

type appReposFields struct {
	questionRepo repository.QuestionRepository
}

type appUseCasesFields struct {
}

func (a *App) initRepos() (*appReposFields, error) {
	db, err := a.Config.Postgres.InitDB()
	if err != nil {
		return nil, err
	}

	f := &appReposFields{}

	f.questionRepo = impl.NewQuestionRepoSqliteImpl(db, a.logger)

	return f, nil
}

func (a *App) initUseCases(repos *appReposFields) *appUseCasesFields {
	u := &appUseCasesFields{}

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
	repos, err := a.initRepos()

	if err != nil {
		logger.Fatal("error init repo", zap.Error(err))
		return err
	}

	a.repos = repos
	a.useCases = a.initUseCases(a.repos)

	bot, err := a.Config.Bot.NewBot()
	if err != nil {
		logger.Fatal("Cannot init telegram bot")
		return err
	}
	a.bot = bot

	return nil
}

func (a *App) Run() {
	logger, err := a.Config.Logger.NewZapLogger()
	if err != nil {
		log.Fatal(err)
		return
	}
	a.logger = logger
	defer func() {
		if err := a.logger.Sync(); err != nil {
			log.Fatal("error logger sync")
		}
	}()

	err = a.Init(logger)
	if err != nil {
		log.Fatal(err)
		return
	}

	u := NewUpdate(0)
	u.Timeout = 60

	updates := a.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
