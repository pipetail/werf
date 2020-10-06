package storage

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/example/stringutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"

	"github.com/werf/logboek"

	"github.com/werf/werf/pkg/container_runtime"
	"github.com/werf/werf/pkg/docker"
	"github.com/werf/werf/pkg/image"
)

const (
	LocalStage_ImageRepoPrefix = "werf-stages-storage/"
	LocalStage_ImageRepoFormat = "werf-stages-storage/%s"
	LocalStage_ImageFormat     = "werf-stages-storage/%s:%s-%d"

	LocalManagedImageRecord_ImageNameFormat = "werf-managed-images/%s"
	LocalManagedImageRecord_ImageFormat     = "werf-managed-images/%s:%s"

	LocalImageMetadataByCommitRecord_ImageNameFormat = "werf-images-metadata-by-commit/%s"
	LocalImageMetadataByCommitRecord_TagFormat       = "%s_%s_%s"

	LocalClientIDRecord_ImageNameFormat = "werf-client-id/%s"
	LocalClientIDRecord_ImageFormat     = "werf-client-id/%s:%s-%d"
)

const ImageDeletionFailedDueToUsedByContainerErrorTip = "Use --force option to remove all containers that are based on deleting werf docker images"

func IsImageDeletionFailedDueToUsingByContainerError(err error) bool {
	return strings.HasSuffix(err.Error(), ImageDeletionFailedDueToUsedByContainerErrorTip)
}

func getDependenciesDigestAndUniqueIDFromLocalStageImageTag(repoStageImageTag string) (string, int64, error) {
	parts := strings.SplitN(repoStageImageTag, "-", 2)
	if uniqueID, err := image.ParseUniqueIDAsTimestamp(parts[1]); err != nil {
		return "", 0, fmt.Errorf("unable to parse unique id %s as timestamp: %s", parts[1], err)
	} else {
		return parts[0], uniqueID, nil
	}
}

type LocalDockerServerStagesStorage struct {
	// Local stages storage is compatible only with docker-server backed runtime
	LocalDockerServerRuntime *container_runtime.LocalDockerServerRuntime
}

func NewLocalDockerServerStagesStorage(localDockerServerRuntime *container_runtime.LocalDockerServerRuntime) *LocalDockerServerStagesStorage {
	return &LocalDockerServerStagesStorage{LocalDockerServerRuntime: localDockerServerRuntime}
}

func (storage *LocalDockerServerStagesStorage) ConstructStageImageName(projectName, dependenciesDigest string, uniqueID int64) string {
	return fmt.Sprintf(LocalStage_ImageFormat, projectName, dependenciesDigest, uniqueID)
}

func (storage *LocalDockerServerStagesStorage) GetStagesIDs(ctx context.Context, projectName string) ([]image.StageID, error) {
	filterSet := localStagesStorageFilterSetBase(projectName)
	images, err := docker.Images(ctx, types.ImageListOptions{Filters: filterSet})
	if err != nil {
		return nil, fmt.Errorf("unable to get docker images: %s", err)
	}

	return convertToStagesList(images)
}

func (storage *LocalDockerServerStagesStorage) DeleteStage(ctx context.Context, stageDescription *image.StageDescription, options DeleteImageOptions) error {
	return deleteRepoImageListInLocalDockerServerStagesStorage(ctx, stageDescription, options.RmiForce)
}

func (storage *LocalDockerServerStagesStorage) FilterStagesAndProcessRelatedData(ctx context.Context, stageDescriptions []*image.StageDescription, options FilterStagesAndProcessRelatedDataOptions) ([]*image.StageDescription, error) {
	return processRelatedContainers(ctx, stageDescriptions, processRelatedContainersOptions{
		skipUsedImages:           options.SkipUsedImage,
		rmContainersThatUseImage: options.RmContainersThatUseImage,
		rmForce:                  options.RmForce,
	})
}

func (storage *LocalDockerServerStagesStorage) CreateRepo(_ context.Context) error {
	return nil
}

func (storage *LocalDockerServerStagesStorage) DeleteRepo(_ context.Context) error {
	return nil
}

func makeLocalManagedImageRecord(projectName, imageName string) string {
	tag := imageName
	if imageName == "" {
		tag = NamelessImageRecordTag
	}

	tag = strings.ReplaceAll(tag, "/", "__slash__")
	tag = strings.ReplaceAll(tag, "+", "__plus__")

	return fmt.Sprintf(LocalManagedImageRecord_ImageFormat, projectName, tag)
}

func (storage *LocalDockerServerStagesStorage) GetStageDescription(ctx context.Context, projectName, dependenciesDigest string, uniqueID int64) (*image.StageDescription, error) {
	stageImageName := storage.ConstructStageImageName(projectName, dependenciesDigest, uniqueID)

	if inspect, err := storage.LocalDockerServerRuntime.GetImageInspect(ctx, stageImageName); err != nil {
		return nil, fmt.Errorf("unable to get image %s inspect: %s", stageImageName, err)
	} else if inspect != nil {
		return &image.StageDescription{
			StageID: &image.StageID{DependenciesDigest: dependenciesDigest, UniqueID: uniqueID},
			Info:    image.NewInfoFromInspect(stageImageName, inspect),
		}, nil
	} else {
		return nil, nil
	}
}

func (storage *LocalDockerServerStagesStorage) AddManagedImage(ctx context.Context, projectName, imageName string) error {
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.AddManagedImage %s %s\n", projectName, imageName)

	if validateImageName(imageName) != nil {
		return nil
	}

	fullImageName := makeLocalManagedImageRecord(projectName, imageName)

	if exsts, err := docker.ImageExist(ctx, fullImageName); err != nil {
		return fmt.Errorf("unable to check existence of image %s: %s", fullImageName, err)
	} else if exsts {
		return nil
	}

	if err := docker.CreateImage(ctx, fullImageName, nil); err != nil {
		return fmt.Errorf("unable to create image %s: %s", fullImageName, err)
	}
	return nil
}

func (storage *LocalDockerServerStagesStorage) RmManagedImage(ctx context.Context, projectName, imageName string) error {
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.RmManagedImage %s %s\n", projectName, imageName)

	fullImageName := makeLocalManagedImageRecord(projectName, imageName)

	if exsts, err := docker.ImageExist(ctx, fullImageName); err != nil {
		return fmt.Errorf("unable to check existence of image %q: %s", fullImageName, err)
	} else if !exsts {
		return nil
	}

	if err := docker.CliRmi(ctx, "--force", fullImageName); err != nil {
		return fmt.Errorf("unable to remove image %q: %s", fullImageName, err)
	}

	return nil
}

func (storage *LocalDockerServerStagesStorage) GetManagedImages(ctx context.Context, projectName string) ([]string, error) {
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.GetManagedImages %s\n", projectName)

	filterSet := filters.NewArgs()
	filterSet.Add("reference", fmt.Sprintf(LocalManagedImageRecord_ImageNameFormat, projectName))

	images, err := docker.Images(ctx, types.ImageListOptions{Filters: filterSet})
	if err != nil {
		return nil, fmt.Errorf("unable to get docker images: %s", err)
	}

	var res []string
	for _, img := range images {
		for _, repoTag := range img.RepoTags {
			_, tag := image.ParseRepositoryAndTag(repoTag)

			imageName := unslugDockerImageTagAsImageName(tag)

			if err := validateImageName(imageName); err != nil {
				continue
			}

			res = append(res, imageName)
		}
	}
	return res, nil
}

func (storage *LocalDockerServerStagesStorage) GetStagesIDsByDependenciesDigest(ctx context.Context, projectName, dependenciesDigest string) ([]image.StageID, error) {
	filterSet := filters.NewArgs()
	filterSet.Add("reference", fmt.Sprintf(LocalStage_ImageRepoFormat, projectName))
	// NOTE dependenciesDigest already depends on build-cache-version
	filterSet.Add("label", fmt.Sprintf("%s=%s", image.WerfStageDependenciesDigestLabel, dependenciesDigest))

	images, err := docker.Images(ctx, types.ImageListOptions{Filters: filterSet})
	if err != nil {
		return nil, fmt.Errorf("unable to get docker images: %s", err)
	}

	return convertToStagesList(images)
}

func (storage *LocalDockerServerStagesStorage) ShouldFetchImage(_ context.Context, _ container_runtime.Image) (bool, error) {
	return false, nil
}

func (storage *LocalDockerServerStagesStorage) FetchImage(_ context.Context, _ container_runtime.Image) error {
	return nil
}

func (storage *LocalDockerServerStagesStorage) StoreImage(ctx context.Context, img container_runtime.Image) error {
	return storage.LocalDockerServerRuntime.TagBuiltImageByName(ctx, img)
}

func (storage *LocalDockerServerStagesStorage) PutImageMetadata(ctx context.Context, projectName, imageName, commit, stageID string) error {
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.PutImageMetadata %s %s %s %s\n", projectName, imageName, commit, stageID)

	fullImageName := makeLocalImageMetadataName(projectName, imageName, commit, stageID)
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.PutImageMetadata full image name: %s\n", fullImageName)

	if exists, err := docker.ImageExist(ctx, fullImageName); err != nil {
		return fmt.Errorf("unable to check existence of image %q: %s", fullImageName, err)
	} else if exists {
		return nil
	}

	if err := docker.CreateImage(ctx, fullImageName, nil); err != nil {
		return fmt.Errorf("unable to create image %q: %s", fullImageName, err)
	}

	logboek.Context(ctx).Info().LogF("Put image %s commit %s stage ID %s\n", imageName, commit, stageID)

	return nil
}

func (storage *LocalDockerServerStagesStorage) RmImageMetadata(ctx context.Context, projectName, imageNameOrID, commit, stageID string) error {
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.RmImageMetadata %s %s %s %s\n", projectName, imageNameOrID, commit, stageID)

	fullImageName, err := storage.selectFullImageMetadataName(ctx, projectName, imageNameOrID, commit, stageID)
	if err != nil {
		return err
	}

	if fullImageName == "" {
		return nil
	}

	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.RmImageMetadata full image name: %s\n", fullImageName)
	if err := docker.CliRmi(ctx, "--force", fullImageName); err != nil {
		return fmt.Errorf("unable to remove image %s: %s", fullImageName, err)
	}

	logboek.Context(ctx).Info().LogF("Removed image %s commit %s stage ID %s\n", imageNameOrID, commit, stageID)

	return nil
}

func (storage *LocalDockerServerStagesStorage) selectFullImageMetadataName(ctx context.Context, projectName, imageNameOrID, commit, stageID string) (string, error) {
	for _, fullImageName := range []string{
		makeLocalImageMetadataName(projectName, imageNameOrID, commit, stageID),
		makeLocalImageMetadataNameByImageID(projectName, imageNameOrID, commit, stageID),
	} {
		if exists, err := docker.ImageExist(ctx, fullImageName); err != nil {
			return "", fmt.Errorf("unable to check existence of image %s: %s", fullImageName, err)
		} else if exists {
			return fullImageName, nil
		}
	}

	return "", nil
}

func (storage *LocalDockerServerStagesStorage) IsImageMetadataExist(ctx context.Context, projectName, imageName, commit, stageID string) (bool, error) {
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.IsImageMetadataExist %s %s %s %s\n", projectName, imageName, commit, stageID)

	fullImageName := makeLocalImageMetadataName(projectName, imageName, commit, stageID)
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.IsImageMetadataExist full image name: %s\n", fullImageName)

	return docker.ImageExist(ctx, fullImageName)
}

func (storage *LocalDockerServerStagesStorage) GetAllAndGroupImageMetadataByImageName(ctx context.Context, projectName string, imageNameList []string) (map[string]map[string][]string, map[string]map[string][]string, error) {
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.GetImageNameStageIDCommitList %s %s\n", projectName)

	filterSet := filters.NewArgs()
	filterSet.Add("reference", fmt.Sprintf(LocalImageMetadataByCommitRecord_ImageNameFormat, projectName))

	images, err := docker.Images(ctx, types.ImageListOptions{Filters: filterSet})
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get docker images: %s", err)
	}

	var tags []string
	for _, img := range images {
		for _, repoTag := range img.RepoTags {
			_, tag := image.ParseRepositoryAndTag(repoTag)
			tags = append(tags, tag)
		}
	}

	return groupImageMetadataTagsByImageName(ctx, imageNameList, tags, "")
}

func makeLocalImageMetadataName(projectName, imageName, commit, stageID string) string {
	return makeLocalImageMetadataNameByImageID(projectName, imageNameID(imageName), commit, stageID)
}

func makeLocalImageMetadataNameByImageID(projectName, imageID, commit, stageID string) string {
	return strings.Join([]string{
		fmt.Sprintf(LocalImageMetadataByCommitRecord_ImageNameFormat, projectName),
		fmt.Sprintf(LocalImageMetadataByCommitRecord_TagFormat, imageID, commit, stageID),
	}, ":")
}

func (storage *LocalDockerServerStagesStorage) String() string {
	return LocalStorageAddress
}

func (storage *LocalDockerServerStagesStorage) Address() string {
	return LocalStorageAddress
}

func (storage *LocalDockerServerStagesStorage) GetClientIDRecords(ctx context.Context, projectName string) ([]*ClientIDRecord, error) {
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.GetClientID for project %s\n", projectName)

	filterSet := filters.NewArgs()
	filterSet.Add("reference", fmt.Sprintf(LocalClientIDRecord_ImageNameFormat, projectName))

	images, err := docker.Images(ctx, types.ImageListOptions{Filters: filterSet})
	if err != nil {
		return nil, fmt.Errorf("unable to get docker images: %s", err)
	}

	var res []*ClientIDRecord
	for _, img := range images {
		for _, repoTag := range img.RepoTags {
			_, tag := image.ParseRepositoryAndTag(repoTag)

			tagParts := strings.SplitN(stringutil.Reverse(tag), "-", 2)
			if len(tagParts) != 2 {
				continue
			}

			clientID, timestampMillisecStr := stringutil.Reverse(tagParts[1]), stringutil.Reverse(tagParts[0])

			timestampMillisec, err := strconv.ParseInt(timestampMillisecStr, 10, 64)
			if err != nil {
				continue
			}

			rec := &ClientIDRecord{ClientID: clientID, TimestampMillisec: timestampMillisec}
			res = append(res, rec)

			logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.GetClientID got clientID record: %s\n", rec)
		}
	}

	return res, nil
}

func (storage *LocalDockerServerStagesStorage) PostClientIDRecord(ctx context.Context, projectName string, rec *ClientIDRecord) error {
	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.PostClientID %s for project %s\n", rec.ClientID, projectName)

	fullImageName := fmt.Sprintf(LocalClientIDRecord_ImageFormat, projectName, rec.ClientID, rec.TimestampMillisec)

	logboek.Context(ctx).Debug().LogF("-- LocalDockerServerStagesStorage.PostClientID full image name: %s\n", fullImageName)

	if exsts, err := docker.ImageExist(ctx, fullImageName); err != nil {
		return fmt.Errorf("unable to check existence of image %q: %s", fullImageName, err)
	} else if exsts {
		return nil
	}

	if err := docker.CreateImage(ctx, fullImageName, map[string]string{}); err != nil {
		return fmt.Errorf("unable to create image %q: %s", fullImageName, err)
	}

	logboek.Context(ctx).Info().LogF("Posted new clientID %q for project %s\n", rec.ClientID, projectName)

	return nil
}

type processRelatedContainersOptions struct {
	skipUsedImages           bool
	rmContainersThatUseImage bool
	rmForce                  bool
}

func processRelatedContainers(ctx context.Context, stageDescriptionList []*image.StageDescription, options processRelatedContainersOptions) ([]*image.StageDescription, error) {
	filterSet := filters.NewArgs()
	for _, stageDescription := range stageDescriptionList {
		filterSet.Add("ancestor", stageDescription.Info.ID)
	}

	containerList, err := containerListByFilterSet(ctx, filterSet)
	if err != nil {
		return nil, err
	}

	var stageDescriptionListToExcept []*image.StageDescription
	var containerListToRemove []types.Container
	for _, container := range containerList {
		for _, stageDescription := range stageDescriptionList {
			imageInfo := stageDescription.Info

			if imageInfo.ID == container.ImageID {
				if options.skipUsedImages {
					logboek.Context(ctx).Default().LogFDetails("Skip image %s (used by container %s)\n", logImageName(imageInfo), logContainerName(container))
					stageDescriptionListToExcept = append(stageDescriptionListToExcept, stageDescription)
				} else if options.rmContainersThatUseImage {
					containerListToRemove = append(containerListToRemove, container)
				} else {
					return nil, fmt.Errorf("cannot remove image %s used by container %s\n%s", logImageName(imageInfo), logContainerName(container), ImageDeletionFailedDueToUsedByContainerErrorTip)
				}
			}
		}
	}

	if err := deleteContainers(ctx, containerListToRemove, options.rmForce); err != nil {
		return nil, err
	}

	return exceptStageDescriptionList(stageDescriptionList, stageDescriptionListToExcept...), nil
}

func containerListByFilterSet(ctx context.Context, filterSet filters.Args) ([]types.Container, error) {
	containersOptions := types.ContainerListOptions{}
	containersOptions.All = true
	containersOptions.Quiet = true
	containersOptions.Filters = filterSet

	return docker.Containers(ctx, containersOptions)
}

func deleteContainers(ctx context.Context, containers []types.Container, rmForce bool) error {
	for _, container := range containers {
		if err := docker.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: rmForce}); err != nil {
			return err
		}
	}

	return nil
}

func exceptStageDescriptionList(stageDescriptionList []*image.StageDescription, stageDescriptionListToExcept ...*image.StageDescription) []*image.StageDescription {
	var result []*image.StageDescription

loop:
	for _, sd1 := range stageDescriptionList {
		for _, sd2 := range stageDescriptionListToExcept {
			if sd2 == sd1 {
				continue loop
			}
		}

		result = append(result, sd1)
	}

	return result
}

func localStagesStorageFilterSetBase(projectName string) filters.Args {
	filterSet := filters.NewArgs()
	filterSet.Add("reference", fmt.Sprintf(LocalStage_ImageRepoFormat, projectName))
	filterSet.Add("label", fmt.Sprintf("%s=%s", image.WerfLabel, projectName))
	filterSet.Add("label", fmt.Sprintf("%s=%s", image.WerfCacheVersionLabel, image.BuildCacheVersion))
	return filterSet
}

func logImageName(image *image.Info) string {
	if image.Name == "<none>:<none>" {
		return image.ID
	} else {
		return image.Name
	}
}

func logContainerName(container types.Container) string {
	name := container.ID
	if len(container.Names) != 0 {
		name = container.Names[0]
	}

	return name
}

func convertToStagesList(imageSummaryList []types.ImageSummary) ([]image.StageID, error) {
	var stagesList []image.StageID

	for _, imageSummary := range imageSummaryList {
		repoTags := imageSummary.RepoTags
		if len(repoTags) == 0 {
			repoTags = append(repoTags, "<none>:<none>")
		}

		for _, repoTag := range repoTags {
			_, tag := image.ParseRepositoryAndTag(repoTag)
			if dependenciesDigest, uniqueID, err := getDependenciesDigestAndUniqueIDFromLocalStageImageTag(tag); err != nil {
				return nil, err
			} else {
				stagesList = append(stagesList, image.StageID{DependenciesDigest: dependenciesDigest, UniqueID: uniqueID})
			}
		}
	}

	return stagesList, nil
}

func deleteRepoImageListInLocalDockerServerStagesStorage(ctx context.Context, stageDescription *image.StageDescription, rmiForce bool) error {
	var imageReferences []string
	imageInfo := stageDescription.Info

	if imageInfo.Name == "" {
		imageReferences = append(imageReferences, imageInfo.ID)
	} else {
		isDanglingImage := imageInfo.Name == "<none>:<none>"
		isTaglessImage := !isDanglingImage && imageInfo.Tag == "<none>"

		if isDanglingImage || isTaglessImage {
			imageReferences = append(imageReferences, imageInfo.ID)
		} else {
			imageReferences = append(imageReferences, imageInfo.Name)
		}
	}

	if err := imageReferencesRemove(ctx, imageReferences, rmiForce); err != nil {
		return err
	}

	return nil
}

func imageReferencesRemove(ctx context.Context, references []string, rmiForce bool) error {
	if len(references) == 0 {
		return nil
	}

	var args []string
	if rmiForce {
		args = append(args, "--force")
	}
	args = append(args, references...)

	if err := docker.CliRmi_LiveOutput(ctx, args...); err != nil {
		return err
	}

	return nil
}
