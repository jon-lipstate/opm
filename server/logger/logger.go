package logger

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	MainLogger     *log.Logger
	SecurityLogger *log.Logger // For 404s, suspicious activity, auth failures
)

// InitLoggers initializes the logging system
func InitLoggers() {
	var mainLogOutput, securityLogOutput io.Writer
	env := os.Getenv("ENV")

	if env == "development" {
		fmt.Println("📝 Sending logs to stdout")
		mainLogOutput = os.Stdout
		securityLogOutput = os.Stdout
	} else {
		// Create logs directory if it doesn't exist
		if err := os.MkdirAll("logs", 0755); err != nil {
			log.Fatalf("Failed to create logs directory: %v", err)
		}

		mainLogOutput = &lumberjack.Logger{
			Filename:   "logs/main.log",
			MaxSize:    25,   // Max megabytes before rotation
			MaxBackups: 10,   // Max number of old files to keep
			MaxAge:     90,   // Max days to retain old files
			Compress:   true, // Compress old files
		}

		securityLogOutput = &lumberjack.Logger{
			Filename:   "logs/security.log",
			MaxSize:    10,   // Smaller since it's mostly suspicious activity
			MaxBackups: 5,    // Keep fewer backups
			MaxAge:     30,   // Shorter retention
			Compress:   true, // Compress since it could be high volume
		}
	}

	MainLogger = log.New(mainLogOutput, "", log.Ldate|log.Ltime|log.Lshortfile)
	SecurityLogger = log.New(securityLogOutput, "", log.Ldate|log.Ltime)
}
