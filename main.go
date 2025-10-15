package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Create Fiber instance with custom config
	app := fiber.New(fiber.Config{
		AppName:          "Fiber Docker App v1.0.0",
		ServerHeader:     "Fiber",
		ReadBufferSize:   16384,           // 16KB - เพิ่มขนาด buffer สำหรับ request headers
		WriteBufferSize:  16384,           // 16KB - เพิ่มขนาด buffer สำหรับ response
		BodyLimit:        4 * 1024 * 1024, // 4MB - จำกัดขนาด request body
		DisableKeepalive: false,           // เปิดใช้ keep-alive
	})

	// Middleware
	app.Use(logger.New())  // Logging middleware
	app.Use(recover.New()) // Recover from panics

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World 🌈",
			"status":  "success",
			"app":     "Fiber Docker App",
		})
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Pong 🏓",
			"status":  "success",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"uptime": "running",
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// Channel to listen for interrupt signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine so we can listen for shutdown signals
	go func() {
		log.Printf("🚀 Server is starting on port %s...", port)
		log.Printf("🌐 Visit: http://localhost:%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("❌ Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-quit
	log.Println("\n🛑 Shutting down server...")

	// Graceful shutdown
	if err := app.Shutdown(); err != nil {
		log.Printf("❌ Error during shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}
