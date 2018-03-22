
FROM golang:1.10.0 as builder
ARG app
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV SRC_DIR=/go/src/github.com/timhugh/ticket_service
WORKDIR $SRC_DIR
COPY . .
RUN cd cmd/$app && go build -a -o /tmp/app


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /tmp/app .
CMD ["./app"]
