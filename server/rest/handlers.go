package rest

import (
	"net/http"
	"strconv"
	"strings"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/domain/boards"
	"fangaoxs.com/go-elasticsearch/internal/domain/goods"
	"fangaoxs.com/go-elasticsearch/internal/infras/errors"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func newHandlers(env environment.Env, logger logger.Logger, goods goods.Goods, boards boards.Boards) (*handlers, error) {
	return &handlers{
		env:    env,
		logger: logger,
		goods:  goods,
		boards: boards,
	}, nil
}

type handlers struct {
	env    environment.Env
	logger logger.Logger

	goods  goods.Goods
	boards boards.Boards
}

func (h *handlers) SearchGoodsByTerm() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET

		q := c.Query("q")
		if q = strings.TrimSpace(q); q == "" {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "invalid q: empty"))
			return
		}

		highlight := false
		if strings.ToLower(c.Query("highlight")) == "true" {
			highlight = true
		}

		page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
		if page == 0 {
			page = 1
		}
		size, _ := strconv.ParseInt(c.Query("size"), 10, 64)
		if size == 0 {
			size = 10
		}

		ctx := c.Request.Context()
		res, err := h.goods.SearchGoodsByTerm(ctx, highlight, q, page, size)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func (h *handlers) SearchGoodsByMatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET

		q := c.Query("q")
		if q = strings.TrimSpace(q); q == "" {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "invalid q: empty"))
			return
		}

		highlight := false
		if strings.ToLower(c.Query("highlight")) == "true" {
			highlight = true
		}

		page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
		if page == 0 {
			page = 1
		}
		size, _ := strconv.ParseInt(c.Query("size"), 10, 64)
		if size == 0 {
			size = 10
		}

		ctx := c.Request.Context()
		res, err := h.goods.SearchGoodsByMatch(ctx, highlight, q, page, size)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func (h *handlers) SearchBoardsByTerm() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET

		q := c.Query("q")
		if q = strings.TrimSpace(q); q == "" {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "invalid q: empty"))
			return
		}

		highlight := false
		if strings.ToLower(c.Query("highlight")) == "true" {
			highlight = true
		}

		page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
		if page == 0 {
			page = 1
		}
		size, _ := strconv.ParseInt(c.Query("size"), 10, 64)
		if size == 0 {
			size = 10
		}
		
		ctx := c.Request.Context()
		res, err := h.boards.SearchBoardsByTerm(ctx, highlight, q, page, size)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func (h *handlers) SearchBoardsByMatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET

		q := c.Query("q")
		if q = strings.TrimSpace(q); q == "" {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "invalid q: empty"))
			return
		}

		highlight := false
		if strings.ToLower(c.Query("highlight")) == "true" {
			highlight = true
		}

		page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
		if page == 0 {
			page = 1
		}
		size, _ := strconv.ParseInt(c.Query("size"), 10, 64)
		if size == 0 {
			size = 10
		}

		ctx := c.Request.Context()
		res, err := h.boards.SearchBoardsByMatch(ctx, highlight, q, page, size)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
