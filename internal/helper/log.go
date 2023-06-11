package helper

import (
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	instances = sync.Map{}
	logPath   string
	level     zap.AtomicLevel
)

func SetupLogger(cfg *Config) (err error) {
	// check log path
	if info, err := os.Stat(cfg.LogPath); err != nil || os.IsNotExist(err) || !info.IsDir() {
		return errors.Wrap(err, "failed to check log path")
	}
	logPath = cfg.LogPath

	// convert log level to zap format
	if level, err = zap.ParseAtomicLevel(cfg.LogLevel); err != nil {
		return errors.Wrap(err, "failed to parse log level")
	}

	return nil
}

func getLoggerCore(key string) (*zap.SugaredLogger, error) {
	f := path.Join(logPath, key+".log")
	encConfig := zap.NewProductionEncoderConfig()
	coreFile := zapcore.NewCore(
		zapcore.NewJSONEncoder(encConfig),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   f,
			MaxSize:    10, // maximum size of a single log file (Mbytes)
			MaxBackups: 10, // maximum number of logs to be saved
			MaxAge:     30, // maximum number of days to keep logs
			Compress:   false,
		}),
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return level.Enabled(l)
		}),
	)
	coreConsole := zapcore.NewCore(
		zapcore.NewJSONEncoder(encConfig),
		zapcore.AddSync(os.Stdout),
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return level.Enabled(l)
		}),
	)
	core := zapcore.NewTee(coreFile, coreConsole)

	// build logger
	zapLog := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	defer zapLog.Sync()
	return zapLog.Sugar(), nil
}

func GetLogger(key string) *zap.SugaredLogger {
	val, _ := instances.LoadOrStore(key, func() *zap.SugaredLogger {
		ins, err := getLoggerCore(key)
		if err != nil {
			panic(err)
		}
		ins.Infof("logger (%v) initialized", key)
		return ins.Named(key)
	})

	return val.(func() *zap.SugaredLogger)()
}
