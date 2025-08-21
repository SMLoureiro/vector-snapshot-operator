# Simple Dockerfile for the manager
FROM golang:1.22 AS build
WORKDIR /workspace
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager ./cmd/manager

FROM gcr.io/distroless/base:latest
WORKDIR /
COPY --from=build /workspace/manager /manager
USER 65532:65532
ENTRYPOINT ["/manager"]
