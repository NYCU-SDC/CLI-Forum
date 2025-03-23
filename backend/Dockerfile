FROM golang:1.24

WORKDIR ./app

COPY . .

RUN go mod tidy && go build -o backend ./cmd/backend/main.go

EXPOSE 8080

CMD ["./backend"]