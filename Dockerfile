FROM golang:1.18.4 AS builder
WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /build/syno-pkg-restart .

###########################################################
# The *final* image

FROM gcr.io/distroless/static
COPY --from=builder /build/syno-pkg-restart /syno-pkg-restart
CMD ["/syno-pkg-restart"]