FROM golang:1.14-alpine AS build

RUN apk update && apk upgrade && \
    apk add --no-cache git

ENV GO111MODULE=on \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . . 

RUN go build -a -installsuffix cgo -o ./out/api .

# After succesfully built the binary

FROM alpine:latest

RUN apk add ca-certificates

COPY --from=build /tmp/app/out/api /app/api

WORKDIR "/app"

EXPOSE 5000

ENTRYPOINT [ "./api" ]
# CMD [ "./api/main.go" ]