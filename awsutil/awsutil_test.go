package awsutil

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/devicefarm"
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

func TestListDevices(t *testing.T) {
	assert := assert.New(t)

	// mock client that returns one device
	mock := &MockClient{}
	mockDevice1 := &devicefarm.Device{Name: aws.String("Test device")}
	mockOutput := &devicefarm.ListDevicesOutput{Devices: []*devicefarm.Device{mockDevice1}}
	mock.enqueue(mockOutput, nil)
	client := &DeviceFarm{mock}

	// blank search should return all devices
	result, err := client.ListDevices("")
	assert.Nil(err)
	assert.Equal([]string{"Test device"}, result)

	// add another device to the list
	mockDevice2 := &devicefarm.Device{Name: aws.String("Another device")}
	mockOutput.Devices = append(mockOutput.Devices, mockDevice2)
	mock.enqueue(mockOutput, nil)

	// search should only return the second device
	result, err = client.ListDevices("another")
	assert.Nil(err)
	assert.Equal([]string{"Another device"}, result)

	// enqueue an error
	mock.enqueue(nil, errors.New("Fake error"))
	result, err = client.ListDevices("")
	assert.NotNil(err)
}
