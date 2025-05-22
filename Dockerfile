FROM golang:latest AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/main.go


FROM alpine:latest

COPY --from=builder /main /bin/main

CMD ["/bin/main"]