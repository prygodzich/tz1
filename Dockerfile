# default value for build dir
FROM golang:1.23.2 AS builder

ENV CGO_ENABLED=0
ADD . /app/
WORKDIR /app

# Можно использовать кеширование зависимостей для ускорения сборок
# --mount=type=cache,target=/go/pkg/mod \
# --mount=type=cache,target=/root/.cache/go-build \
RUN make install-tools
RUN make lint
RUN make test
RUN go build -o targetads cmd/main.go

FROM  alpine:3.19.1
COPY --from=builder /app/targetads /app/targetads

WORKDIR /app
EXPOSE 8090/tcp 
CMD ["/app/targetads"]
