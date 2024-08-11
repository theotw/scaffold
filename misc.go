package scaffold

import (
	"os"

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
