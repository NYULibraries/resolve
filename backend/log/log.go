package log

import (
	"fmt"
	"golang.org/x/exp/slog"
	"math"
	"os"
	"reflect"
)

type Level int

var (
	LevelDebug    = Level(reflect.ValueOf(slog.LevelDebug).Int())
	LevelInfo     = Level(reflect.ValueOf(slog.LevelInfo).Int())
	LevelWarn     = Level(reflect.ValueOf(slog.LevelWarn).Int())
	LevelError    = Level(reflect.ValueOf(slog.LevelWarn).Int())
	LevelDisabled = Level(math.MaxInt)
)

var programLevel = new(slog.LevelVar)
var slogger *slog.Logger

func init() {
	handler := slog.HandlerOptions{Level: programLevel}.NewJSONHandler(os.Stdout)
	slog.SetDefault(slog.New(handler))

	slogger = slog.New(handler)
}

func Fatal(args ...interface{}) {
	Error(fmt.Sprint(args...))
	os.Exit(1)
}

func Error(message string, args ...interface{}) {
	slogger.Error(message, args...)
}

func Info(message string, args ...interface{}) {
	slogger.Info(message, args...)
}

func SetLevel(level Level) {
	programLevel.Set(slog.Level(level))
}
