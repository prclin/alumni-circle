package core

import (
	"fmt"
	"github.com/prclin/alumni-circle/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
initZap 创建zap.SugaredLogger
*/
func initZap() {
	//获取日志输出级别
	levelStr := strings.ToUpper(global.Configuration.Zap.Level)
	var logLevel zapcore.Level
	switch levelStr {
	case "DEBUG":
		logLevel = zapcore.DebugLevel
	case "INFO":
		logLevel = zapcore.InfoLevel
	case "WARN":
		logLevel = zapcore.WarnLevel
	case "ERROR":
		logLevel = zapcore.ErrorLevel
	case "DPANIC":
		logLevel = zapcore.DPanicLevel
	case "PANIC":
		logLevel = zapcore.PanicLevel
	case "FATAL":
		logLevel = zapcore.FatalLevel
	default:
		panic("zap log level invalid,check your configuration file...")
	}
	//构建SugaredLogger核心配置，包括：编码器，输出位置，日志输出级别
	zapCore := zapcore.NewCore(constructEncoder(), zapcore.NewMultiWriteSyncer(constructFileWriteSyncer(), zapcore.AddSync(os.Stdout)), logLevel)
	//构建SugaredLogger并返回
	global.Logger = zap.New(zapCore, zap.AddCaller()).Sugar()
}

/*
constructEncoder 构建zap的日志格式化编码器
*/
func constructEncoder() zapcore.Encoder {
	//初始化encoder配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, e zapcore.PrimitiveArrayEncoder) {
		e.AppendString(t.Local().Format(time.DateTime))
	}
	//构建encoder并返回
	return zapcore.NewJSONEncoder(encoderConfig)
}

/*
constructFileWriteSyncer 构建日志输出同步器，即指定一个日志输出的位置
*/
func constructFileWriteSyncer() zapcore.WriteSyncer {
	//获取系统路径分割符
	pathSeparator := string(filepath.Separator)
	//获取工作目录
	wd, _ := os.Getwd()

	logfilePath := wd + pathSeparator + "log" + pathSeparator + time.Now().Format(time.DateOnly) + ".txt"
	fmt.Println("log files are stored in " + wd + pathSeparator + "log" + pathSeparator)

	//构建lumberjack.Logger，使他作为zap的日志输出位置，再由lumberjack输出到文件以达到拆分日志文件的效果
	lumberjackSyncer := &lumberjack.Logger{
		Filename:   logfilePath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	}

	//指定输出日志到为lumberjack，构建WriteSyncer并返回
	return zapcore.AddSync(lumberjackSyncer)
}
