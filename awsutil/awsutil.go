package awsutil

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"io/ioutil"
	"os"
)

const ENV_ACCESS_KEY = "AWS_ACCESS_KEY_ID"
const ENV_SECRET = "AWS_SECRET_ACCESS_KEY"

func CredsFromEnv() (ok bool, creds *credentials.Credentials) {
	ok = false
	key := os.Getenv(ENV_ACCESS_KEY)
	secret := os.Getenv(ENV_SECRET)
	if len(key) > 0 && len(secret) > 0 {
		ok = true
		creds = credentials.NewStaticCredentials(key, secret, "")
	}
	return
}

type credsJson struct {
	AccessKey string `json:"AWS_ACCESS_KEY_ID"`
	Secret    string `json:"AWS_SECRET_ACCESS_KEY"`
}

func CredsFromFile(filename string) (ok bool, creds *credentials.Credentials) {
	ok = false
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	credsJson := &credsJson{}
	err = json.Unmarshal(bytes, credsJson)
	if err != nil {
		return
	}
	if len(credsJson.AccessKey) > 0 && len(credsJson.Secret) > 0 {
		ok = true
		creds = credentials.NewStaticCredentials(credsJson.AccessKey, credsJson.Secret, "")
	}
	return
}
