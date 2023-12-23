// Package navigator provides structured logging capabilities for the K8sBlackPearl project,
// leveraging the uber-go/zap library. It offers leveled logging with the option to prefix
// messages with emojis for enhanced visual distinction.
//
// A package-level Logger variable is available and protected by a mutex to ensure safe
// concurrent access. This Logger should be initialized and set using SetLogger prior to
// invoking any logging functions. Failing to set the Logger will result in error messages
// being printed to standard output instead of proper log entries.
//
// Logging functions such as LogInfoWithEmoji and LogErrorWithEmoji are available for
// recording informational and error messages, respectively. These functions enhance log
// messages with emojis and a context string for quick identification. They also support
// structured logging by accepting a variable number of zap.Field parameters.
//
// Helper function CreateLogFields is provided to generate a slice of zap.Field from
// specified strings, facilitating the inclusion of consistent structured context within
// logs throughout the application.
//
// # Usage
//
// The Logger must be initialized and set using SetLogger before any logging activity.
// Subsequently, LogInfoWithEmoji and LogErrorWithEmoji can be employed for logging
// messages with structured context.
//
// Example:
//
//	// Initialize the Logger (usually at the start of the application)
//	logger, _ := zap.NewProduction()
//	navigator.SetLogger(logger)
//
//	// Logging with the package functions
//	navigator.LogInfoWithEmoji("🚀", "Application started")
//	navigator.LogErrorWithEmoji("❗️", "An error occurred", zap.Error(err))
//
//	// Structuring log fields
//	fields := navigator.CreateLogFields("navigation", "starry-sea", zap.String("detail", "additional info"))
//	navigator.LogInfoWithEmoji("🧭", "Navigating the stars", fields...)
//
// Important Note:
// Ensure that SetLogger is invoked before any logging functions to prevent nil pointer
// dereferences. If the Logger is nil during a logging attempt, an error message will be
// printed to standard output indicating the absence of a set Logger.
//
// The use of emoji strings in logging functions is optional and for visual enhancement.
// These can be excluded or customized according to the application's logging standards.
//
// Copyright (c) 2023 H0llyW00dzZ
package navigator
