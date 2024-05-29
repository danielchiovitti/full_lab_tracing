package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var lock sync.Mutex
var loggerManagerInstance *LoggerManager

type LoggerManager struct {
	logrus *logrus.Logger
}

func newLogrusInstance() *logrus.Logger {
	logrusInstance := logrus.New()
	logrusInstance.SetLevel(logrus.TraceLevel)
	logrusInstance.SetReportCaller(false)
	formatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logrusInstance.Out = os.Stdout
	logrusInstance.SetFormatter(formatter)
	return logrusInstance
}

func NewLoggerManager() *LoggerManager {
	if loggerManagerInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if loggerManagerInstance == nil {
			loggerManagerInstance = &LoggerManager{
				logrus: newLogrusInstance(),
			}
		}
	}

	return loggerManagerInstance
}

func (l *LoggerManager) Debug(message string) {
	l.logrus.Debug(message)
}

func (l *LoggerManager) Info(message string) {
	l.logrus.Info(message)
}

func (l *LoggerManager) Warn(message string) {
	l.logrus.Warn(message)
}

func (l *LoggerManager) Error(message string) {
	l.logrus.Error(message)
}
