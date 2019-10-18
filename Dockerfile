FROM golang:alpine as builder
RUN	apk add --no-cache \
	ca-certificates
WORKDIR /src/app
COPY . .
RUN CGO_ENABLED=0 go build main.go

FROM scratch
COPY --from=builder /src/app/main /main
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
ENTRYPOINT [ "/main" ]
