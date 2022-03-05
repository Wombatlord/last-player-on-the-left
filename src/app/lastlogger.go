package app

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	logChan       chan string
	file *os.File
}

var (
	loggers map[string]*Logger
)

func makeLogger(logChan chan string, file *os.File) *Logger {
	return &Logger{logChan: logChan, file: file}
}

func (logger *Logger) listen() {
	defer logger.file.Close()
	for {
		select {
		case log := <-logger.logChan:
			fmt.Fprintf(logger.file, "[%s]: MESSAGE: %s\n", time.Now().Format(time.RFC3339), log)
		}
	}	
}

func GetLogChan(loggerName, path string) (chan string, error) {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("the path provided is a directory: %s", path)
	}
	
	logFile, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 
		0644,
	)
	
	if err != nil {
		return nil, err
	}
	
	if _, err := logFile.WriteString("Logging Online!"); err != nil {
		return nil, err
	}

	if loggers == nil {
		loggers = make(map[string]*Logger)
	}
	if logger, ok := loggers[loggerName]; ok {
		return logger.logChan, nil
	} else {
		logChan := make(chan string, 32)
		logger = makeLogger(logChan, logFile)
		loggers[loggerName] = logger

		go func () {logger.listen()}()
		return logger.logChan, nil
	}
}
