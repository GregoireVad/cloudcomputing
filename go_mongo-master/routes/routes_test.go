package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cavdy-play/go_mongo/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performPost(r http.Handler, method, path string, buffer *bytes.Buffer, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func setup() *gin.Engine {
	config.Connect()
	router := gin.Default()
	Routes(router)
	return router
}

func TestRouteWelcome(t *testing.T) {

	message := ""
	r := setup()

	var response map[string]string
	w := performRequest(r, "GET", "/")

	json.Unmarshal([]byte(w.Body.String()), &response)
	message = response["message"]

	assert.Equal(t, "Welcome To API", message)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNotFound(t *testing.T) {

	message := ""
	r := setup()

	var response map[string]string
	w := performRequest(r, "GET", "/notfound")

	json.Unmarshal([]byte(w.Body.String()), &response)
	message = response["message"]

	assert.Equal(t, "Route Not Found", message)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTodos(t *testing.T) {

	message := ""

	r := setup()

	var response map[string]string
	w := performRequest(r, "GET", "/todos")

	json.Unmarshal([]byte(w.Body.String()), &response)
	message = response["message"]

	if w.Code == http.StatusOK {
		assert.Equal(t, "All todos", message)
	} else {
		assert.Equal(t, "Something went wrong", message)
	}

}

func TestTodo(t *testing.T) {

	message := ""

	r := setup()

	var response map[string]string

	var jsonStr = []byte(`{ "Title": "test", "Body" : "test2", "Completed": "Complete"}`)

	w := performPost(r, "POST", "/todo", bytes.NewBuffer(jsonStr), header{
		Key: "Content-Type", Value: "application/json",
	})

	json.Unmarshal([]byte(w.Body.String()), &response)
	message = response["message"]

	if w.Code == http.StatusCreated {
		assert.Equal(t, "Todo created successfully", message)
		assert.Equal(t, http.StatusCreated, w.Code)
	} else {

		assert.Equal(t, "Something went wrong", message)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	}

}

func TestSingleTodo(t *testing.T) {

	message := ""

	r := setup()

	var response map[string]string
	w := performRequest(r, "GET", "/todo/:todoId")

	json.Unmarshal([]byte(w.Body.String()), &response)
	message = response["message"]

	if w.Code == http.StatusOK {
		assert.Equal(t, "Single Todo", message)
	} else {
		assert.Equal(t, "Todo not found", message)
	}

}
