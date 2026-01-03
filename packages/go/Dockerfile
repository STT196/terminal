# ./packages/go/Dockerfile
FROM golang:1.24-alpine AS build
WORKDIR /app
# install build deps if needed: apk add --no-cache build-base
COPY go.mod go.sum ./ 
RUN go mod download
COPY . .
RUN go build -o /out/terminal-ssh ./cmd/ssh

FROM alpine:3.19
WORKDIR /app
RUN adduser -S -D terminal
COPY --from=build /out/terminal-ssh /usr/local/bin/terminal-ssh
# env vars (adjust ports as needed)
ENV PORT=8080
EXPOSE 8080
USER terminal
ENTRYPOINT ["/usr/local/bin/terminal-ssh"]