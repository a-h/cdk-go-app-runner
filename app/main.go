package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	m := http.NewServeMux()
	m.Handle("/panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("oh no")
	}))
	m.Handle("/exit", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		os.Exit(1)
	}))
	m.Handle("/fatal", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Fatal("fatal error")
	}))
	m.Handle("/use-all-memory", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const increment = 1024 * 1024 * 256
		var space []byte
		for {
			// Use 256MB RAM.
			space = append(space, make([]byte, increment)...)
			fmt.Printf("%dMB consumed\n", len(space)/1024/1024)
		}
	}))
	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Incoming request")
		io.WriteString(w, "<html><head><title>Hello</title></head><body><h1>World</h1></body><html>")
	}))
	http.ListenAndServe(":8000", m)
}
