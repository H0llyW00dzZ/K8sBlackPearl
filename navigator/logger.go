package navigator

import (
	"fmt"
	"sync"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"go.uber.org/zap"
)

// Logger is a package-level variable to access the zap logger throughout the worker package.
var Logger *zap.Logger

// mu is used to protect access to the Logger variable to make it safe for concurrent use.
var mu sync.Mutex

// SetLogger sets the logger instance for the package in a thread-safe manner.
func SetLogger(logger *zap.Logger) {
	mu.Lock()
	defer mu.Unlock()
	Logger = logger
}

// LogInfoWithEmoji logs an informational message with a given emoji, context, and fields.
// It checks if the Logger is not nil before logging to prevent panics.
func LogInfoWithEmoji(emoji string, context string, fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()

	if Logger == nil {
		fmt.Printf(language.ErrorLoggerIsNotSet, context)
		return
	}
	Logger.Info(emoji+" "+context, fields...)
}

// logErrorWithEmoji logs an error message with a given emoji, context, and fields.
// It checks if the Logger is not nil before logging to prevent panics.
func LogErrorWithEmoji(emoji string, context string, fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()

	if Logger == nil {
		fmt.Printf(language.ErrorLoggerIsNotSet, context)
		return
	}
	Logger.Error(emoji+" "+context, fields...)
}

// createLogFields creates a slice of zap.Field with the operation and additional info.
// It can be used to add structured context to logs.
func CreateLogFields(sailing string, shipsnamespace string, infos ...string) []zap.Field {
	fields := []zap.Field{
		zap.String("sailing", sailing),
		zap.String("shipsnamespace", shipsnamespace),
	}
	for i, info := range infos {
		fields = append(fields, zap.String(fmt.Sprintf("info%d", i+1), info))
	}
	return fields
}
