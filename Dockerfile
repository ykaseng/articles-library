FROM golang:latest

WORKDIR /go/src/articles-library
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["articles-library", "serve"]