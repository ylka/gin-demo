package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal(file, &recipes)
}

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)

}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, r := range recipes {
		if r.ID == id {
			recipe.ID = r.ID
			recipe.PublishedAt = time.Now()

			recipes[i] = recipe
			c.JSON(http.StatusOK, recipe)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	for i, r := range recipes {
		if r.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "recipe deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
}

func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for _, r := range recipes {
		for _, t := range r.Tags {
			if strings.EqualFold(t, tag) {
				listOfRecipes = append(listOfRecipes, r)
			}
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	r := gin.Default()

	r.POST("/recipes", NewRecipeHandler)
	r.GET("/recipes", ListRecipesHandler)
	r.PUT("/recipes/:id", UpdateRecipeHandler)
	r.DELETE("/recipes/:id", DeleteRecipeHandler)
	r.GET("/recipes/search", SearchRecipesHandler)

	r.Run()
}
