package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	AppEnv    string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPassword string
	DBName    string
	DBSSLMode string
	JWTSecret       string
	JWTRefreshSecret string
	JWTTTL          int // dalam jam, default 24 (access token)
	JWTRefreshTTL   int // dalam jam, default 168 (7 hari)
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	jwtTTL, err := strconv.Atoi(getEnv("JWT_TTL", "24"))
	if err != nil || jwtTTL <= 0 {
		jwtTTL = 24
	}

	jwtRefreshTTL, err := strconv.Atoi(getEnv("JWT_REFRESH_TTL", "168"))
	if err != nil || jwtRefreshTTL <= 0 {
		jwtRefreshTTL = 168
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil || redisDB < 0 {
		redisDB = 0
	}

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "boilerplate"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:         getEnv("JWT_SECRET", "supersecretkey"),
		JWTRefreshSecret:  getEnv("JWT_REFRESH_SECRET", "supersecretkey-refresh"),
		JWTTTL:            jwtTTL,
		JWTRefreshTTL:     jwtRefreshTTL,
		RedisHost:         getEnv("REDIS_HOST", ""),
		RedisPort:         getEnv("REDIS_PORT", ""),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           redisDB,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
