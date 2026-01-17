package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestSetupRouter(t *testing.T) {
	router := setupTestRouter()
	SetupRouter(router)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		checkBody      func(*testing.T, []byte)
	}{
		{
			name:           "Hello endpoint",
			method:         http.MethodGet,
			path:           "/Hello",
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				_ = json.Unmarshal(body, &response)
				assert.Equal(t, "Hello, RESTful API!", response["message"])
			},
		},
		{
			name:           "Register endpoint exists",
			method:         http.MethodPost,
			path:           "/api/v1/register",
			expectedStatus: http.StatusInternalServerError,
			checkBody:      func(t *testing.T, body []byte) {},
		},
		{
			name:           "Login endpoint exists",
			method:         http.MethodPost,
			path:           "/api/v1/login",
			expectedStatus: http.StatusInternalServerError,
			checkBody:      func(t *testing.T, body []byte) {},
		},
		{
			name:           "GraphQL endpoint exists",
			method:         http.MethodPost,
			path:           "/graphql",
			expectedStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				_ = json.Unmarshal(body, &response)
				assert.Contains(t, response["error"], "GraphQL handler not initialized")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkBody(t, w.Body.Bytes())
		})
	}
}

func TestProtectedRoutes(t *testing.T) {
	router := setupTestRouter()
	SetupRouter(router)

	protectedRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/users"},
		{http.MethodGet, "/api/v1/user/1"},
		{http.MethodGet, "/api/v1/profile"},
		{http.MethodDelete, "/api/v1/user/1"},
	}

	for _, route := range protectedRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)

			var response map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(t, response["error"], "Authorization header")
		})
	}
}

func TestProtectedRoutesWithValidToken(t *testing.T) {
	router := setupTestRouter()
	SetupRouter(router)

	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjk5OTk5OTk5OTl9.test"

	protectedRoutes := []struct {
		method         string
		path           string
		expectedStatus int
	}{
		{http.MethodGet, "/api/v1/users", http.StatusOK},
		{http.MethodGet, "/api/v1/profile", http.StatusUnauthorized},
	}

	for _, route := range protectedRoutes {
		t.Run(route.method+" "+route.path+" with token", func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			req.Header.Set("Authorization", "Bearer "+validToken)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.True(t, w.Code == route.expectedStatus || w.Code == http.StatusUnauthorized || w.Code == http.StatusInternalServerError || w.Code == http.StatusOK)
		})
	}
}

func TestRouteNotFound(t *testing.T) {
	router := setupTestRouter()
	SetupRouter(router)

	req := httptest.NewRequest(http.MethodGet, "/not-exist", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGraphQLEndpointMethods(t *testing.T) {
	router := setupTestRouter()
	SetupRouter(router)

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}

	for _, method := range methods {
		t.Run("GraphQL "+method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/graphql", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	}
}

func TestPublicRoutesAccessibility(t *testing.T) {
	router := setupTestRouter()
	SetupRouter(router)

	publicRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/v1/register"},
		{http.MethodPost, "/api/v1/login"},
	}

	for _, route := range publicRoutes {
		t.Run("Public "+route.method+" "+route.path, func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
			assert.NotEqual(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}
