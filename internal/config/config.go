package config

import (
	"local/test-tasks/assisted_team/internal/server"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type MainConf struct {
	DataFilesReturn []string `mapstructure:"files_return"`
	DataFilesOneWay []string `mapstructure:"files_oneway"`
}

type Config struct {
	MainConf       *MainConf          `mapstructure:"main"`
	UserServerConf *server.ServerConf `mapstructure:"user_server"`
}

func Init() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")

	//add also parser from env variables for easier docker access
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	// api_server defaults
	viper.SetDefault("user_server.read_timeout", 10*time.Second)
	viper.SetDefault("user_server.write_timeout", 10*time.Second)
	viper.SetDefault("user_server.max_header_bytes", 1<<20)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
