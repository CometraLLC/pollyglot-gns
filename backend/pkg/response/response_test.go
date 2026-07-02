package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseJSON(t *testing.T) {
	rec := httptest.NewRecorder()

	ResponseJSON(rec, http.StatusCreated, JSON{
		Code:    http.StatusCreated,
		Message: "created",
		Data:    map[string]string{"id": "42"},
		Status:  true,
	})

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var body JSON
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, http.StatusCreated, body.Code)
	assert.Equal(t, "created", body.Message)
	assert.True(t, body.Status)
}

func TestResponseJSONWithNilData(t *testing.T) {
	rec := httptest.NewRecorder()

	ResponseJSON(rec, http.StatusNoContent, nil)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestResponseError(t *testing.T) {
	rec := httptest.NewRecorder()

	ResponseError(rec, http.StatusBadRequest, "name is required")

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var body map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "name is required", body["error"])
}
