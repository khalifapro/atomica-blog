package logging

import (
	"crypto/rand"
	"fmt"
	"os"
	"runtime"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// init initializes the logger
func init() {
	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Only log the warning severity or above.
	// Default log level
	logger.SetLevel(logrus.DebugLevel)

	EnvLogLevel := os.Getenv("LOG_LEVEL")
	if EnvLogLevel == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else if EnvLogLevel == "info" {
		logger.SetLevel(logrus.InfoLevel)
	} else if EnvLogLevel == "warn" {
		logger.SetLevel(logrus.WarnLevel)
	}
}

// WithField log message with field
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

// WithFields logs a message with fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

// Trace returns the source code line and function name (of the calling function)
func Trace() (line string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return fmt.Sprintf("%s,:%d %s\n", frame.File, frame.Line, frame.Function)
}

// GenerateUUID is function to generate our own uuid if the google uuid throws error
func GenerateUUID() string {
	log.Info("entering func generateUUID")
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Error(Trace(), err)
		return ""
	}
	theUUID := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return theUUID
}

// GetRequestID is function to generate uuid as request id if client doesn't pass X-REQUEST-ID request header
func GetRequestID(requestIDParams *string) string {
	log.Debug("entering func getRequestID")
	//generate UUID as request ID if it doesn't exist in request header
	if requestIDParams == nil || *requestIDParams == "" {
		theUUID, err := uuid.NewUUID()
		newUUID := ""
		if err == nil {
			newUUID = theUUID.String()
		} else {
			newUUID = GenerateUUID()
		}
		requestIDParams = &newUUID
	}
	return *requestIDParams
}
