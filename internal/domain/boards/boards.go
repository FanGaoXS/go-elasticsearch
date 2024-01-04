package boards

import (
	"context"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/deps/crawler"
	es "fangaoxs.com/go-elasticsearch/internal/deps/elasticsearch"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"
)

type Board = es.Board

type Boards interface {
	SearchBoardsByTerm(ctx context.Context, highlight bool, keyword string, pageNo, pageSize int64) ([]*Board, error)
	SearchBoardsByMatch(ctx context.Context, highlight bool, keyword string, pageNo, pageSize int64) ([]*Board, error)
}

func New(
	env environment.Env,
	logger logger.Logger,
	c crawler.Client,
	es es.Client,
) (Boards, error) {
	ctx := context.Background()

	boards, err := c.CollectBoards(ctx, crawler.BoardFromBaidu, crawler.BoardTypeRealtime, 0)
	if err != nil {
		return nil, err
	}

	if err = es.InsertBoards(ctx, boards); err != nil {
		return nil, err
	}
	logger.Infof("init %d boards into elasticsearch success", len(boards))

	return &boardsImpl{
		env:    env,
		logger: logger,
		es:     es,
	}, nil
}

type boardsImpl struct {
	env    environment.Env
	logger logger.Logger

	es es.Client
}

func (i *boardsImpl) SearchBoardsByTerm(ctx context.Context, highlight bool, keyword string, pageNo, pageSize int64) ([]*Board, error) {
	return i.es.SearchBoardsByTerm(ctx, highlight, keyword, int(pageNo), int(pageSize))

}

func (i *boardsImpl) SearchBoardsByMatch(ctx context.Context, highlight bool, keyword string, pageNo, pageSize int64) ([]*Board, error) {
	return i.es.SearchBoardsByMatch(ctx, highlight, keyword, int(pageNo), int(pageSize))
}
