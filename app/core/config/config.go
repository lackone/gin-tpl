package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	Name string        `yaml:"name" env:"APP_NAME"`
	Http Http          `yaml:"http"`
	Log  Log           `yaml:"log"`
	DB   map[string]DB `yaml:"db"`
}

type Http struct {
	Listen string `yaml:"listen" env:"APP_HTTP_LISTEN"`
}

type Log struct {
	Level string `yaml:"level" env:"APP_LOG_LEVEL"`
}

// Load 加载配置
func (c *Config) Load(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		return err
	}
	return nil
}

// LoadEnv 加载环境变量
func (c *Config) LoadEnv() {
	c.handleStruct(reflect.ValueOf(c))
}

// handleStruct 处理结构体
func (c *Config) handleStruct(val reflect.Value) {
	v := reflect.Indirect(val)
	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Tag.Get("env")
		c.handleField(key, v.Field(i))
	}
}

// handleField 处理结构体某一项的值
func (c *Config) handleField(key string, val reflect.Value) {
	env, ok := os.LookupEnv(key)

	switch val.Type().Kind() {
	case reflect.String:
		if !ok {
			return
		}
		val.SetString(env)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !ok {
			return
		}
		v, err := strconv.ParseInt(env, 0, val.Type().Bits())
		if err != nil {
			return
		}
		val.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !ok {
			return
		}
		v, err := strconv.ParseUint(env, 0, val.Type().Bits())
		if err != nil {
			return
		}
		val.SetUint(v)
	case reflect.Float32, reflect.Float64:
		if !ok {
			return
		}
		v, err := strconv.ParseFloat(env, val.Type().Bits())
		if err != nil {
			return
		}
		val.SetFloat(v)
	case reflect.Bool:
		if !ok {
			return
		}
		v, err := strconv.ParseBool(env)
		if err != nil {
			return
		}
		val.SetBool(v)
	case reflect.Struct:
		c.handleStruct(val)
	}
}
