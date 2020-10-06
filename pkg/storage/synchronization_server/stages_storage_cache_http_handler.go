package synchronization_server

import (
	"context"
	"net/http"

	"github.com/werf/logboek"
	"github.com/werf/werf/pkg/image"
	"github.com/werf/werf/pkg/util"

	"github.com/werf/werf/pkg/storage"
)

func NewStagesStorageCacheHttpHandler(ctx context.Context, stagesStorageCache storage.StagesStorageCache) *StagesStorageCacheHttpHandler {
	handler := &StagesStorageCacheHttpHandler{
		StagesStorageCache: stagesStorageCache,
		ServeMux:           http.NewServeMux(),
	}
	handler.HandleFunc("/get-all-stages", handler.handleGetAllStages(ctx))
	handler.HandleFunc("/delete-all-stages", handler.handleDeleteAllStages(ctx))
	handler.HandleFunc("/get-stages-by-dependenciesDigest", handler.handleGetStagesByDependenciesDigest(ctx))
	handler.HandleFunc("/store-stages-by-dependenciesDigest", handler.handleStoreStagesByDependenciesDigest(ctx))
	handler.HandleFunc("/delete-stages-by-dependenciesDigest", handler.handleDeleteStagesByDependenciesDigest(ctx))

	return handler
}

type StagesStorageCacheHttpHandler struct {
	*http.ServeMux
	StagesStorageCache storage.StagesStorageCache
}

type GetAllStagesRequest struct {
	ProjectName string `json:"projectName"`
}
type GetAllStagesResponse struct {
	Err    util.SerializableError `json:"err"`
	Found  bool                   `json:"found"`
	Stages []image.StageID        `json:"stages"`
}

func (handler *StagesStorageCacheHttpHandler) handleGetAllStages(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request GetAllStagesRequest
		var response GetAllStagesResponse
		HandleRequest(w, r, &request, &response, func() {
			logboek.Context(ctx).Debug().LogF("StagesStorageCacheHttpHandler -- GetAllStages request %#v\n", request)
			response.Found, response.Stages, response.Err.Error = handler.StagesStorageCache.GetAllStages(ctx, request.ProjectName)
			logboek.Context(ctx).Debug().LogF("StagesStorageCacheHttpHandler -- GetAllStages response %#v\n", response)
		})
	}
}

type DeleteAllStagesRequest struct {
	ProjectName string `json:"projectName"`
}
type DeleteAllStagesResponse struct {
	Err util.SerializableError `json:"err"`
}

func (handler *StagesStorageCacheHttpHandler) handleDeleteAllStages(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request DeleteAllStagesRequest
		var response DeleteAllStagesResponse
		HandleRequest(w, r, &request, &response, func() {
			logboek.Context(ctx).Debug().LogF("StagesStorageCacheHttpHandler -- DeleteAllStages request %#v\n", request)
			response.Err.Error = handler.StagesStorageCache.DeleteAllStages(ctx, request.ProjectName)
			logboek.Context(ctx).Debug().LogF("StagesStorageCacheHttpHandler -- DeleteAllStages response %#v\n", response)
		})
	}
}

type GetStagesByDependenciesDigestRequest struct {
	ProjectName string `json:"projectName"`
	DependenciesDigest   string `json:"dependenciesDigest"`
}
type GetStagesByDependenciesDigestResponse struct {
	Err    util.SerializableError `json:"err"`
	Found  bool                   `json:"found"`
	Stages []image.StageID        `json:"stages"`
}

func (handler *StagesStorageCacheHttpHandler) handleGetStagesByDependenciesDigest(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request GetStagesByDependenciesDigestRequest
		var response GetStagesByDependenciesDigestResponse
		HandleRequest(w, r, &request, &response, func() {
			logboek.Context(ctx).Debug().LogF("StagesStorageCacheHttpHandler -- GetStagesByDependenciesDigest request %#v\n", request)
			response.Found, response.Stages, response.Err.Error = handler.StagesStorageCache.GetStagesByDependenciesDigest(ctx, request.ProjectName, request.DependenciesDigest)
			logboek.Context(ctx).Debug().LogF("StagesStorageCacheHttpHandler -- GetStagesByDependenciesDigest response %#v\n", response)
		})
	}
}

type StoreStagesByDependenciesDigestRequest struct {
	ProjectName string          `json:"projectName"`
	DependenciesDigest   string          `json:"dependenciesDigest"`
	Stages      []image.StageID `json:"stages"`
}
type StoreStagesByDependenciesDigestResponse struct {
	Err util.SerializableError `json:"err"`
}

func (handler *StagesStorageCacheHttpHandler) handleStoreStagesByDependenciesDigest(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request StoreStagesByDependenciesDigestRequest
		var response StoreStagesByDependenciesDigestResponse
		HandleRequest(w, r, &request, &response, func() {
			logboek.Context(ctx).Debug().LogF("StagesStorageCacheHttpHandler -- StoreStagesByDependenciesDigest request %#v\n", request)
			response.Err.Error = handler.StagesStorageCache.StoreStagesByDependenciesDigest(ctx, request.ProjectName, request.DependenciesDigest, request.Stages)
			logboek.Context(ctx).Debug().LogF("StagesStorageCacheHttpHandler -- StoreStagesByDependenciesDigest response %#v\n", response)
		})
	}
}

type DeleteStagesByDependenciesDigestRequest struct {
	ProjectName string `json:"projectName"`
	DependenciesDigest   string `json:"dependenciesDigest"`
}
type DeleteStagesByDependenciesDigestResponse struct {
	Err util.SerializableError `json:"err"`
}

func (handler *StagesStorageCacheHttpHandler) handleDeleteStagesByDependenciesDigest(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request DeleteStagesByDependenciesDigestRequest
		var response DeleteStagesByDependenciesDigestResponse
		HandleRequest(w, r, &request, &response, func() {
			logboek.Context(ctx).Debug().LogF("StagesStorageCacheHttpHandler -- DeleteStagesByDependenciesDigest request %#v\n", request)
			response.Err.Error = handler.StagesStorageCache.DeleteStagesByDependenciesDigest(ctx, request.ProjectName, request.DependenciesDigest)
			logboek.Context(ctx).Debug().LogF("StagesStorageCacheHttpHandler -- DeleteStagesByDependenciesDigest response %#v\n", response)
		})
	}
}
