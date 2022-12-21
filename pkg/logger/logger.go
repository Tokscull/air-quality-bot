package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	*zap.Logger
}

func NewLogger() *Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	//for logging in file
	//	fileEncoder := zapcore.NewJSONEncoder(config)
	//	logFile, _ := os.OpenFile("logs/application.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//	writer := zapcore.AddSync(logFile)

	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	defaultLogLevel := zapcore.InfoLevel

	core := zapcore.NewTee(
		//	zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	return &Logger{
		zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)),
	}
}
