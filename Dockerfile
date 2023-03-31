FROM golang:1.19-alpine AS builder
RUN apk add --no-cache ca-certificates git

WORKDIR /opt/apps

COPY go.mod go.sum ./

RUN go mod download
COPY . . 
RUN go build -o /form-builder-be .

FROM alpine AS release
RUN apk add --no-cache ca-certificates

WORKDIR /form-builder-be
COPY --from=builder /form-builder-be ./main
ENTRYPOINT ["/form-builder-be/main"]