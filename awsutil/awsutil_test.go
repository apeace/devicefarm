package awsutil

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCredsFromEnv(t *testing.T) {
	assert := assert.New(t)

	previousKey := os.Getenv(ENV_ACCESS_KEY)
	previousSecret := os.Getenv(ENV_SECRET)
	defer func() {
		os.Setenv(ENV_ACCESS_KEY, previousKey)
		os.Setenv(ENV_SECRET, previousSecret)
	}()

	// should fail when env vars are unset
	os.Setenv(ENV_ACCESS_KEY, "")
	os.Setenv(ENV_SECRET, "")
	ok, _ := CredsFromEnv()
	assert.False(ok)

	// should succeed when env vars are set
	os.Setenv(ENV_ACCESS_KEY, "access-key")
	os.Setenv(ENV_SECRET, "secret")
	ok, creds := CredsFromEnv()
	assert.True(ok)
	assert.Equal(*credentials.NewStaticCredentials("access-key", "secret", ""), *creds)
}

func TestCredsFromFile(t *testing.T) {
	assert := assert.New(t)

	// should fail because file doesn't exist
	ok, _ := CredsFromFile("./testdata/does-not-exist.json")
	assert.False(ok)

	// should fail because file is not valid JSON
	ok, _ = CredsFromFile("./testdata/creds-invalid.json")
	assert.False(ok)

	// should succeed when given valid file
	ok, creds := CredsFromFile("./testdata/creds.json")
	assert.True(ok)
	assert.Equal(*credentials.NewStaticCredentials("access-key", "secret", ""), *creds)
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)
	creds := credentials.NewStaticCredentials("access-key", "secret", "")
	df := NewClient(creds)
	assert.NotNil(df)
}
