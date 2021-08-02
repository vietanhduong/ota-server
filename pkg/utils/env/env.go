package env

import (
	"github.com/vietanhduong/ota-server/pkg/logger"
	"os"
	"strconv"
)

// GetEnvAsStringOrFallback returns the env variable for the given key
// and falls back to the given defaultValue if not set
func GetEnvAsStringOrFallback(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// GetEnvAsIntOrFallback returns the env variable (parsed as integer) for
// the given key and falls back to the given defaultValue if not set
func GetEnvAsIntOrFallback(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		value, err := strconv.Atoi(v)
		if err != nil {
			logger.Logger.Warnf("Parse env as int failed with error: %v", err)
			return defaultValue
		}
		return value
	}
	return defaultValue
}
