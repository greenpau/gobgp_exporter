#!/bin/bash
set -e
set -x

MYAPP=gobgp-exporter
MYAPP_USER=gobgp_exporter
MYAPP_GROUP=gobgp_exporter
MYAPP_SERVICE=gobgp-exporter
MYAPP_BIN=/usr/bin/gobgp-exporter
MYAPP_DESCRIPTION="Prometheus Exporter for GoBGP"
MYAPP_CONF="/etc/sysconfig/${MYAPP_SERVICE}"

if getent group ${MYAPP_GROUP}  >/dev/null; then
  printf "INFO: ${MYAPP_GROUP} group already exists\n"
else
  printf "INFO: ${MYAPP_GROUP} group does not exist, creating ..."
  groupadd --system ${MYAPP_GROUP}
fi

if getent passwd ${MYAPP_USER} >/dev/null; then
  printf "INFO: ${MYAPP_USER} user already exists\n"
else
  printf "INFO: ${MYAPP_USER} group does not exist, creating ..."
  useradd --system -d /var/lib/${MYAPP} -s /bin/bash -g ${MYAPP_GROUP} ${MYAPP_USER}
fi

mkdir -p /var/lib/${MYAPP}
chown -R ${MYAPP_USER}:${MYAPP_GROUP} /var/lib/${MYAPP}

cat << EOF > /usr/lib/systemd/system/${MYAPP_SERVICE}.service
[Unit]
Description=$MYAPP_DESCRIPTION
After=network.target

[Service]
User=${MYAPP_USER}
Group=${MYAPP_GROUP}
EnvironmentFile=-${MYAPP_CONF}
ExecStart=${MYAPP_BIN} \$OPTIONS
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl is-active --quiet ${MYAPP_SERVICE} && systemctl stop ${MYAPP_SERVICE}
systemctl enable ${MYAPP_SERVICE}
systemctl start ${MYAPP_SERVICE}
if systemctl is-active --quiet ${MYAPP_SERVICE}; then
  printf "INFO: ${MYAPP_SERVICE} service is running\n"
else
  printf "FAIL: ${MYAPP_SERVICE} service is not running\n"
  systemctl status ${MYAPP_SERVICE}
  exit 1
fi
