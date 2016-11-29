package handlers

import (
	"net/http"

	"github.com/src-d/landing/api/config"
	"github.com/src-d/landing/api/github"

	"gopkg.in/gin-gonic/gin.v1"
)

type Repositories interface {
	Main(ctx *gin.Context)
	Other(ctx *gin.Context)
}

type repositories struct {
	main     []config.AllowedRepos
	other    []config.AllowedRepos
	provider github.RepoProvider
}

func NewRepositories(
	conf *config.Config,
	provider github.RepoProvider,
) Repositories {
	return &repositories{
		main:     conf.PinnedRepos.Main,
		other:    conf.PinnedRepos.Other,
		provider: provider,
	}
}

type ReposResponse struct {
	Repos []*github.Repository
}

func (h *repositories) Main(ctx *gin.Context) {
	var result []*github.Repository
	for _, a := range h.main {
		repos, err := h.provider.ByOwner(a.Owner, a.Repos)
		if err != nil {
			status(ctx, http.StatusInternalServerError)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		result = append(result, repos...)
	}

	// BUG: gin won't pass the status with JSON to the cache writer: https://gopkg.in/gin-gonic/gin.v1/pull/625
	status(ctx, http.StatusOK)
	ctx.JSON(http.StatusOK, &ReposResponse{Repos: result})
}

func (h *repositories) Other(ctx *gin.Context) {
	var result []*github.Repository
	for _, a := range h.other {
		for _, name := range a.Repos {
			repo, err := h.provider.ByOwnerAndName(a.Owner, name)
			if err != nil {
				status(ctx, http.StatusInternalServerError)
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			result = append(result, repo)
		}
	}

	status(ctx, http.StatusOK)
	ctx.JSON(http.StatusOK, &ReposResponse{Repos: result})
}