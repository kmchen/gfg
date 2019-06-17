FROM golang:1.11

LABEL maintainer="Kuang-Ming Chen <kuangmingchen0702@gmail.com>"

WORKDIR $GOPATH/src/github.com/gfg

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o gfg .

EXPOSE 8000

# Run the executable
CMD ["gfg"]
