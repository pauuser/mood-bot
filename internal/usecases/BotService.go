package usecases

type BotService interface {
	SendMessage(toChatId int64, message string)
	SendQuestion(toChatId int64, message string, buttons []string)
	Request(callbackId string, data string) error
}
