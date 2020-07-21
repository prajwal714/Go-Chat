package config

import "github.com/spf13/viper"

type Config struct {
	port               int
	secretKey          string
	githubClientID     string
	githubClientSecret string
	callbackURL        string
}

var appConfig *Config

func Load() {
	viper.AutomaticEnv()

	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")

	viper.SetConfigType("yaml")

	viper.ReadInConfig()

	appConfig = &Config{
		port:               getIntOrPanic("APP_PORT"),
		secretKey:          fatalGetString("SECRET_KEY"),
		githubClientID:     fatalGetString("GITHUB_CLIENT_ID"),
		githubClientSecret: fatalGetString("GITHUB_CLIENT_SECRET"),
		callbackURL:        fatalGetString("CALLBACK_URL"),
	}

}

func Port() int {
	return appConfig.port
}

func SecretKey() string {
	return appConfig.secretKey
}

func GithubClientID() string {
	return appConfig.githubClientID
}

func GithubClientSecret() string {
	return appConfig.githubClientSecret
}

func CallbackURL() string {
	return appConfig.callbackURL
}
