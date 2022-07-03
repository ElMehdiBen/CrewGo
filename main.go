package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// talent represents data about a talent with crew.
type talent struct {
	ID        string   `json:"id" bson:"_id"`
	FirstName string   `json:"firstName" bson:"firstName"`
	LastName  string   `json:"lastName" bson:"lastName"`
	Picture   string   `json:"picture" bson:"picture"`
	Job       string   `json:"job" bson:"job"`
	Location  string   `json:"location" bson:"location"`
	LinkedIn  string   `json:"linkedin" bson:"linkedin"`
	Github    string   `json:"github" bson:"github"`
	Twitter   string   `json:"twitter" bson:"twitter"`
	Tags      []string `json:"tags" bson:"tags"`
	Stage     string   `json:"stage" bson:"stage"`
}

func main() {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://crewdb:Opsss112@cluster0.xicjy.mongodb.net/test?authSource=admin&replicaSet=atlas-uhjhzf-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	router := gin.Default()
	router.GET("/talents", getTalents)
	router.POST("/talent", postTalent)

	router.Run("localhost:1112")
}

// getAlbums responds with the list of all albums as JSON.
func getTalents(c *gin.Context) {
	talentsCollection := client.Database("crew").Collection("talents")
	cur, err := talentsCollection.Find(context.TODO(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var talents []talent
	for cur.Next(context.TODO()) {
		var talent talent
		err := cur.Decode(&talent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		talents = append(talents, talent)
	}
	if err := cur.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cur.Close(context.TODO())
	c.JSON(http.StatusOK, talents)
}

func postTalent(c *gin.Context) {
	talentsCollection := client.Database("crew").Collection("talents")
	var talent talent
	if err := c.ShouldBindJSON(&talent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := talentsCollection.InsertOne(context.TODO(), talent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
