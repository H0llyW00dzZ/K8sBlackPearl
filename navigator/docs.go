// Package navigator provides structured logging capabilities tailored for the K8sBlackPearl project,
// utilizing the uber-go/zap library for high-performance, leveled logging. This package enhances
// log messages with emojis for visual distinction and supports structured logging with zap.Field parameters.
//
// A package-level Logger variable is available, which should be set using SetLogger before any logging
// functions are called. If the Logger is not set, logging functions will default to printing error messages
// to standard output to prevent the application from panicking due to a nil Logger.
//
// Functions such as LogInfoWithEmoji and LogErrorWithEmoji are provided for logging informational and error
// messages, respectively. These functions append emojis to the log messages for easier visual scanning in log
// output. They also accept a variable number of zap.Field parameters for structured context logging.
//
// The CreateLogFields helper function is available to generate a slice of zap.Field from specified strings,
// enabling consistent structured context in logs throughout the application.
//
// # Usage
//
// Before logging, the Logger must be initialized and set using SetLogger. After setting up, the logging
// functions such as LogInfoWithEmoji and LogErrorWithEmoji can be used for logging messages with structured
// context.
//
// Example:
//
//	// Initialize the Logger (usually at the start of the application)
//	logger, _ := zap.NewProduction()
//	navigator.SetLogger(logger)
//
//	// Logging with package functions
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
