package main

import (
	"my-ankit-practice/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/get-blog", models.GetBlog)
	r.POST("/create-blog", models.CraetePost)
	r.PUT("/update-blog/:id", models.UpdatePost)
	r.DELETE("/delete-blog/:id", models.DelectePost)
	r.Run(":8001")
}
