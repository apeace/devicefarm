package awsutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/aws/aws-sdk-go/service/devicefarm/devicefarmiface"
	"github.com/ride/devicefarm/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DeviceFarm struct {
	Client          devicefarmiface.DeviceFarmAPI
	Log             util.Logger
	allDevicesCache DeviceList
	initialized     bool
}

func NewClient(creds *credentials.Credentials, log util.Logger) *DeviceFarm {
	sess := session.New(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: creds,
	})
	client := devicefarm.New(sess)
	return &DeviceFarm{client, log, nil, false}
}

func (df *DeviceFarm) SearchDevices(search string, androidOnly bool, iosOnly bool) (devices DeviceList, err error) {
	params := &devicefarm.ListDevicesInput{}
	r, err := df.Client.ListDevices(params)
	if err != nil {
		return
	}
	allDevices := DeviceList(r.Devices)
	allDevices.Sort()
	search = strings.ToLower(search)
	doSearch := len(search) > 0
	for _, device := range allDevices {
		deviceName := *device.Name
		if doSearch && !strings.Contains(strings.ToLower(deviceName), search) {
			continue
		}
		if androidOnly && *device.Platform != devicefarm.DevicePlatformAndroid {
			continue
		}
		if iosOnly && *device.Platform != devicefarm.DevicePlatformIos {
			continue
		}
		devices = append(devices, device)
	}
	return
}

func (df *DeviceFarm) ListDevicePools(projectArn string) ([]*devicefarm.DevicePool, error) {
	params := &devicefarm.ListDevicePoolsInput{Arn: aws.String(projectArn)}
	r, err := df.Client.ListDevicePools(params)
	if err != nil {
		return nil, err
	}
	return r.DevicePools, nil
}

func (df *DeviceFarm) CreateDevicePool(projectArn string, name string, deviceArns []string) (*devicefarm.DevicePool, error) {
	// there will never be an error marshalling a simple slice of strings
	val, _ := json.Marshal(deviceArns)
	params := &devicefarm.CreateDevicePoolInput{
		ProjectArn: aws.String(projectArn),
		Name:       aws.String(name),
		Rules: []*devicefarm.Rule{
			{
				Attribute: aws.String("ARN"),
				Operator:  aws.String("IN"),
				Value:     aws.String(string(val)),
			},
		},
	}
	r, err := df.Client.CreateDevicePool(params)
	if err != nil {
		return nil, err
	}
	return r.DevicePool, nil
}

func (df *DeviceFarm) UpdateDevicePool(pool *devicefarm.DevicePool, deviceArns []string) (*devicefarm.DevicePool, error) {
	// there will never be an error marshalling a simple slice of strings
	val, _ := json.Marshal(deviceArns)
	params := &devicefarm.UpdateDevicePoolInput{
		Arn:  pool.Arn,
		Name: pool.Name,
		Rules: []*devicefarm.Rule{
			{
				Attribute: aws.String("ARN"),
				Operator:  aws.String("IN"),
				Value:     aws.String(string(val)),
			},
		},
	}
	r, err := df.Client.UpdateDevicePool(params)
	if err != nil {
		return nil, err
	}
	return r.DevicePool, nil
}

func (df *DeviceFarm) DevicePoolMatches(pool *devicefarm.DevicePool, deviceArns []string) bool {
	// there will never be an error marshalling a simple slice of strings
	val, _ := json.Marshal(deviceArns)
	for _, rule := range pool.Rules {
		if *rule.Attribute != "ARN" || *rule.Operator != "IN" {
			return false
		}
		if *rule.Value != string(val) {
			return false
		}
	}
	return true
}

func (df *DeviceFarm) UploadToS3(s3Url string, bytes io.ReadSeeker) (err error) {
	req, err := http.NewRequest("PUT", s3Url, bytes)
	if err != nil {
		// TODO: Not sure how to add test coverage for this line
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	return
}

// fileReaderSeeker reads a file into an io.ReadSeeker, see:
// https://github.com/aws/aws-sdk-go/issues/142
// https://github.com/aws/aws-sdk-go/issues/337
// TODO: Get rid of this crap if possible
func (df *DeviceFarm) fileReaderSeeker(filename string) (r io.ReadSeeker, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		// TODO: Not sure how to add test coverage for this line
		return
	}
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)
	file.Read(buffer)
	r = bytes.NewReader(buffer)
	return
}

func (df *DeviceFarm) CreateUpload(projectArn, filename, uploadType, name string) (uploadArn string, err error) {
	// create upload object, get signed S3 URL
	params := &devicefarm.CreateUploadInput{
		Name:        aws.String(name),
		ProjectArn:  aws.String(projectArn),
		Type:        aws.String(uploadType),
		ContentType: aws.String("application/octet-stream"),
	}
	rApp, err := df.Client.CreateUpload(params)
	if err != nil {
		return
	}
	signedUrl := *rApp.Upload.Url

	// get io.ReadSeeker from file
	r, err := df.fileReaderSeeker(filename)
	if err != nil {
		return
	}

	err = df.UploadToS3(signedUrl, r)
	uploadArn = *rApp.Upload.Arn

	return
}

func (df *DeviceFarm) UploadSucceeded(arn string) (bool, error) {
	params := &devicefarm.GetUploadInput{Arn: aws.String(arn)}
	r, err := df.Client.GetUpload(params)
	if err != nil {
		return false, err
	}
	if *r.Upload.Status == devicefarm.UploadStatusFailed {
		return false, fmt.Errorf("Upload failed: %v", arn)
	}
	if *r.Upload.Status == devicefarm.UploadStatusSucceeded {
		return true, nil
	}
	return false, nil
}

func (df *DeviceFarm) WaitForUploadsToSucceed(timeoutMs, delayMs int, arns ...string) error {
	errchan := make(chan error)
	quitchan := make(chan bool)
	go func() {
		var err error
		var succeeded bool
	mainloop:
		for len(arns) > 0 {
			select {
			case <-quitchan:
				break mainloop
			default:
			}
			nextArns := []string{}
			for _, arn := range arns {
				succeeded, err = df.UploadSucceeded(arn)
				if err != nil {
					break mainloop
				}
				if succeeded {
					continue
				}
				nextArns = append(nextArns, arn)
			}
			arns = nextArns
			if len(arns) > 0 && delayMs > 0 {
				time.Sleep(time.Duration(delayMs) * time.Millisecond)
			}
		}
		errchan <- err
	}()
	select {
	case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
		quitchan <- true
		return fmt.Errorf("Timed out")
	case err := <-errchan:
		return err
	}
}

func (df *DeviceFarm) CreateRun(projectArn, poolArn string, test map[string]string) (string, error) {
	testTypeName := test["type"]
	appFile := test["app"]
	appFileBase := filepath.Base(appFile)
	testFile := test["test"]
	testFileBase := filepath.Base(testFile)

	testType, ok := TEST_TYPES[testTypeName]
	if !ok {
		return "", fmt.Errorf("Unknown test type: %v", testTypeName)
	}

	log := df.Log
	log.Println(">> Uploading files...")

	log.Println(appFileBase)
	appArn, err := df.CreateUpload(projectArn, appFile, testType.AppUploadType, appFileBase)
	if err != nil {
		return "", err
	}

	log.Println(testFileBase)
	testArn, err := df.CreateUpload(projectArn, testFile, testType.TestUploadType, testFileBase)
	if err != nil {
		return "", err
	}

	log.Println(">> Waiting for files to be processed...")
	err = df.WaitForUploadsToSucceed(60000, 5000, appArn, testArn)
	if err != nil {
		return "", err
	}

	log.Println(">> Creating test run...")
	params := &devicefarm.ScheduleRunInput{
		DevicePoolArn: aws.String(poolArn),
		ProjectArn:    aws.String(projectArn),
		Test: &devicefarm.ScheduleRunTest{
			Type:           aws.String(testType.TestType),
			TestPackageArn: aws.String(testArn),
		},
		AppArn: aws.String(appArn),
	}
	r, err := df.Client.ScheduleRun(params)
	if err != nil {
		return "", err
	}
	return *r.Run.Arn, nil
}

type TestType struct {
	TestType       string
	AppUploadType  string
	TestUploadType string
}

var TEST_TYPES = map[string]TestType{
	"android_instrumentation": {
		TestType:       devicefarm.TestTypeInstrumentation,
		AppUploadType:  devicefarm.UploadTypeAndroidApp,
		TestUploadType: devicefarm.UploadTypeInstrumentationTestPackage,
	},
	"android_appium_python": {
		TestType:       devicefarm.TestTypeAppiumPython,
		AppUploadType:  devicefarm.UploadTypeAndroidApp,
		TestUploadType: devicefarm.UploadTypeAppiumPythonTestPackage,
	},
}
