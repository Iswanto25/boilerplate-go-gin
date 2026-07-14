package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	AppEnv           string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	JWTSecret        string
	JWTRefreshSecret string
	JWTTTL           int // dalam jam, default 24 (access token)
	JWTRefreshTTL    int // dalam jam, default 168 (7 hari)
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	RedisDB          int
	BCryptRounds     int
	Salt             string
	AllowedOrigins   string
}

func (c *Config) validate() error {
	if c.AppEnv == "production" {
		if c.JWTSecret == "supersecretkey" || c.JWTSecret == "" {
			return fmt.Errorf("JWT_SECRET must be set in production")
		}
		if c.JWTRefreshSecret == "supersecretkey-refresh" || c.JWTRefreshSecret == "" {
			return fmt.Errorf("JWT_REFRESH_SECRET must be set in production")
		}
		if c.DBPassword == "postgres" || c.DBPassword == "" {
			return fmt.Errorf("DB_PASSWORD must be set in production")
		}
	}
	return nil
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

	bcryptRounds, err := strconv.Atoi(getEnv("ROUND", "5"))
	if err != nil || bcryptRounds < 4 || bcryptRounds > 31 {
		bcryptRounds = 5
	}

	cfg := &Config{
		AppPort:          getEnv("PORT", "8080"),
		AppEnv:           getEnv("APP_ENV", "development"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "boilerplate"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		JWTSecret:        getEnv("JWT_SECRET", "supersecretkey"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "supersecretkey-refresh"),
		JWTTTL:           jwtTTL,
		JWTRefreshTTL:    jwtRefreshTTL,
		RedisHost:        getEnv("REDIS_HOST", ""),
		RedisPort:        getEnv("REDIS_PORT", ""),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          redisDB,
		BCryptRounds:     bcryptRounds,
		Salt:             getEnv("SALT", ""),
		AllowedOrigins:   getEnv("ALLOWED_ORIGINS", "*"),
	}

	if err := cfg.validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
