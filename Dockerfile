FROM golang:1.12.0-alpine3.9
RUN apk add git
RUN mkdir -p /app/dbase
ADD server/* /app/
ADD dbase /app/dbase
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]