FROM golang:1.20 as Builder
WORKDIR /source
COPY . /source
RUN go build -o app .

FROM ubuntu:noble
RUN apt update && apt install -y ca-certificates
COPY --from=builder /source/app /app
RUN chmod +x /app

ENV APP_NAME="app"
ENV APP_VERSION="v0.0.1"
ENV LOG_LEVEL="INFO"
ENV REST_LISTEN_ADDR="0.0.0.0:8090"
ENV ES_REST_ADDR="http://localhost:9200"
ENV ES_GOODS_INDEX="goods"
ENV ES_BOARDS_INDEX="boards"
ENV JD_COOKIE="NULL"
ENV TB_COOKIE="NULL"

CMD ["/app"]