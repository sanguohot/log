package log

import (
	"github.com/sanguohot/log/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
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
	switch os.Getenv("LOG_LEVEL") {
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
	default:
		return lvl >= zapcore.InfoLevel
	}
	return false
}

func getEncodeConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	if os.Getenv("LOG_LEVEL") == util.LogLevelDebug {
		config = zap.NewDevelopmentEncoderConfig()
	}
	return config
}

func initConsoleLogConfig() logConfig {
	return initLogConfig(os.Stdout)
}

func initLogConfig(writer io.Writer) logConfig {
	// 实现判断日志等级的interface
	level := IsLogLevelEnable
	writerSync := zapcore.AddSync(writer)
	// 最后创建具体的Logger
	config := getEncodeConfig()
	return logConfig{
		encode: zapcore.NewConsoleEncoder(config),
		sync:   writerSync,
		level:  level,
	}
}

func initFileLogConfig() logConfig {
	logDir := filepath.Join(os.Getenv("LOG_ROOT"), util.LogDirPath)
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logFileEnv := os.Getenv("LOG_FILE")
	if logFileEnv == "" {
		logFileEnv = util.LogFilePath
	}
	logFilePath := filepath.Join(logDir, logFileEnv)
	hook := lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10, // MB
		MaxBackups: 20,
		MaxAge:     7,    //days
		Compress:   true, // disabled by default
		LocalTime:  true,
	}
	return initLogConfig(&hook)
}

// 不要小写的，否则会获取不到其它程序设置的环境变量
// 等其它程序调用初始化，每个进程只需要在入口初始化一次
func Init() {
	InitWithWriter(nil)
}

func InitWithWriter(writer io.Writer) {
	configList := make([]logConfig, 0)
	switch os.Getenv("LOG_TYPE") {
	case util.LogTypeFile:
		configList = append(configList, initFileLogConfig())
	case util.LogTypeGui:
		if writer == nil {
			panic("writer required for gui type")
		}
		configList = append(configList, initLogConfig(writer))
	case util.LogTypeAll:
		configList = append(configList, initConsoleLogConfig())
		configList = append(configList, initFileLogConfig())
		if writer != nil {
			configList = append(configList, initLogConfig(writer))
		}
	default:
		configList = append(configList, initConsoleLogConfig())
	}
	if len(configList) <= 0 {
		return
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
