package config

import (
	"os"
	"strings"
)

type Minio struct {
	Endpoint string
	Bucket   struct {
		User string
	}
}

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string

	Context struct {
		Timeout string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SslMode  string
	}

	OTLPCollector struct {
		Host string
		Port string
	}

	Kafka struct {
		Address []string
		Topic   struct {
			InvestorCreate string
		}
	}
	MinioService Minio
}

func New() *Config {
	var c Config

	// general configuration
	c.APP = getEnv("APP", "dennic_user_service")
	c.Environment = getEnv("ENVIRONMENT", "develop")
	c.LogLevel = getEnv("LOG_LEVEL", "debug")
	c.RPCPort = getEnv("RPC_PORT", ":9070")
	c.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// db configuration
	c.DB.Host = getEnv("POSTGRES_HOST", "postgres")
	c.DB.Port = getEnv("POSTGRES_PORT", "5432")
	c.DB.User = getEnv("POSTGRES_USER", "postgres")
	c.DB.Password = getEnv("POSTGRES_PASSWORD", "20030505")
	c.DB.SslMode = getEnv("POSTGRES_SSLMODE", "disable")
	c.DB.Name = getEnv("POSTGRES_DATABASE", "dennic")

	// otlp collector configuration
	c.OTLPCollector.Host = getEnv("OTLP_COLLECTOR_HOST", "otel-collector")
	c.OTLPCollector.Port = getEnv("OTLP_COLLECTOR_PORT", ":4317")

	// kafka configuration
	c.Kafka.Address = strings.Split(getEnv("KAFKA_ADDRESS", "localhost:29092"), ",")
	c.Kafka.Topic.InvestorCreate = getEnv("KAFKA_TOPIC_INVESTOR_CREATE", "investor.created")

	// Minio
	c.MinioService.Endpoint = getEnv("MINIO_SERVICE_ENDPOINT", "https://minio.dennic.uz")
	c.MinioService.Bucket.User = getEnv("MINIO_SERVICE_BUCKET_USER", "user")

	return &c
}

func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
