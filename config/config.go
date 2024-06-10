package config

import "time"

type Config struct {
	LogEnv   string      `yaml:"env" required:"true"`
	Server   SRVConfig   `yaml:"server" required:"true"`
	Database DBConfig    `yaml:"database" required:"true"`
	Cache    CacheConfig `yaml:"cache" required:"true"`
}

type SRVConfig struct {
	Addr    string        `yaml:"addr" required:"true"`
	Timeout time.Duration `yaml:"timeout" required:"true"`
}

type DBConfig struct {
	URI              string        `yaml:"uri" required:"true"`
	Database         string        `yaml:"database" required:"true"`
	Username         string        `yaml:"username" required:"true"`
	Password         string        `yaml:"password" required:"true"`
	MaxExecutionTime int           `yaml:"max_execution_time" required:"true"`
	DialTimeout      time.Duration `yaml:"dial_timeout" required:"true"`
}

type CacheConfig struct {
	URI string        `yaml:"uri" required:"true"`
	Exp time.Duration `yaml:"exp" required:"true"`
}
