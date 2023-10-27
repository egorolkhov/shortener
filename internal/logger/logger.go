package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func InitLogger() error {
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zap.InfoLevel)

	logger, err := cfg.Build()
	if err != nil {
		return err
	}
	defer logger.Sync()

	Log = logger

	return nil
}
