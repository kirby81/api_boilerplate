package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	// With no router
	srv, err := NewServer(nil)
	assert.Nil(t, srv)
	assert.Error(t, err)

	// Default
	handler := http.NewServeMux()
	srv, err = NewServer(handler)
	assert.NotNil(t, srv)
	assert.Nil(t, err)
}
