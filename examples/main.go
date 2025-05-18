package main

import (
	"bytes"
	"fmt"
	"log"
	"mdhesari/go-multipart-encoder"
	"net/http"
	"os"
)

// UploadRequest represents a file upload with metadata
type UploadRequest struct {
	Username    string `form:"username"`
	Email       string `form:"email"`
	Description string `form:"description"`
	File        []byte `form:"file"`
	Tags        []string
}

func main() {
	// Sample usage: upload a file with metadata
	fileData, err := os.ReadFile("examples/example.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	req := UploadRequest{
		Username:    "johndoe",
		Email:       "john@example.com",
		Description: "An example file upload",
		File:        fileData,
		Tags:        []string{"example", "upload", "golang"},
	}

	// Encode the struct as multipart form-data
	buf, contentType, err := multipart.Encode(req)
	if err != nil {
		log.Fatalf("Failed to encode multipart form: %v", err)
	}

	// Send HTTP request with the multipart form-data
	resp, err := http.Post("https://httpbin.org/post", contentType, buf)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and display the response
	respBuf := new(bytes.Buffer)
	respBuf.ReadFrom(resp.Body)
	fmt.Println("Response:", respBuf.String())
}
