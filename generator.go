package generator

import (
	"fmt"
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
	var signedURLOptions storage.SignedURLOptions
	credentials, err := getGoogleCloudCredential()
	if err != nil {
		return signedURLOptions, err
	}

	signedURLOptions = storage.SignedURLOptions{
		GoogleAccessID: credentials.GoogleAccessID,
		PrivateKey:     []byte(credentials.PrivateKey),
		Method:         "GET",
		Expires:        expires,
	}

	return signedURLOptions, nil
}
