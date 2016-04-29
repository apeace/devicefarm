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

func TestSearchDevices(t *testing.T) {
	assert := assert.New(t)

	// fake devices
	androidDevice := &devicefarm.Device{
		Name:     aws.String("Samsung Galaxy S3"),
		Platform: aws.String(devicefarm.DevicePlatformAndroid),
	}
	iosDevice := &devicefarm.Device{
		Name:     aws.String("Apple iPhone 6S"),
		Platform: aws.String(devicefarm.DevicePlatformIos),
	}

	// mock client and ListDevicesOutput
	mock := &MockClient{}
	client := &DeviceFarm{mock, nil}
	output := &devicefarm.ListDevicesOutput{}

	// enqueue an error
	mock.enqueue(nil, errors.New("Fake error"))
	result, err := client.SearchDevices("", false, false)
	assert.NotNil(err)

	// add both devices and enqueue mock output
	output.Devices = []*devicefarm.Device{androidDevice, iosDevice}
	mock.enqueue(output, nil)

	// blank search should return both devices, sorted
	result, err = client.SearchDevices("", false, false)
	assert.Nil(err)
	assert.Equal(DeviceList{iosDevice, androidDevice}, result)

	// re-enqueue same response
	mock.enqueue(output, nil)

	// search should only return the iphone
	result, err = client.SearchDevices("iphone", false, false)
	assert.Nil(err)
	assert.Equal(DeviceList{iosDevice}, result)

	// re-enqueue same response
	mock.enqueue(output, nil)

	// android filter should only return the android phone
	result, err = client.SearchDevices("", true, false)
	assert.Nil(err)
	assert.Equal(DeviceList{androidDevice}, result)

	// re-enqueue same response
	mock.enqueue(output, nil)

	// ios filter should only return the iphone
	result, err = client.SearchDevices("", false, true)
	assert.Nil(err)
	assert.Equal(DeviceList{iosDevice}, result)
}
