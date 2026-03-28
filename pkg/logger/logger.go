// pkg/logger/logger.go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"os"
	"time"
)

var Log *zap.Logger

func Init() {
	// 按天切割日志文件
	fileWriter := &lumberjack.Logger{
		Filename:   "./logs/ginchat.log", // 日志文件路径
		MaxSize:    100,                  // 每个文件最大 100MB
		MaxBackups: 30,                   // 最多保留 30 个备份
		MaxAge:     30,                   // 最多保留 30 天
		Compress:   true,                 // 压缩旧日志
	}

	// 编码器配置
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // info/error/warn
		EncodeTime:     customTimeEncoder,             // 自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径，如 handler/user.go:32
	}

	// 同时写文件和控制台
	core := zapcore.NewTee(
		// 文件：Info 级别以上
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(fileWriter),
			zapcore.InfoLevel,
		),
		// 控制台：Debug 级别以上（开发用）
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		),
	)

	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// 自定义时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// 对外暴露方法，不用每次都 Log.xxx
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}
