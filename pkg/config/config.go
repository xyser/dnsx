package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	// auto load env
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

// changeEventHandle 配置变更处理器
var changeEventHandle []func(e fsnotify.Event)
var eventLock sync.Mutex
var once sync.Once

// CfgFile 配置文件路径,允许在初始化前,由外部包赋值
var CfgFile string

// Init auto load config
func Init() {
	once.Do(func() {
		// 设置配置文件目录和文件名
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.SetConfigName("dnsx")

		if len(CfgFile) > 0 {
			viper.SetConfigFile(CfgFile)
		}

		// read in environment variables that match
		viper.SetEnvPrefix("dnsx")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		if err := viper.ReadInConfig(); err == nil {
			fmt.Printf("\033[1;30;42m[info]\033[0m using config file %s\n", viper.ConfigFileUsed())
		} else {
			fmt.Printf("\033[1;30;41m[error]\033[0m using config file error %s\n", err.Error())
			os.Exit(1)
		}

		// 监听配置文件变更
		viper.WatchConfig()
		// 调用 config 变更注册
		viper.OnConfigChange(onConfigChange)
	})
}

// RegisterChangeEvent 注册配置变更事件
func RegisterChangeEvent(f func(e fsnotify.Event)) {
	eventLock.Lock()
	defer eventLock.Unlock()

	changeEventHandle = append(changeEventHandle, f)
}

// onConfigChange 循环执行事件调用
func onConfigChange(e fsnotify.Event) {
	for _, f := range changeEventHandle {
		f(e)
	}
}

// Get can retrieve any value given the key to use.
func Get(k string) interface{} { return viper.Get(k) }

// GetBool returns the value associated with the key as a boolean.
func GetBool(k string) bool { return viper.GetBool(k) }

// GetString returns the value associated with the key as a string.
func GetString(key string) string { return viper.GetString(key) }

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int { return viper.GetInt(key) }

// GetInt32 returns the value associated with the key as an integer.
func GetInt32(key string) int32 { return viper.GetInt32(key) }

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 { return viper.GetInt64(key) }

// GetUint returns the value associated with the key as an unsigned integer.
func GetUint(key string) uint { return viper.GetUint(key) }

// GetUint32 returns the value associated with the key as an unsigned integer.
func GetUint32(key string) uint32 { return viper.GetUint32(key) }

// GetUint64 returns the value associated with the key as an unsigned integer.
func GetUint64(key string) uint64 { return viper.GetUint64(key) }

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 { return viper.GetFloat64(key) }

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time { return viper.GetTime(key) }

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration { return viper.GetDuration(key) }

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string { return viper.GetStringSlice(key) }

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} { return viper.GetStringMap(key) }

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string { return viper.GetStringMapString(key) }

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key string) map[string][]string {
	return viper.GetStringMapStringSlice(key)
}

// GetSizeInBytes returns the size of the value associated with the given key
// in bytes.
func GetSizeInBytes(key string) uint { return viper.GetSizeInBytes(key) }

// IsSet checks to see if the key has been set in any of the data locations.
// IsSet is case-insensitive for a key.
func IsSet(key string) bool { return viper.IsSet(key) }

// AllSettings merges all settings and returns them as a map[string]interface{}.
func AllSettings() map[string]interface{} { return viper.AllSettings() }

// AllKeys returns all keys holding a value, regardless of where they are set.
// Nested keys are returned with a v.keyDelim (= ".") separator
func AllKeys() []string { return viper.AllKeys() }
