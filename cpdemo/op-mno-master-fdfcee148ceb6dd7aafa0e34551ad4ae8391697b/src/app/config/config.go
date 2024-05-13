package config

import (
	configs "codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-pkg.git/config"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-pkg.git/mysql8"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type (
	Config struct {
		configs.App  `yaml:"app"`
		configs.HTTP `yaml:"http"`
		configs.Log  `yaml:"logger"`
		Client       `yaml:"client"`
		DB           `yaml:"mysql"`
	}
	Client struct {
		ConnectClientUrl string `env-required:"true" yaml:"connect_client_url" env:"Client_CONNECT_URL"`
	}

	DB struct {
		PoolMax     int                 `env-required:"true" yaml:"pool_max" env:"MYSQL_POOL_MAX"`
		IdleConnMax int                 `env-required:"true" yaml:"idle_conn_max" env:"MYSQL_IDLE_CONN_MAX"`
		MaxIdleTime int                 `env-required:"true" yaml:"max_idle_time" env:"MYSQL_MAX_IDLE_TIME"`
		DsnUrl      mysql8.DBConnString `env-required:"true" yaml:"dsn_url" env:"MYSQL_DSN_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// debug
	fmt.Println(dir)
	//err = cleanenv.ReadConfig("D:\\go-workspace\\src\\op-mno\\src\\app\\config\\"+"/config.yml", cfg)
	//err = cleanenv.ReadConfig("D:\\Environment\\GoWorks\\src\\op-mno\\src\\app\\config\\"+"/config.yml", cfg)
	err = cleanenv.ReadConfig(dir+"/src/app/config/config.yml", cfg)
	//	err = cleanenv.ReadConfig("D:\\cyy\\study\\op-mno-master-fdfcee148ceb6dd7aafa0e34551ad4ae8391697b\\op-mno-master-fdfcee148ceb6dd7aafa0e34551ad4ae8391697b\\src\\app\\config", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
