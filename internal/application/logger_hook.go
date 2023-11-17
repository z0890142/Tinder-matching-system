package app

import (
	"tinderMatchingSystem/pkg/utils/logger"
)

func initLoggerApplicationHook(app *Application) error {
	l, err := logger.New(logger.Options{
		Level:   app.Config.LogLevel,
		Outputs: []string{app.Config.LogFile},
	})

	if err != nil {
		return err
	}

	app.Logger = l.Sugar()
	logger.SetLogger(l)
	return nil
}
