FROM golang:1.17 as builder
ENV CGO_ENABLED 0

ARG BUILD_REF
ARG BUILD_DATE
ARG IMAGE_NAME

RUN mkdir /build
WORKDIR /build

COPY go.* .
RUN go mod download

COPY . .

RUN go build \
    		-ldflags "-X main.buildName=httpd-justforfun \
    				  -X main.buildRef=${BUILD_REF} \
    				  -X main.buildDate=${BUILD_DATE}" \
    		-o httpd.so cmd/httpd/main.go

FROM alpine:latest

RUN echo "Asia/Jakarta" > /etc/localtime

COPY --from=builder /build/httpd.so /


ENTRYPOINT ["/httpd.so"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="${IMAGE_NAME}" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.authors="${AUTHOR}" \
      org.opencontainers.image.source="${SOURCE}" \
      org.opencontainers.image.vendor="${VENDOR}"