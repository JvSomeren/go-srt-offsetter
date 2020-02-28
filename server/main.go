package main

import "log"

func main() {
	srv := CreateServer()

	log.Fatal(srv.ListenAndServe())
}
