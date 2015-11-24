FROM alpine:3.2

EXPOSE 80

RUN echo "@community http://dl-4.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk update
RUN apk add bash git go@community
ADD . /magicbeans

ENV GOPATH /magicbeans
WORKDIR /magicbeans
RUN go build src/magicbeans.go

CMD ["/magicbeans/magicbeans"]
