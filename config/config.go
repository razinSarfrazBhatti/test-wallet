package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
	JWTConfig    JWTConfig
	EthConfig    EthConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type EthConfig struct {
	InfuraURL        string
	USDCContractAddr string
}

var AppConfig Config

func LoadConfig() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Database configuration
	AppConfig.DBConfig = DBConfig{
		Host:     getEnv("MYSQL_DB_HOST", "localhost"),
		Port:     getEnv("MYSQL_DB_PORT", "3306"),
		User:     getEnv("MYSQL_DB_USER", "root"),
		Password: getEnv("MYSQL_DB_PASS", ""),
		Name:     getEnv("MYSQL_DB_NAME", "test_wallet"),
	}

	// Server configuration
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(getEnv("SERVER_IDLE_TIMEOUT", "120"))

	AppConfig.ServerConfig = ServerConfig{
		Port:         getEnv("SERVER_PORT", "8080"),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}

	// JWT configuration
	jwtExpiration, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	AppConfig.JWTConfig = JWTConfig{
		Secret:     getEnv("JWT_SECRET", "your-secret-key"),
		Expiration: time.Duration(jwtExpiration) * time.Hour,
	}

	// Ethereum configuration
	AppConfig.EthConfig = EthConfig{
		InfuraURL:        getEnv("INFURA_URL", ""),
		USDCContractAddr: getEnv("USDC_CONTRACT_ADDRESS", ""),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
