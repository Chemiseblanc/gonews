package main

import (
	"log"

	"github.com/Chemiseblanc/gonews/nntp"
)

func main() {
	srv, err := nntp.NewServer(":119", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.ListenAndServe())
}
