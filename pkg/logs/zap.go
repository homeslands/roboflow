package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(isProd bool) *zap.Logger {
	var conf zap.Config
	conf = NewProductionConfig()

	if !isProd {
		conf = zap.NewDevelopmentConfig()
		conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	conf.DisableStacktrace = true

	log, err := conf.Build()
	if err != nil {
		panic(err)
	}

	return log
}

func NewProductionConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
