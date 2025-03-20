FROM golang:1.24

WORKDIR ./app

COPY . .

RUN go mod download && go build -o backend ./...

EXPOSE 8080

CMD ["./backend"]