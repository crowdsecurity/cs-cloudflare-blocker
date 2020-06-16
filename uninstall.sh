#!/bin/bash

BIN_PATH_INSTALLED="/usr/local/bin/cloudflare-blocker"
CONFIG_DIR="/etc/crowdsec/cloudflare-blocker/"
PID_DIR="/var/run/crowdsec/"
SYSTEMD_PATH_FILE="/etc/systemd/system/cloudflare-blocker.service"

uninstall() {
	systemctl stop cloudflare-blocker
	rm -rf "${CONFIG_DIR}"
	rm -f "${SYSTEMD_PATH_FILE}"
	rm -f "${PID_DIR}cloudflare-blocker.pid"
	rm -f "${BIN_PATH_INSTALLED}"
}

uninstall