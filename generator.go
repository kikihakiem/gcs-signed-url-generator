package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

// GenerateSignedURL generates signed URL for Google Cloud Storage object
func GenerateSignedURL(bucketName, fileName string, expiresInSecond int) (string, error) {
	if expiresInSecond == 0 {
		expiresInSecond = 300 // default expiration is 5 minutes
	}
	expires := time.Now().Add(time.Duration(expiresInSecond) * time.Second)

	opts, err := getSignedURLOptions(expires)
	if err != nil {
		return "", err
	}

	url, err := storage.SignedURL(bucketName, fileName, &opts)
	if err != nil {
		return "", fmt.Errorf("error signing URL: %v", err)
	}

	return url, nil
}

func getSignedURLOptions(expires time.Time) (storage.SignedURLOptions, error) {
	var opts storage.SignedURLOptions
	serviceAccountKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if serviceAccountKey == "" {
		return opts, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
	}

	data, err := ioutil.ReadFile(serviceAccountKey)
	if err != nil {
		return opts, fmt.Errorf("error reading service account key file: %v", err)
	}

	var credentials map[string]string
	json.Unmarshal(data, &credentials)

	opts = storage.SignedURLOptions{
		GoogleAccessID: credentials["client_email"],
		PrivateKey:     []byte(credentials["private_key"]),
		Method:         "GET",
		Expires:        expires,
	}

	return opts, nil
}
