FROM golang:1.17-alpine

WORKDIR $GOPATH/src/github.com/xacnio/aptms-backend

COPY . .

RUN go get -d -v ./...

RUN go build

EXPOSE 80

CMD ["./aptms-backend"]