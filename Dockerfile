FROM golang:alpine as builder 
RUN mkdir /build
WORKDIR /build  
ADD . /build/
RUN go build -o pscan main.go

FROM alpine 
COPY --from=builder /build/pscan /
ENTRYPOINT ["/pscan"]
