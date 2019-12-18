# FROM golang:1.13 AS builder


# EXPOSE 8090

# WORKDIR /go/src/2019_2_IBAT/

# COPY ./cmd ./cmd
# COPY ./pkg ./pkg
# COPY ./go.mod .


# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/chat/chat.go 

FROM base:latest

EXPOSE 8090

ENTRYPOINT ["./chat"]
