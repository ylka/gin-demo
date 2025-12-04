package main

import (
	"context"
	"fmt"
	"gin-demo/recipes-api/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/joho/godotenv"
)

var recipesHandler *handlers.RecipesHandler

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO")))
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
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PWD"),
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	r := gin.Default()

	r.GET("/recipes", recipesHandler.ListRecipesHandler)

	authorized := r.Group("/")
	authorized.Use(handlers.AuthMiddleware())

	authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
	authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	// authorized.GET("/recipes/search", SearchRecipesHandler)

	r.Run()
}
