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
    router.POST("/exercise", createExercise)
    router.PUT("/exercise/:uuid", updateExercise)
    router.DELETE("/exercise/:uuid", deleteExercise)

	router.Run("localhost:8080")
}

// seed the backend with some exercises for demo
var exercises = []exercise{
	{UUID: uuid.New(), Title: "Exercise 1"},
	{UUID: uuid.New(), Title: "Exercise 2"},
	{UUID: uuid.New(), Title: "Exercise 3"},
}

// getExercises responds with the list of all exercises as JSON
func getExercises(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, exercises)
}

func updateExercise(context *gin.Context) {
    uuidStr := context.Param("uuid")
    if err := uuid.Validate(uuidStr); err != nil {
        context.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid uuid"})
        return
    }
    var updateExercise exercise
    if err := context.BindJSON(&updateExercise); err != nil {
        return
    }
    parsedUuid, err := uuid.Parse(uuidStr)
    if(err != nil) {
        context.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid uuid"})
        return
    }
   
    // make sure the exercise uuid cant be changed
    updateExercise.UUID = parsedUuid
    
    for idx, exercise := range exercises {
        if exercise.UUID.String() == uuidStr {
            exercises[idx] = updateExercise
            context.IndentedJSON(http.StatusOK, exercise)
            return
        }
    }
    context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
}

func getExercise(context *gin.Context) {
    uuid := context.Param("uuid")
    for _, exercise := range exercises {
        if exercise.UUID.String() == uuid {
            context.IndentedJSON(http.StatusOK, exercise)
            return
        }
    }
    context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
}

func createExercise(context *gin.Context) {
    var newExercise exercise
    if err := context.BindJSON(&newExercise); err != nil {
        return
    }

    // generate and set new uuid
    newUuid, err := uuid.NewUUID()
    if err != nil {
        return
    }
    newExercise.UUID = newUuid

    exercises = append(exercises, newExercise)
    context.IndentedJSON(http.StatusCreated, newExercise)
}

/* remove the given index from the slice */
func remove(s []exercise, i int) []exercise {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}

func deleteExercise(context *gin.Context) {
    uuid := context.Param("uuid")
    for idx, exercise := range exercises {
        if exercise.UUID.String() == uuid {
            exercises = remove(exercises, idx)
            context.JSON(http.StatusOK, gin.H{"message": "OK"})
            return
        }
    }
    context.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
}

