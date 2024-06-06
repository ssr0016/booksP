package log

import (
	"time"

	"github.com/apex/log"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New3(msg string) *log.Entry {
	logID := xid.New().String()
	return log.WithField(msg, logID)
}

func New(name string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.EncoderConfig.LevelKey = "log_level"
	config.EncoderConfig.StacktraceKey = zapcore.OmitKey

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	log = log.Named(name)
	return log, nil
}
