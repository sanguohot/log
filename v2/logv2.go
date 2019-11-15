package log

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sanguohot/log/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	Sugar  *zap.SugaredLogger
	Logger *zap.Logger
)

type logConfig struct {
	encode zapcore.Encoder
	sync   zapcore.WriteSyncer
	level  func(zapcore.Level) bool
}

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

func IsLogLevelEnable(lvl zapcore.Level) bool {
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
}

func getEncodeConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	if util.LogLevel == util.LogLevelDebug {
		config = zap.NewDevelopmentEncoderConfig()
	}
	return config
}

func initConsoleLogConfig() logConfig {
	fmt.Println("init console log")
	// 实现判断日志等级的interface
	level := IsLogLevelEnable
	consoleSync := zapcore.AddSync(os.Stdout)
	// 最后创建具体的Logger
	config := getEncodeConfig()
	return logConfig{
		encode: zapcore.NewConsoleEncoder(config),
		sync:   consoleSync,
		level:  level,
	}
}

func initFileLogConfig() logConfig {
	fmt.Println("init file log")
	err := os.MkdirAll(util.LogDirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logFilePath := filepath.Join(util.LogDirPath, util.LogFilePath)
	hook := lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10, // MB
		MaxBackups: 20,
		MaxAge:     7,    //days
		Compress:   true, // disabled by default
		LocalTime:  true,
	}
	// 实现判断日志等级的interface
	level := IsLogLevelEnable
	fileSync := zapcore.AddSync(&hook)
	// 最后创建具体的Logger
	config := getEncodeConfig()
	return logConfig{
		encode: zapcore.NewConsoleEncoder(config),
		sync:   fileSync,
		level:  level,
	}
}

func init() {
	fmt.Println(runtime.GOARCH)
	fmt.Println(runtime.GOOS)
	configList := make([]logConfig, 0)
	switch util.LogType {
	case util.LogTypeFile:
		configList = append(configList, initFileLogConfig())
	case util.LogTypeAll:
		configList = append(configList, initConsoleLogConfig())
		configList = append(configList, initFileLogConfig())
	default:
		configList = append(configList, initConsoleLogConfig())
		//if runtime.GOARCH[:3] != "arm" {
			//configList = append(configList, initFileLogConfig())
		//}
	}
	cores := make([]zapcore.Core, 0)
	for _, v := range configList {
		cores = append(cores, zapcore.NewCore(v.encode, v.sync, zap.LevelEnablerFunc(v.level)))
	}
	core := zapcore.NewTee(cores...)
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
