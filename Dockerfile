FROM golang:alpine AS builder
ADD . /go/src/
RUN apk add git
RUN git config --global core.autocrlf false
RUN go get github.com/gorilla/mux
WORKDIR /go/src
RUN go build -o main .

FROM alpine
COPY --from=builder /go/src/main /app/
WORKDIR /app
CMD ["./main"]