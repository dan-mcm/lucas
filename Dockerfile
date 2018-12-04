FROM golang:1.11

MAINTAINER Daniel McMahon <daniel40392@gmail.com>

WORKDIR /opt/lucas

ADD . /opt/lucas

ENV PORT blergh

# installing our golang dependencies
RUN go get -u github.com/gocolly/colly && \
  go get -u github.com/fatih/color && \
  go get -u github.com/lib/pq

EXPOSE 8000

CMD go run lucas.go
