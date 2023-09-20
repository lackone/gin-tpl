package app

import (
	"gin-tpl/app/core/config"
	"gin-tpl/app/core/db"
	"gorm.io/gorm"
	"strings"
	"sync"
)

type App struct {
	Path              string              //项目路径
	Config            *config.Config      //配置
	ConfigDefaultPath string              //默认配置文件路径
	ConfigPath        string              //当前配置文件路径
	DB                map[string]*gorm.DB //数据库连接
	DBLock            sync.RWMutex        //数据库锁
}

// NewApp 实例化一个Application对象
func NewApp(path string) *App {
	return &App{
		Path:   path,
		Config: &config.Config{},
		DB:     map[string]*gorm.DB{},
		DBLock: sync.RWMutex{},
	}
}

// InitConfig 初始化配置
func (app *App) InitConfig() error {
	app.ConfigDefaultPath = strings.TrimRight(app.Path, "/") + "/config/config-default.yml"
	app.ConfigPath = strings.TrimRight(app.Path, "/") + "/config/config.yml"

	err := app.Config.Load(app.ConfigDefaultPath)
	if err != nil {
		return err
	}
	err = app.Config.Load(app.ConfigPath)
	if err != nil {
		return err
	}

	app.Config.LoadEnv()

	return nil
}

// InitDB 初始化数据库
func (app *App) InitDB() {
	app.DBW()
	app.DBR()
}

// DBW 获取写数据库
func (app *App) DBW(keys ...string) *gorm.DB {
	key := "default"
	if len(keys) > 0 {
		key = keys[0]
	}
	dbKey := key + ".write"

	app.DBLock.RLock()
	if _, ok := app.DB[dbKey]; ok {
		app.DBLock.RUnlock()
		return app.DB[dbKey]
	}
	app.DBLock.RUnlock()

	conf, ok := app.Config.DB[key]
	if !ok {
		panic("数据库配置不存在")
	}

	app.DBLock.Lock()
	defer app.DBLock.Unlock()

	if _, exists := app.DB[dbKey]; exists {
		return app.DB[dbKey]
	}

	dbs, err := db.Load(conf.Write)
	if err != nil {
		panic(err)
	}

	app.DB[dbKey] = dbs

	return dbs
}

// DBR 获取读数据库
func (app *App) DBR(keys ...string) *gorm.DB {
	key := "default"
	if len(keys) > 0 {
		key = keys[0]
	}
	dbKey := key + ".read"

	app.DBLock.RLock()
	if _, ok := app.DB[dbKey]; ok {
		app.DBLock.RUnlock()
		return app.DB[dbKey]
	}
	app.DBLock.RUnlock()

	conf, ok := app.Config.DB[key]
	if !ok {
		panic("数据库配置不存在")
	}

	app.DBLock.Lock()
	defer app.DBLock.Unlock()

	if _, exists := app.DB[dbKey]; exists {
		return app.DB[dbKey]
	}

	dbs, err := db.Load(conf.Read)
	if err != nil {
		panic(err)
	}

	app.DB[dbKey] = dbs

	return dbs
}

// Init 初始化
func (app *App) Init() error {
	err := app.InitConfig()
	if err != nil {
		return err
	}
	app.InitDB()
	return nil
}
