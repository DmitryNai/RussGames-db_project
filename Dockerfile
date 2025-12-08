
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go env -w GOPROXY=https://proxy.golang.org,direct && \
    go mod download

COPY . .

RUN cd backend-go && go build -o /russgames ./cmd/server

EXPOSE 8000

CMD ["/russgames"]
