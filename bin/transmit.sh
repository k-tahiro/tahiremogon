#!/bin/bash

readonly DATA_DIR="/sdcard/DCIM/Camera"

sudo /usr/local/bin/bto_ir_cmd -e -t "$1" >/dev/null 2>&1

adb shell touch "${DATA_DIR}/newer"
adb shell input keyevent 82
adb shell input keyevent 3
adb shell am start -n com.huawei.camera/com.huawei.camera
adb shell input keyevent 80
adb shell input keyevent 27

while :
do
  FILENAME=$(adb shell find ${DATA_DIR} -type f -newer ${DATA_DIR}/newer | grep jpg)
  if [ "${FILENAME}" != "" ]; then
    if [ $(echo "${FILENAME}" | wc -l) -eq 1 ]; then
      adb pull "${FILENAME}" "${IMG_DIR}"
      adb shell rm -f "${FILENAME}"
      break
    else
      {
        echo "Unexpected state!!"
        echo "There are too many files."
        echo "${FILENAME}"
      } 1>&2
      exit 1
    fi
  fi
  sleep 1
done
adb shell rm -f "${DATA_DIR}/newer"

basename "${FILENAME}"
