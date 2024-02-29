FROM golang:latest

RUN go version

COPY ./ ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go
RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash
RUN apt-get update
RUN apt-get install migrate

EXPOSE 8081

CMD ["./.bin/app"]