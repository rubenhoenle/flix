package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type exercise struct {
	UUID  uuid.UUID `json:"uuid"`
	Title string    `json:"title"`
}

func main() {
	router := gin.Default()
	router.GET("/exercises", getExercises)
    router.GET("/exercise/:uuid", getExercise)

	router.Run("localhost:8080")
}

// seed the backend with some exercises for demo
var exercises = []exercise{
	{UUID: uuid.New(), Title: "Exercise 1"},
	{UUID: uuid.New(), Title: "Exercise 2"},
	{UUID: uuid.New(), Title: "Exercise 3"},
}

// getExercises responds with the list of all exercises as JSON
func getExercises(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, exercises)
}

func getExercise(c *gin.Context) {
    uuid := c.Param("uuid")
    for _, exercise := range exercises {
        if exercise.UUID.String() == uuid {
            c.IndentedJSON(http.StatusOK, exercise)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "exercise not found"})
}

