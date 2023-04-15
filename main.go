package main

import (
	"fmt"
	"net/http"
	"sentinel-filter/httpserver"
	"sentinel-filter/service"
)

func main() {
	bf := service.NewBloomFilter(10000, 0.01)
	s := httpserver.NewServer(bf)
	http.Handle("/", s)
	fmt.Println("Listening on port 7070...")
	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		panic(err)
	}
}
