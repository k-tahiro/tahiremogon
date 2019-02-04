#!/bin/bash

readonly CAMERA_APP="com.huawei.camera/com.huawei.camera"
readonly DATA_DIR="/sdcard/DCIM/Camera"
readonly LOG_DIR="/var/log/tahiremovision"
readonly ADB_LOG_FILE="${LOG_DIR}/adb.log"
readonly ERR_LOG_FILE="${LOG_DIR}/err.log"

IMG_DIR="${1:-"/var/opt/tahiremogon"}"

{
  adb shell touch "${DATA_DIR}/newer"
  adb shell input keyevent 82 # unlock
  adb shell am start -n "${CAMERA_APP}" # start camera app
  adb shell input keyevent 80 # forcus
  adb shell input keyevent 27 # release the shutter
  adb shell input keyevent 3 # back to home
  adb shell input keyevent 223 # sleep
} >>"${ADB_LOG_FILE}" 2>&1

while :
do
  FILENAME=$(adb shell find ${DATA_DIR} -type f -newer ${DATA_DIR}/newer | grep jpg)
  if [ "${FILENAME}" != "" ]; then
    if [ $(echo "${FILENAME}" | wc -l) -eq 1 ]; then
      adb pull "${FILENAME}" "${IMG_DIR}" >>"${ADB_LOG_FILE}" 2>&1
      adb shell rm -f "${FILENAME}" >>"${ADB_LOG_FILE}" 2>&1
      break
    else
      {
        echo "Unexpected state!!"
        echo "There are too many files."
        echo "${FILENAME}"
      } | tee -a "${ERR_LOG_FILE}" 1>&2
      exit 1
    fi
  fi
  sleep 1
done
adb shell rm -f "${DATA_DIR}/newer" >>"${ADB_LOG_FILE}" 2>&1

echo -n "${IMG_DIR}/${FILENAME}"
