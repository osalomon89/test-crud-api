package log

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/mercadolibre/go-meli-toolkit/goutils/logger"
)

const (
	minTags = 3
)

type ILogger interface {
	Info(source interface{}, tags map[string]string, message string, args ...interface{})
	Warn(source interface{}, tags map[string]string, message string, args ...interface{})
	Error(source interface{}, tags map[string]string, err error, message string, args ...interface{})
	Debug(source interface{}, tags map[string]string, message string, args ...interface{})
	GetRequestID() string
	GetMessage(message string, args ...interface{}) string
	GetTags(source interface{}, tags map[string]string) []string
}

type log struct {
	mutex      sync.Mutex
	requestID  string
	sequenceID int
}

func DefaultLogger() ILogger {
	iLogger := &log{requestID: newRequestID()}
	return iLogger
}

func NewLogger(requestID string) ILogger {
	iLogger := &log{requestID: requestID}

	return iLogger
}

func (theLogger *log) Info(source interface{}, tags map[string]string, message string, args ...interface{}) {
	logger.Info(theLogger.GetMessage(message, args...), theLogger.GetTags(source, tags)...)
}

func (theLogger *log) Warn(source interface{}, tags map[string]string, message string, args ...interface{}) {
	logger.Warn(theLogger.GetMessage(message, args...), theLogger.GetTags(source, tags)...)
}

func (theLogger *log) Error(source interface{}, tags map[string]string, err error,
	message string, args ...interface{}) {
	logger.Error(theLogger.GetMessage(message, args...), err, theLogger.GetTags(source, tags)...)
}

func (theLogger *log) Debug(source interface{}, tags map[string]string, message string, args ...interface{}) {
	logger.Debug(theLogger.GetMessage(message, args...), theLogger.GetTags(source, tags)...)
}

func (theLogger *log) GetRequestID() string {
	return theLogger.requestID
}

func (theLogger *log) GetMessage(message string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(message, args...)
	}

	return message
}

func newRequestID() string {
	id := ""
	logID, err := uuid.NewV4()

	if err == nil {
		id = logID.String()
	}

	return id
}

func getClass(source interface{}) string {
	t := reflect.TypeOf(source)
	if t != nil {
		return t.String()
	}

	return ""
}

func (theLogger *log) GetTags(source interface{}, tags map[string]string) []string {
	var res []string

	i := 0

	if len(tags) == 0 {
		res = make([]string, minTags)
	} else {
		res = make([]string, len(tags)+minTags)
		for key, value := range tags {
			res[i] = fmt.Sprintf("%s:%v", key, value)
			i++
		}
	}

	theLogger.mutex.Lock()
	theLogger.sequenceID++

	res[i] = fmt.Sprintf("Request_ID:%v", theLogger.requestID)
	res[i+1] = fmt.Sprintf("Class:%v", getClass(source))
	res[i+2] = fmt.Sprintf("Sequence_ID:%v", theLogger.sequenceID)

	theLogger.mutex.Unlock()

	return res
}
