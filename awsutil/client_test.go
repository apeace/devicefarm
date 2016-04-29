package awsutil

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
