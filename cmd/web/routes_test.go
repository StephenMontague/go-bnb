package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/stephenmontague/go-bnb/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing; test passed

	default:
		t.Errorf("type is not chi.Mux, but is %t", v)
	}
}
