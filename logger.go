package pam4sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/Sirupsen/logrus"
)

// ILogger is for logging
type ILogger interface {
	Info(message string)
	InfoFL(message string, fn string, line int)
	Warn(message string)
	WarnFL(message string, fn string, line int)
	Debug(message string)
	DebugFL(message string, fn string, line int)
	Error(message string)
	ErrorFL(message string, fn string, line int)
	Print(obj interface{})
}

// Logger is the logger utility with information of request context
type Logger struct {
	log        *logrus.Logger
	simple     bool
	Type       string
	RequestID  string
	TrackingID string
	SourceIP   string
	AppID      string
	HTTPMethod string
	EndPoint   string
}

func isDebugEnv() bool {
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		return true
	}
	return false
}

// NewLoggerSimple return plain text simple logger
func NewLoggerSimple() *Logger {
	log := logrus.New()
	if isDebugEnv() {
		log.SetLevel(logrus.DebugLevel)
	}
	formatter := new(logrus.TextFormatter)
	formatter.DisableTimestamp = true
	formatter.DisableColors = true
	formatter.DisableSorting = true

	log.Formatter = formatter
	multi := io.MultiWriter(os.Stderr)
	log.Out = multi

	return &Logger{
		log:    log,
		simple: true,
	}
}

// NewLogger will create the logger which log context information
func NewLogger(requestType string, requestID string, sourceIP string, httpMethod string, endpoint string, trackingID string, appID string) *Logger {
	log := logrus.New()
	if isDebugEnv() {
		log.SetLevel(logrus.DebugLevel)
	}
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.DisableTimestamp = false
	formatter.DisableColors = false
	formatter.DisableSorting = false
	log.Formatter = formatter
	multi := io.MultiWriter(os.Stderr)
	log.Out = multi

	return &Logger{
		log:        log,
		simple:     false,
		Type:       requestType,
		RequestID:  requestID,
		SourceIP:   sourceIP,
		HTTPMethod: httpMethod,
		EndPoint:   endpoint,
		TrackingID: trackingID,
		AppID:      appID,
	}
}

func (logger *Logger) getLogFields(fn string, line int) logrus.Fields {
	if len(fn) == 0 {
		return logrus.Fields{
			"type":        logger.Type,
			"request_id":  logger.RequestID,
			"tracking_id": logger.TrackingID,
			"source_ip":   logger.SourceIP,
			"app_id":      logger.AppID,
			"http_method": logger.HTTPMethod,
			"endpoint":    logger.EndPoint,
		}
	}
	return logrus.Fields{
		"type":        logger.Type,
		"request_id":  logger.RequestID,
		"tracking_id": logger.TrackingID,
		"source_ip":   logger.SourceIP,
		"app_id":      logger.AppID,
		"http_method": logger.HTTPMethod,
		"endpoint":    logger.EndPoint,
		"function":    fn,
		"line":        line,
	}
}

// Info log information level
func (logger *Logger) Info(message string) {
	logger.InfoFL(message, "", -1)
}

// InfoFL log information level
func (logger *Logger) InfoFL(message string, fn string, line int) {
	if logger.simple {
		logger.log.Info(message)
	} else {
		logger.log.WithFields(logger.getLogFields(fn, line)).Info(message)
	}
}

// Warn log warning level
func (logger *Logger) Warn(message string) {
	logger.WarnFL(message, "", -1)
}

// WarnFL log warnning level
func (logger *Logger) WarnFL(message string, fn string, line int) {
	if logger.simple {
		logger.log.Warn(message)
	} else {
		logger.log.WithFields(logger.getLogFields(fn, line)).Warn(message)
	}
}

// Debug log debug level
func (logger *Logger) Debug(message string) {
	logger.DebugFL(message, "", -1)
}

// DebugFL log debug level
func (logger *Logger) DebugFL(message string, fn string, line int) {
	if logger.simple {
		logger.log.Debug(message)
	} else {
		logger.log.WithFields(logger.getLogFields(fn, line)).Debug(message)
	}
}

// Error log error level
func (logger *Logger) Error(message string) {
	logger.ErrorFL(message, "", -1)
}

// ErrorFL log error level
func (logger *Logger) ErrorFL(message string, fn string, line int) {
	if logger.simple {
		logger.log.Error(message)
	} else {
		logger.log.WithFields(logger.getLogFields(fn, line)).Error(message)
	}
}

// Print log an object
func (logger *Logger) Print(obj interface{}) {
	str := ""
	if obj == nil {
		str = ""
	} else {
		switch reflect.TypeOf(obj).Kind() {
		case reflect.String:
			str = obj.(string)
		case reflect.Struct, reflect.Slice, reflect.Array:
			b, _ := json.MarshalIndent(obj, "", "  ")
			str = string(b)
		default:
			str = fmt.Sprintf("%v", obj)
		}
	}

	logger.log.Print(str)
}
