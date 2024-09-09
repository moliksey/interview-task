# syntax=docker/dockerfile:1

FROM golang:1.23.1

WORKDIR /usr/local/src
RUN go version
ENV GOPATH=/
RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash
RUN apt-get update
RUN apt-get install migrate
RUN apt-get -y install postgresql-client
COPY ./ ./

RUN go mod download
RUN go build -o ./bin/app ./src/main.go
RUN chmod +x wait-for-postgres.sh
RUN chmod +x migrate.sh
CMD ["./bin/app"]
