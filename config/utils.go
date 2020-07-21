package config

import (
	"log"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/spf13/viper"
)

func getIntOrPanic(key string) int {
	checkKey(key)
	v, err := strconv.Atoi(fatalGetString(key))
	if err != nil {
		v, err = strconv.Atoi(os.Getenv(key))
		debug.PrintStack()
		log.Fatalf("Could not parse key: %s, Error: %s", key, err)
	}
	return v
}

func fatalGetString(key string) string {
	checkKey(key)
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}
	return value
}

func getString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}
	return value
}

func checkKey(key string) {
	if !viper.IsSet(key) && os.Getenv(key) == "" {
		log.Print("can't find key", key)
		debug.PrintStack()
		log.Fatalf("%s key is not set", key)
	}
}
