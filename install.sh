#!/usr/bin/env bash
BIN_PATH_INSTALLED="/usr/local/bin/cloudflare-blocker"
BIN_PATH="./cloudflare-blocker"
CONFIG_DIR="/etc/crowdsec/cloudflare-blocker/"
PID_DIR="/var/run/crowdsec/"
SYSTEMD_PATH_FILE="/etc/systemd/system/cloudflare-blocker.service"


install_cloudflare_blocker() {
	install -v -m 755 -D "${BIN_PATH}" "${BIN_PATH_INSTALLED}"
	mkdir -p "${CONFIG_DIR}"
	cp "./config/cloudflare-blocker.yaml" "${CONFIG_DIR}cloudflare-blocker.yaml"
	CFG=${CONFIG_DIR} PID=${PID_DIR} BIN=${BIN_PATH_INSTALLED} envsubst < ./config/cloudflare-blocker.service > "${SYSTEMD_PATH_FILE}"
	systemctl daemon-reload
}


echo "Installing cloudflare-blocker"
install_cloudflare_blocker