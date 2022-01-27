package main

import (
	"os"
)

func SetEnvVariablesUtil() {

	setEnvVariables("DB_HOST", *hostname)
	setEnvVariables("DB_NAME", *dbname)
	setEnvVariables("DB_USER", *user)
	setEnvVariables("DB_PASSWORD", *password)
	setEnvVariables("DB_PORT", *port)
	setEnvVariables("SERVER_ADDR", *server_address)
	setEnvVariables("MINIO_SERVER", *minio_server)
	setEnvVariables("MINIO_USER", *minio_user)
	setEnvVariables("MINIO_PWD", *minio_pwd)
}

func setEnvVariables(key, value string) {
	if (len(value)) > 0 {
		os.Setenv(key, value)
	}
}
