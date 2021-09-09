FROM golang:1.17.0-alpine3.14 as buildstage
WORKDIR /
COPY . .
RUN go get . && go build .

FROM alpine:3.14.2
WORKDIR /
COPY --from=buildstage aenorbot /

CMD ["./aenorbot"]