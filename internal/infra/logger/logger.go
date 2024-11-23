package logger

import (
	"fmt"
	"os"
	"sync"
	"teste/internal/infra/logger/attributes"
	"teste/internal/infra/variables"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	LogLevel       string
}

var (
	once   sync.Once
	option *Option
)

func Init(opt *Option) {
	once.Do(func() {
		option = opt
		initLogsDir()
		zapLogger, err := newZap()
		if err != nil {
			panic(err)
		}
		zap.ReplaceGlobals(zapLogger)
	})
}

func initLogsDir() {
	if _, err := os.Stat(variables.DirLog()); os.IsNotExist(err) {
		err := os.Mkdir(variables.DirLog(), 0755)
		if err != nil {
			panic(fmt.Sprintf("Falha ao criar diretório de logs: %v", err))
		}
	} else {
		if err := os.Chmod(variables.DirLog(), 0755); err != nil {
			panic(fmt.Sprintf("Falha ao definir permissões de gravação: %v", err))
		}
	}
}

func Sync() {
	_ = zap.L().Sync()
}

func newZap() (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(option.LogLevel)); err != nil {
		return nil, err
	}

	debugLogFile := newLogWriter(getLogFileName())
	errorLogFile := newLogWriter(getErrorLogFileName())

	core := zapcore.NewTee(
		zapcore.NewCore(newZapEncoder(), debugLogFile, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})),
		zapcore.NewCore(newZapEncoder(), errorLogFile, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})),
	)

	return zap.New(core), nil
}

func getLogFileName() string {
	return fmt.Sprintf(variables.DirLog()+"%s-debug.log", time.Now().Format("02-01-2006"))
}

func getErrorLogFileName() string {
	return fmt.Sprintf(variables.DirLog()+"%s-error.log", time.Now().Format("02-01-2006"))
}

func newZapEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(newZapEncoderConfig())
}

func newZapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "severity",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("02/01/2006 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func logWithContext(logFunc func(string, ...zap.Field), c *fiber.Ctx, message string, attr attributes.Attributes) {
	if c != nil {
		logFunc(message, makeFiberAttributes(c, attr)...)
	} else {
		logFunc(message, zap.String("error", "fiber context is nil"))
	}
}

func makeFiberAttributes(c *fiber.Ctx, attr attributes.Attributes) []zap.Field {
	if attr == nil {
		attr = attributes.New()
	}
	if cid := c.Locals("cid"); cid == nil {
		attr["cid"] = uuid.New().String()
		c.Locals("cid", attr["cid"])
	}
	return []zap.Field{
		zap.Any("attributes", attr),
		zap.Any("http_request", map[string]interface{}{
			"method": c.Method(),
			"path":   c.Path(),
			"status": c.Response().StatusCode(),
		}),
		zap.Object("resource", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			enc.AddString("service.name", option.ServiceName)
			enc.AddString("service.version", option.ServiceVersion)
			enc.AddString("service.environment", option.Environment)
			return nil
		})),
	}
}

func Debug(c *fiber.Ctx, message string, attr attributes.Attributes) {
	logWithContext(zap.L().Debug, c, message, attr)
}

func Info(c *fiber.Ctx, message string, attr attributes.Attributes) {
	logWithContext(zap.L().Info, c, message, attr)
}

func Warn(c *fiber.Ctx, message string, attr attributes.Attributes) {
	logWithContext(zap.L().Warn, c, message, attr)
}

func Error(c *fiber.Ctx, message string, attr attributes.Attributes) {
	logWithContext(zap.L().Error, c, message, attr)
}

func Fatal(c *fiber.Ctx, message string, attr attributes.Attributes) {
	logWithContext(zap.L().Fatal, c, message, attr)
}

func newLogWriter(filename string) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})
}
