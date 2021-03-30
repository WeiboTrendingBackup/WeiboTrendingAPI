package main

import (
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello.world"))
}

func main() {
	http.HandleFunc("/hello", hello)
	//err := http.ListenAndServe(":8080", nil)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
