package main

import (
	"img-compress-demo/backend/handlers"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const timeFormat = "2006-01-02 15:04:05"

func logf(format string, args ...interface{}) {
	log.Printf("["+timeFormat+"] "+format, append([]interface{}{time.Now()}, args...)...)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		logf("%s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
		logf("<-- %s %s completed in %v with status %d",
			c.Request.Method, c.Request.URL.Path, time.Since(startTime), c.Writer.Status())
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	router := gin.Default()
	router.Use(LoggingMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		AllowMethods: []string{"GET", "POST"},
	}))

	router.POST("/api/compress", handlers.CompressImage)
	router.GET("/api/health", func(c *gin.Context) {
		logf("Health check requested")
		c.JSON(200, gin.H{"status": "ok"})
	})

	logf("========================================")
	logf("Image Compression Server Starting")
	logf("Listening on http://localhost:8080")
	logf("CORS enabled for http://localhost:5173")
	logf("========================================")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("["+timeFormat+"] Failed to start server: %v", time.Now(), err)
	}
}
