package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewArn(t *testing.T) {
	assert := assert.New(t)

	example := "arn:aws:devicefarm:us-west-2::device:5F9CEB47606A4709879003E11BEAFB08"
	arn, err := NewArn(example)
	assert.Nil(err)
	assert.Equal(Arn{
		Partition: "aws",
		Service:   "devicefarm",
		Region:    "us-west-2",
		AccountId: "",
		Resource:  "device:5F9CEB47606A4709879003E11BEAFB08",
	}, *arn)

	arn, err = NewArn("foo")
	assert.Nil(arn)
	assert.NotNil(err)
}
