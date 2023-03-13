package main

import (
	"dbs/routers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("welcome to mongoAPI")
	r := routers.Router()
	fmt.Println("server started")
	// http.ListenAndServe(":4000", r)

	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("listing now...")
}
