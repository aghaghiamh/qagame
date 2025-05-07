package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var once = sync.Once{}

func init() {
	once.Do(func() {

		file := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    10, // megabytes
			MaxBackups: 2,
			MaxAge:     7, // days
		})

		stdoutWriter := zapcore.AddSync(os.Stdout)

		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		fileEncoder := zapcore.NewJSONEncoder(productionCfg)

		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, file, zap.InfoLevel),
			zapcore.NewCore(fileEncoder, stdoutWriter, zap.InfoLevel),
		)

		Logger = zap.New(core)
	})

}