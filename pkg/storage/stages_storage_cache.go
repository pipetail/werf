package storage

import (
	"context"

	"github.com/werf/werf/pkg/image"
)

type StagesStorageCache interface {
	GetAllStages(ctx context.Context, projectName string) (bool, []image.StageID, error)
	DeleteAllStages(ctx context.Context, projectName string) error
	GetStagesByDependenciesDigest(ctx context.Context, projectName, dependenciesDigest string) (bool, []image.StageID, error)
	StoreStagesByDependenciesDigest(ctx context.Context, projectName, dependenciesDigest string, stages []image.StageID) error
	DeleteStagesByDependenciesDigest(ctx context.Context, projectName, dependenciesDigest string) error

	String() string
}
