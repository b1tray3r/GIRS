FROM golang:1.23.1-alpine as builder

WORKDIR /app

ADD . /app/

RUN go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.2.0
RUN go run ./cmd/openapi/generator.go -path ./cmd/server/swagger/
RUN oapi-codegen --config=oapi-codegen.yml cmd/server/swagger/openapi3.json
RUN go mod tidy
RUN go build -o /app/server ./cmd/server/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server /app/server

EXPOSE 8080

CMD ["/app/server"]