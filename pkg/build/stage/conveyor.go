package stage

import (
	"context"

	"github.com/werf/werf/pkg/build/import_server"
)

type Conveyor interface {
	GetImageStageDigest(imageName, stageName string) string
	GetImageDigest(imageName string) string

	GetImageNameForLastImageStage(imageName string) string
	GetImageIDForLastImageStage(imageName string) string

	GetImageNameForImageStage(imageName, stageName string) string
	GetImageIDForImageStage(imageName, stageName string) string

	GetImportServer(ctx context.Context, imageName, stageName string) (import_server.ImportServer, error)
	GetLocalGitRepoVirtualMergeOptions() VirtualMergeOptions

	GetProjectRepoCommit(ctx context.Context) (string, error)
}

type VirtualMergeOptions struct {
	VirtualMerge           bool
	VirtualMergeFromCommit string
	VirtualMergeIntoCommit string
}
