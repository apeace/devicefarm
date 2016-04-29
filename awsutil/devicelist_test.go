package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSort(t *testing.T) {
	assert := assert.New(t)

	device1 := &devicefarm.Device{Name: aws.String("foo")}
	device2 := &devicefarm.Device{Name: aws.String("bar")}

	list := DeviceList{device1, device2}
	list.Sort()

	assert.Equal(DeviceList{device2, device1}, list)
}
