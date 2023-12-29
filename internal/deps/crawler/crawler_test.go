package crawler

import (
	"context"
	"testing"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"
)

func TestCrawler(t *testing.T) {
	env, err := environment.Get()
	if err != nil {
		t.Fatal(err)
	}
	l := logger.New(env)
	c, err := New(env, l)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	goods, err := c.CollectGoods(ctx, GoodFromJinDong, "java", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(goods) == 0 {
		t.Errorf("length of goods should not be empty")
	}

	boards, err := c.CollectBoards(ctx, BoardFromBaidu, BoardTypeRealtime, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(boards) == 0 {
		t.Errorf("length of boards should not be empty")
	}
}
