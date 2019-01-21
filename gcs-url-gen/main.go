package main

import (
	"log"
	"os"

	generator "github.com/kikihakiem/gcs-signed-url-generator"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("missing bucket name or file name")
	}

	var (
		bucketName = args[1]
		fileName   = args[2]
	)

	signedURL, err := generator.GenerateSignedURL(bucketName, fileName, 0)
	if err != nil {
		log.Fatalf("could not generate signed URL: %v", err)
	}

	log.Println("signed URL:", signedURL)
}
