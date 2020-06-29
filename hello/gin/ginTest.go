package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		for k, v := range c.Request.Header {
			fmt.Println(k, v)
		}
		c.Header("test", "jw")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/download", func(c *gin.Context) {
		c.Header("content-disposition", "attachment; filename=jw.txt")
		c.JSON(200, gin.H{
			"message": "jw",
		})
	})
	r.POST("/download", func(c *gin.Context) {
		id := c.Query("postid")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("aa")
		message := c.PostForm("bb")

		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
		c.Header("content-disposition", "attachment; filename=jw.txt")
		c.JSON(200, gin.H{
			"message": "jw",
		})
	})

	r.GET("/forward", func(c *gin.Context) {
		response, err := http.Get("https://www.7-zip.org/a/7z1900.exe")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="test.exe"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})
	r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
