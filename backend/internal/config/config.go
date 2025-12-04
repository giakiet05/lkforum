package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// AppConfig holds the application's configuration
type AppConfig struct {
	Port                 string
	MongoURI             string
	DBName               string
	JWTSecret            string
	JWTIssuer            string
	JWTAudience          string
	TokenTTL             int
	RefreshTokenTTL      int
	FrontendURL          string
	OTPExpirationMinutes int
	SMTP                 SMTPConfig
	Redis                RedisConfig
	Google               GoogleConfig
	Cloudinary           CloudinaryConfig
	Gemini               GeminiConfig
}

// SMTPConfig holds the email server configuration
type SMTPConfig struct {
	Host       string
	Port       int
	User       string
	Pass       string
	SenderName string
}

// RedisConfig holds the Redis server configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// GoogleConfig holds the Google OAuth2 configuration
type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// CloudinaryConfig holds the Cloudinary configuration
type CloudinaryConfig struct {
	CloudName    string
	APIKey       string
	APISecret    string
	UploadFolder string
	UploadPreset string
}

// GeminiConfig holds the Gemini AI configuration
type GeminiConfig struct {
	APIKey              string
	Model               string
	Enabled             bool
	ConfidenceThreshold float64
	Timeout             int
	MaxRetries          int
}

// Cfg is a global variable holding the application's configuration
var Cfg AppConfig

// LoadConfig loads environment variables from .env file and populates the Cfg struct
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found. Using environment variables.")
	}

	//Port
	Cfg.Port = getEnv("PORT", "8080")

	// Database & App
	Cfg.MongoURI = getEnv("MONGO_URI", "mongodb://localhost:27017")
	Cfg.DBName = getEnv("DB_NAME", "lkforum")
	Cfg.FrontendURL = getEnv("FRONTEND_URL", "http://localhost:5173")

	// JWT
	Cfg.JWTSecret = getEnv("JWT_SECRET", "your-secret-key")
	Cfg.JWTIssuer = getEnv("JWT_ISSUER", "lkforum")
	Cfg.JWTAudience = getEnv("JWT_AUDIENCE", "lkforum-users")
	Cfg.TokenTTL = getEnvInt("TOKEN_TTL_MINUTES", 60)
	Cfg.RefreshTokenTTL = getEnvInt("REFRESH_TOKEN_TTL_HOURS", 72)

	// Features
	Cfg.OTPExpirationMinutes = getEnvInt("OTP_EXPIRATION_MINUTES", 15)

	// Services
	Cfg.SMTP.Host = getEnv("SMTP_HOST", "smtp.example.com")
	Cfg.SMTP.Port = getEnvInt("SMTP_PORT", 587)
	Cfg.SMTP.User = getEnv("SMTP_USER", "")
	Cfg.SMTP.Pass = getEnv("SMTP_PASS", "")
	Cfg.SMTP.SenderName = getEnv("SMTP_SENDER_NAME", "LKForum")

	Cfg.Redis.Addr = getEnv("REDIS_ADDR", "localhost:6379")
	Cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	Cfg.Redis.DB = getEnvInt("REDIS_DB", 0)

	Cfg.Google.ClientID = getEnv("GOOGLE_CLIENT_ID", "")
	Cfg.Google.ClientSecret = getEnv("GOOGLE_CLIENT_SECRET", "")
	Cfg.Google.RedirectURL = getEnv("GOOGLE_REDIRECT_URL", "")

	Cfg.Cloudinary.CloudName = getEnv("CLOUDINARY_CLOUD_NAME", "")
	Cfg.Cloudinary.APIKey = getEnv("CLOUDINARY_API_KEY", "")
	Cfg.Cloudinary.APISecret = getEnv("CLOUDINARY_API_SECRET", "")
	Cfg.Cloudinary.UploadFolder = getEnv("CLOUDINARY_FOLDER", "lkforum")
	Cfg.Cloudinary.UploadPreset = getEnv("CLOUDINARY_UPLOAD_PRESET", "lkforum_preset")

	Cfg.Gemini.APIKey = getEnv("GEMINI_API_KEY", "")
	Cfg.Gemini.Model = getEnv("GEMINI_MODEL", "gemini-2.0-flash-lite")
	Cfg.Gemini.Enabled = getEnv("GEMINI_ENABLED", "true") == "true"
	Cfg.Gemini.ConfidenceThreshold = getEnvFloat("GEMINI_CONFIDENCE_THRESHOLD", 0.7)
	Cfg.Gemini.Timeout = getEnvInt("GEMINI_TIMEOUT", 15)
	Cfg.Gemini.MaxRetries = getEnvInt("GEMINI_MAX_RETRIES", 3)

	log.Println("Configuration loaded successfully")
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// Helper function to get float environment variable with a default value
func getEnvFloat(key string, defaultValue float64) float64 {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
			return value
		}
	}
	return defaultValue
}
