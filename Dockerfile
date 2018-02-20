FROM golang:1.9.4-alpine3.7 as build
WORKDIR /go/src/app
COPY . .
RUN go build -o main
# when switching to scratch
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main


# TODO - switch to alpine
# FROM scratch
FROM alpine:3.7
WORKDIR /app
# add certificates
RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/app/main .
COPY conf.json .
CMD ["./main"]
