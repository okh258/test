package main

import (
	"log"
	"net/http"
	"testing"
)

func TestWebSocket(t *testing.T) {
	RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
