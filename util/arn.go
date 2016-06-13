package util

import (
	"fmt"
	"regexp"
	"strings"
)

// ArnRegexp matches an AWS ARN with the following capture groups:
//  1: partition
//  2: service
//  3: region
//  4: account-id
//  5: resource
// Note that the resource may be in one of these forms:
//  resource, resourcetype:resource, resourcetype/resource
// See:
// http://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html#genref-arns
var ArnRegexp *regexp.Regexp = regexp.MustCompile("arn:([^:]+):([^:]+):([^:]+):([^:]*):(.*)")

// An Arn specifies the pieces of an AWS ARN
type Arn struct {
	Partition string
	Service   string
	Region    string
	AccountId string
	Resource  string
}

// NewArn parses an ARN string and returns an Arn struct.
// If the ARN is invalid, it returns an error.
func NewArn(arn string) (*Arn, error) {
	match := ArnRegexp.FindStringSubmatch(arn)
	if len(match) != 6 {
		return nil, fmt.Errorf("Invalid ARN: %v", arn)
	}
	return &Arn{
		Partition: match[1],
		Service:   match[2],
		Region:    match[3],
		AccountId: match[4],
		Resource:  match[5],
	}, nil
}

// String returns a string representation of the ARN
func (arn *Arn) String() string {
	parts := []string{"arn", arn.Partition, arn.Service, arn.Region, arn.AccountId, arn.Resource}
	return strings.Join(parts, ":")
}
