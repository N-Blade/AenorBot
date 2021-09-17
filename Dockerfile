FROM golang:alpine as buildstage
RUN mkdir aenorbot
WORKDIR /aenorbot
COPY . .
RUN go get ./... && go build -o aenorbot cmd/main.go

FROM alpine:3.14.2
WORKDIR /
COPY --from=buildstage aenorbot /

CMD ["./aenorbot"]