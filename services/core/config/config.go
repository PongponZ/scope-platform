package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		ServerPort        string
		RabbitMQURL       string
		MediaConvertQueue string
		MediaDomain       string
		RedisURL          string
		Minio             Minio
	}

	Minio struct {
		Bucket    string
		Endpoint  string
		AccessKey string
		Secert    string
	}
)

func GetConfig() Config {
	godotenv.Load(".env")
	var (
		serverPort        string
		rabbitMQURL       string
		mediaConvertQueue string
		mediaDomain       string
		redisURL          string
		minio             Minio
	)

	serverPort = fmt.Sprintf(":%s", os.Getenv("SEVER_LISTEN_PORT"))

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisURL = fmt.Sprintf("%s:%s", redisHost, redisPort)

	minio.Bucket = os.Getenv("MINIO_BUCKET")
	minio.Endpoint = os.Getenv("MINIO_ENDPOINT")
	minio.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	minio.Secert = os.Getenv("MINIO_SECERT")

	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQPort := os.Getenv("RABBITMQ_PORT")
	rabbitMQUsername := os.Getenv("RABBITMQ_USERNAME")
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitMQURL = fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitMQUsername, rabbitMQPassword, rabbitMQHost, rabbitMQPort)

	mediaConvertQueue = os.Getenv("RABBITMQ_MEDIA_CONVERT_QUEUE_NAME")

	mediaDomain = os.Getenv("MEDIA_PROVIDER_DOMAIN")

	return Config{
		ServerPort:        serverPort,
		RabbitMQURL:       rabbitMQURL,
		MediaConvertQueue: mediaConvertQueue,
		MediaDomain:       mediaDomain,
		RedisURL:          redisURL,
		Minio:             minio,
	}
}
