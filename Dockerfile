FROM golang:1.20.2-alpine3.16 as build
WORKDIR /app
# Copy dependencies list
COPY go.mod go.sum ./
# Build
COPY main.go .
RUN go build -o main main.go
# Copy artifacts to a clean image
FROM alpine:3.16
COPY --from=build /app/main /main
ENTRYPOINT [ "/main" ]