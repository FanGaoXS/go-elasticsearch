package elasticsearch

import (
	"context"
	"testing"

	"fangaoxs.com/go-elasticsearch/environment"
	"fangaoxs.com/go-elasticsearch/internal/entity"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/google/uuid"
)

func TestES(t *testing.T) {
	env, err := environment.Get()
	if err != nil {
		t.Error(err)
	}
	ESListenAddr := env.RestListenAddr
	Index := env.ESIndex

	config := es.Config{
		Addresses: []string{ESListenAddr},
	}
	c, err := es.NewTypedClient(config)
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()

	// 创建index
	_, err = c.Indices.Create(Index).Do(ctx)
	if err != nil {
		t.Error(err)
	}
	defer c.Indices.Delete(Index).Do(ctx)

	// 检查index是否存在
	ok, err := c.Indices.Exists(Index).Do(ctx)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fatal("index应该存在")
	}

	u := &entity.User{
		UUID: uuid.NewString(),
		Name: "狂神说Java",
		Age:  3,
		Desc: "我是一个讲Java的讲师",
		Tags: []string{"直男", "爱旅游", "打游戏", "java"},
	}

	// 创建doc
	_, err = c.Index(Index).Id(u.UUID).Request(u).Do(ctx)
	if err != nil {
		t.Error(err)
	}
	defer c.Delete(Index, u.UUID).Do(ctx)

	got, err := c.Get(Index, u.UUID).Do(ctx)
	if err != nil {
		t.Error(err)
	}
	if !got.Found {
		t.Fatalf("doc: %s 应该存在", u.UUID)
	}

	u = &entity.User{
		UUID: uuid.NewString(),
		Name: "狂神说Go",
		Age:  4,
		Desc: "我是一个讲Go的讲师",
		Tags: []string{"渣男", "爱唱歌", "跳舞", "go"},
	}

	// 创建doc
	_, err = c.Index(Index).Id(u.UUID).Request(u).Do(ctx)
	if err != nil {
		t.Error(err)
	}
	defer c.Delete(Index, u.UUID).Do(ctx)

	got, err = c.Get(Index, u.UUID).Do(ctx)
	if err != nil {
		t.Error(err)
	}
	if !got.Found {
		t.Fatalf("doc: %s 应该存在", u.UUID)
	}

	u = &entity.User{
		UUID: uuid.NewString(),
		Name: "狂神说C++",
		Age:  4,
		Desc: "我是一个讲C++的讲师",
		Tags: []string{"暖男", "爱跳舞", "rap", "c++"},
	}

	// 创建doc
	_, err = c.Index(Index).Id(u.UUID).Request(u).Do(ctx)
	if err != nil {
		t.Error(err)
	}
	defer c.Delete(Index, u.UUID).Do(ctx)

	got, err = c.Get(Index, u.UUID).Do(ctx)
	if err != nil {
		t.Error(err)
	}
	if !got.Found {
		t.Fatalf("doc: %s 应该存在", u.UUID)
	}

	// 查询所有
	r := &search.Request{
		Query: &types.Query{
			MatchAll: &types.MatchAllQuery{},
		},
	}

	res, err := c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// 指定字段值查询
	r = &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"tags": {Query: "男"},
			},
		},
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// 返回指定字段
	r = &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"tags": {Query: "男"},
			},
		},
		Source_: []string{"name", "desc"},
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// 查询排序
	r = &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"tags": {Query: "男"},
			},
		},
		Sort: []types.SortCombinations{
			map[string]string{"age": "asc"},
		},
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// 查询分页
	from := 0
	size := 2
	r = &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"tags": {Query: "男"},
			},
		},
		From: &from,
		Size: &size,
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// and 查询
	r = &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{Match: map[string]types.MatchQuery{
						"name": {Query: "go"},
					}},
					{Match: map[string]types.MatchQuery{
						"age": {Query: "4"},
					}},
				},
			},
		},
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// or 查询
	r = &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Should: []types.Query{
					{Match: map[string]types.MatchQuery{
						"name": {Query: "go"},
					}},
					{Match: map[string]types.MatchQuery{
						"age": {Query: "4"},
					}},
				},
			},
		},
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// not 查询
	r = &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				MustNot: []types.Query{
					{Match: map[string]types.MatchQuery{
						"name": {Query: "go"},
					}},
					{Match: map[string]types.MatchQuery{
						"age": {Query: "4"},
					}},
				},
			},
		},
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// filter
	var gte types.Float64 = 3
	var lt types.Float64 = 5
	r = &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{Match: map[string]types.MatchQuery{
						"name": {Query: "狂神"},
					}},
				},
				Filter: []types.Query{
					{Range: map[string]types.RangeQuery{
						"age": types.NumberRangeQuery{
							Gte: &gte,
							Lt:  &lt,
						},
					}},
				},
			},
		},
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// 模糊查询
	r = &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"tags": {Query: "rap go"},
			},
		},
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// 分词查询（倒排索引）
	r = &search.Request{
		Query: &types.Query{
			Term: map[string]types.TermQuery{
				"tags": {Value: "rap"},
			},
		},
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// 查询结果高亮
	r = &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"tags": {Query: "rap"},
			},
		},
		Highlight: &types.Highlight{
			Fields: map[string]types.HighlightField{
				"tags": {},
			},
		},
	}

	res, err = c.Search().Index(Index).Request(r).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_ = res
}
