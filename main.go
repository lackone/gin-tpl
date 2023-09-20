package main

import (
	"flag"
	"gin-tpl/app"
	"gin-tpl/app/http"
	"path/filepath"
)

var (
	path string   // 项目路径
	a    *app.App //项目应用
)

func main() {
	flag.StringVar(&path, "path", "", "项目路径")
	flag.Parse()

	if len(path) == 0 {
		path = "./"
	}

	path, err := filepath.Abs(path)
	if err != nil {
		panic("路径获取失败" + err.Error())
	}

	a = app.NewApp(path)
	err = a.Init()
	if err != nil {
		panic("初始化失败" + err.Error())
	}

	err = http.NewServer(a)
	if err != nil {
		panic("服务启动失败" + err.Error())
	}
}
