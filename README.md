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
...
```

## search API

```http request
GET localhost:8090/api/v1/search/goods/term?q={keyword}&page={page}&size={size}highlight={isHighlight}

GET localhost:8090/api/v1/search/goods/match?q={keyword}&page={page}&size={size}highlight={isHighlight}

GET localhost:8090/api/v1/search/boards/term?q={keyword}&page={page}&size={size}highlight={isHighlight}

GET localhost:8090/api/v1/search/boards/match?q={keyword}&page={page}&size={size}highlight={isHighlight}
```


## docker

### build docker image
```bash
make docker-build
```

### run docker image
```bash
docker run go-es
```

### run with docker-compose
```bash
make docker-build
docker-compose up
```