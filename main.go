package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo client database connection, as global variable for simplicity
// Later : validate concurrency handling of the go mongo driver
// Later : validate concurrency safety of the mongo driver
var client *mongo.Client

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

	router := gin.Default()
	router.GET("/talents", getTalents)
	router.POST("/talent", postTalent)

	router.Run("0.0.0.0:1112")
}

func connectMongoDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://crewdb:Opsss112@cluster0.xicjy.mongodb.net/test?authSource=admin&replicaSet=atlas-uhjhzf-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// getTalents responds with the list of all talents as JSON.
func getTalents(c *gin.Context) {
	client = connectMongoDB()
	talentsCollection := client.Database("crew").Collection("talents")
	cur, err := talentsCollection.Find(context.TODO(), bson.D{})
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

// postTalent get a json body for a Talent and responds with the InsertedId or an error.
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
