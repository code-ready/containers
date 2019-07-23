package logging

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/code-ready/crc/pkg/crc/constants"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	logfile       *os.File
	LogLevel      string
	originalHooks = logrus.LevelHooks{}
)

type fileHook struct {
	file      io.Writer
	formatter logrus.Formatter
	level     logrus.Level
}

func newFileHook(file io.Writer, level logrus.Level, formatter logrus.Formatter) *fileHook {
	return &fileHook{
		file:      file,
		formatter: formatter,
		level:     level,
	}
}

func (h fileHook) Levels() []logrus.Level {
	var levels []logrus.Level
	for _, level := range logrus.AllLevels {
		if level <= h.level {
			levels = append(levels, level)
		}
	}

	return levels
}

func (h *fileHook) Fire(entry *logrus.Entry) error {
	line, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = h.file.Write(line)
	return err
}

func OpenLogFile() (*os.File, error) {
	l, err := os.OpenFile(filepath.Join(constants.LogFilePath), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func setupFileHook() error {
	logfile, err := OpenLogFile()
	if err != nil {
		return err
	}

	logrus.AddHook(newFileHook(logfile, logrus.TraceLevel, &logrus.TextFormatter{
		DisableColors:          true,
		DisableTimestamp:       false,
		FullTimestamp:          true,
		DisableLevelTruncation: false,
	}))
	return nil
}

func SetupFileHook() {
	err := setupFileHook()
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "Failed to open logfile"))
	}
}

func RemoveFileHook() {
	CloseLogFile()
	logrus.StandardLogger().ReplaceHooks(originalHooks)
}

func CloseLogFile() {
	logfile.Close()
}

func CloseLogging() {
	CloseLogFile()
	logrus.StandardLogger().ReplaceHooks(make(logrus.LevelHooks))
}

func InitLogrus(logLevel string) {
	logrus.SetOutput(ioutil.Discard)
	logrus.SetLevel(logrus.TraceLevel)
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}

	logrus.AddHook(newFileHook(os.Stdout, level, &logrus.TextFormatter{
		// Setting ForceColors is necessary because logrus.TextFormatter determines
		// whether or not to enable colors by looking at the output of the logger.
		// In this case, the output is ioutil.Discard, which is not a terminal.
		// Overriding it here allows the same check to be done, but against the
		// hook's output instead of the logger's output.
		ForceColors:            terminal.IsTerminal(int(os.Stderr.Fd())),
		DisableTimestamp:       true,
		DisableLevelTruncation: false,
	}))

	for k, v := range logrus.StandardLogger().Hooks {
		originalHooks[k] = v
	}
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Infof(s string, args ...interface{}) {
	logrus.Infof(s, args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Warnf(s string, args ...interface{}) {
	logrus.Warnf(s, args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func Fatalf(s string, args ...interface{}) {
	logrus.Fatalf(s, args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Errorf(s string, args ...interface{}) {
	logrus.Errorf(s, args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Debugf(s string, args ...interface{}) {
	logrus.Debugf(s, args...)
}
