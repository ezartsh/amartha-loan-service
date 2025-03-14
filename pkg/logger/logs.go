package logger

import (
	"fmt"
	"io"
	"loan-service/config"
	"loan-service/utils"
	"math/rand"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogVariable = map[string]any

type LogHandler struct {
	Instance zerolog.Logger
}

type logStack struct {
	logger LogHandler
	Writer io.Closer
}

var LogStack = map[string]logStack{}

func ClearLogStack() {
	LogStack = map[string]logStack{}
}

func Register() {
	InitAppLog()
}

func NewFile(fileName string) LogHandler {
	filePath := utils.ConcatPaths(config.Env.LogDirectoryPath, fileName)
	if val, ok := LogStack[filePath]; ok {
		return val.logger
	}

	fileLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    50, //
		MaxBackups: 10,
		MaxAge:     14,
		Compress:   true,
		LocalTime:  true,
	}

	logger := zerolog.New(fileLogger).With().Timestamp().Logger()

	handler := LogHandler{logger}

	LogStack[filePath] = logStack{
		logger: handler,
		Writer: fileLogger,
	}

	return handler
}

func NewFileConsoleStructured(fileName string) LogHandler {
	filePath := utils.ConcatPaths(config.Env.LogDirectoryPath, fileName)
	if val, ok := LogStack[filePath]; ok {
		return val.logger
	}

	var stdOut = os.Stdout

	fileLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    50, //
		MaxBackups: 10,
		MaxAge:     14,
		Compress:   true,
		LocalTime:  true,
	}

	multiOutput := io.MultiWriter(stdOut, fileLogger)

	logger := zerolog.New(multiOutput).With().Timestamp().Logger()

	handler := LogHandler{logger}

	LogStack[filePath] = logStack{
		logger: handler,
		Writer: fileLogger,
	}

	randomizer := rand.New(rand.NewSource(10))
	randomInt := fmt.Sprintf("%d", randomizer.Uint32())
	LogStack[randomInt] = logStack{
		logger: handler,
		Writer: stdOut,
	}

	return handler
}

func (l LogHandler) Warn(m ...any) {
	var messages = l.messageToString(m)
	l.Instance.Warn().Msg(strings.Join(messages, " "))
}

func (l LogHandler) Warnf(format string, m ...any) {
	l.Instance.Warn().Msgf(format, m...)
}

func (l LogHandler) Error(e error, m ...any) {
	var messages = l.messageToString(m)
	l.Instance.Error().Err(e).Msg(strings.Join(messages, " "))
}

func (l LogHandler) Errorf(e error, format string, m ...any) {
	l.Instance.Error().Err(e).Msgf(format, m...)
}

func (l LogHandler) Fatal(e error, m ...any) {
	var messages = l.messageToString(m)
	l.Instance.Fatal().Err(e).Msg(strings.Join(messages, " "))
}

func (l LogHandler) Fatalf(e error, format string, m ...any) {
	l.Instance.Fatal().Err(e).Msgf(format, m...)
}

func (l LogHandler) Info(m ...any) {
	var messages = l.messageToString(m)
	l.Instance.Info().Msg(strings.Join(messages, " "))
}

func (l LogHandler) Infof(format string, m ...any) {
	l.Instance.Info().Msgf(format, m...)
}

func (l LogHandler) Debug(m ...any) {
	var messages = l.messageToString(m)
	l.Instance.Debug().Msg(strings.Join(messages, " "))
}

func (l LogHandler) Debugf(format string, m ...any) {
	l.Instance.Debug().Msgf(format, m...)
}

func (l LogHandler) InfoWithVariables(vars LogVariable, m ...any) {
	var messages = l.messageToString(m)
	logEvent := l.Instance
	for key, value := range vars {
		variableValue := fmt.Sprintf("%v", value)
		logEvent = logEvent.With().Str(key, variableValue).Logger()
	}
	logEvent.Info().Msg(strings.Join(messages, " "))
}

func (l LogHandler) DebugWithVariables(vars LogVariable, m ...any) {
	var messages = l.messageToString(m)
	logEvent := l.Instance
	for key, value := range vars {
		variableValue := fmt.Sprintf("%v", value)
		logEvent = logEvent.With().Str(key, variableValue).Logger()
	}
	logEvent.Debug().Msg(strings.Join(messages, " "))
}

func (l LogHandler) messageToString(m []any) []string {
	var messages []string
	for _, msg := range m {
		messages = append(messages, fmt.Sprintf("%v", msg))
	}
	return messages
}
