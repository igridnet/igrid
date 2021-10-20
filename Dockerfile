FROM golang:1.16-alpine as builder
WORKDIR /igrid
COPY . .
RUN go clean -modcache
RUN go mod tidy
RUN go mod download
RUN cd cmd && go build -o igrid

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /igrid/
COPY --from=builder /igrid ./
ARG IGRID_SERVER_PORT
EXPOSE ${IGRID_SERVER_PORT}
ENTRYPOINT ["cmd/igrid"]
CMD ["--help"]