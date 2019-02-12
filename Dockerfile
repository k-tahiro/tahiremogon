FROM resin/rpi-raspbian:jessie

ARG GOLANG_SRC="https://dl.google.com/go/go1.11.5.linux-armv6l.tar.gz"

RUN wget "${GOLANG_SRC}" \
  && tar -C /usr/local -xzf $(basename "${GOLANG_SRC}")
COPY bin bin
COPY handler handler
COPY middleware middleware
COPY model model
COPY server.go server.go
RUN go build server.go

ENTRYPOINT [ "server" ]
