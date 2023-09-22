package db

import (
	"gin-tpl/app/core/config"
	"github.com/go-sql-driver/mysql"
	"github.com/natefinch/lumberjack"
	"gorm.io/driver/clickhouse"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

func LoadDB(c config.DBWriteRead, l config.DBLog) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return nil, err
	}
	readTimeout, err := time.ParseDuration(c.ReadTimeout)
	if err != nil {
		return nil, err
	}
	writeTimeout, err := time.ParseDuration(c.WriteTimeout)
	if err != nil {
		return nil, err
	}
	connMaxLifetime, err := time.ParseDuration(c.ConnMaxLifetime)
	if err != nil {
		return nil, err
	}
	location, err := time.LoadLocation(c.Loc)
	if err != nil {
		return nil, err
	}

	conf := &mysql.Config{
		User:                 c.Username,
		Passwd:               c.Password,
		Net:                  c.Protocol,
		Addr:                 net.JoinHostPort(c.Host, c.Port),
		DBName:               c.Database,
		Collation:            c.Collation,
		Loc:                  location,
		Timeout:              timeout,
		ReadTimeout:          readTimeout,
		WriteTimeout:         writeTimeout,
		ParseTime:            c.ParseTime,
		AllowNativePasswords: c.AllowNativePasswords,
	}

	dsn := conf.FormatDSN()

	var newLogger logger.Interface

	if l.Enable {
		var writer logger.Writer
		var level logger.LogLevel

		switch l.Type {
		case "stdout":
			writer = log.New(os.Stdout, "\r\n", log.LstdFlags)
		case "file":
			path, _ := filepath.Abs(l.Path)
			w := &lumberjack.Logger{
				Filename:   path,
				MaxSize:    50,
				MaxBackups: 10,
				MaxAge:     10,
				Compress:   false,
			}
			writer = log.New(w, "\r\n", log.LstdFlags)
		default:
			writer = log.New(os.Stdout, "\r\n", log.LstdFlags)
		}

		switch l.Level {
		case "silent":
			level = logger.Silent
		case "error":
			level = logger.Error
		case "warn":
			level = logger.Warn
		case "info":
			level = logger.Info
		}

		newLogger = logger.New(
			writer,
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  level,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
				Colorful:                  false,
			},
		)
	}

	switch c.Driver {
	case "mysql":
		db, err = gorm.Open(gmysql.Open(dsn), &gorm.Config{Logger: newLogger})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: newLogger})
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{Logger: newLogger})
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{Logger: newLogger})
	}

	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	return db, err
}
