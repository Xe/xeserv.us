FROM golang:1.4.2

ENV VIRTUAL_HOST xeserv.us

ADD . /go/src/github.com/Xe/xeserv.us
RUN go get github.com/Xe/xeserv.us/...

WORKDIR /go/src/github.com/Xe/xeserv.us
CMD xeserv.us

EXPOSE 3000
