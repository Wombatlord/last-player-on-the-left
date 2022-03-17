package app

import (
	"log"
	"os"
)

type Logger struct {
	logChan chan string
	name    string
}

var (
	loggers map[string]*Logger
	config  *ConfigFile
	err     error
)

func makeLogger(logChan chan string, loggerName string) *Logger {
	return &Logger{logChan: logChan, name: loggerName}
}

func (logger *Logger) listen() {
	for {
		select {
		case logString := <-logger.logChan:
			log.Printf("[%s] MESSAGE: %s", logger.name, logString)
		}
	}
}

func GetLogChan(loggerName string) chan string {
	config, err = LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	logFile, err := os.OpenFile(config.Config.Logs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)

	if loggers == nil {
		loggers = make(map[string]*Logger)
	}
	if logger, ok := loggers[loggerName]; ok {
		return logger.logChan
	} else {
		logChan := make(chan string, 32)
		logger = makeLogger(logChan, loggerName)
		loggers[loggerName] = logger

		go func() { logger.listen() }()
		logChan <- "Logging Connected!"
		return logger.logChan
	}
}
