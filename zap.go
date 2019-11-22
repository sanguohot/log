package log

import (
	"github.com/sanguohot/log/util"
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
	atom zap.AtomicLevel
)

func GetLogPath() string {
	return filepath.Join(util.LogRoot, util.LogDirPath, util.LogFilePath)
}

func ServerTypeIsProd() bool {
	if util.ServerType == util.ServerTypeProd {
		return true
	}
	return false
}

func init() {
	var (
		config     zapcore.EncoderConfig
		stackLevel zapcore.Level
	)
	atom = zap.NewAtomicLevel()
	fileSync := zapcore.AddSync(&lumberjack.Logger{
		Filename:   GetLogPath(),
		MaxSize:    50, // MB
		MaxBackups: 20,
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
		atom.SetLevel(zap.InfoLevel)
	} else {
		config = zap.NewDevelopmentEncoderConfig()
		stackLevel = zap.WarnLevel
		atom.SetLevel(zap.DebugLevel)
	}
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	//core := zapcore.NewCore(
	//	zapcore.NewJSONEncoder(config),
	//	zapcore.NewMultiWriteSyncer(consoleSync, fileSync),
	//	atom, //debug,info,warn,error
	//)
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), consoleSync, atom), // 日志同步到控制台
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), fileSync, atom),    // 日志同步到app.log
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
