package main

import (
	"fmt"
	"net/http"
	"sentinel-filter/service"
)

func main() {
	bf := service.NewBloomFilter(10000, 0.01)
	http.Handle("/", bf)
	fmt.Println("Listening on port 7070...")
	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		panic(err)
	}
}
