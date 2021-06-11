FROM golang:1.16-alpine3.13 as builder 
ENV GOPROXY=https://goproxy.io
WORKDIR /build
ADD . /build/
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o pscan /build/

FROM scratch
COPY --from=builder /build/pscan /
ENTRYPOINT ["/pscan"]
