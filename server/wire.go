//go:build wireinject
// +build wireinject

package server

import (
	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/deps/crawler"
	"fangaoxs.com/go-elasticsearch/internal/deps/elasticsearch"
	"fangaoxs.com/go-elasticsearch/internal/domain/boards"
	"fangaoxs.com/go-elasticsearch/internal/domain/goods"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func initServer(env environment.Env, logger logger.Logger, engine *gin.Engine) (*Server, error) {
	panic(wire.Build(
		crawler.New,
		elasticsearch.New,
		goods.New,
		boards.New,
		newServer,
	))
}
