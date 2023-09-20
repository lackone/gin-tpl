package config

type DBWriteRead struct {
	Driver               string `yaml:"driver"`
	Host                 string `yaml:"host"`
	Port                 string `yaml:"port"`
	Username             string `yaml:"username"`
	Password             string `yaml:"password"`
	Database             string `yaml:"database"`
	Loc                  string `yaml:"loc"`
	Charset              string `yaml:"charset"`
	Collation            string `yaml:"collation"`
	Protocol             string `yaml:"protocol"`
	Timeout              string `yaml:"timeout"`
	WriteTimeout         string `yaml:"write_timeout"`
	ReadTimeout          string `yaml:"read_timeout"`
	ParseTime            bool   `yaml:"parse_time"`
	AllowNativePasswords bool   `yaml:"allow_native_passwords"`
	MaxIdleConn          int    `yaml:"max_idle_conn"`
	MaxOpenConn          int    `yaml:"max_open_conn"`
	ConnMaxLifetime      string `yaml:"conn_max_lifetime"`
}

type DB struct {
	Write DBWriteRead `yaml:"write"`
	Read  DBWriteRead `yaml:"read"`
}
