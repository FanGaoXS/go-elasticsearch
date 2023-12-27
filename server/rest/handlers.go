package rest

import (
	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/domain/goods"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func newHandlers(env environment.Env, logger logger.Logger, goods goods.Goods) (*handlers, error) {
	return &handlers{
		env:    env,
		logger: logger,
		goods:  goods,
	}, nil
}

type handlers struct {
	env    environment.Env
	logger logger.Logger

	goods goods.Goods
}

func (h *handlers) SearchGoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET

	}
}
