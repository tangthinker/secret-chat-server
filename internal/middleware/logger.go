package middleware

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2"
)

func LoggerInFile(filepath string) fiber.Handler {
	loggingFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return logger.New(logger.Config{
		Format:     "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     loggingFile,
	})
}

func LoggerInConsole() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} ${method} ${path} - ${header:X-Forwarded-For} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     os.Stdout,
	})
}
