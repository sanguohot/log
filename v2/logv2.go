package log

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sanguohot/log/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	Sugar  *zap.SugaredLogger
	Logger *zap.Logger
)

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func pathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func init() {
	err := os.MkdirAll(util.LogDirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logFilePath := filepath.Join(util.LogDirPath, util.LogFilePath)
	hook := lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    50, // MB
		MaxBackups: 3,
		MaxAge:     3,    //days
		Compress:   true, // disabled by default
		LocalTime:  true,
	}

	// 实现判断日志等级的interface
	level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		switch util.LogLevel {
		case util.LogLevelOff:
			return false
		case util.LogLevelDebug:
			return lvl >= zapcore.DebugLevel
		case util.LogLevelInfo:
			return lvl >= zapcore.InfoLevel
		case util.LogLevelWarn:
			return lvl >= zapcore.WarnLevel
		case util.LogLevelError:
			return lvl >= zapcore.ErrorLevel
		case util.LogLevelFatal:
			return lvl >= zapcore.DPanicLevel
		}
		if util.LogLevel == "" {
			return lvl >= zapcore.InfoLevel
		}
		return false
	})
	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	writer := getWriter(filepath.Join(util.LogDirPath, util.LinkFilePath))
	linkSync := zapcore.AddSync(writer)
	consoleSync := zapcore.AddSync(os.Stdout)
	fileSync := zapcore.AddSync(&hook)
	// 最后创建具体的Logger
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	if util.LogLevel == util.LogLevelDebug {
		config = zap.NewDevelopmentEncoderConfig()
	}
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), consoleSync, level), // 日志同步到控制台
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), fileSync, level),    // 日志同步到app.log
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), linkSync, level),    // 日志同步到link.log(symbol link)，用于切割文件
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	defer Logger.Sync() // flushes buffer, if any
	Sugar = Logger.Sugar()
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每2小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*2),
	)

	if err != nil {
		panic(err)
	}
	return hook
}
