package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AWS      AWSConfig
	JWT      JWTConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	S3Bucket        string
	S3Endpoint      string
	SQSQueueURL     string
}

type JWTConfig struct {
	Secret                string
	ExpiresIn             time.Duration
	RefreshTokenExpiresIn time.Duration
}

type UploadConfig struct {
	Path          string
	MaxUploadSize int64
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	refreshTokenExpiresIn, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRES_IN", "72h"))
	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "ecommerce_shop"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		AWS: AWSConfig{
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", "test"),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", "test"),
			Region:          getEnv("AWS_REGION", "us-east-1"),
			S3Bucket:        getEnv("AWS_S3_BUCKET", "ecommerce-uploads"),
			S3Endpoint:      getEnv("AWS_S3_ENDPOINT", "http://localhost:9000"),
			SQSQueueURL:     getEnv("AWS_SQS_QUEUE_URL", "http://localhost:4566/000000000000/ecommerce-shop"),
		},
		JWT: JWTConfig{
			Secret:                getEnv("JWT_SECRET", "secret"),
			ExpiresIn:             jwtExpiresIn,
			RefreshTokenExpiresIn: refreshTokenExpiresIn,
		},
		Upload: UploadConfig{
			Path:          getEnv("UPLOAD_PATH", "./uploads"),
			MaxUploadSize: maxUploadSize,
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
