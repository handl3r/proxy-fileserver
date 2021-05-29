package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"proxy-fileserver/common/log/hooks"
	"proxy-fileserver/configs"
	"strings"
)

var (
	localEnv   = "local"
	devEnv     = "dev"
	productEnv = "prod"
)
var _Hook *hooks.TelegramHook

type Logger struct {
	zap *zap.SugaredLogger
}

func (l *Logger) Info(args ...interface{}) {
	l.zap.Info(args...)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.zap.Infof(msg, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.zap.Warn(args...)
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.zap.Warnf(msg, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.zap.Error(args...)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.zap.Errorf(msg, args...)
}

func (l Logger) Fatal(args ...interface{}) {
	panic("implement me")
}

func (l Logger) Fatalf(msg string, args ...interface{}) {
	panic("implement me")
}

func NewLogger() (*Logger, error) {
	var zapSugarLogger *zap.SugaredLogger
	conf := configs.Get()
	options := make([]zap.Option, 0)
	options = append(options, zap.AddCallerSkip(2))
	if conf.TelegramBotConfig.BotToken != "" {
		_Hook = hooks.NewTelegramHook(conf.TelegramBotConfig.BaseURL, conf.TelegramBotConfig.BotToken, conf.TelegramBotConfig.ChannelID)
		hookOption := zap.Hooks(func(e zapcore.Entry) error {
			if e.Level != zapcore.ErrorLevel && e.Level != zapcore.WarnLevel {
				return nil
			}
			message := fmt.Sprintf("<u>[SERVER-%s]\n[LEVEL:%s]</u>\n<u>[MESSAGE]</u> %s\n<u>[STACK]</u>\n %s", conf.Env,
				strings.ToUpper(e.Level.String()), e.Message, e.Stack)
			go func() {
				_ = _Hook.SendMessage(message)
			}()
			return nil
		})
		options = append(options, hookOption)
	}

	if conf.Env == localEnv || conf.Env == devEnv {
		zapLogger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		zapSugarLogger = zapLogger.WithOptions(options...).Sugar()
	} else if conf.Env == productEnv {
		zapLogger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		zapSugarLogger = zapLogger.WithOptions(options...).Sugar()
	}
	return &Logger{
		zap: zapSugarLogger,
	}, nil
}
