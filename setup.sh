#!/bin/bash

echo "Start build"
go build
echo "End build"

echo "Start install"
sudo cp -rp bin/* /usr/local/bin/
echo "End install"

cat <<__EOF__
Start server using command below

  IMG_DIR=/var/opt/tahiremogon/image nohup ./tahiremogon &
__EOF__
