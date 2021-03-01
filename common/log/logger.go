package log

import (
	"go.uber.org/zap"
	"proxy-fileserver/configs"
)

var (
	localEnv   = "local"
	devEnv     = "dev"
	productEnv = "prod"
)

type Logger struct {
	zap *zap.SugaredLogger
}

func (l *Logger) Info(args ...interface{}) {
	l.zap.Info(args...)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.zap.Infof(msg, args...)
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
	if conf.Env == localEnv || conf.Env == devEnv {
		zapLogger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		zapSugarLogger = zapLogger.Sugar()
	} else if conf.Env == productEnv {
		zapLogger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		zapSugarLogger = zapLogger.Sugar()
	}
	return &Logger{
		zap: zapSugarLogger,
	}, nil
}
