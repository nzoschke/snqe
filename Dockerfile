FROM gliderlabs/alpine:edge

RUN apk-install ca-certificates go

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /go/src/github.com/convox/snqe
COPY . /go/src/github.com/convox/snqe
RUN go install ./...

CMD ["snqe"]