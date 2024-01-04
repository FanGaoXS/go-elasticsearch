package server

import (
	"context"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/domain/boards"
	"fangaoxs.com/go-elasticsearch/internal/domain/goods"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"
	"fangaoxs.com/go-elasticsearch/server/rest"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func New(env environment.Env, logger logger.Logger) (*Server, error) {
	httpServer := gin.New()
	gin.ForceConsoleColor()
	httpServer.Use(gin.Logger())

	return initServer(env, logger, httpServer)
}

func newServer(
	env environment.Env,
	logger logger.Logger,
	httpServer *gin.Engine,
	goods goods.Goods,
	boards boards.Boards,
) (*Server, error) {
	restServer, err := rest.New(env, logger, httpServer, goods, boards)
	if err != nil {
		return nil, err
	}

	return &Server{
		env:        env,
		logger:     logger,
		restServer: restServer,
	}, nil
}

type Server struct {
	env    environment.Env
	logger logger.Logger

	restServer *rest.Server
}

func (s *Server) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		s.logger.Infof("rest server listen on %s", s.env.RestListenAddr)
		err := s.restServer.ListenAndServe()
		if err != nil {
			return err
		}
		s.logger.Info("rest server stopped")
		return nil
	})

	go func() {
		select {
		case <-ctx.Done():
			s.Close()
		}
	}()

	defer s.Close()

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Close() error {
	return s.restServer.Close()
}
