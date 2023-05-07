package logging

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

var (
	Info   = logger.Infow
	Infof  = logger.Infof
	Warn   = logger.Warnw
	Warnf  = logger.Warnf
	Error  = logger.Errorw
	Errorf = logger.Errorf
)

type Config struct {
	OutputPath string
}

func Init(lcfg Config) error {
	cfg := zap.NewProductionConfig()
	if lcfg.OutputPath != `` {
		cfg.OutputPaths = []string{
			lcfg.OutputPath,
		}
	}
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.DisableStacktrace = true

	l, err := cfg.Build()
	if err != nil {
		return err
	}
	logger = l.Sugar()

	Info = logger.Infow
	Infof = logger.Infof
	Warn = logger.Warnw
	Warnf = logger.Warnf
	Error = logger.Errorw
	Errorf = logger.Errorf

	return nil
}
