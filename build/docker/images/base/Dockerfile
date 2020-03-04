FROM golang:1.13 AS builder

EXPOSE 8080

WORKDIR /go/src/2019_2_IBAT/

COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY ./go.mod .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth ./cmd/auth/auth.go 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o recommends ./cmd/recommends/recommends.go 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notifications ./cmd/notifications/notifications.go 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chat ./cmd/chat/chat.go 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main/main.go 
