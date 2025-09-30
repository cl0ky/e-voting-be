package env

import (
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	AllowOrigins []string
	AllowMethods []string
	GinMode      string
	Port         uint64
	DBHost       string
	DBPort       uint64
	DBUser       string
	DBPassword   string
	DBName       string
	DBString     string
)

func GetEnv() {
	AllowOrigins = strings.Split(os.Getenv("ALLOW_ORIGINS"), ",")
	AllowMethods = strings.Split(os.Getenv("ALLOW_METHODS"), ",")
	GinMode = os.Getenv("GIN_MODE")
	Port = parseToUint(os.Getenv("PORT"), 8003)

	// Database Configuration
	DBHost = os.Getenv("PG_HOST")
	DBPort = parseToUint(os.Getenv("PG_PORT"), 5432)
	DBUser = os.Getenv("PG_USER")
	DBPassword = os.Getenv("PG_PASSWORD")
	DBName = os.Getenv("PG_DATABASE")
	DBString = os.Getenv("DB_STRING")
}

func parseToUint(val string, def ...uint64) uint64 {
	if val == "" {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}

	parsed, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		log.Fatal("failed to parse string to uint:", err)
	}
	return parsed
}
