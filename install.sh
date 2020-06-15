#!/usr/bin/env bash
BIN_PATH_INSTALLED="/usr/local/bin/cloudfare-blocker"
BIN_PATH="./cloudfare-blocker"
CONFIG_DIR="/etc/crowdsec/cloudfare-blocker/"
PID_DIR="/var/run/crowdsec/"
SYSTEMD_PATH_FILE="/etc/systemd/system/cloudfare-blocker.service"


install_cloudfare_blocker() {
	install -v -m 755 -D "${BIN_PATH}" "${BIN_PATH_INSTALLED}"
	mkdir -p "${CONFIG_DIR}"
	cp "./config/cloudfare-blocker.yaml" "${CONFIG_DIR}cloudfare-blocker.yaml"
	CFG=${CONFIG_DIR} PID=${PID_DIR} BIN=${BIN_PATH_INSTALLED} envsubst < ./config/cloudfare-blocker.service > "${SYSTEMD_PATH_FILE}"
	systemctl daemon-reload
	systemctl start cloudfare-blocker
}


echo "Installing cloudfare-blocker"
install_cloudfare_blocker