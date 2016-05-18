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

func mockClient() (*DeviceFarm, *MockClient) {
	mock := &MockClient{}
	client := &DeviceFarm{mock, util.NilLogger, nil, false}
	return client, mock
}

func TestSearchDevices(t *testing.T) {
	assert := assert.New(t)
	client, mock := mockClient()

	// enqueue mock output with two devices
	output := &devicefarm.ListDevicesOutput{}
	output.Devices = []*devicefarm.Device{androidDevice, iosDevice}
	mock.enqueue(output, nil)

	// blank search should return both devices, sorted
	result, err := client.SearchDevices("", false, false)
	assert.Nil(err)
	assert.Equal(DeviceList{iosDevice, androidDevice}, result)

	// search should only return the iphone
	mock.enqueue(output, nil)
	result, err = client.SearchDevices("iphone", false, false)
	assert.Nil(err)
	assert.Equal(DeviceList{iosDevice}, result)

	// android filter should only return the android phone
	mock.enqueue(output, nil)
	result, err = client.SearchDevices("", true, false)
	assert.Nil(err)
	assert.Equal(DeviceList{androidDevice}, result)

	// ios filter should only return the iphone
	mock.enqueue(output, nil)
	result, err = client.SearchDevices("", false, true)
	assert.Nil(err)
	assert.Equal(DeviceList{iosDevice}, result)

	// should fail due to mock error
	mock.enqueue(nil, errors.New("fake error"))
	result, err = client.SearchDevices("", false, false)
	assert.NotNil(err)
	assert.Nil(result)
}

func TestListDevicePools(t *testing.T) {
	assert := assert.New(t)
	client, mock := mockClient()

	// should succeed
	output := &devicefarm.ListDevicePoolsOutput{
		DevicePools: []*devicefarm.DevicePool{
			&devicefarm.DevicePool{
				Arn:         aws.String("foo"),
				Description: aws.String("foo"),
				Name:        aws.String("foo"),
			},
		},
	}
	mock.enqueue(output, nil)
	pools, err := client.ListDevicePools("foo")
	assert.Nil(err)
	assert.Equal(output.DevicePools, pools)

	// should fail
	mock.enqueue(nil, errors.New("fake error"))
	pools, err = client.ListDevicePools("foo")
	assert.NotNil(err)
	assert.Nil(pools)
}

func TestCreateDevicePool(t *testing.T) {
	assert := assert.New(t)
	client, mock := mockClient()

	// enqueue mock device pool output
	output := &devicefarm.CreateDevicePoolOutput{
		DevicePool: &devicefarm.DevicePool{
			Arn: aws.String("poolarn"),
		},
	}
	mock.enqueue(output, nil)

	// should succeed and return device pool
	pool, err := client.CreateDevicePool("arn", "name", []string{"foo"})
	assert.Nil(err)
	assert.Equal(*output.DevicePool, *pool)

	// check input given to mock.CreateDevicePool()
	expectedInput := devicefarm.CreateDevicePoolInput{
		ProjectArn: aws.String("arn"),
		Name:       aws.String("name"),
		Rules: []*devicefarm.Rule{
			{
				Attribute: aws.String("ARN"),
				Operator:  aws.String("IN"),
				Value:     aws.String("[\"foo\"]"),
			},
		},
	}
	actualInput := (mock.Inputs()[0][0]).(*devicefarm.CreateDevicePoolInput)
	assert.Equal(expectedInput, *actualInput)

	// should fail
	mock.enqueue(nil, errors.New("fake error"))
	pool, err = client.CreateDevicePool("arn", "name", []string{"foo"})
	assert.NotNil(err)
	assert.Nil(pool)
}

func TestUpdateDevicePool(t *testing.T) {
	assert := assert.New(t)
	client, mock := mockClient()

	// enqueue mock device pool output
	pool := &devicefarm.DevicePool{
		Arn:  aws.String("poolarn"),
		Name: aws.String("poolname"),
	}
	output := &devicefarm.UpdateDevicePoolOutput{
		DevicePool: pool,
	}
	mock.enqueue(output, nil)

	// should succeed and return device pool
	updatedPool, err := client.UpdateDevicePool(pool, []string{"foo"})
	assert.Nil(err)
	assert.Equal(*pool, *updatedPool)

	// check input given to mock.UpdateDevicePool()
	expectedInput := devicefarm.UpdateDevicePoolInput{
		Arn:  aws.String("poolarn"),
		Name: aws.String("poolname"),
		Rules: []*devicefarm.Rule{
			{
				Attribute: aws.String("ARN"),
				Operator:  aws.String("IN"),
				Value:     aws.String("[\"foo\"]"),
			},
		},
	}
	actualInput := (mock.Inputs()[0][0]).(*devicefarm.UpdateDevicePoolInput)
	assert.Equal(expectedInput, *actualInput)

	// should fail
	mock.enqueue(nil, errors.New("fake error"))
	pool, err = client.UpdateDevicePool(pool, []string{"foo"})
	assert.NotNil(err)
	assert.Nil(pool)
}

func TestWaitForUploadsToSucceed(t *testing.T) {
	assert := assert.New(t)
	client, mock := mockClient()

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
