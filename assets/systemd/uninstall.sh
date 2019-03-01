#!/bin/bash
set -e
#set -x

MYAPP=gobgp-exporter
MYAPP_USER=gobgp_exporter
MYAPP_GROUP=gobgp_exporter
MYAPP_SERVICE=${MYAPP}
MYAPP_BIN=/usr/bin/${MYAPP}
MYAPP_DESCRIPTION="Prometheus Exporter for Networking"
MYAPP_SYSCONF="/etc/sysconfig/${MYAPP_SERVICE}"

systemctl stop ${MYAPP_SERVICE}
systemctl disable ${MYAPP_SERVICE}
if systemctl is-active --quiet ${MYAPP_SERVICE}; then
  printf "FAIL: ${MYAPP_SERVICE} service is running\n"
  exit 1
else
  printf "INFO: ${MYAPP_SERVICE} service is not running\n"
fi

rm -rf ${MYAPP_BIN}

if [ -e ${MYAPP_SYSCONF} ]; then
  mv ${MYAPP_SYSCONF} ${MYAPP_SYSCONF}.`date +"%Y%m%d.%H%M%S"`
fi
