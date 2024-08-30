package scaffold

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

func GetEnv(name, defaultValue string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return defaultValue
}
func NewUUID() string {
	return uuid.NewString()
}
func GetEnvBool(name string, defaultValue bool) bool {
	if value, ok := os.LookupEnv(name); ok {
		value = strings.ToLower(value)
		return value == "true" || value == "yes" || value == "1"
	}
	return defaultValue
}
