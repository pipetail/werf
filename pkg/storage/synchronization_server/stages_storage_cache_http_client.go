package synchronization_server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/werf/werf/pkg/image"
)

func NewStagesStorageCacheHttpClient(url string) *StagesStorageCacheHttpClient {
	return &StagesStorageCacheHttpClient{
		URL:        url,
		HttpClient: &http.Client{},
	}
}

type StagesStorageCacheHttpClient struct {
	URL        string
	HttpClient *http.Client
}

func (client *StagesStorageCacheHttpClient) String() string {
	return fmt.Sprintf("http-client %s", client.URL)
}

func (client *StagesStorageCacheHttpClient) GetAllStages(_ context.Context, projectName string) (bool, []image.StageID, error) {
	var request = GetAllStagesRequest{projectName}
	var response GetAllStagesResponse
	if err := PerformPost(client.HttpClient, fmt.Sprintf("%s/%s", client.URL, "get-all-stages"), request, &response); err != nil {
		return false, nil, err
	}
	return response.Found, response.Stages, response.Err.Error
}

func (client *StagesStorageCacheHttpClient) DeleteAllStages(_ context.Context, projectName string) error {
	var request = DeleteAllStagesRequest{projectName}
	var response DeleteAllStagesResponse
	if err := PerformPost(client.HttpClient, fmt.Sprintf("%s/%s", client.URL, "delete-all-stages"), request, &response); err != nil {
		return err
	}
	return response.Err.Error
}

func (client *StagesStorageCacheHttpClient) GetStagesByDependenciesDigest(_ context.Context, projectName, dependenciesDigest string) (bool, []image.StageID, error) {
	var request = GetStagesByDependenciesDigestRequest{projectName, dependenciesDigest}
	var response GetStagesByDependenciesDigestResponse
	if err := PerformPost(client.HttpClient, fmt.Sprintf("%s/%s", client.URL, "get-stages-by-dependenciesDigest"), request, &response); err != nil {
		return false, nil, err
	}
	return response.Found, response.Stages, response.Err.Error
}

func (client *StagesStorageCacheHttpClient) StoreStagesByDependenciesDigest(_ context.Context, projectName, dependenciesDigest string, stages []image.StageID) error {
	var request = StoreStagesByDependenciesDigestRequest{projectName, dependenciesDigest, stages}
	var response StoreStagesByDependenciesDigestResponse
	if err := PerformPost(client.HttpClient, fmt.Sprintf("%s/%s", client.URL, "store-stages-by-dependenciesDigest"), request, &response); err != nil {
		return err
	}
	return response.Err.Error
}

func (client *StagesStorageCacheHttpClient) DeleteStagesByDependenciesDigest(_ context.Context, projectName, dependenciesDigest string) error {
	var request = DeleteStagesByDependenciesDigestRequest{projectName, dependenciesDigest}
	var response DeleteStagesByDependenciesDigestResponse
	if err := PerformPost(client.HttpClient, fmt.Sprintf("%s/%s", client.URL, "delete-stages-by-dependenciesDigest"), request, &response); err != nil {
		return err
	}
	return response.Err.Error
}
