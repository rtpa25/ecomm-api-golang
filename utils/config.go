package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver           string `mapstructure:"DB_DRIVER"`
	DBSource           string `mapstructure:"DB_SOURCE"`
	ServerAddress      string `mapstructure:"SERVER_ADDRESS"`
	WebsiteAddress     string `mapstructure:"WEBSITE_ADDRESS"`
	ConnectionUri      string `mapstructure:"CONNECTION_URI"`
	ApiKey             string `mapstructure:"API_KEY"`
	ServerDomainLocal  string `mapstructure:"SERVER_DOMAIN_LOCAL"`
	ServerDomainProd   string `mapstructure:"SERVER_DOMAIN_PROD"`
	WebsiteDomainLocal string `mapstructure:"WEBSITE_DOMAIN_LOCAL"`
	WebsiteDomainProd  string `mapstructure:"WEBSITE_DOMAIN_PROD"`
	GoEnv              string `mapstructure:"GO_ENV"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
