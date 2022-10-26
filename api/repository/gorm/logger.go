package gorm

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	LogLevel                 logger.LogLevel
	infoStr, warnStr, errStr string
}

func (GormLogger) Printf(v ...interface{}) {
	switch v[0] {
	case "sql":
		logrus.WithFields(
			logrus.Fields{
				"module":        "gorm",
				"type":          "sql",
				"rows_returned": v[5],
				"src":           v[1],
				"values":        v[4],
				"duration":      v[2],
			},
		).Info(v[3])
		//case "log":
	default:
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}

func (l GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logrus.Info(msg)
	if l.LogLevel >= logger.Info {
		l.Printf(data...)
	}
}

func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logrus.Warn(msg)
	if l.LogLevel >= logger.Warn {
		l.Printf(data...)
	}
}

func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logrus.Error(msg)
	if l.LogLevel >= logger.Error {
		l.Printf(data...)
	}
}

func (GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
}
