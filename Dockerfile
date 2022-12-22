FROM golang:1.19-alpine AS build
#RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:3.17.0
#FROM scratch
WORKDIR /dist

COPY --from=build /app/main ./
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
COPY --from=build /app/assets ./assets/
#COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV ZONEINFO=/zoneinfo.zip

ENTRYPOINT ["./main"]

