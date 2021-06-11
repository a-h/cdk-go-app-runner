package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Incoming request")
		io.WriteString(w, "<html><head><title>Hello</title></head><body><h1>World</h1></body><html>")
	}))
}
