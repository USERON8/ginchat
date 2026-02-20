package global

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化Zap日志
func InitLogger() error {
	// 1.  lumberjack日志切割
	hook := lumberjack.Logger{
		Filename:   Config.Log.Filename,
		MaxSize:    Config.Log.MaxSize,
		MaxBackups: Config.Log.MaxBackups,
		MaxAge:     Config.Log.MaxAge,
		Compress:   Config.Log.Compress,
	}

	// 2. 日志级别
	var level zapcore.Level
	switch Config.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 3. 编码器配置（生产级格式）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 4. 同时输出到文件和控制台
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(&hook), level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level),
	)

	// 5. 初始化全局Log
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(Log)

	return nil
}
