package httpcontext

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const ReqIdKey = "reqId"
const ReqIdHeaderName = "X-Req-Id"
const LoggerKey = "logger"

func BackgroundFromContext(c *gin.Context) context.Context {
	var reqId string

	if reqId = c.GetHeader(ReqIdHeaderName); reqId == "" {
		reqId = uuid.New().String()
	}

	return useLogger(context.WithValue(c, ReqIdKey, reqId))
}

func useLogger(ctx context.Context) context.Context {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if reqId := ctx.Value(ReqIdKey).(string); reqId == "" {
		// panics
		logrus.Fatal("No request id associated with request")
	} else {
		log := logrus.WithField(ReqIdKey, reqId)
		ctx = context.WithValue(ctx, LoggerKey, log)
	}

	return ctx
}

func GetLogger(ctx context.Context) *logrus.Entry {
	log := ctx.Value(LoggerKey)

	if log == nil {
		logrus.Fatal("Logger is missing in the context") // panics
	}

	return log.(*logrus.Entry)
}
