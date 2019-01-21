package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/storage"
)

var (
	signedURLOptions *storage.SignedURLOptions

	initOptions    sync.Once
	optionsInitErr error
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

	url, err := storage.SignedURL(bucketName, fileName, opts)
	if err != nil {
		return "", fmt.Errorf("error signing URL: %v", err)
	}

	return url, nil
}

func getSignedURLOptions(expires time.Time) (*storage.SignedURLOptions, error) {
	initOptions.Do(func() {
		if signedURLOptions == nil {
			serviceAccountKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
			if serviceAccountKey == "" {
				optionsInitErr = fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
				return
			}

			data, err := ioutil.ReadFile(serviceAccountKey)
			if err != nil {
				optionsInitErr = fmt.Errorf("error reading service account key file: %v", err)
				return
			}

			var credentials map[string]string
			err = json.Unmarshal(data, &credentials)
			if err != nil {
				optionsInitErr = fmt.Errorf("error parsing service account credentials: %v", err)
			}

			signedURLOptions = &storage.SignedURLOptions{
				GoogleAccessID: credentials["client_email"],
				PrivateKey:     []byte(credentials["private_key"]),
				Method:         "GET",
				Expires:        expires,
			}
		}
	})

	return signedURLOptions, optionsInitErr
}
