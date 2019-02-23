FROM golang:latest

# RUN go get github.com/zodahu/bitmark-map
WORKDIR $GOPATH/src/github.com/zodahu/bitmark-map
COPY . $GOPATH/src/github.com/zodahu/bitmark-map

RUN go get github.com/gin-gonic/gin
RUN go get github.com/gin-gonic/contrib/static
RUN go get github.com/gin-contrib/cors
RUN go get github.com/boltdb/bolt
RUN go build .

EXPOSE 8080
ENTRYPOINT ["./bitmark-map"]