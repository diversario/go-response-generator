package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	//"github.com/diversario/go-echo-server/server"
	"github.com/gin-gonic/gin"
	"github.com/zsais/go-gin-prometheus"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

func GenRandomBytes(size int, c *gin.Context) {
	seed := rand.Int()
	randSource := rand.NewSource(int64(seed))
	randGen := rand.New(randSource)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic: ", err)
			c.Abort()
		}
	}()

	c.DataFromReader(200, int64(size), "application/octet-stream", randGen, map[string]string{})

	return
}

func main() {
	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"ts": fmt.Sprintf("%s", time.Now().Format(time.RFC3339Nano)),
			"ip": c.ClientIP(),
		})
	})

	r.GET("/headers", func(c *gin.Context) {
		resp := gin.H{
			"ts": fmt.Sprintf("%s", time.Now().Format(time.RFC3339Nano)),
			"ip": c.ClientIP(),
			"headers": c.Request.Header.Clone(),
		}

		c.JSON(200, resp)
	})

	r.GET("/response", func(c *gin.Context) {
		size := c.DefaultQuery("bytes", "10K")
		var sizeBytes int

		quantity, err := strconv.Atoi(size[:len(size)-1])

		if err != nil {
			c.JSON(200, gin.H{
				"error": err.Error(),
			})

			return
		}

		unit := size[len(size)-1]

		switch unit {
		case 'K':
			sizeBytes = quantity * KB
		case 'M':
			sizeBytes = quantity * MB
		case 'G':
			sizeBytes = quantity * GB
		default:
			sizeBytes = quantity
		}

		GenRandomBytes(sizeBytes, c)
	})

	r.GET("/env", func(c *gin.Context) {
		name := c.DefaultQuery("name", "")

		val := os.Getenv(name)

		c.JSON(200, gin.H{
			name: val,
			"ip": c.ClientIP(),
			"ts": fmt.Sprintf("%s", time.Now().Format(time.RFC3339Nano)),
		})
	})

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//go func() {
	//	server.Run()
	//}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}




