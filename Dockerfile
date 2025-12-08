FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /russgames ./cmd/server
EXPOSE 8000
CMD ["/russgames"]