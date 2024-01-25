package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type Logger struct {
	logger *slog.Logger
	file   *os.File
}

func (logger *Logger) Init(outputFile string) {
	replacer := func(groups []string, a slog.Attr) slog.Attr {

		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}

	var file, err = os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logger.file = file
	if err != nil {
		error := os.MkdirAll(outputFile, os.ModePerm)
		if error != nil {
			fmt.Println(fmt.Errorf("Error creating directory %s, error message %s", outputFile, error.Error()))
		}
		var file, err = os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		logger.file = file
		if err != nil {
			var textHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true, ReplaceAttr: replacer})
			logger.logger = slog.New(textHandler)
			logger.Info("Logging to stdout because file " + outputFile + " could not be opened")
		}
	}
	var textHandler = slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true, ReplaceAttr: replacer})
	logger.logger = slog.New(textHandler)
}

func (logger *Logger) Debug(msg string, args ...any) {
	logger.logger.Debug(msg, args...)
}

func (logger *Logger) Info(msg string, args ...any) {
	logger.logger.Info(msg, args...)
}

func (logger *Logger) Error(msg string, args ...any) {
	logger.logger.Error(msg, args...)
}

func (logger *Logger) CloseFile() {
	if logger.file != nil {
		logger.file.Close()
	}
}
