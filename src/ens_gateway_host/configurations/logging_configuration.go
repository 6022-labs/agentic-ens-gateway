package configurations

import (
	"os"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func ConfigureLogging(container *dig.Container) {
	logLevel := os.Getenv("LOG_LEVEL")

	zapConfig := zap.NewProductionConfig()

	switch logLevel {
	case "debug":
		zapConfig.Level.SetLevel(zap.DebugLevel)
	case "info", "information":
		zapConfig.Level.SetLevel(zap.InfoLevel)
	case "warn", "warning":
		zapConfig.Level.SetLevel(zap.WarnLevel)
	case "error":
		zapConfig.Level.SetLevel(zap.ErrorLevel)
	case "fatal":
		zapConfig.Level.SetLevel(zap.FatalLevel)
	case "panic":
		zapConfig.Level.SetLevel(zap.PanicLevel)
	default:
		zapConfig.Level.SetLevel(zap.InfoLevel)
	}

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() *zap.Logger {
		return logger
	})
	if err != nil {
		logger.Fatal("failed to provide logger", zap.Error(err))
	}
}
