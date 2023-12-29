package boards

import (
	"context"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/deps/crawler"
	es "fangaoxs.com/go-elasticsearch/internal/deps/elasticsearch"
	"fangaoxs.com/go-elasticsearch/internal/infras/errors"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"
)

type Board = crawler.Board

type Boards interface {
	SearchBoards(ctx context.Context, highlight bool, searchType es.SearchType, keyword string, pageNo, pageSize int64) ([]*Board, error)
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

func (i *boardsImpl) SearchBoards(ctx context.Context, highlight bool, searchType es.SearchType, keyword string, pageNo, pageSize int64) ([]*Board, error) {
	if searchType == es.SearchTypeInvalid {
		return nil, errors.Newf(errors.InvalidArgument, nil, "unsupported search type")
	}

	if searchType == es.SearchTypeTerm {
		return i.es.SearchBoardsByTerm(ctx, highlight, keyword, int(pageNo), int(pageSize))
	}

	return i.es.SearchBoardsByMatch(ctx, highlight, keyword, int(pageNo), int(pageSize))
}
