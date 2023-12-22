package navigator

import (
	"fmt"
	"sync"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"go.uber.org/zap"
)

// LogFieldOption is a type that represents a function which returns a zap.Field.
// It's used to pass custom field generators to functions that create a slice of zap.Fields for logging.
// This allows for deferred creation of zap.Fields, which can be useful when the value of the field
// is not immediately available or when you want to include custom logic to determine the value of the field.
//
// Example usage:
//
//	fields := navigator.CreateLogFields("sailing", "the high seas", navigator.WithAnyZapField(zap.Int("treasure chests", 5)))
//	navigator.LogInfoWithEmoji("🏴‍☠️", "Pirates ahoy!", fields...)
type LogFieldOption func() zap.Field

// Logger is a package-level variable to access the zap logger throughout the worker package.
var Logger *zap.Logger

// mu is used to protect access to the Logger variable to make it safe for concurrent use.
var mu sync.Mutex

// SetLogger sets the logger instance for the package in a thread-safe manner.
func SetLogger(logger *zap.Logger) {
	mu.Lock()
	Logger = logger
	mu.Unlock()
}

// LogInfoWithEmoji logs an informational message with a given emoji, context, and fields.
// It checks if the Logger is not nil before logging to prevent panics.
func LogInfoWithEmoji(emoji string, context string, fields ...zap.Field) {
	mu.Lock()
	logger := Logger
	mu.Unlock()

	if logger == nil {
		fmt.Printf(language.ErrorLoggerIsNotSet, context)
		return
	}
	logger.Info(emoji+" "+context, fields...)
}

// logErrorWithEmoji logs an error message with a given emoji, context, and fields.
// It checks if the Logger is not nil before logging to prevent panics.
func LogErrorWithEmoji(emoji string, context string, fields ...zap.Field) {
	mu.Lock()
	logger := Logger
	mu.Unlock()

	if logger == nil {
		fmt.Printf(language.ErrorLoggerIsNotSet, context)
		return
	}
	logger.Info(emoji+" "+context, fields...)
}

// WithAnyZapField creates a LogFieldOption that encapsulates a zap.Field for deferred addition to a log entry.
// This function is particularly handy when you have a custom field to add to your log that isn't
// already covered by existing "With*" functions. It allows for a more flexible and dynamic approach
// to logging, akin to how a pirate might prefer the freedom to navigate the open seas.
//
// Example:
//
// Let's say we want to log an event related to a pirate's treasure map. We have a custom binary field
// that represents the map, and we want to include this in our log fields. We can use WithAnyZapField
// to add this custom field to our logs as follows:
//
//	treasureMap := []byte{0x0A, 0x0B, 0x0C, 0x0D} // This represents our treasure map in binary form.
//	fields := navigator.CreateLogFields("sailing", "find treasure", navigator.WithAnyZapField(zap.Binary("treasureMap", treasureMap)))
//	navigator.LogInfoWithEmoji("🏴‍☠️", "Found a treasure map!", fields...)
func WithAnyZapField(field zap.Field) LogFieldOption {
	return func() zap.Field {
		return field
	}
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
