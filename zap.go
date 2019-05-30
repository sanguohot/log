package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

var (
	Sugar  *zap.SugaredLogger
	Logger *zap.Logger
	// Atom.SetLevel(zap.DebugLevel) 程序运行时动态级别
	Atom           zap.AtomicLevel
	ServerType     = os.Getenv("SERVER_TYPE")
	serverTypeProd = "production"
)

func GetLogPath() string {
	return filepath.Join("./", "app.log")
}

func ServerTypeIsProd() bool {
	if ServerType == serverTypeProd {
		return true
	}
	return false
}

func init() {
	var (
		config     zapcore.EncoderConfig
		stackLevel zapcore.Level
	)
	Atom = zap.NewAtomicLevel()
	fileSync := zapcore.AddSync(&lumberjack.Logger{
		Filename:   GetLogPath(),
		MaxSize:    500, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
		LocalTime:  true,
		Compress:   true,
	})
	consoleSync := zapcore.AddSync(os.Stdout)
	// 默认开发者Encoder，包含函数调用信息
	// 可以根据环境变量调整
	// 根据当前环境和日志级别（warn以上）自动打印调用栈信息
	if ServerTypeIsProd() {
		config = zap.NewProductionEncoderConfig()
		stackLevel = zap.ErrorLevel
		Atom.SetLevel(zap.InfoLevel)
	} else {
		config = zap.NewDevelopmentEncoderConfig()
		stackLevel = zap.WarnLevel
		Atom.SetLevel(zap.DebugLevel)
	}
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.NewMultiWriteSyncer(consoleSync, fileSync),
		Atom, //debug,info,warn,error
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(stackLevel))
	defer Logger.Sync() // flushes buffer, if any
	Sugar = Logger.Sugar()
}

func Info(args ...interface{}) {
	Sugar.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Sugar.Infof(template, args...)
}

func Warn(args ...interface{}) {
	Sugar.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	Sugar.Warnf(template, args...)
}

func Debug(args ...interface{}) {
	Sugar.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	Sugar.Debugf(template, args...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Errorf(template string, args ...interface{}) {
	Sugar.Errorf(template, args...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

func Fatalf(template string, args ...interface{}) {
	Sugar.Fatalf(template, args...)
}
