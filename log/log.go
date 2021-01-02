package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	Logger = NewLogger("")
	Level  string
)

func NewLogger(in string) log.Logger {
	lvl := autoLevel(in)

	l := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	l = log.With(l, "ts", log.DefaultTimestampUTC)

	switch lvl {
	case "debug":
		l = log.With(l, "caller", log.Caller(6))
		l = level.NewFilter(l, level.AllowDebug())
	case "error":
		l = level.NewFilter(l, level.AllowError())
	default:
		l = level.NewFilter(l, level.AllowInfo())
	}
	Level = lvl
	return l
}

// we do this env overwriting here because even if viper take cares of it,
// we want env overwriting in tests to use debug level in tests
// 1. environment variable SHYNC_LOGLEVEL
// 2. config
// 3. defaults to "info"
func autoLevel(v string) string {
	e := strings.ToLower(os.Getenv("SHYNC_LOGLEVEL"))
	if e == "debug" || e == "info" || e == "error" {
		return e
	}
	if v == "debug" || v == "info" || v == "error" {
		return v
	}
	e = "info"
	return e
}

func IsDebug() bool {
	return Level == "debug"
}

type msgWriter struct {
	log.Logger
}

func (w msgWriter) Write(p []byte) (int, error) {
	err := w.Logger.Log("msg", string(p))
	return len(p), err
}

func Writer() io.Writer {
	return msgWriter{Logger}
}

func Debug(keyvals ...interface{}) error {
	return level.Debug(Logger).Log(keyvals...)
}

func Debugf(format string, a ...interface{}) error {
	return level.Debug(Logger).Log("msg", fmt.Sprintf(format, a...))
}

func Debugfn(format string, a ...interface{}) {
	level.Debug(Logger).Log("msg", fmt.Sprintf(format, a...))
}

func Info(keyvals ...interface{}) error {
	return level.Info(Logger).Log(keyvals...)
}

func Infof(format string, a ...interface{}) error {
	return level.Info(Logger).Log("msg", fmt.Sprintf(format, a...))
}

func Error(keyvals ...interface{}) error {
	return level.Error(Logger).Log(keyvals...)
}

func Errorf(format string, a ...interface{}) error {
	return level.Error(Logger).Log("msg", fmt.Sprintf(format, a...))
}
