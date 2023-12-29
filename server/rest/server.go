package rest

import (
	"context"
	"net/http"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/domain/goods"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func New(
	env environment.Env,
	logger logger.Logger,
	router *gin.Engine,
	goods goods.Goods,
) (*Server, error) {
	h, err := newHandlers(env, logger, goods)
	if err != nil {
		return nil, err
	}

	v1 := router.Group("api/v1")
	c := v1.Group("search")
	{
		c.GET("goods", h.SearchGoods())
		c.GET("hotpots", h.SearchHotpots())
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
