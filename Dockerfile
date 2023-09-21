FROM golang:1.18 AS builder
WORKDIR /go/src/better-admin-backend-service
COPY . .
RUN go mod download
RUN go install -ldflags '-w -extldflags "-static"'

# make application docker image use alpine
FROM alpine:3.10
# using timezone
ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Seoul
RUN apk add -U tzdata

WORKDIR /go/bin/
# copy config files to image
COPY --from=builder /go/src/better-admin-backend-service/config/*.json ./config/
COPY --from=builder /go/src/better-admin-backend-service/authorization ./authorization
# copy execute file to image
COPY --from=builder /go/bin/better-admin-backend-service .
EXPOSE 2016
CMD ["./better-admin-backend-service"]
