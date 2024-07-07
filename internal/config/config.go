package config

import (
	"bytes"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Config struct {
	Server Server
	//Kafka         service.Kafka
	Logger   Logger
	Postgres Postgres
}

type Server struct {
	Port                  string
	Development           bool
	JWTKey                string
	TrustedReverseProxies []string
}

type Logger struct {
	Level string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	// pool settings
	MaxOpenConns int
	MaxIdleConns int

	MigrationsDirPath string
}

var configFileName = "config.yaml"

func ReadConfig(configObj *Config) error {
	if err := exportConfig("", configFileName); err != nil {
		return err
	}

	return parseConfig(configObj)
}

func parseConfig(cfg *Config) error {
	for _, key := range getKeysFromConfig(cfg) {
		err := viper.BindEnv(key)
		if err != nil {
			return err
		}
	}

	err := viper.Unmarshal(cfg)
	if err != nil {
		return err
	}

	return nil
}

func getKeysFromConfig(cfg interface{}) []string {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return nil
	}

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(b)); err != nil {
		return nil
	}

	return v.AllKeys()
}

func exportConfig(configPath string, configFileName string) error {
	viper.SetConfigType("yaml")
	if configPath == "" {
		if os.Getenv("MODE") == "DOCKER" {
			viper.AddConfigPath("/configs")
		} else {
			//development mode
			viper.AddConfigPath("configs")
		}
	} else {
		viper.AddConfigPath(configPath)
	}
	viper.SetConfigName(configFileName)

	//read config from env
	viper.AutomaticEnv()
	viper.SetEnvPrefix("x-test")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetTypeByDefaultValue(true)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; fallback to env
			return nil
		} else {
			// Config file was found but another error was produced
			return err
		}
	}
	return nil
}
