package main

import (
	"log"
	"os"
)

// GetEnv returns a set environment variable based on a passed key or
// if the variable isn't set the passed fallback
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	srv := CreateServer()

	log.Fatal(srv.ListenAndServe())
}
