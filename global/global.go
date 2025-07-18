package global

import (
	"go.uber.org/zap"
)

var (
	Log *zap.Logger
)

func Getlog() {

	Log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Log.Info("init log", zap.String("level", "info"))
	Log.Debug("This is a debug message")
	Log.Warn("This is a warning message")
	Log.Error("This is an error message")
	Log.Fatal("This is a fatal message")

}
