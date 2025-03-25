package config

import (
	"log"
	"os"
)

func GetEnv() string {
	return getEnvironmentValue("ENV")
}
func GetTarantoolUser() string {
	return getEnvironmentValue("TARANTOOL_USER_NAME")
}
func GetTarantoolPassword() string {
	return getEnvironmentValue("TARANTOOL_USER_PASSWORD")
}
func GetTarantoolHost() string {
	return getEnvironmentValue("TARANTOOL_HOST")
}
func GetApplicationPort() string {
	port := getEnvironmentValue("APPLICATION_PORT")
	return port
}
func GetSecretKey() string {
	return getEnvironmentValue("SECRET_KEY")
}

func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}
	return os.Getenv(key)
}
