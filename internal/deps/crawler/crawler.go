package crawler

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/infras/errors"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"

	"github.com/gocolly/colly"
)

type Client interface {
	CollectGoods(keyword string, pageNo int) ([]*Good, error)
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
}

func (c *crawlerImpl) CollectGoods(keyword string, pageNo int) ([]*Good, error) {
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
			ImageURI: "https:" + imageURI,
			Link:     "https:" + link,
			Title:    strings.ReplaceAll(strings.ReplaceAll(title, "\n", ""), "\t", ""),
			Shop:     shop,
			Price:    price,
			From:     GoodFromJinDong,
		}
		res = append(res, g)
	})

	h := http.Header{}
	//h.Set("Origin", "https://search.jd.com")
	//h.Set("Referer", "https://search.jd.com/")
	h.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	h.Set("Cookie", c.env.JDCookie)
	collector.Request(http.MethodGet, u.String(), nil, nil, h)

	return res, nil
}

type Board struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Index string `json:"index"`
	Link  string `json:"link"`

	Type BoardType `json:"type"`
	From BoardFrom `json:"from"`
}

func (c *crawlerImpl) CollectBoardBaidu(typ BoardType) ([]*Board, error) {
	panic("implement me")
}

type TreadingRepo struct{}

func (c *crawlerImpl) CollectTreadingRepoGithub() ([]*TreadingRepo, error) {
	panic("implement me")
}
