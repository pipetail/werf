package storage

import (
	"context"
	"fmt"

	"github.com/werf/lockgate"
	"github.com/werf/logboek"
	"github.com/werf/werf/pkg/werf"
)

func NewGenericLockManager(locker lockgate.Locker) *GenericLockManager {
	return &GenericLockManager{Locker: locker}
}

type GenericLockManager struct {
	// Single Locker for all projects
	Locker lockgate.Locker
}

func (manager *GenericLockManager) LockStage(ctx context.Context, projectName, dependenciesDigest string) (LockHandle, error) {
	_, lock, err := manager.Locker.Acquire(genericStageLockName(projectName, dependenciesDigest), werf.SetupLockerDefaultOptions(ctx, lockgate.AcquireOptions{}))
	return LockHandle{LockgateHandle: lock, ProjectName: projectName}, err
}

func (manager *GenericLockManager) LockStageCache(ctx context.Context, projectName, dependenciesDigest string) (LockHandle, error) {
	_, lock, err := manager.Locker.Acquire(genericStageCacheLockName(projectName, dependenciesDigest), werf.SetupLockerDefaultOptions(ctx, lockgate.AcquireOptions{}))
	return LockHandle{LockgateHandle: lock, ProjectName: projectName}, err
}

func (manager *GenericLockManager) LockImage(ctx context.Context, projectName, imageName string) (LockHandle, error) {
	_, lock, err := manager.Locker.Acquire(genericImageLockName(imageName), werf.SetupLockerDefaultOptions(ctx, lockgate.AcquireOptions{}))
	return LockHandle{LockgateHandle: lock, ProjectName: projectName}, err
}

func (manager *GenericLockManager) LockStagesAndImages(ctx context.Context, projectName string, opts LockStagesAndImagesOptions) (LockHandle, error) {
	_, lock, err := manager.Locker.Acquire(genericStagesAndImagesLockName(projectName), werf.SetupLockerDefaultOptions(ctx, lockgate.AcquireOptions{Shared: opts.GetOrCreateImagesOnly}))
	return LockHandle{LockgateHandle: lock, ProjectName: projectName}, err
}

func (manager *GenericLockManager) Unlock(ctx context.Context, lock LockHandle) error {
	err := manager.Locker.Release(lock.LockgateHandle)
	if err != nil {
		logboek.Context(ctx).Error().LogF("ERROR: unable to release lock for %q: %s\n", lock.LockgateHandle.LockName, err)
	}
	return err
}

func genericStageLockName(projectName, dependenciesDigest string) string {
	return fmt.Sprintf("%s.%s", projectName, dependenciesDigest)
}

func genericStageCacheLockName(projectName, dependenciesDigest string) string {
	return fmt.Sprintf("%s.%s.cache", projectName, dependenciesDigest)
}

func genericImageLockName(imageName string) string {
	return fmt.Sprintf("%s.image", imageName)
}

func genericStagesAndImagesLockName(projectName string) string {
	return fmt.Sprintf("%s.stages_and_images", projectName)
}
