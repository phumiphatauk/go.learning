package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server           ServerConfig     `mapstructure:"server"`
	Databasepostgres Databasepostgres `mapstructure:"databasepostgres"`
}

type ServerConfig struct {
	Port     uint   `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Databasepostgres struct {
	Host     string `mapstructure:"host"`
	Port     uint   `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	Appname  string `mapstructure:"appname"`
}

func LoadConfig() (config Config, err error) {
	var c Config
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return c, err
	}

	viper.Unmarshal(&c)

	c.Server = ServerConfig{
		Port:     getEnvInteger("server.port", c.Server.Port),
		User:     getEnv("server.user", c.Server.User),
		Password: getEnv("server.password", c.Server.Password),
	}

	fmt.Printf("Port after %d\n", c.Server.Port)

	return c, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInteger(key string, defaultValue uint) uint {
	if strValue := os.Getenv(key); strValue != "" {
		var intValue uint
		_, err := fmt.Sscan(strValue, &intValue)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true"
	}
	return defaultValue
}
