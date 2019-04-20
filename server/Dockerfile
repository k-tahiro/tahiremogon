FROM balenalib/raspberry-pi2-golang

RUN apt update \
  && apt install autoconf libtool automake build-essential git libudev-dev wget \
  && wget https://github.com/libusb/libusb/releases/download/v1.0.22/libusb-1.0.22.tar.bz2 \
  && tar jxf libusb-1.0.22.tar.bz2 \
  && pushd libusb-1.0.22 \
  && ./configure \
  && make \
  && make install \
  && popd \
  && git clone https://github.com/kjmkznr/bto_ir_cmd.git \
  && pushd bto_ir_cmd \
  && make \
  && mv bto_ir_cmd /usr/local/bin/

RUN go get -u github.com/labstack/echo/... \
  && go get -u github.com/gocraft/dbr \
  && go get github.com/mattn/go-sqlite3

COPY handler handler
COPY middleware middleware
COPY model model
COPY server.go server.go
RUN go build server.go

CMD [ "./server" ]
