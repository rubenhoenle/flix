package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllExercises(t *testing.T) {
	router := setupRouter()

	assert.Equal(t, 0, len(exercises))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/exercises", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[]", w.Body.String())

	// TODO: add new exercises to the slice and test the updated endpoint response
}

func TestDeleteExercise(t *testing.T) {
	router := setupRouter()

	// assert the slice is empty before
	assert.Equal(t, 0, len(exercises))

	// create a new exercise and add it to the slice
	var e exercise
	e.UUID = uuid.New()
	e.Title = "delete-test-exercise"
	exercises = append(exercises, e)
	assert.Equal(t, 1, len(exercises))

	// delete the exercise via the API
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/exercise/%s", e.UUID.String()), nil)
	router.ServeHTTP(w, req)

	// make sure the API responds correctly
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"OK\"}", w.Body.String())

	// assert the slice is empty again afterwards
	assert.Equal(t, 0, len(exercises))
}
