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

func (e Err) Log(logger *zap.Logger, lvl zapcore.Level) {
	fields := []zap.Field{}
	msg := ""

	msgStrs := []string{}
	if e.Name != "" {
		msgStrs = append(msgStrs, e.Name)
	}
	if e.Description != "" {
		msgStrs = append(msgStrs, e.Description)
	}
	msg = strings.Join(msgStrs, " - ")

	if e.InternalErr != nil {
		fields = append(fields, zap.Error(e.InternalErr))
	}

	if e.Payload != nil {
		e.Payload = map[string]any{}
	}
	for k, v := range e.Payload {
		fields = append(fields, zap.Any(k, v))
	}

	if e.TypeErr != nil {
		fields = append(fields, zap.String("err_type", e.TypeErr.Error()))
	}

	if !shouldSplit && e.Stack != nil {
		fmt.Println(string(e.Stack))
	}

	logger.Log(lvl, msg, fields...)

	if shouldSplit && e.Stack != nil {
		fmt.Println(string(e.Stack))
	}
}
