package goods

import (
	"context"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/deps/crawler"
	"fangaoxs.com/go-elasticsearch/internal/deps/elasticsearch"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"
)

type Goods interface {
	SearchGoods(ctx context.Context, keyword string, pageNo, pageSize int64) ([]*Good, error)
}

func New(
	env environment.Env,
	logger logger.Logger,
	crawler crawler.Client,
	es elasticsearch.Client,
) (Goods, error) {
	ctx := context.Background()

	goods := make([]*Good, 0, 300)
	// init data
	for i := 0; i < 5; i++ {
		gds, err := crawler.CollectGoods("java", i+1)
		if err != nil {
			return nil, err
		}
		goods = append(goods, gds...)

		gds, err = crawler.CollectGoods("vue", i+1)
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

	es elasticsearch.Client
}

type Good = elasticsearch.Good

func (i *goodsImpl) SearchGoods(ctx context.Context, keyword string, pageNo, pageSize int64) ([]*Good, error) {
	return i.es.SearchGoods(ctx, keyword, int(pageNo), int(pageSize))
}
