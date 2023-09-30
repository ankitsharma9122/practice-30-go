package models

import (
	"context"
	"my-ankit-practice/controller"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func init() {
	collection = controller.ConnectionDbInstance()
}

type Blog struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Writer    string             `json:"writer,omitempty" bson:"writer,omitempty"`
	Content   string             `json:"content,omitempty" bson:"content,omitempty"`
	Likes     int64              `json:"likes,omitempty" bson:"likes,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

func GetBlog(c *gin.Context) {
	var data []Blog
	cursor, _ := collection.Find(context.Background(), bson.M{})
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var blog Blog
		cursor.Decode(context.Background())
		data = append(data, blog)
	}
	c.JSON(http.StatusCreated, gin.H{"message": data})
}

func CraetePost(c *gin.Context) {
	var blog Blog
	err := c.ShouldBindJSON(&blog)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	blog.CreatedAt = time.Now()
	blog.Likes = 0
	result, err := collection.InsertOne(context.Background(), blog)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": result})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post Blog
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ID}
	update := bson.M{"$set": bson.M{"writer": post.Writer, "content": post.Content, "likes": post.Likes}}
	rst, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": rst})

}

func DelectePost(c *gin.Context) {
	id := c.Param("id")
	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ID}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog post deleted successfully"})
}
