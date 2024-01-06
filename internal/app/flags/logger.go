package flags

import (
	"os"
	"time"
)

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerFlags struct {
	TimeKey     string   `mapstructure:"time_key"`
	Level       string   `mapstructure:"level"`
	OutputPaths []string `mapstructure:"output_paths"`
}

func (l *LoggerFlags) NewZapLogger() (*zap.Logger, error) {
	loggerConfig := zap.NewProductionConfig()

	var logLevel zapcore.Level
	switch l.Level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warning":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	case "panic":
		logLevel = zapcore.PanicLevel
	case "fatal":
		logLevel = zapcore.FatalLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	loggerConfig.Level = zap.NewAtomicLevelAt(logLevel)
	loggerConfig.EncoderConfig.TimeKey = l.TimeKey
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.OutputPaths = l.OutputPaths

	encoderCfg := zap.NewProductionEncoderConfig()
	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

	ws, _ := os.OpenFile("log.log", os.O_CREATE, 0666)
	bws := &zapcore.BufferedWriteSyncer{
		WS:            ws,
		Size:          1024,        // размер буфера перед записью в ws
		FlushInterval: time.Minute, // время, после которого будет запись
	}
	defer bws.Stop()
	core := zapcore.NewCore(jsonEncoder, bws, zapcore.InfoLevel)

	logger := zap.New(core)

	logger.Info("Successful log build!")

	return loggerConfig.Build()
}
