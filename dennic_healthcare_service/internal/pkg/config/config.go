package config

import (
	"os"
	"strings"
)

type Minio struct {
	Endpoint string
	Bucket   struct {
		Department     string
		Doctor         string
		Reasons        string
		Specialization string
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
			Healthcare string
		}
	}
	MinioService Minio
}

func New() *Config {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "dennic_healthcare_service")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.RPCPort = getEnv("RPC_PORT", ":9080")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "postgres")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "20030505")
	config.DB.SslMode = getEnv("POSTGRES_SSLMODE", "disable")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "dennic")

	// otlp collector configuration
	config.OTLPCollector.Host = getEnv("OTLP_COLLECTOR_HOST", "otel-collector")
	config.OTLPCollector.Port = getEnv("OTLP_COLLECTOR_PORT", ":4317")

	// kafka configuration
	config.Kafka.Address = strings.Split(getEnv("KAFKA_ADDRESS", "localhost:29092"), ",")
	config.Kafka.Topic.Healthcare = getEnv("KAFKA_TOPIC_HEALTHCARE_CREATE", "user.created")

	// Minio
	config.MinioService.Endpoint = getEnv("MINIO_SERVICE_ENDPOINT", "https://minio.dennic.uz")
	config.MinioService.Bucket.Department = getEnv("MINIO_SERVICE_BUCKET_DEPARTMENT", "department")
	config.MinioService.Bucket.Doctor = getEnv("MINIO_SERVICE_BUCKET_DOCTOR", "doctor")
	config.MinioService.Bucket.Reasons = getEnv("MINIO_SERVICE_BUCKET_REASONS", "reasons")
	config.MinioService.Bucket.Specialization = getEnv("MINIO_SERVICE_BUCKET_SPECIALIZATION", "specialization")

	return &config
}

func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
