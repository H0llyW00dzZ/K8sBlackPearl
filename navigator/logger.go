package navigator

import (
	"fmt"
	"sync"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// tryLog attempts to log a message if the rate limiter allows it.
func tryLog(logFunc func(string, ...zap.Field), emoji string, context string, fields ...zap.Field) {
	if logLimiter.Allow() {
		logFunc(emoji+" "+context, fields...)
	}
}

// logMessage logs a message with the provided log function.
func logMessage(logFunc func(string, ...zap.Field), message string, fields ...zap.Field) {
	logFunc(message, fields...)
}

// LogWithEmoji logs a message with a given level, emoji, context, and fields.
// It checks if the Logger is not nil before logging to prevent panics.
// The rateLimited flag determines whether the log should be rate limited.
func LogWithEmoji(level zapcore.Level, emoji string, context string, rateLimited bool, fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()

	if Logger == nil {
		fmt.Printf(language.ErrorLoggerIsNotSet, context)
		return
	}

	if rateLimited && !logLimiter.Allow() {
		return
	}

	message := emoji + " " + context
	logByLevel(level, message, fields...)
}

// logByLevel logs the message by the appropriate level.
func logByLevel(level zapcore.Level, message string, fields ...zap.Field) {
	switch level {
	case zapcore.InfoLevel:
		logMessage(Logger.Info, message, fields...)
	case zapcore.ErrorLevel:
		// Note: Temporarily, errors are logged at the info level for testing purposes.
		// This is to ensure visibility during the development phase where the global logger
		// is shared across multiple tasks and workers. Each worker and their respective tasks
		// are synchronized to use this logger without conflicts.
		logMessage(Logger.Info, message, fields...)
	default:
		// Output an error message if an unsupported log level is encountered.
		// The 'Unsupportedloglevel' variable should be defined in the 'language' package
		// and contain an appropriate error message template.
		fmt.Printf(language.Unsupportedloglevel, level)
	}
}

// LogInfoWithEmojiRateLimited logs an informational message with rate limiting.
func LogInfoWithEmojiRateLimited(emoji string, context string, fields ...zap.Field) {
	LogWithEmoji(zapcore.InfoLevel, emoji, context, true, fields...)
}

// LogErrorWithEmojiRateLimited logs an error message with rate limiting.
func LogErrorWithEmojiRateLimited(emoji string, context string, fields ...zap.Field) {
	LogWithEmoji(zapcore.ErrorLevel, emoji, context, true, fields...)
}

// SetLogger sets the logger instance for the package in a thread-safe manner.
func SetLogger(logger *zap.Logger) {
	mu.Lock()
	Logger = logger
	mu.Unlock()
}

// LogInfoWithEmoji logs an informational message with a given emoji, context, and fields.
// It checks if the Logger is not nil before logging to prevent panics.
func LogInfoWithEmoji(emoji string, context string, fields ...zap.Field) {
	LogWithEmoji(zapcore.InfoLevel, emoji, context, false, fields...)
}

// logErrorWithEmoji logs an error message with a given emoji, context, and fields.
// It checks if the Logger is not nil before logging to prevent panics.
func LogErrorWithEmoji(emoji string, context string, fields ...zap.Field) {
	LogWithEmoji(zapcore.ErrorLevel, emoji, context, false, fields...)
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
func CreateLogFields(sailing string, shipsnamespace string, fieldOpts ...LogFieldOption) []zap.Field {
	fields := []zap.Field{
		zap.String("sailing", sailing),
		zap.String("shipsnamespace", shipsnamespace),
	}
	for _, opt := range fieldOpts {
		fields = append(fields, opt())
	}
	return fields
}
