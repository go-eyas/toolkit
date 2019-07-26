package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SugaredLogger *zap.SugaredLogger
var Logger *zap.Logger

// LogConfig 日志配置
type LogConfig struct {
	Level        string        // 日志级别
	Path         string        // 路径
	Name         string        // 文件名称
	Console      bool          // 是否输出到控制台
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

	level := new(zapcore.Level)
	err := level.Set(conf.Level)
	if err != nil {
		return err
	}

	lv := *level

	cores := []zapcore.Core{}

	if lv <= zapcore.DebugLevel {
		debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl <= zapcore.DebugLevel
		})

		debugWriter, err := getWriter(conf.Path+"/"+conf.Name+"_debug.", conf)
		if err != nil {
			return err
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel))
	}
	if lv <= zapcore.InfoLevel {
		infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.WarnLevel && lvl > zapcore.DebugLevel
		})

		infoWriter, err := getWriter(conf.Path+"/"+conf.Name+"_info", conf)
		if err != nil {
			return err
		}

		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel))
	}

	if lv <= zapcore.WarnLevel {
		warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.WarnLevel
		})
		warnWriter, err := getWriter(conf.Path+"/"+conf.Name+"_error", conf)
		if err != nil {
			return err
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel))
	} else {
		// 级别是 error 以上
		errorLevel := lv
		errorWriter, err := getWriter(conf.Path+"/"+conf.Name+"_error", conf)
		if err != nil {
			return err
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel))
	}

	if conf.Console {
		consoleWriter := os.Stdout
		//var consoleCore zapcore.Core
		//if conf.DebugConsole {
		//	consoleCore = zapcore.NewCore(encoder, zapcore.AddSync(consoleWriter), lv)
		//} else {
		//	consoleCore = zapcore.NewCore(encoder, zapcore.AddSync(consoleWriter), lv)
		//}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(consoleWriter), lv))
	}

	// 最后创建具体的Logger
	core := zapcore.NewTee(cores...)

	Logger = zap.New(core)
	SugaredLogger = Logger.Sugar()

	return nil
}

func getWriter(filename string, conf *LogConfig) (io.Writer, error) {
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H.log", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filepath.Base(filename)+".log"),
		rotatelogs.WithMaxAge(conf.MaxAge),
		rotatelogs.WithRotationTime(conf.RotationTime),
	)

	if err != nil {
		return nil, err
	}
	return hook, nil
}
