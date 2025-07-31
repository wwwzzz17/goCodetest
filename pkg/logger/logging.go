package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	logger = zap.New(core)
}

func Debug(msg string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(msg, args...))
}
func Info(msg string, args ...interface{}) {
	logger.Info(fmt.Sprintf(msg, args...))
}
func Warn(msg string, args ...interface{}) {
	logger.Warn(fmt.Sprintf(msg, args...))
}
func Error(msg string, args ...interface{}) {
	logger.Error(fmt.Sprintf(msg, args...))
}
func Fatal(msg string, args ...interface{}) {
	logger.Fatal(fmt.Sprintf(msg, args...))
	os.Exit(1)
}
func Panic(msg string, args ...interface{}) {
	logger.Panic(fmt.Sprintf(msg, args...))
	os.Exit(1)
}
