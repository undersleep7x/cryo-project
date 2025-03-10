package main

import (
	"testing"

	"github.com/undersleep7x/cryptowallet-v0.1/app"
)

func TestStartServer(t *testing.T) {
	server := startServer()

	// Ensure the server is not nil
	if server == nil {
		t.Fatal("Expected server to be initialized, got nil")
	}

	// Check if the server address is correctly set
	expectedAddr := ":" + app.Config.App.Port
	if server.Addr != expectedAddr {
		t.Errorf("Expected server address to be %s, got %s", expectedAddr, server.Addr)
	}

	// Ensure the handler is set
	if server.Handler == nil {
		t.Fatal("Expected server handler to be initialized, got nil")
	}
}