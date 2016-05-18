package awsutil

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/ride/devicefarm/util"
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

func mockClient(t *testing.T) (*DeviceFarm, *MockClient) {
	// mock client and ListDevicesOutput
	mock := &MockClient{}
	client := &DeviceFarm{mock, util.NilLogger, nil, false}
	output := &devicefarm.ListDevicesOutput{}

	// add both devices and enqueue mock output
	output.Devices = []*devicefarm.Device{androidDevice, iosDevice}
	mock.enqueue(output, nil)

	// init should succeed
	err := client.Init()
	assert.Nil(t, err)

	return client, mock
}

func TestInit(t *testing.T) {
	assert := assert.New(t)

	// test error case
	mock := &MockClient{}
	client := &DeviceFarm{mock, util.NilLogger, nil, false}
	mock.enqueue(nil, errors.New("Fake error"))
	err := client.Init()
	assert.NotNil(err)

	// test success case
	client, _ = mockClient(t)

	// init should succeed another time
	err = client.Init()
	assert.Nil(err)
}

func TestSearchDevices(t *testing.T) {
	assert := assert.New(t)
	client, _ := mockClient(t)

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

func TestWaitForUploadsToSucceed(t *testing.T) {
	assert := assert.New(t)
	client, mock := mockClient(t)

	// should succeed immediately
	output := &devicefarm.GetUploadOutput{
		Upload: &devicefarm.Upload{
			Arn:    aws.String("arn123"),
			Status: aws.String(devicefarm.UploadStatusSucceeded),
		},
	}
	mock.enqueue(output, nil)
	err := client.WaitForUploadsToSucceed(1000, 0, "arn123")
	assert.Nil(err)

	// should succeed on the third iteration
	output = &devicefarm.GetUploadOutput{
		Upload: &devicefarm.Upload{
			Arn:    aws.String("arn123"),
			Status: aws.String(devicefarm.UploadStatusInitialized),
		},
	}
	mock.enqueue(output, nil)
	mock.enqueue(output, nil)
	output = &devicefarm.GetUploadOutput{
		Upload: &devicefarm.Upload{
			Arn:    aws.String("arn123"),
			Status: aws.String(devicefarm.UploadStatusSucceeded),
		},
	}
	mock.enqueue(output, nil)
	err = client.WaitForUploadsToSucceed(1000, 0, "arn123")
	assert.Nil(err)

	// should fail because upload failed
	output = &devicefarm.GetUploadOutput{
		Upload: &devicefarm.Upload{
			Arn:    aws.String("arn123"),
			Status: aws.String(devicefarm.UploadStatusFailed),
		},
	}
	mock.enqueue(output, nil)
	err = client.WaitForUploadsToSucceed(1000, 0, "arn123")
	assert.NotNil(err)

	// should fail because of request error
	mock.enqueue(nil, errors.New("Fake error"))
	err = client.WaitForUploadsToSucceed(1000, 0, "arn123")
	assert.NotNil(err)

	// should fail because of timeout
	output = &devicefarm.GetUploadOutput{
		Upload: &devicefarm.Upload{
			Arn:    aws.String("arn123"),
			Status: aws.String(devicefarm.UploadStatusInitialized),
		},
	}
	mock.enqueue(output, nil)
	err = client.WaitForUploadsToSucceed(1, 2, "arn123")
	assert.NotNil(err)
}
