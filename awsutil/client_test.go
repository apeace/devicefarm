package awsutil

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/ride/devicefarm/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

// see client_mock_test.go for MockClient implementation
func mockClient() (*DeviceFarm, *MockClient) {
	mock := &MockClient{}
	client := &DeviceFarm{mock, util.NilLogger, nil, false}
	return client, mock
}

// MockHandler implements http.Handler. It is used to verify
// the payload of PUT requests
type MockHandler struct {
	t               *testing.T
	expectedPayload string
}

func (h *MockHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	assert := assert.New(h.t)
	assert.Equal(req.Method, http.MethodPut)
	body, err := ioutil.ReadAll(req.Body)
	assert.Nil(err)
	assert.Equal(h.expectedPayload, string(body))
	res.WriteHeader(http.StatusCreated)
}

// mockS3 creates a mock S3 HTTP server which only accepts PUT
// requests, and verifies the payload of the PUT requests
func mockS3(t *testing.T, expectedPayload string) (url string, ln net.Listener, err error) {
	assert := assert.New(t)
	// listen on random port
	ln, err = net.Listen("tcp", "localhost:0")
	assert.Nil(err)
	url = "http://" + ln.Addr().String() + "/"
	server := &http.Server{
		Handler:        &MockHandler{t, expectedPayload},
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		server.Serve(ln)
	}()
	return
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)
	creds := credentials.NewStaticCredentials("foo", "bar", "baz")
	log := util.NilLogger
	client := NewClient(creds, log)
	assert.NotNil(client)
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

	// should fail due to error
	mock.enqueue(nil, fmt.Errorf("fake error"))
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

	// should fail due to error
	mock.enqueue(nil, fmt.Errorf("fake error"))
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

	// should fail due to error
	mock.enqueue(nil, fmt.Errorf("fake error"))
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

	// should fail due to error
	mock.enqueue(nil, fmt.Errorf("fake error"))
	pool, err = client.UpdateDevicePool(pool, []string{"foo"})
	assert.NotNil(err)
	assert.Nil(pool)
}

func TestDevicePoolMatches(t *testing.T) {
	assert := assert.New(t)
	client, _ := mockClient()

	pool := &devicefarm.DevicePool{
		Rules: []*devicefarm.Rule{
			{
				Attribute: aws.String("ARN"),
				Operator:  aws.String("IN"),
				Value:     aws.String("[\"foo\"]"),
			},
		},
	}

	// should match
	result := client.DevicePoolMatches(pool, []string{"foo"})
	assert.True(result)

	// should not match
	result = client.DevicePoolMatches(pool, []string{"foo", "bar"})
	assert.False(result)

	pool = &devicefarm.DevicePool{
		Rules: []*devicefarm.Rule{
			{
				Attribute: aws.String("FOO"),
				Operator:  aws.String("BAR"),
				Value:     aws.String("[\"foo\"]"),
			},
		},
	}

	// should not match
	result = client.DevicePoolMatches(pool, []string{"foo"})
	assert.False(result)
}

func TestUploadToS3(t *testing.T) {
	assert := assert.New(t)
	client, _ := mockClient()

	url, ln, err := mockS3(t, "foo")
	defer ln.Close()

	// should succeed
	err = client.UploadToS3(url, strings.NewReader("foo"))
	assert.Nil(err)

	// should fail because 'fakeurl' does not exist
	err = client.UploadToS3("fakeurl", nil)
	assert.NotNil(err)
}

func TestCreateUpload(t *testing.T) {
	assert := assert.New(t)
	client, mock := mockClient()

	// create temporary directory with file "foo.txt"
	// (containing the string "foo")
	tmpDir, err := ioutil.TempDir("", "devicefarm")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)
	filename := filepath.Join(tmpDir, "foo.txt")

	// create mock S3 server
	url, ln, err := mockS3(t, "Foo\n")
	defer ln.Close()

	// mock upload output
	output := &devicefarm.CreateUploadOutput{
		Upload: &devicefarm.Upload{
			Arn: aws.String("uploadArn"),
			Url: aws.String(url),
		},
	}

	// should fail because foo.txt does not exist
	mock.enqueue(output, nil)
	_, err = client.CreateUpload("projectArn", filename, "uploadType", "name")
	assert.NotNil(err)

	// should succeed
	util.CopyFile("testdata/foo.txt", filename)
	mock.enqueue(output, nil)
	uploadArn, err := client.CreateUpload("projectArn", filename, "uploadType", "name")
	assert.Nil(err)
	assert.Equal("uploadArn", uploadArn)

	// should fail due to error
	mock.enqueue(nil, fmt.Errorf("fake error"))
	_, err = client.CreateUpload("projectArn", filename, "uploadType", "name")
	assert.NotNil(err)
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

	// should fail due to error
	mock.enqueue(nil, fmt.Errorf("Fake error"))
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
