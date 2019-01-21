package main

import (
	"log"
	"os"

	generator "github.com/kikihakiem/gcs-signed-url-generator"
)

func main() {
	var (
		bucketName = os.Args[1]
		fileName   = os.Args[2]
	)

	signedURL := generator.GenerateSignedURL(bucketName, fileName, 0)
	log.Println("Signed URL:", signedURL)
}
