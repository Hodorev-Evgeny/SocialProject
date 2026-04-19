FROM golang:1.26.2-bookworm

WORKDIR /app

COPY ../.. .

RUN go mod tidy
RUN go build -o /app/cmd/trackerapp/exe /app/cmd/trackerapp

CMD ["/app/exe"]