package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
	m.Handle("/dynamodb/read", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("read: creating session")
		sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})
		if err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}
		client := dynamodb.New(sess)
		tableName := os.Getenv("TABLE_NAME")
		log.Printf("tableName: %q", tableName)
		_, err = client.Query(&dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			KeyConditionExpression: aws.String("#pk = :pk"),
			ExpressionAttributeNames: map[string]*string{
				"#pk": aws.String("pk"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":pk": {S: aws.String("pk")},
			},
		})
		if err != nil {
			log.Printf("read error: %v", err)
			http.Error(w, "failed to read", http.StatusInternalServerError)
			return
		}
		io.WriteString(w, "<html><head><title>Read data</title></head><body><h1>Data read</h1></body><html>")
	}))
	m.Handle("/dynamodb/write", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("write: creating session")
		sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})
		if err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}
		client := dynamodb.New(sess)
		tableName := os.Getenv("TABLE_NAME")
		log.Printf("tableName: %q", tableName)
		_, err = client.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]*dynamodb.AttributeValue{
				"pk": {S: aws.String("pk")},
				"sk": {S: aws.String(time.Now().String())},
			},
		})
		if err != nil {
			log.Printf("write error: %v", err)
			http.Error(w, "failed to write", http.StatusInternalServerError)
			return
		}
		io.WriteString(w, "<html><head><title>Written data</title></head><body><h1>Data written</h1></body><html>")
	}))
	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Incoming request")
		io.WriteString(w, "<html><head><title>Hello</title></head><body><h1>World</h1></body><html>")
	}))
	http.ListenAndServe(":8000", m)
}
