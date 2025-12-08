FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go env -w GOPROXY=https://proxy.golang.org,direct && \
    go build -mod=vendor -o /russgames ./backend-go/cmd/server

EXPOSE 8000
CMD ["/russgames"]
