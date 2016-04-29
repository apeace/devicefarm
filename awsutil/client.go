package awsutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/aws/aws-sdk-go/service/devicefarm/devicefarmiface"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type DeviceFarm struct {
	Client          devicefarmiface.DeviceFarmAPI
	allDevicesCache DeviceList
}

func NewClient(creds *credentials.Credentials) *DeviceFarm {
	sess := session.New(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: creds,
	})
	client := devicefarm.New(sess)
	return &DeviceFarm{client, nil}
}

func (df *DeviceFarm) ListDevicesCached() (DeviceList, error) {
	if df.allDevicesCache != nil {
		return df.allDevicesCache, nil
	}
	params := &devicefarm.ListDevicesInput{}
	r, err := df.Client.ListDevices(params)
	if err != nil {
		return nil, err
	}
	list := DeviceList(r.Devices)
	list.Sort()
	df.allDevicesCache = list
	return df.allDevicesCache, nil
}

func (df *DeviceFarm) DevicesLookup() (map[string]*devicefarm.Device, error) {
	allDevices, err := df.ListDevicesCached()
	if err != nil {
		return nil, err
	}
	lookup := map[string]*devicefarm.Device{}
	for _, device := range allDevices {
		lookup[*device.Name] = device
		lookup[*device.Arn] = device
	}
	return lookup, nil
}

func (df *DeviceFarm) SearchDevices(search string, androidOnly bool, iosOnly bool) (devices DeviceList, err error) {
	allDevices, err := df.ListDevicesCached()
	if err != nil {
		return
	}
	devices = DeviceList{}
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

func (df *DeviceFarm) ListDevicePools(arn string) ([]*devicefarm.DevicePool, error) {
	params := &devicefarm.ListDevicePoolsInput{Arn: aws.String(arn)}
	r, err := df.Client.ListDevicePools(params)
	if err != nil {
		return nil, err
	}
	return r.DevicePools, nil
}

func (df *DeviceFarm) NamesToArns(names []string) ([]string, error) {
	lookup, err := df.DevicesLookup()
	if err != nil {
		return nil, err
	}
	arns := []string{}
	for _, name := range names {
		device, ok := lookup[name]
		if !ok {
			return nil, errors.New("No such device: " + name)
		}
		arns = append(arns, *device.Arn)
	}
	return arns, nil
}

func (df *DeviceFarm) CreateDevicePool(arn string, name string, deviceNames []string) (*devicefarm.DevicePool, error) {
	arns, err := df.NamesToArns(deviceNames)
	if err != nil {
		return nil, err
	}
	val, err := json.Marshal(arns)
	if err != nil {
		return nil, err
	}
	params := &devicefarm.CreateDevicePoolInput{
		ProjectArn: aws.String(arn),
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

func (df *DeviceFarm) UpdateDevicePool(pool *devicefarm.DevicePool, deviceNames []string) (*devicefarm.DevicePool, error) {
	arns, err := df.NamesToArns(deviceNames)
	if err != nil {
		return nil, err
	}
	val, err := json.Marshal(arns)
	if err != nil {
		return nil, err
	}
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

func (df *DeviceFarm) DevicePoolMatches(pool *devicefarm.DevicePool, deviceNames []string) (bool, error) {
	arns, err := df.NamesToArns(deviceNames)
	if err != nil {
		return false, err
	}
	val, err := json.Marshal(arns)
	if err != nil {
		return false, err
	}
	for _, rule := range pool.Rules {
		if *rule.Attribute != "ARN" || *rule.Operator != "IN" {
			return false, nil
		}
		if *rule.Value != string(val) {
			return false, nil
		}
	}
	return true, nil
}

func (df *DeviceFarm) CreateUpload(projectArn, filename, uploadType, name string) (string, error) {
	// open file and read into io.ReaderSeeker
	// to avoid 501 Not Implemented Transfer-Encoding
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	// create upload object, get signed S3 URL
	params := &devicefarm.CreateUploadInput{
		Name:        aws.String(name),
		ProjectArn:  aws.String(projectArn),
		Type:        aws.String(uploadType),
		ContentType: aws.String("application/octet-stream"),
	}
	rApp, err := df.Client.CreateUpload(params)
	if err != nil {
		return "", err
	}
	uploadUrl := *rApp.Upload.Url

	// upload to S3
	req, err := http.NewRequest("PUT", uploadUrl, fileBytes)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return *rApp.Upload.Arn, nil
}

func (df *DeviceFarm) CreateRun(projectArn, poolArn, apk, apkInstrumentation string) (string, error) {
	log.Println(">> Uploading files...")
	log.Println(apk)
	appArn, err := df.CreateUpload(projectArn, apk, "ANDROID_APP", "app.apk")
	if err != nil {
		return "", err
	}
	log.Println(apkInstrumentation)
	instArn, err := df.CreateUpload(projectArn, apkInstrumentation, "INSTRUMENTATION_TEST_PACKAGE", "instrumentation.apk")
	if err != nil {
		return "", err
	}
	log.Println(">> Waiting for files to be processed...")
	time.Sleep(60 * time.Second)
	log.Println(">> Creating test run...")
	params := &devicefarm.ScheduleRunInput{
		DevicePoolArn: aws.String(poolArn),
		ProjectArn:    aws.String(projectArn),
		Test: &devicefarm.ScheduleRunTest{
			Type:           aws.String("INSTRUMENTATION"),
			TestPackageArn: aws.String(instArn),
		},
		AppArn: aws.String(appArn),
	}
	r, err := df.Client.ScheduleRun(params)
	if err != nil {
		return "", err
	}
	return *r.Run.Arn, nil
}
