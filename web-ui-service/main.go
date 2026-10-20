package main

import (
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./public"))

	log.Println("Web UI http://localhost:3000 adresinde başlatılıyor")

	err := http.ListenAndServe(":3000", fileServer)
	if err != nil {
		log.Fatal(err)
	}
}