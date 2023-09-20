package db

import (
	"gin-tpl/app/core/config"
	"github.com/go-sql-driver/mysql"
	"gorm.io/driver/clickhouse"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"net"
	"time"
)

func Load(c config.DBWriteRead) (*gorm.DB, error) {
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

	switch c.Driver {
	case "mysql":
		db, err = gorm.Open(gmysql.Open(dsn), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
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
