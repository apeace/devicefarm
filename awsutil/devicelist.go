package awsutil

import (
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"sort"
)

type DeviceList []*devicefarm.Device

func (list DeviceList) Len() int {
	return len(list)
}

func (list DeviceList) Less(i, j int) bool {
	return *list[i].Name < *list[j].Name
}

func (list DeviceList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list DeviceList) Sort() {
	sort.Sort(list)
}
