package elasticsearch

import (
	"context"
	"encoding/json"
	"strings"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/deps/crawler"
	"fangaoxs.com/go-elasticsearch/internal/infras/errors"
	"fangaoxs.com/go-elasticsearch/internal/infras/logger"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type SearchType int

const (
	SearchTypeInvalid SearchType = iota
	SearchTypeTerm
	SearchTypeMatch
)

type Good = crawler.Good
type Board = crawler.Board

type Client interface {
	InsertGoods(ctx context.Context, goods []*Good) error
	SearchGoodsByTerm(ctx context.Context, isHighlight bool, keyword string, pageNo, pageSize int) ([]*Good, error)
	SearchGoodsByMatch(ctx context.Context, isHighlight bool, keyword string, pageNo, pageSize int) ([]*Good, error)

	InsertBoards(ctx context.Context, boards []*Board) error
	SearchBoardsByTerm(ctx context.Context, isHighlight bool, keyword string, pageNo, pageSize int) ([]*Board, error)
	SearchBoardsByMatch(ctx context.Context, isHighlight bool, keyword string, pageNo, pageSize int) ([]*Board, error)
}

func New(env environment.Env, logger logger.Logger) (Client, error) {
	addr := env.ESRestAddr

	config := es.Config{
		Addresses: []string{addr},
	}
	client, err := es.NewTypedClient(config)
	if err != nil {
		return nil, errors.New(errors.Internal, nil, "init typed elasticsearch client failed")
	}

	e := &esImpl{
		env:    env,
		logger: logger,
		es:     client,
	}
	ctx := context.Background()
	if strings.EqualFold(env.ESGoodsIndex, env.ESBoardsIndex) {
		return nil, errors.Newf(errors.Internal, nil, "goods index: %s is the same as board index: %s", env.ESGoodsIndex, env.ESBoardsIndex)
	}
	if err = e.initGoodsIndex(ctx); err != nil {
		return nil, err
	}
	if err = e.initBoardsIndex(ctx); err != nil {
		return nil, err
	}

	return e, nil
}

type esImpl struct {
	env    environment.Env
	logger logger.Logger

	es *es.TypedClient
}

func (e *esImpl) initGoodsIndex(ctx context.Context) error {
	index := e.env.ESGoodsIndex

	ok, err := e.es.Indices.Exists(index).Do(ctx)
	if err != nil {
		return errors.Newf(errors.Internal, err, "get whether index: %s exists failed", index)
	}
	if ok {
		if _, err = e.es.Indices.Delete(index).Do(ctx); err != nil {
			return errors.Newf(errors.Internal, err, "delete index: %s failed", index)
		}
	}

	var (
		analyzer = "ik_smart"
		mappings = &types.TypeMapping{
			Properties: map[string]types.Property{
				"title": types.TextProperty{Analyzer: &analyzer},
			},
		} // analyzer: ik_smart on title field
	)

	if _, err = e.es.Indices.Create(index).Mappings(mappings).Do(ctx); err != nil {
		return errors.Newf(errors.Internal, err, "create index: %s failed", index)
	}

	return nil
}

func (e *esImpl) InsertGoods(ctx context.Context, goods []*Good) error {
	for _, g := range goods {
		if _, err := e.es.Index(e.env.ESGoodsIndex).Request(g).Do(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (e *esImpl) SearchGoodsByTerm(ctx context.Context, isHighlight bool, keyword string, pageNo, pageSize int) ([]*Good, error) {
	var hl *types.Highlight
	if isHighlight {
		hl = &types.Highlight{
			Fields: map[string]types.HighlightField{
				"title": {},
			},
		}
	}
	from := (pageNo - 1) * pageSize
	r := &search.Request{
		Query: &types.Query{
			Term: map[string]types.TermQuery{
				"title": {
					Value: keyword,
				},
			},
		},
		From:      &from,
		Size:      &pageSize,
		Highlight: hl,
	}
	resp, err := e.es.Search().Index(e.env.ESGoodsIndex).Request(r).Do(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*Good, 0, pageSize)
	for _, hit := range resp.Hits.Hits {
		var g Good
		if err = json.Unmarshal(hit.Source_, &g); err != nil {
			return nil, err
		}

		if isHighlight {
			g.Title = hit.Highlight["title"][0]

		}

		res = append(res, &g)
	}

	return res, nil
}

func (e *esImpl) SearchGoodsByMatch(ctx context.Context, isHighlight bool, keyword string, pageNo, pageSize int) ([]*Good, error) {
	var hl *types.Highlight
	if isHighlight {
		hl = &types.Highlight{
			Fields: map[string]types.HighlightField{
				"title": {},
			},
		}
	}
	from := (pageNo - 1) * pageSize
	r := &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"title": {
					Query: keyword,
				},
			},
		},
		From:      &from,
		Size:      &pageSize,
		Highlight: hl,
	}
	resp, err := e.es.Search().Index(e.env.ESGoodsIndex).Request(r).Do(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*Good, 0, pageSize)
	for _, hit := range resp.Hits.Hits {
		var g Good
		if err = json.Unmarshal(hit.Source_, &g); err != nil {
			return nil, err
		}

		if isHighlight {
			g.Title = hit.Highlight["title"][0]

		}

		res = append(res, &g)
	}

	return res, nil
}

func (e *esImpl) initBoardsIndex(ctx context.Context) error {
	index := e.env.ESBoardsIndex

	ok, err := e.es.Indices.Exists(index).Do(ctx)
	if err != nil {
		return errors.Newf(errors.Internal, err, "get whether index: %s exists failed", index)
	}
	if ok {
		if _, err = e.es.Indices.Delete(index).Do(ctx); err != nil {
			return errors.Newf(errors.Internal, err, "delete index: %s failed", index)
		}
	}

	var (
		analyzer = "ik_smart"
		mappings = &types.TypeMapping{
			Properties: map[string]types.Property{
				"title": types.TextProperty{Analyzer: &analyzer},
			},
		} // analyzer: ik_smart on title field
	)

	if _, err = e.es.Indices.Create(index).Mappings(mappings).Do(ctx); err != nil {
		return errors.Newf(errors.Internal, err, "create index: %s failed", index)
	}

	return nil
}

func (e *esImpl) InsertBoards(ctx context.Context, boards []*Board) error {
	for _, b := range boards {
		if _, err := e.es.Index(e.env.ESBoardsIndex).Request(b).Do(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (e *esImpl) SearchBoardsByTerm(ctx context.Context, isHighlight bool, keyword string, pageNo, pageSize int) ([]*Board, error) {
	var hl *types.Highlight
	if isHighlight {
		hl = &types.Highlight{
			Fields: map[string]types.HighlightField{
				"title": {},
			},
		}
	}
	from := (pageNo - 1) * pageSize
	r := &search.Request{
		Query: &types.Query{
			Term: map[string]types.TermQuery{
				"title": {
					Value: keyword,
				},
			},
		},
		From:      &from,
		Size:      &pageSize,
		Highlight: hl,
	}
	resp, err := e.es.Search().Index(e.env.ESBoardsIndex).Request(r).Do(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*Board, 0, pageSize)
	for _, hit := range resp.Hits.Hits {
		var b Board
		if err = json.Unmarshal(hit.Source_, &b); err != nil {
			return nil, err
		}

		if isHighlight {
			b.Title = hit.Highlight["title"][0]

		}

		res = append(res, &b)
	}

	return res, nil
}

func (e *esImpl) SearchBoardsByMatch(ctx context.Context, isHighlight bool, keyword string, pageNo, pageSize int) ([]*Board, error) {
	var hl *types.Highlight
	if isHighlight {
		hl = &types.Highlight{
			Fields: map[string]types.HighlightField{
				"title": {},
			},
		}
	}
	from := (pageNo - 1) * pageSize
	r := &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"title": {
					Query: keyword,
				},
			},
		},
		From:      &from,
		Size:      &pageSize,
		Highlight: hl,
	}
	resp, err := e.es.Search().Index(e.env.ESBoardsIndex).Request(r).Do(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*Board, 0, pageSize)
	for _, hit := range resp.Hits.Hits {
		var b Board
		if err = json.Unmarshal(hit.Source_, &b); err != nil {
			return nil, err
		}

		if isHighlight {
			b.Title = hit.Highlight["title"][0]

		}

		res = append(res, &b)
	}

	return res, nil
}
