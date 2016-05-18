package awsutil

/*

This is a mock of "github.com/aws/aws-sdk-go/service/devicefarm/devicefarmiface".DeviceFarmAPI.

It exists so that we can test functionality of the awsutil.DeviceFarm wrapper without making
any real requests to AWS.

I am only implementing methods as needed by tests, the rest produce a panic. As you can see
from the enqueue() and dequeue() methods, the idea is for tests to enqueue arbitrary values
which will then be dequeued and returned in FIFO order by methods.

See awsutil_test.go in this directory for usage.

*/

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/devicefarm"
)

type MockClient struct {
	responses [][]interface{}
}

func (client *MockClient) enqueue(values ...interface{}) {
	client.responses = append(client.responses, values)
}

func (client *MockClient) dequeue() []interface{} {
	if len(client.responses) == 0 {
		panic("Nothing in MockClient queue")
	}
	response := client.responses[0]
	client.responses = client.responses[1:]
	return response
}

func (client *MockClient) CreateDevicePoolRequest(*devicefarm.CreateDevicePoolInput) (*request.Request, *devicefarm.CreateDevicePoolOutput) {
	panic("Not implemented")
}

func (client *MockClient) CreateDevicePool(*devicefarm.CreateDevicePoolInput) (*devicefarm.CreateDevicePoolOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) CreateProjectRequest(*devicefarm.CreateProjectInput) (*request.Request, *devicefarm.CreateProjectOutput) {
	panic("Not implemented")
}

func (client *MockClient) CreateProject(*devicefarm.CreateProjectInput) (*devicefarm.CreateProjectOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) CreateUploadRequest(*devicefarm.CreateUploadInput) (*request.Request, *devicefarm.CreateUploadOutput) {
	panic("Not implemented")
}

func (client *MockClient) CreateUpload(*devicefarm.CreateUploadInput) (*devicefarm.CreateUploadOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) DeleteDevicePoolRequest(*devicefarm.DeleteDevicePoolInput) (*request.Request, *devicefarm.DeleteDevicePoolOutput) {
	panic("Not implemented")
}

func (client *MockClient) DeleteDevicePool(*devicefarm.DeleteDevicePoolInput) (*devicefarm.DeleteDevicePoolOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) DeleteProjectRequest(*devicefarm.DeleteProjectInput) (*request.Request, *devicefarm.DeleteProjectOutput) {
	panic("Not implemented")
}

func (client *MockClient) DeleteProject(*devicefarm.DeleteProjectInput) (*devicefarm.DeleteProjectOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) DeleteRunRequest(*devicefarm.DeleteRunInput) (*request.Request, *devicefarm.DeleteRunOutput) {
	panic("Not implemented")
}

func (client *MockClient) DeleteRun(*devicefarm.DeleteRunInput) (*devicefarm.DeleteRunOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) DeleteUploadRequest(*devicefarm.DeleteUploadInput) (*request.Request, *devicefarm.DeleteUploadOutput) {
	panic("Not implemented")
}

func (client *MockClient) DeleteUpload(*devicefarm.DeleteUploadInput) (*devicefarm.DeleteUploadOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetAccountSettingsRequest(*devicefarm.GetAccountSettingsInput) (*request.Request, *devicefarm.GetAccountSettingsOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetAccountSettings(*devicefarm.GetAccountSettingsInput) (*devicefarm.GetAccountSettingsOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetDeviceRequest(*devicefarm.GetDeviceInput) (*request.Request, *devicefarm.GetDeviceOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetDevice(*devicefarm.GetDeviceInput) (*devicefarm.GetDeviceOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetDevicePoolRequest(*devicefarm.GetDevicePoolInput) (*request.Request, *devicefarm.GetDevicePoolOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetDevicePool(*devicefarm.GetDevicePoolInput) (*devicefarm.GetDevicePoolOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetDevicePoolCompatibilityRequest(*devicefarm.GetDevicePoolCompatibilityInput) (*request.Request, *devicefarm.GetDevicePoolCompatibilityOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetDevicePoolCompatibility(*devicefarm.GetDevicePoolCompatibilityInput) (*devicefarm.GetDevicePoolCompatibilityOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetJobRequest(*devicefarm.GetJobInput) (*request.Request, *devicefarm.GetJobOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetJob(*devicefarm.GetJobInput) (*devicefarm.GetJobOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetOfferingStatusRequest(*devicefarm.GetOfferingStatusInput) (*request.Request, *devicefarm.GetOfferingStatusOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetOfferingStatus(*devicefarm.GetOfferingStatusInput) (*devicefarm.GetOfferingStatusOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetOfferingStatusPages(*devicefarm.GetOfferingStatusInput, func(*devicefarm.GetOfferingStatusOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) GetProjectRequest(*devicefarm.GetProjectInput) (*request.Request, *devicefarm.GetProjectOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetProject(*devicefarm.GetProjectInput) (*devicefarm.GetProjectOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetRunRequest(*devicefarm.GetRunInput) (*request.Request, *devicefarm.GetRunOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetRun(*devicefarm.GetRunInput) (*devicefarm.GetRunOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetSuiteRequest(*devicefarm.GetSuiteInput) (*request.Request, *devicefarm.GetSuiteOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetSuite(*devicefarm.GetSuiteInput) (*devicefarm.GetSuiteOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetTestRequest(*devicefarm.GetTestInput) (*request.Request, *devicefarm.GetTestOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetTest(*devicefarm.GetTestInput) (*devicefarm.GetTestOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) GetUploadRequest(*devicefarm.GetUploadInput) (*request.Request, *devicefarm.GetUploadOutput) {
	panic("Not implemented")
}

func (client *MockClient) GetUpload(*devicefarm.GetUploadInput) (*devicefarm.GetUploadOutput, error) {
	response := client.dequeue()
	var out *devicefarm.GetUploadOutput
	if response[0] != nil {
		out = response[0].(*devicefarm.GetUploadOutput)
	}
	var err error
	if response[1] != nil {
		err = response[1].(error)
	}
	return out, err
}

func (client *MockClient) ListArtifactsRequest(*devicefarm.ListArtifactsInput) (*request.Request, *devicefarm.ListArtifactsOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListArtifacts(*devicefarm.ListArtifactsInput) (*devicefarm.ListArtifactsOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListArtifactsPages(*devicefarm.ListArtifactsInput, func(*devicefarm.ListArtifactsOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListDevicePoolsRequest(*devicefarm.ListDevicePoolsInput) (*request.Request, *devicefarm.ListDevicePoolsOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListDevicePools(*devicefarm.ListDevicePoolsInput) (*devicefarm.ListDevicePoolsOutput, error) {
	response := client.dequeue()
	var out *devicefarm.ListDevicePoolsOutput
	if response[0] != nil {
		out = response[0].(*devicefarm.ListDevicePoolsOutput)
	}
	var err error
	if response[1] != nil {
		err = response[1].(error)
	}
	return out, err
}

func (client *MockClient) ListDevicePoolsPages(*devicefarm.ListDevicePoolsInput, func(*devicefarm.ListDevicePoolsOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListDevicesRequest(*devicefarm.ListDevicesInput) (*request.Request, *devicefarm.ListDevicesOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListDevices(*devicefarm.ListDevicesInput) (*devicefarm.ListDevicesOutput, error) {
	response := client.dequeue()
	var out *devicefarm.ListDevicesOutput
	if response[0] != nil {
		out = response[0].(*devicefarm.ListDevicesOutput)
	}
	var err error
	if response[1] != nil {
		err = response[1].(error)
	}
	return out, err
}

func (client *MockClient) ListDevicesPages(*devicefarm.ListDevicesInput, func(*devicefarm.ListDevicesOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListJobsRequest(*devicefarm.ListJobsInput) (*request.Request, *devicefarm.ListJobsOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListJobs(*devicefarm.ListJobsInput) (*devicefarm.ListJobsOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListJobsPages(*devicefarm.ListJobsInput, func(*devicefarm.ListJobsOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListOfferingTransactionsRequest(*devicefarm.ListOfferingTransactionsInput) (*request.Request, *devicefarm.ListOfferingTransactionsOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListOfferingTransactions(*devicefarm.ListOfferingTransactionsInput) (*devicefarm.ListOfferingTransactionsOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListOfferingTransactionsPages(*devicefarm.ListOfferingTransactionsInput, func(*devicefarm.ListOfferingTransactionsOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListOfferingsRequest(*devicefarm.ListOfferingsInput) (*request.Request, *devicefarm.ListOfferingsOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListOfferings(*devicefarm.ListOfferingsInput) (*devicefarm.ListOfferingsOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListOfferingsPages(*devicefarm.ListOfferingsInput, func(*devicefarm.ListOfferingsOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListProjectsRequest(*devicefarm.ListProjectsInput) (*request.Request, *devicefarm.ListProjectsOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListProjects(*devicefarm.ListProjectsInput) (*devicefarm.ListProjectsOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListProjectsPages(*devicefarm.ListProjectsInput, func(*devicefarm.ListProjectsOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListRunsRequest(*devicefarm.ListRunsInput) (*request.Request, *devicefarm.ListRunsOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListRuns(*devicefarm.ListRunsInput) (*devicefarm.ListRunsOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListRunsPages(*devicefarm.ListRunsInput, func(*devicefarm.ListRunsOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListSamplesRequest(*devicefarm.ListSamplesInput) (*request.Request, *devicefarm.ListSamplesOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListSamples(*devicefarm.ListSamplesInput) (*devicefarm.ListSamplesOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListSamplesPages(*devicefarm.ListSamplesInput, func(*devicefarm.ListSamplesOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListSuitesRequest(*devicefarm.ListSuitesInput) (*request.Request, *devicefarm.ListSuitesOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListSuites(*devicefarm.ListSuitesInput) (*devicefarm.ListSuitesOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListSuitesPages(*devicefarm.ListSuitesInput, func(*devicefarm.ListSuitesOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListTestsRequest(*devicefarm.ListTestsInput) (*request.Request, *devicefarm.ListTestsOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListTests(*devicefarm.ListTestsInput) (*devicefarm.ListTestsOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListTestsPages(*devicefarm.ListTestsInput, func(*devicefarm.ListTestsOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListUniqueProblemsRequest(*devicefarm.ListUniqueProblemsInput) (*request.Request, *devicefarm.ListUniqueProblemsOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListUniqueProblems(*devicefarm.ListUniqueProblemsInput) (*devicefarm.ListUniqueProblemsOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListUniqueProblemsPages(*devicefarm.ListUniqueProblemsInput, func(*devicefarm.ListUniqueProblemsOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) ListUploadsRequest(*devicefarm.ListUploadsInput) (*request.Request, *devicefarm.ListUploadsOutput) {
	panic("Not implemented")
}

func (client *MockClient) ListUploads(*devicefarm.ListUploadsInput) (*devicefarm.ListUploadsOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ListUploadsPages(*devicefarm.ListUploadsInput, func(*devicefarm.ListUploadsOutput, bool) bool) error {
	panic("Not implemented")
}

func (client *MockClient) PurchaseOfferingRequest(*devicefarm.PurchaseOfferingInput) (*request.Request, *devicefarm.PurchaseOfferingOutput) {
	panic("Not implemented")
}

func (client *MockClient) PurchaseOffering(*devicefarm.PurchaseOfferingInput) (*devicefarm.PurchaseOfferingOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) RenewOfferingRequest(*devicefarm.RenewOfferingInput) (*request.Request, *devicefarm.RenewOfferingOutput) {
	panic("Not implemented")
}

func (client *MockClient) RenewOffering(*devicefarm.RenewOfferingInput) (*devicefarm.RenewOfferingOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) ScheduleRunRequest(*devicefarm.ScheduleRunInput) (*request.Request, *devicefarm.ScheduleRunOutput) {
	panic("Not implemented")
}

func (client *MockClient) ScheduleRun(*devicefarm.ScheduleRunInput) (*devicefarm.ScheduleRunOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) StopRunRequest(*devicefarm.StopRunInput) (*request.Request, *devicefarm.StopRunOutput) {
	panic("Not implemented")
}

func (client *MockClient) StopRun(*devicefarm.StopRunInput) (*devicefarm.StopRunOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) UpdateDevicePoolRequest(*devicefarm.UpdateDevicePoolInput) (*request.Request, *devicefarm.UpdateDevicePoolOutput) {
	panic("Not implemented")
}

func (client *MockClient) UpdateDevicePool(*devicefarm.UpdateDevicePoolInput) (*devicefarm.UpdateDevicePoolOutput, error) {
	panic("Not implemented")
}

func (client *MockClient) UpdateProjectRequest(*devicefarm.UpdateProjectInput) (*request.Request, *devicefarm.UpdateProjectOutput) {
	panic("Not implemented")
}

func (client *MockClient) UpdateProject(*devicefarm.UpdateProjectInput) (*devicefarm.UpdateProjectOutput, error) {
	panic("Not implemented")
}
