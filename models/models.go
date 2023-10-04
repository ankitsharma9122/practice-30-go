package models

import (
	"context"
	"my-ankit-practice/controller"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var mu sync.Mutex

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
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	skip := (page - 1) * limit
	mu.Lock()
	defer mu.Unlock()
	var data []Blog
	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var blog Blog
		cursor.Decode(context.Background())
		data = append(data, blog)
	}
	c.JSON(http.StatusOK, gin.H{"message": data})
}

func CraetePost(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

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
	mu.Lock()
	defer mu.Unlock()
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
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": rst})
}

func DelectePost(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	id := c.Param("id")
	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ID}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog post deleted successfully"})
}
