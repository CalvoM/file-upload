package main

import (
	"github.com/CalvoM/file-upload/pkg/auth"
	"github.com/CalvoM/file-upload/pkg/server"
	"log"
)

func main() {
	defer auth.CleanUp()
	srv := server.GetNewServer()
	log.Fatal(srv.ListenAndServe())
}
