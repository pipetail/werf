package storage

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/example/stringutil"

	"github.com/werf/logboek"

	"github.com/werf/werf/pkg/container_runtime"
	"github.com/werf/werf/pkg/docker_registry"
	"github.com/werf/werf/pkg/image"
	"github.com/werf/werf/pkg/util"
)

const (
	RepoStage_ImageFormat = "%s:%s-%d"

	RepoManagedImageRecord_ImageTagPrefix  = "managed-image-"
	RepoManagedImageRecord_ImageNameFormat = "%s:managed-image-%s"

	RepoImageMetadataByCommitRecord_ImageTagPrefix = "meta-"
	RepoImageMetadataByCommitRecord_TagFormat      = "meta-%s_%s_%s"

	RepoClientIDRecrod_ImageTagPrefix  = "client-id-"
	RepoClientIDRecrod_ImageNameFormat = "%s:client-id-%s-%d"

	UnexpectedTagFormatErrorPrefix = "unexpected tag format"
)

func getDependenciesDigestAndUniqueIDFromRepoStageImageTag(repoStageImageTag string) (string, int64, error) {
	parts := strings.SplitN(repoStageImageTag, "-", 2)

	if len(parts) != 2 {
		return "", 0, fmt.Errorf("%s %s", UnexpectedTagFormatErrorPrefix, repoStageImageTag)
	}

	if uniqueID, err := image.ParseUniqueIDAsTimestamp(parts[1]); err != nil {
		return "", 0, fmt.Errorf("%s %s: unable to parse unique id %s as timestamp: %s", UnexpectedTagFormatErrorPrefix, repoStageImageTag, parts[1], err)
	} else {
		return parts[0], uniqueID, nil
	}
}

func isUnexpectedTagFormatError(err error) bool {
	return strings.HasPrefix(err.Error(), UnexpectedTagFormatErrorPrefix)
}

type RepoStagesStorage struct {
	RepoAddress      string
	DockerRegistry   docker_registry.DockerRegistry
	ContainerRuntime container_runtime.ContainerRuntime
}

type RepoStagesStorageOptions struct {
	docker_registry.DockerRegistryOptions
	Implementation string
}

func NewRepoStagesStorage(repoAddress string, containerRuntime container_runtime.ContainerRuntime, options RepoStagesStorageOptions) (*RepoStagesStorage, error) {
	implementation := options.Implementation

	dockerRegistry, err := docker_registry.NewDockerRegistry(repoAddress, implementation, options.DockerRegistryOptions)
	if err != nil {
		return nil, fmt.Errorf("error creating docker registry accessor for repo %q: %s", repoAddress, err)
	}

	return &RepoStagesStorage{
		RepoAddress:      repoAddress,
		DockerRegistry:   dockerRegistry,
		ContainerRuntime: containerRuntime,
	}, nil
}

func (storage *RepoStagesStorage) ConstructStageImageName(_, dependenciesDigest string, uniqueID int64) string {
	return fmt.Sprintf(RepoStage_ImageFormat, storage.RepoAddress, dependenciesDigest, uniqueID)
}

func (storage *RepoStagesStorage) GetStagesIDs(ctx context.Context, projectName string) ([]image.StageID, error) {
	var res []image.StageID

	if tags, err := storage.DockerRegistry.Tags(ctx, storage.RepoAddress); err != nil {
		return nil, fmt.Errorf("unable to fetch tags for repo %q: %s", storage.RepoAddress, err)
	} else {
		logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.GetRepoImagesByDependenciesDigest fetched tags for %q: %#v\n", storage.RepoAddress, tags)

		for _, tag := range tags {
			if strings.HasPrefix(tag, RepoManagedImageRecord_ImageTagPrefix) || strings.HasPrefix(tag, RepoImageMetadataByCommitRecord_ImageTagPrefix) {
				continue
			}

			if dependenciesDigest, uniqueID, err := getDependenciesDigestAndUniqueIDFromRepoStageImageTag(tag); err != nil {
				if isUnexpectedTagFormatError(err) {
					logboek.Context(ctx).Debug().LogLn(err.Error())
					continue
				}
				return nil, err
			} else {
				res = append(res, image.StageID{DependenciesDigest: dependenciesDigest, UniqueID: uniqueID})

				logboek.Context(ctx).Debug().LogF("Selected stage by dependenciesDigest %q uniqueID %d\n", dependenciesDigest, uniqueID)
			}
		}

		return res, nil
	}
}

func (storage *RepoStagesStorage) DeleteStage(ctx context.Context, stageDescription *image.StageDescription, _ DeleteImageOptions) error {
	return storage.DockerRegistry.DeleteRepoImage(ctx, stageDescription.Info)
}

func (storage *RepoStagesStorage) FilterStagesAndProcessRelatedData(_ context.Context, stageDescriptions []*image.StageDescription, _ FilterStagesAndProcessRelatedDataOptions) ([]*image.StageDescription, error) {
	return stageDescriptions, nil
}

func (storage *RepoStagesStorage) CreateRepo(ctx context.Context) error {
	return storage.DockerRegistry.CreateRepo(ctx, storage.RepoAddress)
}

func (storage *RepoStagesStorage) DeleteRepo(ctx context.Context) error {
	return storage.DockerRegistry.DeleteRepo(ctx, storage.RepoAddress)
}

func (storage *RepoStagesStorage) GetStagesIDsByDependenciesDigest(ctx context.Context, projectName, dependenciesDigest string) ([]image.StageID, error) {
	var res []image.StageID

	if tags, err := storage.DockerRegistry.Tags(ctx, storage.RepoAddress); err != nil {
		return nil, fmt.Errorf("unable to fetch tags for repo %q: %s", storage.RepoAddress, err)
	} else {
		logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.GetRepoImagesByDependenciesDigest fetched tags for %q: %#v\n", storage.RepoAddress, tags)
		for _, tag := range tags {
			if !strings.HasPrefix(tag, dependenciesDigest) {
				logboek.Context(ctx).Debug().LogF("Discard tag %q: should have prefix %q\n", tag, dependenciesDigest)
				continue
			}
			if _, uniqueID, err := getDependenciesDigestAndUniqueIDFromRepoStageImageTag(tag); err != nil {
				if isUnexpectedTagFormatError(err) {
					logboek.Context(ctx).Debug().LogLn(err.Error())
					continue
				}
				return nil, err
			} else {
				logboek.Context(ctx).Debug().LogF("Tag %q is suitable for dependenciesDigest %q\n", tag, dependenciesDigest)
				res = append(res, image.StageID{DependenciesDigest: dependenciesDigest, UniqueID: uniqueID})
			}
		}
	}

	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.GetRepoImagesByDependenciesDigest result for %q: %#v\n", storage.RepoAddress, res)

	return res, nil
}

func (storage *RepoStagesStorage) GetStageDescription(ctx context.Context, projectName, dependenciesDigest string, uniqueID int64) (*image.StageDescription, error) {
	stageImageName := storage.ConstructStageImageName(projectName, dependenciesDigest, uniqueID)

	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage GetStageDescription %s %s %d\n", projectName, dependenciesDigest, uniqueID)
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage stageImageName = %q\n", stageImageName)

	if imgInfo, err := storage.DockerRegistry.TryGetRepoImage(ctx, stageImageName); err != nil {
		return nil, err
	} else if imgInfo != nil {
		return &image.StageDescription{
			StageID: &image.StageID{DependenciesDigest: dependenciesDigest, UniqueID: uniqueID},
			Info:    imgInfo,
		}, nil
	}
	return nil, nil
}

func (storage *RepoStagesStorage) AddManagedImage(ctx context.Context, projectName, imageName string) error {
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.AddManagedImage %s %s\n", projectName, imageName)

	if validateImageName(imageName) != nil {
		return nil
	}

	fullImageName := makeRepoManagedImageRecord(storage.RepoAddress, imageName)
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.AddManagedImage full image name: %s\n", fullImageName)

	if isExists, err := storage.DockerRegistry.IsRepoImageExists(ctx, fullImageName); err != nil {
		return err
	} else if isExists {
		logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.AddManagedImage record %q is exists => exiting\n", fullImageName)
		return nil
	}

	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.AddManagedImage record %q does not exist => creating record\n", fullImageName)

	if err := storage.DockerRegistry.PushImage(ctx, fullImageName, nil); err != nil {
		return fmt.Errorf("unable to push image %s: %s", fullImageName, err)
	}

	return nil
}

func (storage *RepoStagesStorage) RmManagedImage(ctx context.Context, projectName, imageName string) error {
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.RmManagedImage %s %s\n", projectName, imageName)

	fullImageName := makeRepoManagedImageRecord(storage.RepoAddress, imageName)

	if imgInfo, err := storage.DockerRegistry.TryGetRepoImage(ctx, fullImageName); err != nil {
		return fmt.Errorf("unable to get repo image %q info: %s", fullImageName, err)
	} else if imgInfo == nil {
		logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.RmManagedImage record %q does not exist => exiting\n", fullImageName)
		return nil
	} else {
		if err := storage.DockerRegistry.DeleteRepoImage(ctx, imgInfo); err != nil {
			return fmt.Errorf("unable to delete image %q from repo: %s", fullImageName, err)
		}
	}

	return nil
}

func (storage *RepoStagesStorage) GetManagedImages(ctx context.Context, projectName string) ([]string, error) {
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.GetManagedImages %s\n", projectName)

	var res []string

	if tags, err := storage.DockerRegistry.Tags(ctx, storage.RepoAddress); err != nil {
		return nil, fmt.Errorf("unable to get repo %s tags: %s", storage.RepoAddress, err)
	} else {
		for _, tag := range tags {
			if !strings.HasPrefix(tag, RepoManagedImageRecord_ImageTagPrefix) {
				continue
			}

			managedImageName := unslugDockerImageTagAsImageName(strings.TrimPrefix(tag, RepoManagedImageRecord_ImageTagPrefix))

			if validateImageName(managedImageName) != nil {
				continue
			}

			res = append(res, managedImageName)
		}
	}

	return res, nil
}

func (storage *RepoStagesStorage) FetchImage(ctx context.Context, img container_runtime.Image) error {
	switch containerRuntime := storage.ContainerRuntime.(type) {
	case *container_runtime.LocalDockerServerRuntime:
		return containerRuntime.PullImageFromRegistry(ctx, img)
	default:
		// TODO: case *container_runtime.LocalHostRuntime:
		panic("not implemented")
	}
}

func (storage *RepoStagesStorage) StoreImage(ctx context.Context, img container_runtime.Image) error {
	switch containerRuntime := storage.ContainerRuntime.(type) {
	case *container_runtime.LocalDockerServerRuntime:
		dockerImage := img.(*container_runtime.DockerImage)

		if dockerImage.Image.GetBuiltId() != "" {
			return containerRuntime.PushBuiltImage(ctx, img)
		} else {
			return containerRuntime.PushImage(ctx, img)
		}

	default:
		// TODO: case *container_runtime.LocalHostRuntime:
		panic("not implemented")
	}
}

func (storage *RepoStagesStorage) ShouldFetchImage(_ context.Context, img container_runtime.Image) (bool, error) {
	switch storage.ContainerRuntime.(type) {
	case *container_runtime.LocalDockerServerRuntime:
		dockerImage := img.(*container_runtime.DockerImage)
		return !dockerImage.Image.IsExistsLocally(), nil
	default:
		panic("not implemented")
	}
}

func (storage *RepoStagesStorage) PutImageMetadata(ctx context.Context, projectName, imageName, commit, stageID string) error {
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.PutImageMetadata %s %s %s %s\n", projectName, imageName, commit, stageID)

	fullImageName := makeRepoImageMetadataName(storage.RepoAddress, imageName, commit, stageID)
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.PutImageMetadata full image name: %s\n", fullImageName)

	if err := storage.DockerRegistry.PushImage(ctx, fullImageName, nil); err != nil {
		return fmt.Errorf("unable to push image %s with metadata: %s", fullImageName, err)
	}
	logboek.Context(ctx).Info().LogF("Put image %s commit %s stage ID %s\n", imageName, commit, stageID)

	return nil
}

func (storage *RepoStagesStorage) RmImageMetadata(ctx context.Context, projectName, imageNameOrID, commit, stageID string) error {
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.RmImageMetadata %s %s %s %s\n", projectName, imageNameOrID, commit, stageID)

	img, err := storage.selectMetadataNameImage(ctx, imageNameOrID, commit, stageID)
	if err != nil {
		return err
	}

	if img == nil {
		return nil
	}
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.RmImageMetadata full image name: %s\n", img.Tag)

	if err := storage.DockerRegistry.DeleteRepoImage(ctx, img); err != nil {
		return fmt.Errorf("unable to remove repo image %s: %s", img.Tag, err)
	}

	logboek.Context(ctx).Info().LogF("Removed image %s commit %s stage ID %s\n", imageNameOrID, commit, stageID)

	return nil
}

func (storage *RepoStagesStorage) selectMetadataNameImage(ctx context.Context, imageNameOrID, commit, stageID string) (*image.Info, error) {
	for _, fullImageName := range []string{
		makeRepoImageMetadataName(storage.RepoAddress, imageNameOrID, commit, stageID),
		makeRepoImageMetadataNameByImageID(storage.RepoAddress, imageNameOrID, commit, stageID),
	} {
		if img, err := storage.DockerRegistry.TryGetRepoImage(ctx, fullImageName); err != nil {
			return nil, fmt.Errorf("unable to get repo image %s: %s", fullImageName, err)
		} else if img != nil {
			return img, nil
		}
	}

	return nil, nil
}

func (storage *RepoStagesStorage) IsImageMetadataExist(ctx context.Context, projectName, imageName, commit, stageID string) (bool, error) {
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.IsImageMetadataExist %s %s %s %s\n", projectName, imageName, commit, stageID)

	fullImageName := makeRepoImageMetadataName(storage.RepoAddress, imageName, commit, stageID)
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.IsImageMetadataExist full image name: %s\n", fullImageName)

	img, err := storage.DockerRegistry.TryGetRepoImage(ctx, fullImageName)
	return img != nil, err
}

func (storage *RepoStagesStorage) GetAllAndGroupImageMetadataByImageName(ctx context.Context, projectName string, imageNameList []string) (map[string]map[string][]string, map[string]map[string][]string, error) {
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.GetImageNameStageIDCommitList %s %s\n", projectName)

	tags, err := storage.DockerRegistry.Tags(ctx, storage.RepoAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get repo %s tags: %s", storage.RepoAddress, err)
	}

	return groupImageMetadataTagsByImageName(ctx, imageNameList, tags, RepoImageMetadataByCommitRecord_ImageTagPrefix)
}

func groupImageMetadataTagsByImageName(ctx context.Context, imageNameList []string, tags []string, imageTagPrefix string) (map[string]map[string][]string, map[string]map[string][]string, error) {
	imageNameNameByID := map[string]string{}
	for _, imageName := range imageNameList {
		imageNameNameByID[imageNameID(imageName)] = imageName
	}

	result := map[string]map[string][]string{}
	resultNotManagedImageName := map[string]map[string][]string{}
	for _, tag := range tags {
		var res map[string]map[string][]string

		if !strings.HasPrefix(tag, imageTagPrefix) {
			continue
		}

		sluggedImageAndCommit := strings.TrimPrefix(tag, imageTagPrefix)
		sluggedImageAndCommitParts := strings.Split(sluggedImageAndCommit, "_")
		if len(sluggedImageAndCommitParts) != 3 {
			// unexpected
			continue
		}

		tagImageNameID := sluggedImageAndCommitParts[0]
		tagCommit := sluggedImageAndCommitParts[1]
		tagStageID := sluggedImageAndCommitParts[2]

		logboek.Context(ctx).Debug().LogF("Found image ID %s commit %s stage ID %s\n", tagImageNameID, tagCommit, tagStageID)

		imageName, ok := imageNameNameByID[tagImageNameID]
		if !ok {
			res = resultNotManagedImageName
			imageName = tagImageNameID
		} else {
			res = result
		}

		stageIDCommitList, ok := res[imageName]
		if !ok {
			stageIDCommitList = map[string][]string{}
		}

		commitList, ok := stageIDCommitList[tagStageID]
		if !ok {
			commitList = []string{}
		}

		commitList = append(commitList, tagCommit)
		stageIDCommitList[tagStageID] = commitList
		res[imageName] = stageIDCommitList
	}

	return result, resultNotManagedImageName, nil
}

func makeRepoImageMetadataName(repoAddress, imageNameOrID, commit, stageID string) string {
	return makeRepoImageMetadataNameByImageID(repoAddress, imageNameID(imageNameOrID), commit, stageID)
}

func makeRepoImageMetadataNameByImageID(repoAddress, imageID, commit, stageID string) string {
	return strings.Join([]string{
		repoAddress,
		fmt.Sprintf(RepoImageMetadataByCommitRecord_TagFormat, imageID, commit, stageID),
	}, ":")
}

func (storage *RepoStagesStorage) String() string {
	return storage.RepoAddress
}

func (storage *RepoStagesStorage) Address() string {
	return storage.RepoAddress
}

func makeRepoManagedImageRecord(repoAddress, imageName string) string {
	return fmt.Sprintf(RepoManagedImageRecord_ImageNameFormat, repoAddress, slugImageNameAsDockerImageTag(imageName))
}

func imageNameID(imageName string) string {
	return util.MurmurHash(imageName)
}

func slugImageNameAsDockerImageTag(imageName string) string {
	res := imageName
	res = strings.ReplaceAll(res, "/", "__slash__")
	res = strings.ReplaceAll(res, "+", "__plus__")

	if imageName == "" {
		res = NamelessImageRecordTag
	}

	return res
}

func unslugDockerImageTagAsImageName(tag string) string {
	res := tag
	res = strings.ReplaceAll(res, "__slash__", "/")
	res = strings.ReplaceAll(res, "__plus__", "+")

	if res == NamelessImageRecordTag {
		res = ""
	}

	return res
}

func validateImageName(name string) error {
	if strings.ToLower(name) != name {
		return fmt.Errorf("no upcase symbols allowed")
	}
	return nil
}

func (storage *RepoStagesStorage) GetClientIDRecords(ctx context.Context, projectName string) ([]*ClientIDRecord, error) {
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.GetClientIDRecords for project %s\n", projectName)

	var res []*ClientIDRecord

	if tags, err := storage.DockerRegistry.Tags(ctx, storage.RepoAddress); err != nil {
		return nil, fmt.Errorf("unable to get repo %s tags: %s", storage.RepoAddress, err)
	} else {
		for _, tag := range tags {
			if !strings.HasPrefix(tag, RepoClientIDRecrod_ImageTagPrefix) {
				continue
			}

			tagWithoutPrefix := strings.TrimPrefix(tag, RepoClientIDRecrod_ImageTagPrefix)
			dataParts := strings.SplitN(stringutil.Reverse(tagWithoutPrefix), "-", 2)
			if len(dataParts) != 2 {
				continue
			}

			clientID, timestampMillisecStr := stringutil.Reverse(dataParts[1]), stringutil.Reverse(dataParts[0])

			timestampMillisec, err := strconv.ParseInt(timestampMillisecStr, 10, 64)
			if err != nil {
				continue
			}

			rec := &ClientIDRecord{ClientID: clientID, TimestampMillisec: timestampMillisec}
			res = append(res, rec)

			logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.GetClientIDRecords got clientID record: %s\n", rec)
		}
	}

	return res, nil
}

func (storage *RepoStagesStorage) PostClientIDRecord(ctx context.Context, projectName string, rec *ClientIDRecord) error {
	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.PostClientID %s for project %s\n", rec.ClientID, projectName)

	fullImageName := fmt.Sprintf(RepoClientIDRecrod_ImageNameFormat, storage.RepoAddress, rec.ClientID, rec.TimestampMillisec)

	logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.PostClientID full image name: %s\n", fullImageName)

	if isExists, err := storage.DockerRegistry.IsRepoImageExists(ctx, fullImageName); err != nil {
		return err
	} else if isExists {
		logboek.Context(ctx).Debug().LogF("-- RepoStagesStorage.AddManagedImage record %q is exists => exiting\n", fullImageName)
		return nil
	}

	if err := storage.DockerRegistry.PushImage(ctx, fullImageName, nil); err != nil {
		return fmt.Errorf("unable to push image %s: %s", fullImageName, err)
	}

	logboek.Context(ctx).Info().LogF("Posted new clientID %q for project %s\n", rec.ClientID, projectName)

	return nil
}
