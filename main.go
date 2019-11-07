package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"

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

	blk := make([]byte, size)
	_, err := randGen.Read(blk)

	data := bytes.NewReader(blk)
	dataLength := int64(data.Len())
	_ = dataLength

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.Data(200, "application/octet-stream", blk)
	}

	return
}

func main() {
	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
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

	fmt.Println(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}




