package oceanlog

import (
	"context"
	"fmt"
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"path"
)

const (
	traceIDKey    = "trace_id"
	spanIDKey     = "span_id"
	traceFlagsKey = "trace_flags"
	logEventKey   = "log"
	logIDKey      = "log_id"
	logJson       = "json"
	logConsole    = "console"
)

//var Oceanlog hlog.FullLogger

func InitOceanLog(LogFileName, logFormat string, level Level) *DefaultLogger {
	// Provides compression and deletion
	lumberjackLogger := getLumberjackLogger(LogFileName)
	iw := io.MultiWriter(lumberjackLogger, os.Stdout) // os.Stdout, logger.Gin.Writer()

	if logFormat != logJson {
		iw = NewConsole(iw)
	}
	// For logrus detailed settings, please refer to https://github.com/hertz-contrib/logger/tree/main/logrus and https://github.com/sirupsen/logrus
	ologger := New(
		WithOutput(iw),   // allows to specify output
		WithLevel(level), // option with log level
		WithFormattedTimestamp("2006-01-02 15:04:05"), // option with timestamp
		//WithTimestamp(),                               // option with timestamp
		//WithFields(map[string]interface{}{})
		//WithCallerSkipFrameCount(6), // 自动记录日志调用位置
		//WithCaller(),                                  // 自动记录日志调用位置
		// ...
	)
	ologger.SetOutput(iw)
	ologger.SetLevel(level)

	//hlog.SetLogger(ologger)
	//Oceanlog = ologger
	return ologger
}

func getLumberjackLogger(fileName string) *lumberjack.Logger {
	if err := InitOutToFile(fileName); err != nil {
		panic(err)
	}
	// Provides compression and deletion
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    20,   // A file can be up to 20M.
		MaxBackups: 5,    // Save up to 5 files at the same time.
		MaxAge:     10,   // A file can exist for a maximum of 10 days.
		Compress:   true, // Compression with gzip.
	}
	return lumberjackLogger
}

func (c *LogConf) GetHzLog(ctx context.Context) *hertzlogrus.Logger {
	lumberjackLogger := GetLumberjackLogger(c)
	var writers []io.Writer
	if c.Fileout {
		writers = append(writers, lumberjackLogger)
	}
	if c.Stdout {
		writers = append(writers, os.Stdout)
	}
	iw := io.MultiWriter(writers...) // os.Stdout, logger.Gin.Writer()
	lo := hertzlogrus.NewLogger(hertzlogrus.WithLogger(logrus.New()))
	temp := lo.Logger()
	// 设置日志格式为json格式
	//temp := hlog.DefaultLogger{}
	if c.Formatter == "json" {
		temp.SetFormatter(&logrus.JSONFormatter{})
	} else {
		temp.SetFormatter(&logrus.TextFormatter{})
	}
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	temp.SetOutput(iw)

	// 设置日志级别为warn以上
	if lev, err := logrus.ParseLevel(c.Level); err == nil {
		temp.SetLevel(lev)
	} else {
		temp.SetLevel(logrus.InfoLevel)
	}

	return lo
}

func (c *LogConf) GetLogrusLog() *logrus.Logger {
	lumberjackLogger := GetLumberjackLogger(c)
	var writers []io.Writer
	if c.Fileout {
		writers = append(writers, lumberjackLogger)
	}
	if c.Stdout {
		writers = append(writers, os.Stdout)
	}
	iw := io.MultiWriter(writers...) // os.Stdout, logger.Gin.Writer()
	// 设置日志格式为json格式
	temp := logrus.Logger{}
	if c.Formatter == "json" {
		temp.SetFormatter(&logrus.JSONFormatter{})
	} else {
		temp.SetFormatter(&logrus.TextFormatter{})
	}
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	temp.SetOutput(iw)

	// 设置日志级别为warn以上
	if lev, err := logrus.ParseLevel(c.Level); err == nil {
		temp.SetLevel(lev)
	} else {
		temp.SetLevel(logrus.InfoLevel)
	}
	return &temp
}

func defaultLumberjackLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		MaxSize:    20,   // A file can be up to 20M.
		MaxBackups: 5,    // Save up to 5 files at the same time.
		MaxAge:     10,   // A file can exist for a maximum of 10 days.
		Compress:   true, // Compression with gzip.
	}
}

func GetLumberjackLogger(c *LogConf) *lumberjack.Logger {
	if err := InitOutToFile(c.LogFileName); err != nil {
		panic(err)
	}
	return c.Lumberjack
}

func GetFileIO(LogFileName string) *os.File {
	_ = InitOutToFile(LogFileName)
	f, err := os.OpenFile(LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	//defer f.Close()
	return f
}

func InitOutToFile(logFilePath string) error {
	//logFilePath := path.Join(logPath, logFileName)
	logDirPath := path.Dir(logFilePath)
	if _, err := os.Stat(logDirPath); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(logDirPath, os.ModePerm)
			if err != nil {
				msg := fmt.Sprintf("mkdir log dir %s error: %s\n", logDirPath, err)
				log.Println(msg)
				return err
			}
		}
	}
	if _, err := os.Stat(logFilePath); err != nil {
		if _, err := os.Create(logFilePath); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func NewDefaultLogger(LogFileName, level string, opts ...Option) *LogConf {
	var filename = "./log/std.log"
	if LogFileName == "" {
		LogFileName = filename
	}
	cfg := &LogConf{
		LogFileName: LogFileName,
		Level:       level,
		Stdout:      true,
		Fileout:     true,
		Lumberjack:  defaultLumberjackLogger(),
	}
	cfg.Lumberjack.Filename = LogFileName
	// apply options
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return cfg
}

type LogConf struct {
	LogFileName string `json:"log_file_name"` // ./log/std.log
	Formatter   string // json、text
	Stdout      bool   // 日志控制台输出
	Fileout     bool   // 日志文件输出
	Level       string
	Lumberjack  *lumberjack.Logger
}

// Option logger options
type Option interface {
	apply(cfg *LogConf)
}

type option func(cfg *LogConf)

func (fn option) apply(cfg *LogConf) {
	fn(cfg)
}

// WithLumberjackLogger configures logger
func WithLumberjackLogger(logger *lumberjack.Logger) Option {
	return option(func(cfg *LogConf) {
		cfg.Lumberjack = logger
	})
}
