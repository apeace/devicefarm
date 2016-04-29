package awsutil

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/stretchr/testify/assert"
	"testing"
)

// fake devices
var androidDevice *devicefarm.Device = &devicefarm.Device{
	Name:     aws.String("Samsung Galaxy S3"),
	Platform: aws.String(devicefarm.DevicePlatformAndroid),
	Arn:      aws.String("arn123"),
}
var iosDevice *devicefarm.Device = &devicefarm.Device{
	Name:     aws.String("Apple iPhone 6S"),
	Platform: aws.String(devicefarm.DevicePlatformIos),
	Arn:      aws.String("arn456"),
}

func mockClient(t *testing.T) *DeviceFarm {
	// mock client and ListDevicesOutput
	mock := &MockClient{}
	client := &DeviceFarm{mock, nil, false}
	output := &devicefarm.ListDevicesOutput{}

	// add both devices and enqueue mock output
	output.Devices = []*devicefarm.Device{androidDevice, iosDevice}
	mock.enqueue(output, nil)

	// init should succeed
	err := client.Init()
	assert.Nil(t, err)

	return client
}

func TestInit(t *testing.T) {
	assert := assert.New(t)

	// test error case
	mock := &MockClient{}
	client := &DeviceFarm{mock, nil, false}
	mock.enqueue(nil, errors.New("Fake error"))
	err := client.Init()
	assert.NotNil(err)

	// test success case
	client = mockClient(t)

	// init should succeed another time
	err = client.Init()
	assert.Nil(err)
}

func TestDevicesLookup(t *testing.T) {
	assert := assert.New(t)
	client := mockClient(t)
	lookup := client.DevicesLookup()
	assert.Equal(len(lookup), 4)
	assert.Equal(lookup["Samsung Galaxy S3"], androidDevice)
	assert.Equal(lookup["arn123"], androidDevice)
	assert.Equal(lookup["Apple iPhone 6S"], iosDevice)
	assert.Equal(lookup["arn456"], iosDevice)
}

func TestSearchDevices(t *testing.T) {
	assert := assert.New(t)
	client := mockClient(t)

	// blank search should return both devices, sorted
	result := client.SearchDevices("", false, false)
	assert.Equal(DeviceList{iosDevice, androidDevice}, result)

	// search should only return the iphone
	result = client.SearchDevices("iphone", false, false)
	assert.Equal(DeviceList{iosDevice}, result)

	// android filter should only return the android phone
	result = client.SearchDevices("", true, false)
	assert.Equal(DeviceList{androidDevice}, result)

	// ios filter should only return the iphone
	result = client.SearchDevices("", false, true)
	assert.Equal(DeviceList{iosDevice}, result)
}
