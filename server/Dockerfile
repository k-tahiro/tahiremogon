FROM balenalib/raspberry-pi2-golang

RUN go get -u github.com/labstack/echo/... \
  && go get -u github.com/gocraft/dbr \
  && go get github.com/mattn/go-sqlite3

COPY handler handler
COPY middleware middleware
COPY model model
COPY server.go server.go
RUN go build server.go

ARG DB_FILE
COPY "${DB_FILE}" ./

CMD [ "./server" ]
