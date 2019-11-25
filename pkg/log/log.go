package log

import (
	"context"
	"os"
	"sync"
	"time"

	"dnsx/pkg/config"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once

const TraceId = "trace_id"

// 日志级别
var levelType = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
}

// zlog logger 标准日志
type zlog struct {
	*zap.Logger
}

var zlogger = new(zlog)

// Logger new Logger
func Logger() *zlog {
	return zlogger
}

// level 日志级别操作
var level = zap.NewAtomicLevel()

// Setup init Logger
func Init() {
	once.Do(func() {
		handle := lumberjack.Logger{
			Filename:   getLogfilePath(),                // 日志文件路径
			MaxSize:    viper.GetInt("log.max_size"),    // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: viper.GetInt("log.max_backups"), // 日志文件最多保存多少个备份
			MaxAge:     viper.GetInt("log.max_age"),     // 文件最多保存多少天
			Compress:   true,                            // 是否压缩
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "category",
			CallerKey:      "line",
			MessageKey:     "msg",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		}

		SetLevel(config.GetString("log.log_level"))
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&handle)),
			level,
		)
		zlogger.Logger = zap.New(core, zap.AddCaller(), zap.Development()).
			With(zap.String("app_name", config.GetString("app.name")))

		// 注册配置变更事件
		config.RegisterChangeEvent(func(e fsnotify.Event) {
			SetLevel(config.GetString("log.log_level"))
		})
	})
}

// WithCTX 从上下文中获取 trace-id 并在日志中加入 trace-id 字段
func (l zlog) WithCTX(c context.Context) zlog {
	id, ok := c.Value(TraceId).(string)
	if !ok {
		id = ""
	}
	l.Logger = l.With(zap.String(TraceId, id))
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

// timeEncoder 日志时间格式化
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
