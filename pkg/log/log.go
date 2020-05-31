package log

import (
	"context"
	"os"
	"sync"

	"dnsx/pkg/config"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once

const TraceID = "trace_id"

// 日志级别
var levelType = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
}

// Logger logger 标准日志
type Logger struct {
	*zap.Logger
}

var logger = new(Logger)

// Logger new Logger
func New() *Logger {
	return logger
}

// level 日志级别操作
var level = zap.NewAtomicLevel()

// Setup init Logger
func Init() {
	once.Do(func() {

		SetLevel(config.GetString("log.log_level"))
		core := zapcore.NewCore(
			encoder(),
			zapcore.NewMultiWriteSyncer(writers()...),
			level,
		)
		logger.Logger = zap.New(core, zap.AddCaller(), zap.Development()).
			With(zap.String("app_name", config.GetString("app.name")))

		// 注册配置变更事件
		config.RegisterChangeEvent(func(e fsnotify.Event) {
			SetLevel(config.GetString("log.log_level"))
		})
	})
}

// WithCTX 从上下文中获取 trace-id 并在日志中加入 trace-id 字段
func (l Logger) WithContext(c context.Context) Logger {
	id, ok := c.Value(TraceID).(string)
	if !ok {
		id = ""
	}
	l.Logger = l.With(zap.String(TraceID, id))
	return l
}

// SetLevel 设置日志级别
func SetLevel(name string) {
	var l zapcore.Level
	if v, ok := levelType[name]; ok {
		l = v
	} else {
		l = zap.InfoLevel
	}
	if l == GetLevel() {
		return
	}
	level.SetLevel(l)
}

// GetLevel 获取当前日志级别
func GetLevel() zapcore.Level {
	return level.Level()
}

// getLogfilePath 获取日志文件全路径
func getLogfilePath() string {
	return config.GetString("log.log_path") + config.GetString("log.log_file_name") + ".log"
}

// writers 日志输出
func writers() (ws []zapcore.WriteSyncer) {
	handle := lumberjack.Logger{
		Filename:   getLogfilePath(),
		MaxSize:    viper.GetInt("log.max_size"),
		MaxBackups: viper.GetInt("log.max_backups"),
		MaxAge:     viper.GetInt("log.max_age"),
		Compress:   true,
	}
	ws = []zapcore.WriteSyncer{
		zapcore.AddSync(&handle),
	}
	if viper.GetBool("log.stdout") {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}
	return
}

// encoder 日志编码
func encoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "category",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	if config.GetString("log.encoder") == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}
