package main

import (
	"fmt"
	"log"
	"strings"
)

type LogLevel int

const (
	LevelError LogLevel = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

var currentLogLevel = LevelInfo

func initLogger(levelStr string) {
	levelStr = strings.ToLower(strings.TrimSpace(levelStr))
	switch levelStr {
	case "error":
		currentLogLevel = LevelError
	case "warn", "warning":
		currentLogLevel = LevelWarn
	case "info", "":
		currentLogLevel = LevelInfo
	case "debug":
		currentLogLevel = LevelDebug
	default:
		currentLogLevel = LevelInfo
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func logError(format string, args ...any) {
	if currentLogLevel >= LevelError {
		log.Printf("[ERROR] "+format, args...)
	}
}

func logWarn(format string, args ...any) {
	if currentLogLevel >= LevelWarn {
		log.Printf("[WARN] "+format, args...)
	}
}

func logInfo(format string, args ...any) {
	if currentLogLevel >= LevelInfo {
		log.Printf("[INFO] "+format, args...)
	}
}

func logDebug(format string, args ...any) {
	if currentLogLevel >= LevelDebug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// small helper for user-facing fatal messages (in main)
func fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	logError("%s", msg)
}
