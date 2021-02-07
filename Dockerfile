FROM golang

ADD . /go/src/

WORKDIR /app

COPY . /app

EXPOSE 4000

RUN go build -o app cmd/web/*

ENTRYPOINT  /app/app