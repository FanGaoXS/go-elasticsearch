package elasticsearch

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/deps/crawler"
	"fangaoxs.com/go-elasticsearch/internal/infras/errors"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"

	es "github.com/elastic/go-elasticsearch/v8"
)

type Client interface {
	InsertGoods(ctx context.Context, goods []*Good) error
}

func New(env environment.Env, logger logger.Logger) (Client, error) {
	addr := env.ESRestAddr
	index := env.ESIndex

	config := es.Config{
		Addresses: []string{addr},
	}
	c, err := es.NewTypedClient(config)
	if err != nil {
		return nil, errors.New(errors.Internal, nil, "init typed elasticsearch client failed")
	}

	ctx := context.Background()
	ok, err := c.Indices.Exists(index).Do(ctx)
	if err != nil {
		return nil, errors.Newf(errors.Internal, err, "get whether index: %s exists failed", index)
	}
	if ok {
		if _, err = c.Indices.Delete(index).Do(ctx); err != nil {
			return nil, errors.Newf(errors.Internal, err, "delete index: %s failed", index)
		}
	}
	if _, err = c.Indices.Create(index).Do(ctx); err != nil {
		return nil, errors.Newf(errors.Internal, err, "create index: %s failed", index)
	}

	return &esImpl{
		env:    env,
		logger: logger,
		es:     c,
		index:  index,
	}, nil
}

type esImpl struct {
	env    environment.Env
	logger logger.Logger

	es    *es.TypedClient
	index string
}

type Good = crawler.Good

func (e *esImpl) InsertGoods(ctx context.Context, goods []*Good) error {
	for _, g := range goods {
		if _, err := e.es.Index(e.index).Request(g).Do(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (e *esImpl) SearchGoods(ctx context.Context) ([]*Good, error) {
	res := make([]*Good, 0)

	r := &search.Request{}
	resp, err := e.es.Search().Index(e.index).Request(r).Do(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
