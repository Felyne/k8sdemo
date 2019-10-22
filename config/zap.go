// package zap handles creating zap logger
package config

import (
	"encoding/json"
	"github.com/jfeng45/k8sdemo/tool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type loggerWrapper struct {
	lw *zap.SugaredLogger
}


func (logger *loggerWrapper) Errorf(format string, args ...interface{}) {
	logger.lw.Errorf(format, args)
}
func (logger *loggerWrapper) Fatalf(format string, args ...interface{}) {
	logger.lw.Fatalf(format, args)
}
func (logger *loggerWrapper) Fatal(args ...interface{}) {
	logger.lw.Fatal(args)
}
func (logger *loggerWrapper) Infof(format string, args ...interface{}) {
	logger.lw.Infof(format, args)
}
func (logger *loggerWrapper) Warnf(format string, args ...interface{}) {
	logger.lw.Warnf(format, args)
}
func (logger *loggerWrapper) Debugf(format string, args ...interface{}) {
	logger.lw.Debugf(format, args)
}
func (logger *loggerWrapper) Printf(format string, args ...interface{}) {
	logger.lw.Infof(format, args)
}
func (logger *loggerWrapper) Println(args ...interface{}) {
	logger.lw.Info(args, "\n")
}

func RegisterZapLog() error {
	zLogger, err := initLog()
	if err != nil {
		return errors.Wrap(err, "RegisterZapLog")
	}
	defer zLogger.Sync()
	zSugarlog := zLogger.Sugar()
	zSugarlog.Info()

	//This is for loggerWrapper implementation
	//appLogger.SetLogger(&loggerWrapper{zaplog})

	tool.SetLogger(zSugarlog)
	return nil
}

// initLog create logger
func initLog() (zap.Logger, error) {
	rawJSON := []byte(`{
	 "level": "info",
     "Development": true,
      "DisableCaller": false,
	 "encoding": "console",
	 "outputPaths": ["stdout", "../logs/demo.log"],
	 "errorOutputPaths": ["stderr"],
	 "encoderConfig": {
		"timeKey":        "ts",
		"levelKey":       "level",
		"messageKey":     "msg",
         "nameKey":        "name",
		"stacktraceKey":  "stacktrace",
         "callerKey":      "caller",
		"lineEnding":     "\n\t",
        "timeEncoder":     "time",
		"levelEncoder":    "lowercaseLevel",
        "durationEncoder": "stringDuration",
         "callerEncoder":   "shortCaller"
	 }
	}`)

	var cfg zap.Config
	var zLogger *zap.Logger
	//standard configuration
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return *zLogger, errors.Wrap(err, "Unmarshal")
	}
	//customize it from configuration file
	err := customizeLogFromConfig(&cfg)
	if err != nil {
		return *zLogger, errors.Wrap(err, "cfg.Build()")
	}
	zLogger, err = cfg.Build()
	if err != nil {
		return *zLogger, errors.Wrap(err, "cfg.Build()")
	}

	zLogger.Debug("logger construction succeeded")
	return *zLogger, nil
}

// customizeLogFromConfig customize log based on parameters from configuration file
func customizeLogFromConfig(cfg *zap.Config) error {
	//cfg.DisableCaller = !lc.EnableCaller

	// set log level
	l := zap.NewAtomicLevel().Level()
	err := l.Set("debug")
	if err != nil {
		return errors.Wrap(err, "")
	}
	cfg.Level.SetLevel(l)

	return nil
}
