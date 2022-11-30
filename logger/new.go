package logger

import (
	"github.com/regismelgaco/go-sdks/erring"
	"go.uber.org/zap"
)

func New(isDev bool) (*zap.Logger, error) {
	var (
		l   *zap.Logger
		err error
	)
	if isDev {
		lCfg := zap.NewDevelopmentConfig()
		lCfg.DisableStacktrace = true
		lCfg.DisableCaller = true

		l, err = lCfg.Build()

		erring.SplitStackFromLogs()
	} else {
		l, err = zap.NewProduction()
	}

	return l, err
}
