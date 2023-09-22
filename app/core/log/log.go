package log

import (
	"fmt"
	"gin-tpl/app/core/config"
	"gin-tpl/app/utils"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"strings"
)

func LoadLog(c config.Log, env string) (*zerolog.Logger, error) {
	var log zerolog.Logger
	var err error
	var level zerolog.Level

	timeFormat := "2006-01-02 15:04:05"
	zerolog.TimeFieldFormat = timeFormat

	switch c.Type {
	case "stdout":
		w := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}
		w.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		w.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		w.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		w.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s;", i)
		}
		log = zerolog.New(w)
	case "file":
		path, _ := filepath.Abs(c.Path)
		w := &lumberjack.Logger{
			Filename:   path,
			MaxSize:    50,
			MaxBackups: 10,
			MaxAge:     10,
			Compress:   false,
		}
		log = zerolog.New(w)
	}

	switch c.Level {
	case "panic":
		level = zerolog.PanicLevel
	case "fatal":
		level = zerolog.FatalLevel
	case "error":
		level = zerolog.ErrorLevel
	case "warn":
		level = zerolog.WarnLevel
	case "info":
		level = zerolog.InfoLevel
	case "debug":
		level = zerolog.DebugLevel
	case "trace":
		level = zerolog.TraceLevel
	}

	hostname, _ := os.Hostname()
	localIp, _ := utils.GetInternalIP()

	log = log.Level(level).With().Timestamp().Str("hostname", hostname).Str("local_ip", localIp).Str("env", env).Logger()

	return &log, err
}
