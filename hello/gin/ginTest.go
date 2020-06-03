package main

import (
	"fmt"

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
	r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
