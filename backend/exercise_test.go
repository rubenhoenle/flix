package main

import (
	"fmt"
    "bytes"
    "strings"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGetAllExercises(t *testing.T) {
	router := setupRouter()

	// assert the slice is empty before
    exercises = []exercise{}
	assert.Equal(t, 0, len(exercises))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/exercises", nil)
	router.ServeHTTP(w, req)

    // assert the endpoint returns no exercises
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[]", w.Body.String())

    // add 5 exercises to the slice
	for i := 1; i <= 5; i++ {
		var e exercise
		e.UUID = uuid.New()
		e.Title = fmt.Sprintf("Exercise %d", i)
		exercises = append(exercises, e)
	}
    
    // assert the slice contains 5 exercises
    assert.Equal(t, 5, len(exercises))
	
    w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/exercises", nil)
	router.ServeHTTP(w, req)

    var sb strings.Builder
    sb.WriteString("[")
    for idx, e := range exercises {
        sb.WriteString(fmt.Sprintf(`{"uuid":"%s","title":"Exercise %d"}`, e.UUID.String(), idx + 1))
        if idx + 1 < len(exercises) {
            sb.WriteString(",")
        } 
    }
    sb.WriteString("]")

    // assert the endpoint returns all 5 exercises
    assert.Equal(t, 200, w.Code)
    require.JSONEq(t, sb.String(), w.Body.String())
}

func TestGetExercise(t *testing.T) {
	router := setupRouter()

	// assert the slice is empty before
    exercises = []exercise{}
	assert.Equal(t, 0, len(exercises))
		
	var e exercise
    e.UUID = uuid.New()
	e.Title = "exercise-get-test"
	exercises = append(exercises, e)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/exercise/%s", e.UUID.String()), nil)
	router.ServeHTTP(w, req)
    
    // assert the slice contains an exercise
    assert.Equal(t, 1, len(exercises))
	
    // assert the endpoint responds correctly
	assert.Equal(t, 200, w.Code)
    require.JSONEq(t, fmt.Sprintf(`{"uuid":"%s","title":"exercise-get-test"}`, exercises[0].UUID.String()), w.Body.String())
}

func TestDeleteExercise(t *testing.T) {
	router := setupRouter()

	// assert the slice is empty before
    exercises = []exercise{}
	assert.Equal(t, 0, len(exercises))

	// create a new exercise and add it to the slice
	var e exercise
	e.UUID = uuid.New()
	e.Title = "delete-test-exercise"
	exercises = append(exercises, e)
	assert.Equal(t, 1, len(exercises))

	// delete the exercise via the API
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/exercise/%s", e.UUID.String()), nil)
	router.ServeHTTP(w, req)

	// make sure the API responds correctly
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"OK\"}", w.Body.String())

	// assert the slice is empty again afterwards
	assert.Equal(t, 0, len(exercises))
}

func TestCreateExercise(t *testing.T) {
    router := setupRouter()

	// assert the slice is empty before
    exercises = []exercise{}
	assert.Equal(t, 0, len(exercises))

	// create an exercise via the API
	w := httptest.NewRecorder()
    body := []byte(`{"title": "exercise-create-test"}`)
	req, _ := http.NewRequest(http.MethodPost, "/exercise", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	// make sure the API responds correctly
	assert.Equal(t, 201, w.Code)
    require.JSONEq(t, fmt.Sprintf(`{"uuid":"%s","title":"exercise-create-test"}`, exercises[0].UUID.String()), w.Body.String())

	// assert the slice contains one exercise afterwards
	assert.Equal(t, 1, len(exercises))
}

func TestUpdateExercise(t *testing.T) {
    router := setupRouter()

	// assert the slice is empty before
    exercises = []exercise{}
	assert.Equal(t, 0, len(exercises))

    exerciseUuid := uuid.New().String()

    // create a exercise which will be updated later
    e := exercise{UUID: uuid.MustParse(exerciseUuid), Title: "exercise-update-test-old-title"}
    exercises = append(exercises, e)
	assert.Equal(t, 1, len(exercises))
	
    w := httptest.NewRecorder()
    // also pass a new uuid, to make sure it is not possible to change the uuid of the exercise via the update endpoint
    body := fmt.Sprintf(`{"title": "exercise-create-test-new-title", "uuid": "%s"}`, uuid.New().String())
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/exercise/%s", e.UUID.String()), strings.NewReader(body))
	router.ServeHTTP(w, req)
	
    // make sure the API responds correctly
	assert.Equal(t, 200, w.Code)
    //assert.Equal(t, fmt.Sprintf(`{"uuid":"%s","title":"exercise-create-test"}`, e.UUID.String(), e.Title), w.Body.String())

    assert.Equal(t, "exercise-create-test-new-title", exercises[0].Title)
    assert.Equal(t, exerciseUuid, exercises[0].UUID.String())

	// assert the slice contains still exactly one exercise afterwards
	assert.Equal(t, 1, len(exercises))
}

func TestGetUpdateDeleteNonExistentExercise(t *testing.T) {
    router := setupRouter()

	// assert the slice is empty before
    exercises = []exercise{}
	assert.Equal(t, 0, len(exercises))

    missingUuid := uuid.New().String()
    url := fmt.Sprintf("/exercise/%s", missingUuid)
   
    // Delete
    w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
    
    // Get
    w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, url, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
    
    // Update
    w = httptest.NewRecorder()
    body := `{"title": "update-non-existent-exercise-test"}`
	req, _ = http.NewRequest(http.MethodPut, url, strings.NewReader(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)

	// assert the slice is still empty afterwards
	assert.Equal(t, 0, len(exercises))
}

