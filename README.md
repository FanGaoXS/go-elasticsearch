# go-elasticsearch
a sample demo for elasticsearch with golang.

## environment

```shell
cp env.example .env
```

## dependencies

```shell
github.com/elastic/go-elasticsearch/v8
github.com/gin-gonic/gin
github.com/gocolly/colly
github.com/google/wire
```

## search API

`GET localhost:8090/api/v1/search/goods?q={keyword}&type={searchType}&highlight={isHighlight}`
