package tracing

import (
	"encoding/json"
	"os"
)

type Config struct {
	// EnableFile determines if the log should be writed to local file.
	EnableFile bool `json:"enablefile" yaml:"enablefile"`

	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string `json:"filename" yaml:"filename"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxsize" yaml:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxage" yaml:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"localtime" yaml:"localtime"`

	// Compress determines if the rotated log files should be compressed
	// using gzip.
	Compress bool `json:"compress" yaml:"compress"`

	// EnableConsole determines if the log should be displayed in stderr.
	EnableConsole bool `json:"enableconsole" yaml:"enableconsole"`

	// EnableCaller determines if the log should contain the caller
	EnableCaller bool `json:"enablecaller" yaml:"enablecaller"`

	// EnableSourceIP determines if the log should contain the sourceIP
	EnableSourceIP bool `json:"enablesourceip" yaml:"enablesourceip"`

	// EnablePID determines if the log should contain the PID
	EnablePID bool `json:"enablePID" yaml:"enablePID"`

	// log level in log file
	FileLevel string `json:"filelevel" yaml:"filelevel"`

	// log level in console
	ConsoleLevel string `json:"consolelevel" yaml:"consolelevel"`

	// encoding in log file. Valid values are "json" and
	// "console"
	FileEncodeing string `json:"fileencoding" yaml:"fileencoding"`

	// encoding in console. Valid values are "json" and
	// "console"
	ConsoleEncodeing string `json:"consoleencoding" yaml:"consoleencoding"`

	// application name
	// default is app
	AppName string `json:"appname" yaml:"appname"`

	// SourceEth determine which eth to get SourceIp
	// defautl is en0
	SourceEth string `json:"sourceeth" yaml:"sourceeth"`

	// DisableTraceID disable trace id
	DisableTraceID bool `json:"disable_trace_id" yaml:"disable_trace_id"`

	// GlobalCallerSkip increases the number of callers skipped
	GlobalCallerSkip int `json:"-" yaml:"-"`
}

func NewDevelopmentConfig(appname string, filename string) Config {
	return Config{
		EnableFile:       true,
		Filename:         filename,
		MaxSize:          0,
		MaxAge:           0,
		MaxBackups:       0,
		LocalTime:        true,
		Compress:         false,
		EnableConsole:    true,
		EnableCaller:     true,
		EnableSourceIP:   true,
		EnablePID:        true,
		FileLevel:        "debug",
		ConsoleLevel:     "debug",
		FileEncodeing:    "json",
		ConsoleEncodeing: "console",
		AppName:          appname,
		SourceEth:        "en0",
		DisableTraceID:   false,
	}
}

func NewProductionConfig(appname string, filename string) Config {
	return Config{
		EnableFile:       true,
		Filename:         filename,
		MaxSize:          0,
		MaxAge:           0,
		MaxBackups:       0,
		LocalTime:        true,
		Compress:         false,
		EnableConsole:    false,
		EnableCaller:     true,
		EnableSourceIP:   true,
		EnablePID:        true,
		FileLevel:        "info",
		ConsoleLevel:     "info",
		FileEncodeing:    "json",
		ConsoleEncodeing: "console",
		AppName:          appname,
		SourceEth:        "en0",
		DisableTraceID:   false,
	}
}

func NewStdConfig() Config {
	return Config{
		EnableFile:       false,
		Filename:         "",
		MaxSize:          0,
		MaxAge:           0,
		MaxBackups:       0,
		LocalTime:        true,
		Compress:         false,
		EnableConsole:    true,
		EnableCaller:     true,
		EnableSourceIP:   false,
		EnablePID:        true,
		FileLevel:        "info",
		ConsoleLevel:     "debug",
		FileEncodeing:    "json",
		ConsoleEncodeing: "console",
		AppName:          "",
		SourceEth:        "en0",
		DisableTraceID:   false,
	}
}

func NewConfig(filename string) (Config, error) {
	file, _ := os.Open(filename)
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)

	return config, err
}

type Option func(*Config)

func WithGlobalCallerSkip(n int) Option {
	return func(c *Config) {
		c.GlobalCallerSkip = n
	}
}
