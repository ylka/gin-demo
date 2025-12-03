package main

import (
	"context"
	"fmt"
	"gin-demo/recipes-api/handlers"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesHandler *handlers.RecipesHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:Glb%4012345@10.100.1.217:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	collection := client.Database("demo1203").Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "10.100.1.217:6379",
		Password: "8URXPL2x3HZMi7xoGTdk3Upc",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	r := gin.Default()

	r.POST("/recipes", recipesHandler.NewRecipeHandler)
	r.GET("/recipes", recipesHandler.ListRecipesHandler)
	r.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	r.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	// r.GET("/recipes/search", SearchRecipesHandler)

	r.Run()
}
