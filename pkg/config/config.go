package config

import (
	"fmt"
	"strings"

	"github.com/iamolegga/enviper"
	"github.com/spf13/viper"
)

type Config struct {
	Production bool `mapstructure:"production"`

	// logger
	LogLevel    string `mapstructure:"log_level"`
	LogResponse bool   `mapstructure:"log_response"`

	// base url
	BaseURL string `mapstructure:"base_url"`
	Port    string `mapstructure:"port"`

	// MySQL
	MySQLHost     string `mapstructure:"mysql_host"`
	MySQLUser     string `mapstructure:"mysql_user"`
	MySQLPass     string `mapstructure:"mysql_pass"`
	MySQLDatabase string `mapstructure:"mysql_database"`

	// HMAC Secret
	HMACSecret string `mapstructure:"hmac_secret"`

	// Auth
	GoogleClientID     string `mapstructure:"google_client_id"`
	GoogleClientSecret string `mapstructure:"google_client_secret"`
	GoogleClientState  string `mapstructure:"google_client_state"`
	GoogleRedirectURL  string `mapstructure:"google_redirect_url"`

	FacebookClientID     string `mapstructure:"facebook_client_id"`
	FacebookClientSecret string `mapstructure:"facebook_client_secret"`
	FacebookClientState  string `mapstructure:"facebook_client_state"`
	FacebookRedirectURL  string `mapstructure:"facebook_redirect_url"`

	GithubClientID     string `mapstructure:"github_client_id"`
	GithubClientSecret string `mapstructure:"github_client_secret"`
	GithubClientState  string `mapstructure:"github_client_state"`
	GithubRedirectURL  string `mapstructure:"github_redirect_url"`
}

func New() *Config {
	viper := enviper.New(viper.New())

	// initialize config
	viper.SetConfigName(".apicfg")
	viper.AddConfigPath("../..")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("GOREST")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err == nil {
		viper.WatchConfig()
	}

	viper.SetDefault("production", false)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("base_url", "http://127.0.0.1")
	viper.SetDefault("port", "8080")
	viper.SetDefault("mysql_host", "127.0.0.1:3306")
	viper.SetDefault("mysql_user", "root")
	viper.SetDefault("mysql_pass", "")
	viper.SetDefault("mysql_database", "gorest_db")

	// read config
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Sprintf("failed to read config: %s", err.Error()))
	}

	return &config
}
