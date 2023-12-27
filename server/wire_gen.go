// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/deps/crawler"
	"fangaoxs.com/go-elasticsearch/internal/deps/elasticsearch"
	"fangaoxs.com/go-elasticsearch/internal/domain/goods"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func initServer(env environment.Env, logger2 logger.Logger, engine *gin.Engine) (*Server, error) {
	client, err := crawler.New(env, logger2)
	if err != nil {
		return nil, err
	}
	elasticsearchClient, err := elasticsearch.New(env, logger2)
	if err != nil {
		return nil, err
	}
	goodsGoods, err := goods.New(env, logger2, client, elasticsearchClient)
	if err != nil {
		return nil, err
	}
	server, err := newServer(env, logger2, engine, goodsGoods)
	if err != nil {
		return nil, err
	}
	return server, nil
}