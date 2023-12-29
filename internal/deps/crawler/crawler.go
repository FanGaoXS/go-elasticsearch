package crawler

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/infras/errors"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"

	"github.com/gocolly/colly"
)

type Client interface {
	CollectGoods(ctx context.Context, from GoodFrom, keyword string, pageNo int) ([]*Good, error)
	CollectBoards(ctx context.Context, from BoardFrom, typ BoardType, pageNo int) ([]*Board, error)
}

func New(env environment.Env, logger logger.Logger) (Client, error) {
	return &crawlerImpl{
		env:    env,
		logger: logger,
	}, nil
}

type crawlerImpl struct {
	env    environment.Env
	logger logger.Logger
}

type Good struct {
	ImageURI string   `json:"image_uri"`
	Link     string   `json:"link"`
	Title    string   `json:"title"`
	Shop     string   `json:"shop"`
	Price    string   `json:"price"`
	From     GoodFrom `json:"from"`

	CreatedAt time.Time `json:"created_at"`
}

func (c *crawlerImpl) CollectGoods(ctx context.Context, from GoodFrom, keyword string, pageNo int) ([]*Good, error) {
	if from == GoodFromInvalid {
		return nil, errors.New(errors.InvalidArgument, nil, "invalid good from")
	}
	if from == GoodFromTaoBao {
		return nil, errors.New(errors.InvalidArgument, nil, "invalid good from: unsupported taobao")
	}

	return c.collectGoodsFromJD(ctx, keyword, pageNo)
}

func (c *crawlerImpl) collectGoodsFromJD(ctx context.Context, keyword string, pageNo int) ([]*Good, error) {
	if pageNo <= 0 {
		return nil, errors.Newf(errors.InvalidArgument, nil, "invalid page number: %d", pageNo)
	}

	collector := colly.NewCollector(
		colly.MaxDepth(1),
	)

	u, _ := url.Parse("search.jd.com/Search")

	values := u.Query()
	values.Set("keyword", keyword)
	values.Set("page", strconv.Itoa(pageNo))
	u.RawQuery = values.Encode()
	u.Scheme = "https"

	res := make([]*Good, 0)
	collector.OnHTML("li[class=gl-item]", func(e *colly.HTMLElement) {
		imageURI := e.ChildAttr("div[class=p-img] > a > img", "data-lazy-img")
		link := e.ChildAttr("div[class=p-img] > a", "href")
		title := strings.TrimSpace(e.ChildText("div[class^=p-name]"))
		shop := strings.TrimSpace(e.ChildText("div[class^=p-shop]"))
		price := e.ChildText("div[class=p-price]")

		g := &Good{
			ImageURI:  "https:" + imageURI,
			Link:      "https:" + link,
			Title:     strings.ReplaceAll(strings.ReplaceAll(title, "\n", ""), "\t", ""),
			Shop:      shop,
			Price:     price,
			From:      GoodFromJinDong,
			CreatedAt: time.Now(),
		}
		res = append(res, g)
	})

	h := http.Header{}
	h.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	h.Set("Cookie", c.env.JDCookie)
	if err := collector.Request(http.MethodGet, u.String(), nil, nil, h); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *crawlerImpl) collectGoodsFromTaobao(ctx context.Context, keyword string, pageNo int) ([]*Good, error) {
	panic("implement me")
}

type Board struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Index int    `json:"index"`
	Link  string `json:"link"`

	Type      BoardType `json:"type"`
	From      BoardFrom `json:"from"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *crawlerImpl) CollectBoards(ctx context.Context, from BoardFrom, typ BoardType, pageNo int) ([]*Board, error) {
	if typ == BoardTypeInvalid {
		return nil, errors.New(errors.InvalidArgument, nil, "unsupported board type")
	}

	return c.collectBoardsFromBaidu(ctx, typ, pageNo)
}

func (c *crawlerImpl) collectBoardsFromBaidu(ctx context.Context, typ BoardType, pageNo int) ([]*Board, error) {
	collector := colly.NewCollector(
		colly.MaxDepth(1),
	)

	u, _ := url.Parse("top.baidu.com/board")

	values := u.Query()
	values.Set("tab", ToBoardTypeString(typ))
	//values.Set("page", strconv.Itoa(pageNo))
	u.RawQuery = values.Encode()
	u.Scheme = "https"

	res := make([]*Board, 0)
	collector.OnHTML("div[class^=category]", func(e *colly.HTMLElement) {
		link := e.ChildAttr("div[class^=content] > a", "href")
		title := strings.TrimSpace(e.ChildText("div[class^=content] div[class^=c-single-text]"))
		desc := strings.TrimSpace(e.ChildText("div[class^=content] div[class*=large]"))
		desc = strings.ReplaceAll(desc, "查看更多>", "")
		indexStr := strings.TrimSpace(e.ChildText("div[class^=hot-index]"))
		index, _ := strconv.Atoi(indexStr)

		res = append(res, &Board{
			Title:     title,
			Desc:      desc,
			Index:     index,
			Link:      link,
			Type:      BoardTypeRealtime,
			From:      BoardFromBaidu,
			CreatedAt: time.Now(),
		})
	})

	h := http.Header{}
	if err := collector.Request(http.MethodGet, u.String(), nil, nil, h); err != nil {
		return nil, err
	}
	return res, nil
}

type TreadingRepo struct{}

func (c *crawlerImpl) CollectTreadingRepoGithub() ([]*TreadingRepo, error) {
	panic("implement me")
}
