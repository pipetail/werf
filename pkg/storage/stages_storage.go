package storage

import (
	"context"
	"fmt"

	"github.com/werf/werf/pkg/container_runtime"
	"github.com/werf/werf/pkg/image"
)

const (
	LocalStorageAddress             = ":local"
	DefaultKubernetesStorageAddress = "kubernetes://werf-synchronization"
	NamelessImageRecordTag          = "__nameless__"
)

type StagesStorage interface {
	GetStagesIDs(ctx context.Context, projectName string) ([]image.StageID, error)
	GetStagesIDsBySignature(ctx context.Context, projectName, signature string) ([]image.StageID, error)
	GetStageDescription(ctx context.Context, projectName, signature string, uniqueID int64) (*image.StageDescription, error)
	DeleteStage(ctx context.Context, stageDescription *image.StageDescription, options DeleteImageOptions) error
	FilterStagesAndProcessRelatedData(ctx context.Context, stageDescriptions []*image.StageDescription, options FilterStagesAndProcessRelatedDataOptions) ([]*image.StageDescription, error)

	ConstructStageImageName(projectName, signature string, uniqueID int64) string

	// FetchImage will create a local image in the container-runtime
	FetchImage(ctx context.Context, img container_runtime.Image) error
	// StoreImage will store a local image into the container-runtime, local built image should exist prior running store
	StoreImage(ctx context.Context, img container_runtime.Image) error
	ShouldFetchImage(ctx context.Context, img container_runtime.Image) (bool, error)

	CreateRepo(ctx context.Context) error
	DeleteRepo(ctx context.Context) error

	AddManagedImage(ctx context.Context, projectName, imageName string) error
	RmManagedImage(ctx context.Context, projectName, imageName string) error
	GetManagedImages(ctx context.Context, projectName string) ([]string, error)

	PutImageMetadata(ctx context.Context, projectName, imageName, commit, stageID string) error
	RmImageMetadata(ctx context.Context, projectName, imageNameOrID, commit, stageID string) error
	IsImageMetadataExist(ctx context.Context, projectName, imageName, commit, stageID string) (bool, error)
	GetAllAndGroupImageMetadataByImageName(ctx context.Context, projectName string, imageNameList []string) (map[string]map[string][]string, map[string]map[string][]string, error)

	GetClientIDRecords(ctx context.Context, projectName string) ([]*ClientIDRecord, error)
	PostClientIDRecord(ctx context.Context, projectName string, rec *ClientIDRecord) error

	String() string
	Address() string
}

type ClientIDRecord struct {
	ClientID          string
	TimestampMillisec int64
}

func (rec *ClientIDRecord) String() string {
	return fmt.Sprintf("clientID:%s tsMillisec:%d", rec.ClientID, rec.TimestampMillisec)
}

type ImageMetadata struct {
	Digest string
}

type StagesStorageOptions struct {
	RepoStagesStorageOptions
}

func NewStagesStorage(stagesStorageAddress string, containerRuntime container_runtime.ContainerRuntime, options StagesStorageOptions) (StagesStorage, error) {
	if stagesStorageAddress == LocalStorageAddress {
		return NewLocalDockerServerStagesStorage(containerRuntime.(*container_runtime.LocalDockerServerRuntime)), nil
	} else { // Docker registry based stages storage
		return NewRepoStagesStorage(stagesStorageAddress, containerRuntime, options.RepoStagesStorageOptions)
	}
}
