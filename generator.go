package generator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

// GenerateSignedURL generates signed URL for Google Cloud Storage object
func GenerateSignedURL(bucketName, fileName string, expiresInSecond int) string {
	if expiresInSecond == 0 {
		expiresInSecond = 300 // default expiration is 5 minutes
	}
	expires := time.Now().Add(time.Duration(expiresInSecond) * time.Second)

	opts := getSignedURLOptions(expires)
	url, err := storage.SignedURL(bucketName, fileName, &opts)
	if err != nil {
		log.Fatalf("error signing URL: %v", err)
	}

	return url
}

func getSignedURLOptions(expires time.Time) storage.SignedURLOptions {
	serviceAccountKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if serviceAccountKey == "" {
		log.Fatalf("missing service account key path. Forgot to set GOOGLE_APPLICATION_CREDENTIALS environment variable?")
	}

	data, err := ioutil.ReadFile(serviceAccountKey)
	if err != nil {
		log.Fatalf("error reading service account key file: %v", err)
	}

	var credentials map[string]string
	json.Unmarshal(data, &credentials)

	return storage.SignedURLOptions{
		GoogleAccessID: credentials["client_email"],
		PrivateKey:     []byte(credentials["private_key"]),
		Method:         "GET",
		Expires:        expires,
	}
}
