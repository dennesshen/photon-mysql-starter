package mysql

type ConnectData struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Database struct {
		Master     ConnectData `yaml:"master"`
		Slave      ConnectData `yaml:"slave"`
		Connection struct {
			MaxIdleConns      int `yaml:"maxIdleConns"`
			MaxOpenConns      int `yaml:"maxOpenConns"`
			MaxLifetimeSecond int `yaml:"maxLifetimeSecond"`
		} `yaml:"connection"`
		Database string `yaml:"database"`
	} `yaml:"database"`
}

var config Config
