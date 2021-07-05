package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
	}
}

type CronLogger struct {
}

func (cl CronLogger) Info(msg string, keysAndValues ...interface{}) {

	// is := []interface{}{msg}
	// is = append(is, keysAndValues...)

	// log.ErrS(is...)
}

func (cl CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {

	is := []zap.Field{zap.Error(err)}

	// for k, v := range keysAndValues {
	// 	is = append(is, zap.Any(k.(string), v))
	// }

	Logger.Error(msg, is...)
}
