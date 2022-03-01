package helper

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	hostname          = kingpin.Flag("host", "Hostname").String()
	dbname            = kingpin.Flag("db", "Database name").String()
	user              = kingpin.Flag("user", "Username").String()
	password          = kingpin.Flag("pwd", "Password").String()
	port              = kingpin.Flag("port", "Database port").String()
	server_address    = kingpin.Flag("server", "Server_IP:Port").String()
	minio_server      = kingpin.Flag("minio_server", "Minio Server_IP:Port").String()
	minio_pwd         = kingpin.Flag("minio_pwd", "Minio Password").String()
	minio_user        = kingpin.Flag("minio_user", "Minio user").String()
	redis_addr        = kingpin.Flag("redis_server", "Redis Server_IP:Port").String()
	redis_pwd         = kingpin.Flag("redis_pwd", "Redis Password").String()
	rabbitmq_server   = kingpin.Flag("rabbitmq_server", "RabbitMQ Server").String()
	rabbitmq_user     = kingpin.Flag("rabbitmq_user", "RabbitMQ Username").String()
	rabbitmq_pwd      = kingpin.Flag("rabbitmq_pwd", "RabbitMQ Password").String()
	face_worker_queue = kingpin.Flag("face_queue", "RabbitMQ Queue Name for Face match worker").String()
)

func SetEnvVariablesUtil() {

	kingpin.Parse()

	setEnvVariables("DB_HOST", *hostname)
	setEnvVariables("DB_NAME", *dbname)
	setEnvVariables("DB_USER", *user)
	setEnvVariables("DB_PASSWORD", *password)
	setEnvVariables("DB_PORT", *port)
	setEnvVariables("SERVER_ADDR", *server_address)
	setEnvVariables("MINIO_SERVER", *minio_server)
	setEnvVariables("MINIO_USER", *minio_user)
	setEnvVariables("MINIO_PWD", *minio_pwd)
	setEnvVariables("REDIS_SERVER", *redis_addr)
	setEnvVariables("REDIS_PASSWORD", *redis_pwd)
	setEnvVariables("RABBITMQ_SERVER", *rabbitmq_server)
	setEnvVariables("RABBITMQ_USER", *rabbitmq_user)
	setEnvVariables("RABBITMQ_PWD", *rabbitmq_pwd)
	setEnvVariables("FACE_WORKER_QUEUE", *face_worker_queue)
}

func setEnvVariables(key, value string) {
	if (len(value)) > 0 {
		os.Setenv(key, value)
	}
}
