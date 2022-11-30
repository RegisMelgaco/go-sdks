package erring

import (
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	splitOnce   sync.Once
	shouldSplit bool
)

func SplitStackFromLogs() {
	splitOnce.Do(func() {
		shouldSplit = true
	})
}

func (err Err) Log(logger *zap.Logger, lvl zapcore.Level) Err {
	fields := []zap.Field{}
	msg := ""

	msgStrs := []string{}
	if err.Name != "" {
		msgStrs = append(msgStrs, err.Name)
	}
	if err.Description != "" {
		msgStrs = append(msgStrs, err.Description)
	}
	msg = strings.Join(msgStrs, " - ")

	if err.InternalErr != nil {
		fields = append(fields, zap.Error(err.InternalErr))
	}

	if err.Payload != nil {
		err.Payload = map[string]any{}
	}
	for k, v := range err.Payload {
		fields = append(fields, zap.Any(k, v))
	}

	if err.TypeErr != nil {
		fields = append(fields, zap.String("err_type", err.TypeErr.Error()))
	}

	if !shouldSplit && err.Stack != nil {
		fmt.Println(string(err.Stack))
	}

	logger.Log(lvl, msg, fields...)

	if shouldSplit && err.Stack != nil {
		fmt.Println(string(err.Stack))
	}

	return err
}
