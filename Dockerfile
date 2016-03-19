FROM gliderlabs/alpine:edge

RUN apk-install go

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /go/src/github.com/convox/snqe
COPY . /go/src/github.com/convox/snqe
RUN go install ./...

CMD ["snqe"]