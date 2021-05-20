FROM golang:alpine as builder 
RUN mkdir /build
WORKDIR /build  
ADD . /build/

RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o pscan .

FROM scratch
COPY --from=builder /build/pscan /
ENTRYPOINT ["/pscan"]
