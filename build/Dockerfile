FROM alpine:3.8

RUN apk upgrade --update --no-cache

USER nobody

ADD build/_output/bin/golang-operator /usr/local/bin/golang-operator
