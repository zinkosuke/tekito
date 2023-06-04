package logging

import "go.uber.org/zap"

type Logger = *zap.Logger

func getLogger() Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger
}

var Log = getLogger()
