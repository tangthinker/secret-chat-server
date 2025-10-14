package core

import "github.com/spf13/viper"

type Config struct {
	configPath string
}

func NewConfig(configPath string) *Config {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	return &Config{
		configPath: configPath,
	}
}

func (c *Config) GetString(key string) string {
	return viper.GetString(key)
}

func (c *Config) GetInt(key string) int {
	return viper.GetInt(key)
}
