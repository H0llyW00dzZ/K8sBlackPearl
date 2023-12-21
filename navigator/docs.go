// Package navigator provides logging functionality to the K8sBlackPearl project.
// It utilizes the uber-go/zap logging library to provide structured, leveled, and
// optionally emoji-prefixed logging capabilities.
//
// The package exposes a package-level Logger variable that is protected by a mutex
// to ensure thread-safe access and modification. This Logger must be set using the
// SetLogger function before any of the logging functions are called. If the Logger
// is not set, the logging functions will output an error message to standard output.
//
// The LogInfoWithEmoji and LogErrorWithEmoji functions provide a way to log
// informational and error messages, respectively. They both require an emoji string
// and a context string as part of the message to add a visual identifier to the logs.
// These functions also accept a variadic number of zap.Field parameters to include
// structured data in the log messages.
//
// The CreateLogFields function is a helper that creates a slice of zap.Field from
// given information strings. This can be used to add consistent structured context
// to logs across different parts of the application.
//
// # Usage
//
// Before logging can occur, the Logger must be initialized and set using the
// SetLogger function. Once set, the LogInfoWithEmoji and LogErrorWithEmoji functions
// can be used to log messages with structured data.
//
// Example:
//
//	// Initialize the Logger (typically done at application start)
//	logger, _ := zap.NewProduction()
//	navigator.SetLogger(logger)
//
//	// Use the logging functions
//	navigator.LogInfoWithEmoji("🚀", "Application started")
//	navigator.LogErrorWithEmoji("❗️", "An error occurred", zap.Error(err))
//
//	// Create structured log fields
//	fields := navigator.CreateLogFields("navigation", "starry-sea", "additional info")
//	navigator.LogInfoWithEmoji("🧭", "Navigating the stars", fields...)
//
// Note:
// It is crucial that the SetLogger function is called before any logging functions
// to avoid nil pointer dereferences. If the Logger is nil when a logging function
// is called, a message indicating that the Logger is not set will be printed to
// standard output.
//
// The emoji strings used in the logging functions are purely for visual effect and
// can be omitted or replaced based on the logging preferences of the application.
//
// Copyright (c) 2023 H0llyW00dzZ
package navigator
