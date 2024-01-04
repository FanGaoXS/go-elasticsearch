package rest

import (
	"context"
	"net/http"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/domain/boards"
	"fangaoxs.com/go-elasticsearch/internal/domain/goods"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func New(
	env environment.Env,
	logger logger.Logger,
	router *gin.Engine,
	goods goods.Goods,
	boards boards.Boards,
) (*Server, error) {
	h, err := newHandlers(env, logger, goods, boards)
	if err != nil {
		return nil, err
	}

	v1 := router.Group("api/v1")
	c := v1.Group("search")
	{
		c.GET("goods/term", h.SearchGoodsByTerm())
		c.GET("goods/match", h.SearchGoodsByMatch())
		c.GET("boards/term", h.SearchBoardsByTerm())
		c.GET("boards/match", h.SearchBoardsByMatch())
	}

	s := &http.Server{
		Addr:    env.RestListenAddr,
		Handler: router,
	}
	return &Server{
		server: s,
	}, nil
}

type Server struct {
	server *http.Server
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	return s.server.Shutdown(context.Background())
}
