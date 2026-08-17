package config

import (
	"os"
	"strconv"
	"time"
)

// Config 应用配置：全部来自环境变量
type Config struct {
	AppEnv        string
	HTTPPort      string
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPassword    string
	RedisAddr     string
	RedisPassword string
	JWTSecret     string
	JWTTTL        time.Duration
	UploadDir     string
	AdminEmail    string
	AdminPassword string
	DebugCodeMode bool // 演示模式：发送验证码接口返回 debug_code
}

// Load 从环境变量加载配置（含默认值）
func Load() *Config {
	return &Config{
		AppEnv:        getEnv("APP_ENV", "dev"),
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "127.0.0.1"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBName:        getEnv("DB_NAME", "techblog_db"),
		DBUser:        getEnv("DB_USER", "techblog_user"),
		DBPassword:    getEnv("DB_PASSWORD", "techblog_pwd"),
		RedisAddr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", "change_me_to_a_long_random_string"),
		JWTTTL:        time.Duration(getEnvInt("JWT_TTL_HOURS", 72)) * time.Hour,
		UploadDir:     getEnv("UPLOAD_DIR", "./data/uploads"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@techblog.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123456"),
		DebugCodeMode: getEnv("DEBUG_CODE_MODE", "true") == "true",
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
