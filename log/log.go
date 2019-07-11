package log

import (
	"fmt"
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

// LogConfig 日志配置
type LogConfig struct {
	Path         string        // 路径
	Name         string        // 文件名称
	Console      bool          // 是否输出到控制台
	DebugConsole bool          // 是否把调试日志输出到控制台
	MaxAge       time.Duration // 保存多久的日志，默认15天
	RotationTime time.Duration // 多久分割一次日志
	Caller       bool          // 是否打印文件行号
}

var printCaller = false

// Init 初始化日志库
func Init(conf *LogConfig) error {
	// 默认保存最近15天日志
	if conf.MaxAge == 0 {
		conf.MaxAge = time.Hour * 24 * 15
	}
	if conf.RotationTime == 0 {
		conf.RotationTime = time.Hour
	}
	printCaller = conf.Caller
	return newLog(conf)
}

func newLog(conf *LogConfig) error {
	// 建立日志目录
	if err := os.MkdirAll(conf.Path+"/", os.ModePerm); err != nil {
		fmt.Println("init log path error.")
		return err
	}
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现两个判断日志等级的interface
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.DebugLevel
	})

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl > zapcore.DebugLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	// 获取 debug、info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	debugWriter, err := getWriter(conf.Path+"/"+conf.Name+"_debug.log", conf)
	if err != nil {
		return err
	}
	infoWriter, err := getWriter(conf.Path+"/"+conf.Name+"_info.log", conf)
	if err != nil {
		return err
	}
	warnWriter, err := getWriter(conf.Path+"/"+conf.Name+"_error.log", conf)
	if err != nil {
		return err
	}

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
	}

	if conf.Console {
		consoleWriter := os.Stdout
		var consoleCore zapcore.Core
		if conf.DebugConsole {
			consoleCore = zapcore.NewCore(encoder, zapcore.AddSync(consoleWriter), zapcore.DebugLevel)
		} else {
			consoleCore = zapcore.NewCore(encoder, zapcore.AddSync(consoleWriter), zapcore.InfoLevel)
		}
		cores = append(cores, consoleCore)
	}

	// 最后创建具体的Logger
	core := zapcore.NewTee(cores...)

	logger := zap.New(core)
	Logger = logger.Sugar()

	return nil
}

func getWriter(filename string, conf *LogConfig) (io.Writer, error) {
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(conf.MaxAge),
		rotatelogs.WithRotationTime(conf.RotationTime),
	)

	if err != nil {
		return nil, err
	}
	return hook, nil
}
