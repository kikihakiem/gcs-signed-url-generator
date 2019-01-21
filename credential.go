package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var (
	gcloudCredential = &googleCloudCredential{}

	initCredential sync.Once
	initErr        error
)

type googleCloudCredential struct {
	GoogleAccessID string `json:"client_email"`
	PrivateKey     string `json:"private_key"`
}

func getGoogleCloudCredential() (*googleCloudCredential, error) {
	initCredential.Do(func() {
		serviceAccountKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if serviceAccountKey == "" {
			initErr = fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
			return
		}

		data, err := ioutil.ReadFile(serviceAccountKey)
		if err != nil {
			initErr = fmt.Errorf("error reading service account key file: %v", err)
			return
		}
		err = json.Unmarshal(data, gcloudCredential)
		if err != nil {
			initErr = fmt.Errorf("error parsing service account credentials: %v", err)
		}

	})

	return gcloudCredential, initErr
}
