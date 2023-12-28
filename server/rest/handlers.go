package rest

import (
	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/domain/goods"
	"fangaoxs.com/go-elasticsearch/internal/infras/errors"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"
	"net/http"
	"strconv"
	"strings"

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

		q := c.Query("q")
		if q = strings.TrimSpace(q); q == "" {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "invalid q: empty"))
			return
		}
		var page int64 = 1
		var size int64 = 10
		var err error
		if pageStr, ok := c.GetQuery("page"); ok {
			page, err = strconv.ParseInt(pageStr, 10, 64)
			if err != nil {
				WrapGinError(c, err)
				return
			}
		}
		if sizeStr, ok := c.GetQuery("size"); ok {
			size, err = strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				WrapGinError(c, err)
				return
			}
		}
		ctx := c.Request.Context()
		res, err := h.goods.SearchGoods(ctx, q, page, size)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
