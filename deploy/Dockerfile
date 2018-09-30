FROM golang:latest

ADD . /app
WORKDIR /app

RUN go get -d -v ./...
RUN go build -o main .

ENV DATA_DIR /app/data

ADD docker_entrypoint.sh /docker_entrypoint.sh

EXPOSE 3000

RUN chmod a+x /docker_entrypoint.sh
ENTRYPOINT ["/docker_entrypoint.sh"]