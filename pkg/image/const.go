package image

const (
	WerfLabel                      = "werf"
	WerfVersionLabel               = "werf-version"
	WerfCacheVersionLabel          = "werf-cache-version"
	WerfImageLabel                 = "werf-image"
	WerfImageNameLabel             = "werf-image-name"
	WerfDockerImageName            = "werf-docker-image-name"
	WerfStageDependenciesDigestLabel        = "werf-stage-dependenciesDigest"
	WerfStageDigestLabel = "werf-stage-digest"
	WerfProjectRepoCommitLabel     = "werf-project-repo-commit"
	WerfDigestLabel      = "werf-digest"

	WerfMountTmpDirLabel          = "werf-mount-type-tmp-dir"
	WerfMountBuildDirLabel        = "werf-mount-type-build-dir"
	WerfMountCustomDirLabelPrefix = "werf-mount-type-custom-dir-"

	WerfImportLabelPrefix = "werf-import-"

	BuildCacheVersion = "1.1"

	StageContainerNamePrefix = "werf.build."
)
