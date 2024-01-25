package goods

import (
	"context"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/deps/crawler"
	es "fangaoxs.com/go-elasticsearch/internal/deps/elasticsearch"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"
)

type Good = es.Good

type Goods interface {
	SearchGoodsByTerm(ctx context.Context, highlight bool, keyword string, pageNo, pageSize int64) ([]*Good, error)
	SearchGoodsByMatch(ctx context.Context, highlight bool, keyword string, pageNo, pageSize int64) ([]*Good, error)
}

func New(
	env environment.Env,
	logger logger.Logger,
	c crawler.Client,
	es es.Client,
) (Goods, error) {
	ctx := context.Background()

	logger.Info("collecting goods from JD...")
	goods := make([]*Good, 0)
	// init data
	for i := 0; i < 1; i++ {
		gds, err := c.CollectGoods(ctx, crawler.GoodFromJinDong, "java", i+1)
		if err != nil {
			return nil, err
		}
		goods = append(goods, gds...)

		gds, err = c.CollectGoods(ctx, crawler.GoodFromJinDong, "vue", i+1)
		if err != nil {
			return nil, err
		}
		goods = append(goods, gds...)
	}

	if err := es.InsertGoods(ctx, goods); err != nil {
		return nil, err
	}
	logger.Infof("init %d goods into elasticsearch success", len(goods))

	return &goodsImpl{
		env:    env,
		logger: logger,
		es:     es,
	}, nil
}

type goodsImpl struct {
	env    environment.Env
	logger logger.Logger

	es es.Client
}

func (i *goodsImpl) SearchGoodsByTerm(ctx context.Context, highlight bool, keyword string, pageNo, pageSize int64) ([]*Good, error) {
	return i.es.SearchGoodsByTerm(ctx, highlight, keyword, int(pageNo), int(pageSize))
}

func (i *goodsImpl) SearchGoodsByMatch(ctx context.Context, highlight bool, keyword string, pageNo, pageSize int64) ([]*Good, error) {
	return i.es.SearchGoodsByMatch(ctx, highlight, keyword, int(pageNo), int(pageSize))
}
